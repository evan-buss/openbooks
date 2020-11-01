package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/evan-buss/openbooks/irc"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"

	// Load the static content
	_ "github.com/evan-buss/openbooks/server/statik"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Connection represents a connection to the websocket client and the irc server
type Connection struct {
	irc *irc.Conn
	ws  *websocket.Conn
	sync.Mutex
}

// Conn is the current connection between the browser
// and server and server and irc channel
var Conn Connection

// Start instantiates the web server and opens the browser
func Start(irc *irc.Conn, port string) {

	Conn = Connection{
		irc: irc,
	}

	staticFs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(staticFs))
	http.HandleFunc("/ws", wsHandler)

	openbrowser("http://127.0.0.1:" + port + "/")

	log.Fatal(http.ListenAndServe(":"+port, nil))
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

	Conn.ws = ws // Set global WS variable

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

// A wrapper for websocket.WriteJSON that ensures a single writer
// and handles errors
func writeJSON(obj interface{}) {
	Conn.Lock()
	err := Conn.ws.WriteJSON(obj)
	if err != nil {
		log.Println("Error writing JSON to websocket: ", err)
	}
	Conn.Unlock()
}
