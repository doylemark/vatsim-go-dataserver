package fsd

import (
	"bufio"
	"net"

	"github.com/doylemark/vatsim-go-dataserver/log"
	"github.com/spf13/viper"
)

// ConnectToFSD establishes a connection to FSD over TCP socket
func ConnectToFSD() {
	log.FSDLogger.Print("Connecting to Socket")
	conn, err := net.Dial("tcp", viper.GetString("fsd.url"))

	if err != nil {
		log.FSDLogger.Fatal("Error connecting to FSD\n", err)
	}

	log.FSDLogger.Print("Connected to Socket")

	defer func() {
		conn.Close()
		log.FSDLogger.Print("Closing Connection")
	}()

	_, err = conn.Write([]byte(`$"SYNC:*:` + viper.GetString("fsd.name") + `:B1:1:\r\n"`))

	if err != nil {
		log.FSDLogger.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		if scanner.Text() == "#You are not allowed on this port." {
			log.FSDLogger.Print("Connection Blocked on Port")
			break
		}

		log.FSDLogger.Print(scanner.Text())
	}
}
