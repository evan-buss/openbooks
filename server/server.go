package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/evan-buss/openbooks/irc"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// IRC is a global variable to access the current IRC connection
var IRC *irc.Conn

// WS is a global variable to access the current Websocket Connection
var WS *websocket.Conn

// Start instantiates the web server and opens the browser
func Start(irc *irc.Conn, port int) {
	IRC = irc

	// Access the SPA bundled in the binary
	box := packr.New("ReactApp", "./app/build")

	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/ws", wsHandler)

	openbrowser("http://localhost:" + strconv.Itoa(port) + "/")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

// wsHandler handles upgrading and connecting websocket requests
func wsHandler(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(req *http.Request) bool {
		return true
	}

	// upgrade the connection into a websocket
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	WS = ws // Set global WS variable

	log.Println("Client connected: " + req.RemoteAddr)
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println(ws.RemoteAddr().String())
		return nil
	})

	reader(ws)
}

// reader listens to all incoming messages on an websocket connection and delegates them or response with error
func reader(conn *websocket.Conn) {
	for {
		var message Request
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			err := conn.WriteJSON(ErrorResponse{
				Error:   ERROR,
				Details: err.Error(),
			})
			if err != nil {
				log.Println("Error sending error JSON:", err)
			}
			return
		}

		go asyncResponse(conn, message)
	}
}

// asyncResponse handles the given JSON messages via the messageRouter or sends and error
func asyncResponse(conn *websocket.Conn, message Request) {
	result, err := messageRouter(message)
	if err != nil {
		log.Println(err)
		// Write the error to the socket if bad
		err := conn.WriteJSON(ErrorResponse{
			Error:   message.RequestType,
			Details: err.Error(),
		})
		if err != nil {
			log.Println("asyncResponse error: ", err)
		}
		return
	}

	if err := conn.WriteJSON(result); err != nil {
		log.Println(err)
		return
	}
}
