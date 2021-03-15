package main

import (
	"bufio"
	"log"
	"net"

	"github.com/spf13/viper"
)

func connectToServer() {
	log.Print("[FSD]: Connecting to Socket")
	conn, err := net.Dial("tcp", viper.GetString("fsd.url"))

	if err != nil {
		log.Fatal("Error connecting to FSD\n", err)
	}

	log.Print("[FSD]: Connected to Socket")

	defer func() {
		conn.Close()
		log.Print("[FSD]: Closing Connection")
	}()

	_, err = conn.Write([]byte(`$"SYNC:*:` + viper.GetString("fsd.name") + `:B1:1:\r\n"`))

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		if scanner.Text() == "#You are not allowed on this port." {
			log.Print("[FSD]: Connection Blocked on Port")
			break
		}

		log.Print(scanner.Text())
	}

	if scanner.Err() != nil {
		log.Print(scanner.Err().Error())
	}
}
