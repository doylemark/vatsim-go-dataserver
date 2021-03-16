package hub

import (
	"fmt"
	"strings"
	"time"
)

type message struct {
	callsigns []string
	client    *Client
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients map[*Client][]string

	// Inbound messages from the clients.
	update chan *message

	// Register client
	register chan *Client

	// Unregister client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		update:     make(chan *message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client][]string),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = []string{""}
			go hub.manageClient(client)

		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}

		case msg := <-hub.update:
			hub.clients[msg.client] = msg.callsigns
			fmt.Println(&msg.callsigns)
		}
	}
}

func (hub *Hub) manageClient(client *Client) {
	i := 0

	for {
		if hub.clients[client] == nil {
			fmt.Println("Client Disconnected")
			break
		}

		msg := "Reporting Aircraft:" + strings.Join(hub.clients[client], " ") + "\nClients Connected:" + fmt.Sprint(len(hub.clients))

		client.send <- []byte(msg)
		i++
		time.Sleep(time.Millisecond * 500)
	}
}
