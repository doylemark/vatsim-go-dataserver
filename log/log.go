package log

import (
	"log"
	"os"
)

var (
	FSDLogger *log.Logger
)

// InitLog creates app loggers
func InitLog() {
	file, err := os.Create("log.txt")

	if err != nil {
		log.Print("Error creating log file \n", err)
	}

	FSDLogger = log.New(file, "[FSD]: ", log.Ldate|log.Ltime)
}
