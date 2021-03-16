package main

import "github.com/doylemark/vatsim-go-dataserver/hub"

func main() {
	initLog()
	initConfig()
	go connectToFSD()
	hub.HandleConnections()
}
