package irc

import (
	"log"
	"net"
)

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

// Connect connects to the given server at port 6667
func (i *Conn) Connect(address string) {
	conn, err := net.Dial("tcp", address+":6667")

	if err != nil {
		log.Fatal("IRC Connection Error", err)
	}

	i.Conn = conn

	user := "USER " + i.Username + " " + i.Username + " " + i.Username + " :" + i.realname + "\r\n"
	nick := "NICK " + i.Username + "\r\n"

	i.Write([]byte(user))
	i.Write([]byte(nick))
}

// Disconnect closes connection to the IRC server
func (i *Conn) Disconnect() {
	i.Write([]byte("QUIT :Goodbye\r\n"))
	i.Conn.Close()
}

// SendMessage sends the given message string to the connected IRC server
func (i *Conn) SendMessage(message string) {
	i.Write([]byte("PRIVMSG #" + i.channel + " :" + message + "\r\n"))
}

// JoinChannel joins the channel given by channel string
func (i *Conn) JoinChannel(channel string) {
	i.channel = channel
	_, err := i.Write([]byte("JOIN #" + channel + "\r\n"))
	if err != nil {
		log.Fatal(err)
	}
}

// GetUsers sends a NAMES request to the IRC server
func (i *Conn) GetUsers(channel string) {
	i.Write([]byte("NAMES #" + channel + "\r\n"))
}

// PONG sends a PONG message to the server, often used after a PING request
func (i *Conn) PONG(server string) {
	i.Write([]byte("PONG " + server + "\r\n"))
}

// IsConnected returns true if the IRC connection is not null
func (i *Conn) IsConnected() bool {
	return i.Conn != nil
}
