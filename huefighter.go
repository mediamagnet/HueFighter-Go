package main

import (
	"Huefighter-go/config"
	"github.com/jrm780/gotirc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	var cfg config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatal("unable to decode into struct, %v", err)
	}

}
