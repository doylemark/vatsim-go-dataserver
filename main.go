package main

import (
	"github.com/doylemark/vatsim-go-dataserver/fsd"
	"github.com/doylemark/vatsim-go-dataserver/hub"
	"github.com/doylemark/vatsim-go-dataserver/log"
)

func main() {
	log.InitLog()
	initConfig()
	go fsd.ConnectToFSD()
	hub.HandleConnections()
}
