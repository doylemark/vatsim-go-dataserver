package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type msg struct {
	User string `json:"user"`
	Text string `json:"text"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func startHub() {
	http.HandleFunc("/", handleConnection)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("New Client Connected")

	for i := 0; i < 100; i++ {
		ws.WriteJSON(msg{User: "Mark", Text: fmt.Sprint(i + 1)})
		time.Sleep(time.Millisecond * 500)
	}

	if err != nil {
		log.Println(err)
	}
}
