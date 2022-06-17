package hue

import (
	"Huefighter-go/config"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/amimof/huego"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Fighter does the thing that makes the colors change
func FightDMX() {

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
	// light, err := bridge.GetLights()
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

	// TODO: moderator role and broadcaster role add for commands
	// TODO: Change connect message to group instead of total

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Infof("[%v, %v] %v \n", message.User.DisplayName, message.User.Color, message.Message)
		// log.Infoln(message.User.Badges)
		c, _ := colorful.Hex(message.User.Color)
		/* w.SetContent(widget.NewVBox(
			cwidget,
			widget.NewLabel(message.User.Color),
		)) */
		cx, cy, cy2 := c.Xyy()
		log.Debugf("color CIE: %v, %v, %v \n", cx, cy, cy2)

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
		moderator := message.User.Badges["moderator"]
		broadcaster := message.User.Badges["broadcaster"]
		msgary := strings.Split(message.Message, " ")
		if moderator == 1 || broadcaster == 1 {
			switch {
			case strings.Contains(msgary[0], "!alert"):
				for i := 0; i < 15; i++ {
					_, err := bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
						On:             true,
						Bri:            uint8(255),
						Sat:            uint8(255),
						TransitionTime: 5,
					})
					if err != nil {
						log.Warn("Could not set group state")
					}
					time.Sleep(500 * time.Millisecond)
					_, err = bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
						On:             false,
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
			case strings.Contains(msgary[0], "!reset"):
				var white []float32
				white = append(white, 0.34510, 0.35811)
				log.Printf("Resetting")
				time.Sleep(1 * time.Second)
				_, err = bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
					On:             true,
					Bri:            255,
					Xy:             white,
					TransitionTime: 0,
				})
				if err != nil {
					log.Warn("Unable to set light state.")
				}
			case strings.Contains(msgary[0], "!lightson"):
				log.Info("Turning lights on")
				_, err = bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
					On:             true,
					TransitionTime: 1,
				})
				if err != nil {
					log.Warn("Unable to set group state")
				}
			case strings.Contains(msgary[0], "!lightsoff"):
				log.Info("Turning lights off")
				_, err = bridge.SetGroupState(cfg.Bridge.GroupNumber, huego.State{
					On:             false,
					TransitionTime: 1,
				})
				if err != nil {
					log.Warn("Unable to set group state")
				}
			}
		}
		// group, _ := bridge.GetGroup(cfg.Bridge.GroupNumber)
		// log.Printf("HSB: %v %v %v", group.State.Hue, group.State.Sat, group.State.Bri)
		// testHCL := colorful.Hcl(float64(group.State.Hue), float64(group.State.Sat), float64(group.State.Bri))
		// log.Printf("HCL: %v", testHCL)

	})

	client.OnConnect(func() {
		client.Join(cfg.IRC.Channel)
		log.Infof("%s connected and monitoring %v lights", cfg.Twitch.User, len(cfg.Bridge.LightGroup))
		client.Say(cfg.IRC.Channel, fmt.Sprintf("HueFighter-Go Connected and working with %v lights.", len(cfg.Bridge.LightGroup)))
	})

	err = client.Connect()
	if err != nil {
		log.Panic(err)
	}
}
