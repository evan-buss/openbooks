package main

import (
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			conn.WriteJSON(ErrorMessage{
				Error:   -1,
				Details: err.Error(),
			})
			return
		}

		go asyncResponse(conn, message)
	}
}

// asyncResponse handles the given JSON messages via the messageRouter
// The messageRouter writes the results
func asyncResponse(conn *websocket.Conn, message Message) {
	result, err := messageRouter(message)
	if err != nil {
		log.Println(err)
		// Write the error to the socket if bad
		conn.WriteJSON(ErrorMessage{
			Error:   message.MessageType,
			Details: err.Error(),
		})
		return
	}

	if err := conn.WriteJSON(result); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	box := packr.New("ReactApp", "./web/build")

	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/ws", wsEndpoint)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
