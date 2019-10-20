package server

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/evan-buss/openbooks/irc"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Start instantiates the web server and opens the browser
func Start(irc *irc.Conn, port int) {
	// This should ideally be called from the server
	connect(irc)

	box := packr.New("ReactApp", "./app/build")

	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/ws", wsEndpoint)

	openbrowser("http://localhost:" + strconv.Itoa(port) + "/")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

// wsEndpoint handles upgrading and connecting websocket requests
func wsEndpoint(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(req *http.Request) bool {
		return true
	}

	// upgrade the connection into a websocket
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client connected: " + req.RemoteAddr)
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println(ws.RemoteAddr().String())
		return nil
	})

	reader(ws)
}

// Reader listens to all incoming messages and delegates them to
// the response handler if there is not an error parsing the json
func reader(conn *websocket.Conn) {
	for {
		var message Request
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			conn.WriteJSON(ErrorMessage{
				Error:   ERROR,
				Details: err.Error(),
			})
			return
		}

		go asyncResponse(conn, message)
	}
}

// asyncResponse handles the given JSON messages via the messageRouter
// The messageRouter writes the results
func asyncResponse(conn *websocket.Conn, message Request) {
	result, err := messageRouter(message)
	if err != nil {
		log.Println(err)
		// Write the error to the socket if bad
		conn.WriteJSON(ErrorMessage{
			Error:   message.RequestType,
			Details: err.Error(),
		})
		return
	}

	if err := conn.WriteJSON(result); err != nil {
		log.Println(err)
		return
	}
}

func connect(irc *irc.Conn) {
	irc.Connect("irc.irchighway.net")
	// Wait before joining the ebooks room
	// Often you recieve a private message from the server
	time.Sleep(time.Second * 2)
	irc.JoinChannel("ebooks")
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
