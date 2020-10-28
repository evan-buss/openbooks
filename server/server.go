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

// Connection represents a connection to the websocket client and the irc server
type Connection struct {
	irc *irc.Conn
	ws  *websocket.Conn
	sync.Mutex
}

// Conn is the current connection between the browser
// and server and server and irc channel
var Conn Connection

//hacky method of seeing numberofconnections in order to close
//irc session when not client is connected
var numberOfConnections = 0
var timerForIrcDeath *time.Timer
var timeTillDeath = 120;

// Start instantiates the web server and opens the browser
func Start(irc *irc.Conn, port string) {

	Conn = Connection{
		irc: irc,
	}

	// Access the SPA bundled in the binary
	box := packr.New("ReactApp", "./app/build")

	http.Handle("/", http.FileServer(box))
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

	numberOfConnections++
	log.Println(strconv.Itoa(numberOfConnections)+" clients are connected!")
	if (numberOfConnections == 1 && timerForIrcDeath != nil){
		timerForIrcDeath.Stop()
		log.Println("The IRC connection timer has been stopped, yay!")
	}
	reader(ws)
}

// reader listens to all incoming messages on an websocket connection and delegates them or response with error
func reader(conn *websocket.Conn) {
	for {
		var message Request
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			if (strings.Contains(err.Error(), "close") && strings.Contains(err.Error(), "websocket")){
				numberOfConnections--
				log.Println("Client disconnected! "+strconv.Itoa(numberOfConnections)+" remaining!")
				if (numberOfConnections == 0){
					log.Println("All clients disconnected, waiting for "+strconv.Itoa(timeTillDeath)+" seconds before closing irc connection :(")
					timerForIrcDeath = time.AfterFunc(time.Second * time.Duration(timeTillDeath), func() {
						log.Println("Death.")
						ircConn.ChangeState(false)
						ircConn.Disconnect()
					})
					//timerForIrcDeath
					//ircCon.changeState(false)
					//here is where we do the disconnecting jazz
				}
			}

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

func resetIRC(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Disconnecting from irc server")
	ircConn.ChangeState(false)
	ircConn.Disconnect()
}
