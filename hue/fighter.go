package hue

import (
	"Huefighter-go/config"
	"fmt"
	"github.com/amimof/huego"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Fighter() {

	var LColor []float32

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	var cfg config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	bridge := huego.New(cfg.Bridge.IP, cfg.Bridge.User)
	light, err := bridge.GetLights()
	if err != nil {
		log.Panic(err)
	}

	c, _ := colorful.Hex("#ffffff")
	cx, cy, cy2 := c.Xyy()
	LColor = append(LColor, float32(cx))
	LColor = append(LColor, float32(cy))
	log.Debug(cy2)
	LState, err := bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
		On:  true,
		Bri: 255,
		Xy:  LColor,
	})

	if err != nil {
		log.Error(err)
	}
	log.Debugln(LState)
	LColor = make([]float32, 0)

	client := twitch.NewClient(cfg.Twitch.User, cfg.Twitch.OAuth)


	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Infof("[%v, %v] %v \n", message.User.DisplayName, message.User.Color, message.Message)
		c, _ := colorful.Hex(message.User.Color)
		cx, cy, cy2 := c.Xyy()
		log.Infof("color CIE: %v, %v, %v \n", cx, cy, cy2)

		LColor = append(LColor, float32(cx))
		LColor = append(LColor, float32(cy))
		in := cfg.Bridge.LightGroup
		randomIndex := rand.Intn(len(in))
		pick := in[randomIndex]
		light, _ := strconv.Atoi(pick)
		LState, err := bridge.SetLightState(light, huego.State{
			On: true,
			// Bri:            uint8(cy2),
			Bri:            255,
			Xy:             LColor,
			TransitionTime: 1,
		})
		if err != nil {
			log.Error(err)
		}
		log.Debugln(LState)
		LColor = make([]float32, 0)
		if strings.Contains(message.Message, "!alert") {
			for i := 0; i < 15; i++ {
				_, err := bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
					On: true,
					Bri: uint8(255),
					Sat: uint8(255),
					TransitionTime: 5,
				})
				if err != nil {
					log.Warn("Could not set group state")
				}
				time.Sleep(500 * time.Millisecond)
				_, err = bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
					On: false,
					TransitionTime: 5,
				})
				if err != nil {
					log.Warn("Could not set group state")
				}
				time.Sleep(500 * time.Millisecond)
			}
			_, _ = bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
				On: true,
			})
		}
	})

	client.OnConnect(func() {
		client.Join(cfg.IRC.Channel)
		log.Infof("%s connected and monitoring %v lights", cfg.Twitch.User, len(light))
		client.Say(cfg.IRC.Channel, fmt.Sprintf("HueFighter-Go Connected and working with %v lights.", len(light)))
	})

	err = client.Connect()
	if err != nil {
		log.Panic(err)
	}
}
