package main

import (
	"strings"

	"github.com/spf13/viper"
)

// InitConfig initialises environment variables
func InitConfig() {
	initDefaults()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func initDefaults() {
	viper.SetDefault("fsd.url", "")
}
