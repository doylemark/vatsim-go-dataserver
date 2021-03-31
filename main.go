package main

import (
	"github.com/doylemark/vatsim-go-dataserver/fsd"
	"github.com/doylemark/vatsim-go-dataserver/hub"
	"github.com/doylemark/vatsim-go-dataserver/log"
	"github.com/doylemark/vatsim-go-dataserver/store"
)

func main() {
	log.InitLog()
	url, clientName, serverName := initConfig()

	store := store.NewStore()
	go store.Run()

	go fsd.ConnectToFSD(store, url, clientName, serverName)
	hub.HandleConnections(store)
}
