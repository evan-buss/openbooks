package dcc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
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
			&Download{Filename: "SearchOok_results_for__hp_lovecraft.txt.zip", IP: net.ParseIP("92.3.199.54").To4(), Port: 2043, Size: 784},
		},
		{
			":Search!Search@ihw-4q5hcb.dyn.suddenlink.net PRIVMSG evan_bot :DCC SEND SearchBot_results_for__stephen_king_the_stand.txt.zip 2907707975 4342 1116",
			&Download{Filename: "SearchBot_results_for__stephen_king_the_stand.txt.zip", IP: net.ParseIP("173.80.26.71").To4(), Port: 4342, Size: 1116},
		},
		{
			`:DV8!HandyAndy@ihw-39fkft.ip-164-132-173.eu PRIVMSG negative-bishop-1 :DCC SEND "Douglas Adams - [HITCHHIKER'S GUIDE TO THE GALAXY & THE 01] - Hitchhiker's Guide to the Galaxy & The (v5.0) (EPUB).rar" 2760158537 2050 2321788`,
			&Download{Filename: "Douglas Adams - [HITCHHIKER'S GUIDE TO THE GALAXY & THE 01] - Hitchhiker's Guide to the Galaxy & The (v5.0) (EPUB).rar", IP: net.ParseIP("164.132.173.73").To4(), Port: 2050, Size: 2321788},
		},
	}

	for _, table := range tables {
		download, err := ParseString(table.search)
		require.NoError(t, err)
		assert.Equal(t, table.download, download)
	}
}

func TestIpConversion(t *testing.T) {
	dccStr := delimiterChar + `DCC SEND "SearchBot_results_for__stephen_king_the_stand.txt.zip" 3232235777 4342 1116`

	download, err := ParseString(dccStr)
	if err != nil {
		t.Error(err)
	}

	if download.Filename != "SearchBot_results_for__stephen_king_the_stand.txt.zip" {
		t.Error("invalid filename")
	}

	if download.IP.String() != "192.168.1.1" {
		t.Errorf("invalid ip got: %s expected: %s", download.IP.To4(), net.ParseIP("192.168.1.1").To4())
	}

	if download.Port != 4342 {
		t.Error("invalid port")
	}

	if download.Size != 1116 {
		t.Error("invalid size")
	}

	if download.String() != dccStr {
		t.Errorf("invalid String() got: %s expected %s", download.String(), dccStr)
	}
}
