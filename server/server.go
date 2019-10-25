package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/evan-buss/openbooks/irc"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ircConn is a global variable to access the current IRC connection
var ircConn *irc.Conn

// wsConn is a global variable to access the current Websocket Connection
var wsConn *websocket.Conn

// mutex ensures that only a single thread is writing to wsConn
var mutex sync.Mutex

// Start instantiates the web server and opens the browser
func Start(irc *irc.Conn) {
	ircConn = irc

	// Access the SPA bundled in the binary
	box := packr.New("ReactApp", "./app/build")

	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/ws", wsHandler)

	openbrowser("http://localhost:5228" + "/")

	log.Fatal(http.ListenAndServe(":5228", nil))
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

	wsConn = ws // Set global WS variable

	log.Println("Client connected: " + req.RemoteAddr)
	ws.SetCloseHandler(func(code int, text string) error {
		log.Println("Close Handler: " + ws.RemoteAddr().String())
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
			writeJSON(ErrorResponse{
				Error:   ERROR,
				Details: err.Error(),
			})
			return
		}
		go messageRouter(message)
	}
}

// A wrapper for websocket.WriteJSON that handles errors implicitly
func writeJSON(obj interface{}) {
	mutex.Lock()
	err := wsConn.WriteJSON(obj)
	if err != nil {
		log.Println("Error writing JSON to websocket: ", err)
	}
	mutex.Unlock()
}

// writeJSON(ServersResponse{
// 	MessageType: SERVERS,
// 	Servers:     core.Servers.Servers,
// })
