package main

import (
	"log"

	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Couldn't read config file\n", err)
	}
}

func initDefaults() {
	viper.SetDefault("fsd.url", "")
	viper.SetDefault("fsd.name", "")
}
