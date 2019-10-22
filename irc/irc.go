package irc

import (
	"log"
	"net"
	"strings"
)

// Conn represents an IRC connection to a server
type Conn struct {
	irc      net.Conn
	channel  string
	username string
	realname string
	Logging  bool
}

// New creates a new IRC connection to the server using the supplied username and realname
func New(username, realname string) *Conn {
	irc := Conn{}
	irc.channel = ""
	irc.username = username
	irc.realname = realname

	return &irc
}

// Connect connects to the given server at port 6667
func (i *Conn) Connect(address string) {
	conn, err := net.Dial("tcp", address+":6667")

	if err != nil {
		log.Fatal("IRC Connection Error", err)
	}
	i.irc = conn

	user := "USER " + i.username + " " + i.username + " " + i.username + " :" + i.realname + "\r\n"
	nick := "NICK " + i.username + "\r\n"

	i.irc.Write([]byte(user))
	i.irc.Write([]byte(nick))
}

// Disconnect closes connection to the IRC server
func (i *Conn) Disconnect() {
	i.irc.Write([]byte("QUIT :Goodbye\r\n"))
}

// SendMessage sends the given message string to the connected IRC server
func (i *Conn) SendMessage(message string) {
	i.irc.Write([]byte("PRIVMSG #" + i.channel + " :" + message + "\r\n"))
}

// JoinChannel joins the channel given by channel
// NOTE: You must explicitly join a channel before you can start sending messages
func (i *Conn) JoinChannel(channel string) {
	i.channel = channel
	_, err := i.irc.Write([]byte("JOIN #" + channel + "\r\n"))
	if err != nil {
		log.Fatal(err)
	}
}

// GetMessage listens to the IRC connection and returns any message found
func (i *Conn) GetMessage() (text string) {
	buf := make([]byte, 512)
	n, err := i.irc.Read(buf)

	if err != nil {
		log.Println("Probably a connection error.")
		log.Fatal("Get Message Error: ", err)
	}

	text = string(buf[:n])

	if strings.Contains(text, "PING") {
		i.irc.Write([]byte("PONG :pingis\r\n"))
	}
	return
}

func (i *Conn) GetUsers(channel string) {
	i.irc.Write([]byte("NAMES #" + channel + "\r\n"))
}

func (i *Conn) IsConnected() bool {
	return i.irc != nil
}
