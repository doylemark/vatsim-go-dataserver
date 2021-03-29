package fsd

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/doylemark/vatsim-go-dataserver/log"
	"github.com/doylemark/vatsim-go-dataserver/store"
	"github.com/spf13/viper"
)

// ConnectToFSD establishes a connection to FSD over TCP socket
func ConnectToFSD(st *store.Store) {
	log.FSDLogger.Print("Connecting to Socket")

	ip := viper.GetString("fsd.url")
	conn, err := net.Dial("tcp", ip)

	if err != nil {
		log.FSDLogger.Fatal("Error connecting to FSD\n", err)
	}

	log.FSDLogger.Print("Connected to Socket")

	defer func() {
		conn.Close()
		log.FSDLogger.Print("Closing Connection")
	}()

	name := viper.GetString("fsd.name")

	_, err = conn.Write([]byte(`$"SYNC:*:` + name + `:B1:1:\r\n"`))

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
		go handlePacket(scanner.Text(), st)
	}

	fmt.Println(scanner.Err())
}

const (
	TYPE_INDEX = 8
)

func handlePacket(packet string, st *store.Store) {
	fields := strings.Split(packet, ":")
	log.FSDLogger.Println(fields)

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

	if fields[TYPE_INDEX] == "1" {
		const (
			CID_INDEX      = 5
			SERVER_INDEX   = 6
			CALLSIGN_INDEX = 7
			RATING_INDEX   = 9
			NAME_INDEX     = 11
		)

		pilot := &store.Pilot{
			Cid:      fields[CID_INDEX],
			Server:   fields[SERVER_INDEX],
			Type:     fields[TYPE_INDEX],
			Rating:   fields[RATING_INDEX],
			RealName: fields[NAME_INDEX],
			Callsign: fields[CALLSIGN_INDEX],
		}

		st.AddPilot <- pilot
	}
}

func updatePilot(fields []string, st *store.Store) {
	if len(fields) != 15 {
		return
	}

	const (
		CALLSIGN_INDEX    = 6
		TRANSPONDER_INDEX = 7
		LATITUDE_INDEX    = 9
		LONGITUDE_INDEX   = 10
		ALTITUDE_INDEX    = 11
		SPEED_INDEX       = 12
	)

	fs, errs := parseFloats([]string{fields[LATITUDE_INDEX], fields[LONGITUDE_INDEX], fields[ALTITUDE_INDEX], fields[SPEED_INDEX]})

	if len(errs) != 0 {
		return
	}

	update := &store.PositionUpdate{
		Callsign:    fields[CALLSIGN_INDEX],
		Transponder: fields[TRANSPONDER_INDEX],
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
