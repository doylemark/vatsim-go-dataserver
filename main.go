package main

import (
	"github.com/doylemark/vatsim-go-dataserver/fsd"
	"github.com/doylemark/vatsim-go-dataserver/hub"
	"github.com/doylemark/vatsim-go-dataserver/log"
	"github.com/doylemark/vatsim-go-dataserver/store"
)

func main() {
	log.InitLog()
	initConfig()

	store := store.NewStore()
	go store.Run()

	go fsd.ConnectToFSD(store)
	hub.HandleConnections(store)
}
