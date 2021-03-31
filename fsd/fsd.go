package fsd

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/doylemark/vatsim-go-dataserver/log"
	"github.com/doylemark/vatsim-go-dataserver/store"
)

// ConnectToFSD establishes a connection to FSD over TCP socket
func ConnectToFSD(st *store.Store, url string, clientName string, serverName string) {
	log.FSDLogger.Print("Connecting to Socket")

	conn, err := net.Dial("tcp", url)

	if err != nil {
		log.FSDLogger.Fatal("Error connecting to FSD\n", err)
	}

	defer conn.Close()

	log.FSDLogger.Print("Connected to Socket")

	// Sync with other FSD servers
	_, err = conn.Write([]byte(`"SYNC:*:` + serverName + `:B1:1:\r\n"`))

	if err != nil {
		log.FSDLogger.Fatal(err)
	}

	go handleConn(conn, serverName, clientName)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		log.FSDLogger.Print(scanner.Text())
		go handlePacket(scanner.Text(), st, conn)
	}

	if scanner.Err() != nil {
		log.FSDLogger.Print(scanner.Err())
	}
}

func handleConn(c net.Conn, serverName string, clientName string) {
	transferCount := 0

	for {
		// must perform at least every 30s to prevent disconnection from server (can be any valid write)
		time.Sleep(time.Second * 30)
		transferCount += 1

		// request all atc data
		msg := "AD:*:" + serverName + ":B" + fmt.Sprint(transferCount) + ":1:" + clientName + ":99999:1:100:12:0.00000:0.00000:0\r\n"

		c.Write([]byte(msg))
	}
}

func handlePacket(packet string, st *store.Store, c net.Conn) {
	fields := strings.Split(packet, ":")

	switch fields[0] {
	case "ADDCLIENT":
		addClient(fields, st)
	case "RMCLIENT":
		st.RemovePilot <- fields[5]
	case "PD":
		updatePilot(fields, st)
	}
}

func addClient(fields []string, st *store.Store) {
	if len(fields) != 14 {
		return
	}

	const (
		TYPE = 8
	)

	if fields[TYPE] == "1" {

		const (
			CID      = 5
			SERVER   = 6
			CALLSIGN = 7
			RATING   = 9
			NAME     = 11
		)

		pilot := &store.Pilot{
			Cid:      fields[CID],
			Server:   fields[SERVER],
			Type:     fields[TYPE],
			Rating:   fields[RATING],
			RealName: fields[NAME],
			Callsign: fields[CALLSIGN],
		}

		st.AddPilot <- pilot
	}
}

func updatePilot(fields []string, st *store.Store) {
	if len(fields) != 15 {
		return
	}

	const (
		CALLSIGN    = 6
		TRANSPONDER = 7
		LATITUDE    = 9
		LONGITUDE   = 10
		ALTITUDE    = 11
		SPEED       = 12
	)

	fs, errs := parseFloats([]string{fields[LATITUDE], fields[LONGITUDE], fields[ALTITUDE], fields[SPEED]})

	if len(errs) != 0 {
		return
	}

	update := &store.PositionUpdate{
		Callsign:    fields[CALLSIGN],
		Transponder: fields[TRANSPONDER],
		Latitude:    fs[0],
		Longitude:   fs[1],
		Altitude:    fs[2],
		Speed:       fs[3],
	}

	st.UpdatePilot <- update
}

func parseFloats(s []string) ([]float64, []error) {
	var output []float64
	var errs []error

	for _, str := range s {
		f, err := strconv.ParseFloat(str, 64)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		output = append(output, f)
	}

	return output, errs
}

var Ratings = map[int]string{
	-1: "INAC",
	0:  "SUS",
	1:  "OBS",
	2:  "S1",
	3:  "S2",
	4:  "S3",
	5:  "C1",
	6:  "C2",
	7:  "C3",
	8:  "I1",
	9:  "I2",
	10: "I3",
	11: "SUP",
	12: "ADM",
}
