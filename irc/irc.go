package irc

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
)

var ErrTLSHandshake = errors.New("TLS handshake failed. The server may not support TLS")

// Conn represents an IRC connection to a server
type Conn struct {
	net.Conn
	channel  string
	Username string
	realname string
}

// New creates a new IRC connection to the server using the supplied username and realname
func New(username, realname string) *Conn {
	irc := &Conn{
		channel:  "",
		Username: username,
		realname: realname,
	}

	return irc
}

func (i *Conn) Connect(address string, enableTLS bool) error {
	// If switching servers, we will already have an active connection.
	// A single IRC{} instance is used per client.
	if i.IsConnected() {
		i.Disconnect()
		slog.Info("Disconnected from IRC server.")
	}

	var conn net.Conn
	var err error
	if enableTLS {
		conn, err = tls.Dial("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = net.Dial("tcp", address)
	}

	if err != nil {
		if err.Error() == "tls: first record does not look like a TLS handshake" {
			return ErrTLSHandshake
		}
		return err
	}

	i.Conn = conn

	user := "USER " + i.Username + " 0 * :" + i.Username + "\r\n"
	nick := "NICK " + i.Username + "\r\n"

	i.Write([]byte(user))
	i.Write([]byte(nick))
	return nil
}

// Disconnect closes connection to the IRC server
func (i *Conn) Disconnect() {
	if !i.IsConnected() {
		return
	}
	i.Write([]byte("QUIT :Goodbye\r\n"))
	i.Conn.Close()
}

// SendMessage sends the given message string to the connected IRC server
func (i *Conn) SendMessage(message string) {
	if !i.IsConnected() {
		return
	}
	i.Write([]byte("PRIVMSG #" + i.channel + " :" + message + "\r\n"))
}

// SendNotice sends a notice message to the specified user
func (i *Conn) SendNotice(user string, message string) {
	if !i.IsConnected() {
		return
	}
	i.Write([]byte("NOTICE " + user + " :" + message + "\r\n"))
}

// JoinChannel joins the channel given by channel string
func (i *Conn) JoinChannel(channel string) {
	if !i.IsConnected() {
		return
	}
	i.channel = channel
	i.Write([]byte("JOIN #" + channel + "\r\n"))
}

// GetUsers sends a NAMES request to the IRC server
func (i *Conn) GetUsers(channel string) {
	if !i.IsConnected() {
		return
	}
	i.Write([]byte("NAMES #" + channel + "\r\n"))
}

// Pong sends a Pong message to the server, often used after a PING request
func (i *Conn) Pong(server string) {
	if !i.IsConnected() {
		return
	}
	i.Write([]byte("PONG " + server + "\r\n"))
}

// IsConnected returns true if the IRC connection is not null
func (i *Conn) IsConnected() bool {
	return i.Conn != nil
}
