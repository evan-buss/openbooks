package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/util"
	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// Unique ID for the client
	uuid uuid.UUID

	sse *sse.Server

	// Signal to indicate the connection should be terminated.
	// disconnect chan struct{}

	// Message to send to the client ws connection
	send chan interface{}

	// Individual IRC connection per connected client.
	irc *irc.Conn

	log *log.Logger

	// Context is used to signal when this client should close.
	ctx context.Context
}

// TODO: we don't really need clients anymore as SSE handles it for us and we aren't allowing
// multiple IRC connections simulatenously.
// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (server *server) writePump(c *Client) {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				continue
			}

			byteMessage, err := json.Marshal(message)
			if err != nil {
				c.log.Printf("Error marshalling message to JSON: %s\n", err)
				return
			}
			c.sse.Publish(c.uuid.String(), &sse.Event{
				Data: byteMessage,
			})
		case <-c.ctx.Done():
			return
		}
	}
}

func (server *server) connectIRC(client *Client) {
	// The IRC connection could be re-used if it is already connected
	if !client.irc.IsConnected() {
		err := core.Join(client.irc, server.config.Server, server.config.EnableTLS)
		if err != nil {
			client.log.Println(err)
			client.send <- newErrorResponse("Unable to connect to IRC server.")
			return
		}

		handler := server.NewIrcEventHandler(client)

		if server.config.Log {
			logger, _, err := util.CreateLogFile(client.irc.Username, server.config.DownloadDir)
			if err != nil {
				server.log.Println(err)
			}
			handler[core.Message] = func(text string) { logger.Println(text) }
		}

		go core.StartReader(client.ctx, client.irc, handler)

		client.send <- ConnectionResponse{
			StatusResponse: StatusResponse{
				MessageType:      CONNECT,
				NotificationType: SUCCESS,
				Title:            "Welcome, connection established.",
				Detail:           fmt.Sprintf("IRC username %s", client.irc.Username),
			},
			Name: client.irc.Username,
		}

		return
	}

	client.send <- ConnectionResponse{
		StatusResponse: StatusResponse{
			MessageType:      CONNECT,
			NotificationType: NOTIFY,
			Title:            "Welcome back, re-using open IRC connection.",
			Detail:           fmt.Sprintf("IRC username %s", client.irc.Username),
		},
		Name: client.irc.Username,
	}
}
