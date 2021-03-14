package main

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"
)

// ConnectToServer connects to VATSIM FSD over TCP
func ConnectToServer() (conn net.Conn) {
	fmt.Println("Connecting to FSD 🚀\n", viper.GetString("FSD_URL"))

	conn, err := net.Dial("tcp", viper.GetString("fsd.url"))

	if err != nil {
		log.Fatal("Error connecting to FSD\n", err)
	}

	_, err = conn.Write([]byte(`$"SYNC:*:` + viper.GetString("fsd.name") + `:B1:1:\r\n"`))

	if err != nil {
		log.Fatal(err)
	}

	return conn
}
