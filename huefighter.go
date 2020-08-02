package main

import (
	"Huefighter-go/config"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	var cfg config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatalf("unable to decode into struct, %v", err)
	}

	client := twitch.NewClient(cfg.Twitch.User, cfg.Twitch.OAuth)
	fmt.Println(cfg.Twitch.User)
	// client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		logrus.Infof("[%s, %s] %s \n", message.User, message.User.Color, message.Message)
	})

	client.Join("#mediamagnet")
	client.OnConnect()

	err = client.Connect()
	if err != nil {
		logrus.Panic(err)
	}
}


