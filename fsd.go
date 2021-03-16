package main

import (
	"bufio"
	"log"
	"net"

	"github.com/spf13/viper"
)

func connectToFSD() {
	FSDLogger.Print("Connecting to Socket")
	conn, err := net.Dial("tcp", viper.GetString("fsd.url"))

	if err != nil {
		log.Fatal("Error connecting to FSD\n", err)
	}

	FSDLogger.Print("Connected to Socket")

	defer func() {
		conn.Close()
		FSDLogger.Print("Closing Connection")
	}()

	_, err = conn.Write([]byte(`$"SYNC:*:` + viper.GetString("fsd.name") + `:B1:1:\r\n"`))

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		if scanner.Text() == "#You are not allowed on this port." {
			FSDLogger.Print("Connection Blocked on Port")
			break
		}

		FSDLogger.Print(scanner.Text())
	}
}
