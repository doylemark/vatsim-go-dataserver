package main

import (
	"bufio"
	"fmt"
	"io"
)

func main() {
	InitConfig()
	conn := ConnectToServer()

	for {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')

		if err == io.EOF {
			conn.Close()
			return
		}

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(msg))
	}

}
