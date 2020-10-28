package server

import (
	"log"
	"net/http"
	"sync"
	"fmt"

	"github.com/evan-buss/openbooks/irc"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/websocket"

	"strconv"
	"strings"
	"time"
	"os"
	"path/filepath"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TODO: Maybe use a struct instead of global variables.
// ircConn is a global variable to access the current IRC connection
var ircConn *irc.Conn

//hacky method of seeing numberofconnections in order to close
//irc session when not client is connected
var numberOfConnections = 0

var timerForIrcDeath *time.Timer
var timeTillDeath = 120;

// wsConn is a global variable to access the current Websocket Connection
var wsConn *websocket.Conn

// mutex ensures that only a single thread is writing to wsConn
var mutex sync.Mutex

// Start instantiates the web server and opens the browser
func Start(irc *irc.Conn, port string) {
	ircConn = irc

	// Access the SPA bundled in the binary
	box := packr.New("ReactApp", "./app/build")
	ebookFolder := http.Dir(filepath.Join(getProgramDir(), "downloadedEbooks"))

	http.Handle("/", http.FileServer(box))
	http.Handle("/attic/", http.StripPrefix("/attic/", http.FileServer(ebookFolder))) //if stripprefix wasn't used the file server would try to find the book in $PROGRAMDIR/downloadedEbooks/attic
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/reset", resetIRC)

	//openbrowser("http://127.0.0.1:" + port + "/")

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

	wsConn = ws // Set global WS variable

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
	mutex.Lock()
	err := wsConn.WriteJSON(obj)
	if err != nil {
		log.Println("Error writing JSON to websocket: ", err)
	}
	mutex.Unlock()
}

func resetIRC(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Disconnecting from irc server")
	ircConn.ChangeState(false)
	ircConn.Disconnect()
}

func getProgramDir() string{
	ex, err := os.Executable()
    if err != nil {
        panic(err)
	}
	//fmt.Println("The program dir is: "+filepath.Dir(ex))
    return(filepath.Dir(ex))
}
