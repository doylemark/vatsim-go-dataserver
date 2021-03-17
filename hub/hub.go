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

// hub maintains the active set of clients and tracks their active
// aircraft
type hub struct {
	clients map[*Client][]string

	// Inbound Callsign Updates
	update chan *message

	// Register client
	register chan *Client

	// Unregister client
	unregister chan *Client
}

func newHub() *hub {
	return &hub{
		update:     make(chan *message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client][]string),
	}
}

func (hub *hub) run() {
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

func (hub *hub) manageClient(client *Client) {
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
