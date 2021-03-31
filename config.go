package main

import (
	"log"

	"github.com/spf13/viper"
)

func initConfig() (string, string, string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Couldn't read config file\n", err)
	}

	url := viper.GetString("fsd.url")
	clientName := viper.GetString("fsd.clientName")
	serverName := viper.GetString("fsd.serverName")

	return url, clientName, serverName
}
