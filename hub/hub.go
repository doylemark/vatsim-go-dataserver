package hub

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/doylemark/vatsim-go-dataserver/store"
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

func (hub *hub) run(st *store.Store) {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = []string{}
			go hub.manageClient(client, st)

		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}

		case msg := <-hub.update:
			hub.clients[msg.client] = msg.callsigns
		}
	}
}

func (hub *hub) manageClient(client *Client, st *store.Store) {
	for {
		if hub.clients[client] == nil {
			fmt.Println("Client Disconnected")
			break
		}

		var output []store.Pilot

		for _, callsign := range hub.clients[client] {
			if _, ok := st.Pilots[callsign]; ok {
				output = append(output, st.Pilots[callsign])
			}
		}

		data, err := json.Marshal(output)

		if err != nil {
			fmt.Println(err)
		}

		client.send <- []byte(data)
		time.Sleep(time.Second * 5)
	}
}
