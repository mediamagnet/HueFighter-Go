package main

import (
	"Huefighter-go/config"
	"github.com/amimof/huego"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	var LColor []float32

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

	bridge := huego.New(cfg.Bridge.IP, cfg.Bridge.User)
	light, err := bridge.GetLights()
	if err != nil {
		logrus.Panic(err)
	}
	bridge.SetLightState(3, huego.State{
		On: true,
		Bri: 255,
		Sat: 255,
		HueInc: 255,
	})

	client := twitch.NewClient(cfg.Twitch.User, cfg.Twitch.OAuth)
	// client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// logrus.Printf("[%v, %v] %v \n", message.User.DisplayName, message.User.Color, message.Message)
		c, _ := colorful.Hex(message.User.Color)
		cx, cy, cy2 := c.Xyy()
		// logrus.Printf("color CIE: %v, %v, %v \n", cx, cy, cy2)

		LColor = append(LColor, float32(cx))
		LColor = append(LColor, float32(cy))
		LState, err := bridge.SetLightState(3, huego.State{
			On:             true,
			Bri:            uint8(cy2),
			Xy:             LColor,
			TransitionTime: 1,
		})
		if err != nil {
			logrus.Error(err)
		}
		logrus.Debugln(LState)
		LColor = make([]float32, 0)
	})

	client.OnConnect(func() {
		client.Join(cfg.IRC.Channel)
		logrus.Infof("%s connected and monitoring %v lights", cfg.Twitch.User, len(light))
		// client.Say(cfg.IRC.Channel, fmt.Sprintf("HueFighter-Go Connected and working with %v lights.", len(light)))
	})


	err = client.Connect()
	if err != nil {
		logrus.Panic(err)
	}
}


