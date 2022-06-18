package mock

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type IrcServer struct {
	Port string
	log  *log.Logger
}

func (irc *IrcServer) Start(ready chan<- struct{}) {
	irc.log = log.New(os.Stdout, "MOCK SERVER: ", 0)

	server, err := net.Listen("tcp", irc.Port)
	if err != nil {
		panic(err)
	}
	irc.log.Println("Listening on " + irc.Port)
	ready <- struct{}{}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go irc.handler(conn)
	}
}

func (irc *IrcServer) handler(conn net.Conn) {
	irc.log.Printf("Connection received from %s", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)

	irc.serverHandler(conn)

	irc.sendVersionRequest(conn)

	for scanner.Scan() {
		request := scanner.Text()
		if err := scanner.Err(); err != nil {
			irc.log.Println(err)
		}

		irc.log.Printf("Request Received: %s\n", request)

		if strings.Contains(request, "@search") {
			go irc.searchHandler(request, conn)
		}

		if strings.Contains(request, "!") {
			go irc.downloadHandler(request, conn)
		}
	}

	irc.log.Println("Connection closed.")
}

func (irc *IrcServer) sendVersionRequest(conn net.Conn) {
	irc.log.Println("Sending CTCP Version inquiry.")
	fmt.Fprintf(conn, ":mock_server PRIVMSG evan_28 :\x01VERSION\x01\r\n")
}

func (irc *IrcServer) serverHandler(conn net.Conn) {
	fmt.Fprintf(conn, "353 ~DV8 ~Horla +server1 ~server2 ~evan_irc\r\n")
	fmt.Fprintf(conn, "end_list 366\r\n")
}

func (irc *IrcServer) searchHandler(request string, conn net.Conn) {
	irc.log.Printf("Sending search results.")
	fmt.Fprint(conn, "NOTICE: Search returned 27 matches\r\n")
	fmt.Fprint(conn, ":SearchOok!ook@only.ook PRIVMSG evan_28 :DCC SEND SearchOok_results_for__the_great_gatsby.txt.zip 2130706433 6668 1184\r\n")
}

func (irc *IrcServer) downloadHandler(request string, conn net.Conn) {
	irc.log.Println("Sending book file.")
	time.Sleep(time.Second * 4)
	fmt.Fprint(conn, ":SearchOok!ook@only.ook PRIVMSG evan_28 :DCC SEND great-gatsby.epub 2130706433 6669 358887\r\n")
}
