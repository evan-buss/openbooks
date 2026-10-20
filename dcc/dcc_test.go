package dcc

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/evan-buss/openbooks/mock"
)

// TestStringParsing makes sure that data is properly extracted from the DCC
// response string. (filename, IP conversion, port, and size)
func TestStringParsing(t *testing.T) {
	tables := []struct {
		search   string
		download *Download
	}{
		{
			":SearchOok!ook@only.ook PRIVMSG evan_28 :DCC SEND SearchOok_results_for__hp_lovecraft.txt.zip 1543751478 2043 784",
			&Download{Filename: "SearchOok_results_for__hp_lovecraft.txt.zip", IP: "92.3.199.54", Port: "2043", Size: 784},
		},
		{
			":Search!Search@ihw-4q5hcb.dyn.suddenlink.net PRIVMSG evan_bot :DCC SEND SearchBot_results_for__stephen_king_the_stand.txt.zip 2907707975 4342 1116",
			&Download{Filename: "SearchBot_results_for__stephen_king_the_stand.txt.zip", IP: "173.80.26.71", Port: "4342", Size: 1116},
		},
		{
			`:DV8!HandyAndy@ihw-39fkft.ip-164-132-173.eu PRIVMSG negative-bishop-1 :DCC SEND "Douglas Adams - [HITCHHIKER'S GUIDE TO THE GALAXY & THE 01] - Hitchhiker's Guide to the Galaxy & The (v5.0) (EPUB).rar" 2760158537 2050 2321788`,
			&Download{Filename: "Douglas Adams - [HITCHHIKER'S GUIDE TO THE GALAXY & THE 01] - Hitchhiker's Guide to the Galaxy & The (v5.0) (EPUB).rar", IP: "164.132.173.73", Port: "2050", Size: 2321788},
		},
	}

	for _, table := range tables {
		download, err := ParseString(table.search)
		require.NoError(t, err)
		assert.Equal(t, table.download, download)
	}
}

func TestDownload(t *testing.T) {
	text := "Test dcc download content."

	textDownload := Download{
		Filename: "test.txt",
		IP:       "localhost",
		Port:     "6969",
		Size:     int64(len(text)),
	}

	reader := bytes.NewReader([]byte(text))
	server := mock.DccServer{
		Port:   ":" + textDownload.Port,
		Reader: reader,
	}

	ready := make(chan struct{}, 1)
	go server.Start(ready)
	<-ready

	t.Log("After server start")

	received := new(mock.WriteCloser)
	err := textDownload.Download(received)
	require.NoError(t, err)
	assert.Equal(t, text, string(received.Data))
}

// TestDownloadAcknowledgesSenderWindow reproduces the real world stall where a
// DCC sender with a limited window stops permanently if the receiver never
// acknowledges the bytes it has received. Without acknowledgements the transfer
// hangs once the window is full.
func TestDownloadAcknowledgesSenderWindow(t *testing.T) {
	const window = 100 * 1024
	const total = window*2 + 512

	data := make([]byte, total)
	for i := range data {
		data[i] = byte(i % 251)
	}

	textDownload := Download{
		Filename: "windowed.bin",
		IP:       "localhost",
		Port:     "6971",
		Size:     int64(total),
	}

	server := mock.AckDccServer{
		Port:       ":" + textDownload.Port,
		Data:       data,
		WindowSize: window,
	}

	ready := make(chan struct{}, 1)
	go server.Start(ready)
	<-ready

	received := new(mock.WriteCloser)

	done := make(chan error, 1)
	go func() {
		done <- textDownload.Download(received)
	}()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		t.Fatal("download stalled: receiver never acknowledged the sender window")
	}

	assert.Equal(t, total, len(received.Data))
	assert.True(t, bytes.Equal(data, received.Data))
}

// scriptedReceiver plays back a fixed sequence of Read results so the transfer
// loop can be tested without a real socket.
type scriptedReceiver struct {
	frames []frame
	index  int
	acks   bytes.Buffer
}

type frame struct {
	data []byte
	err  error
}

func (s *scriptedReceiver) Read(p []byte) (int, error) {
	if s.index >= len(s.frames) {
		return 0, io.EOF
	}
	f := s.frames[s.index]
	s.index++
	return copy(p, f.data), f.err
}

func (s *scriptedReceiver) Write(p []byte) (int, error) { return s.acks.Write(p) }
func (s *scriptedReceiver) SetReadDeadline(time.Time) error {
	return nil
}

// TestReceiveKeepsDataReturnedWithError covers the case where the sender resets
// the connection in the same Read that carries the final bytes. The tail must
// not be dropped; the transfer is complete.
func TestReceiveKeepsDataReturnedWithError(t *testing.T) {
	payload := bytes.Repeat([]byte("A"), 5000)

	conn := &scriptedReceiver{frames: []frame{
		{data: payload[:4096]},
		{data: payload[4096:], err: io.ErrUnexpectedEOF},
	}}

	var out bytes.Buffer
	n, err := receive(conn, &out, int64(len(payload)))

	require.NoError(t, err)
	assert.EqualValues(t, len(payload), n)
	assert.True(t, bytes.Equal(payload, out.Bytes()))
}

// TestReceiveReportsIncompleteTransfer makes sure a genuine mid-transfer error
// is still surfaced together with the byte count that actually arrived.
func TestReceiveReportsIncompleteTransfer(t *testing.T) {
	conn := &scriptedReceiver{frames: []frame{
		{data: bytes.Repeat([]byte("B"), 1000), err: io.ErrUnexpectedEOF},
	}}

	var out bytes.Buffer
	n, err := receive(conn, &out, 5000)

	require.Error(t, err)
	assert.EqualValues(t, 1000, n)
}
