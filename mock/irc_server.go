package mock

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var logger *log.Logger

type IrcServer struct {
	Port string
}

func (irc *IrcServer) Start(ready chan<- struct{}) {
	logger = log.New(os.Stdout, "MOCK SERVER: ", 0)

	server, err := net.Listen("tcp", irc.Port)
	if err != nil {
		panic(err)
	}
	logger.Println("Listening on " + irc.Port)
	ready <- struct{}{}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	logger.Printf("Connection received from %s", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		request := scanner.Text()
		if err := scanner.Err(); err != nil {
			logger.Println(err)
		}

		logger.Printf("Request Received: %s\n", request)

		if strings.Contains(request, "@search") {
			go searchHandler(request, conn)
		}

		if strings.Contains(request, "!") {
			go downloadHandler(request, conn)
		}
	}

	logger.Println("Connection closed.")
}

func searchHandler(request string, conn net.Conn) {
	logger.Printf("Sending search results.")
	fmt.Fprint(conn, "NOTICE: Search returned 27 matches\r\n")
	fmt.Fprint(conn, ":SearchOok!ook@only.ook PRIVMSG evan_28 :DCC SEND SearchOok_results_for__the_great_gatsby.txt.zip 2130706433 6668 1184\r\n")
}

func downloadHandler(request string, conn net.Conn) {
	logger.Println("Sending book file.")
	fmt.Fprint(conn, ":SearchOok!ook@only.ook PRIVMSG evan_28 :DCC SEND great-gatsby.epub 2130706433 6669 358887\r\n")
}
