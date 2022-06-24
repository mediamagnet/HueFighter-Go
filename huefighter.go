package main

import (
	"Huefighter-go/config"
	"Huefighter-go/hue"
	"flag"
	"fmt"
	"github.com/amimof/huego"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
  "github.com/akualab/dmx"

)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {

	listCFG := flag.Bool("list", false, "List detected Lights and exit")
	groupSet := flag.Bool("new", false, "Setup a group of lights")
	groupName := flag.String("name", "", "Name of the light group")
	groupID := flag.Int("id", 0, "ID of the Light group")
	configFile := flag.String("config", "config.toml", "location of config file.")
  dmxDev := flag.String("dev", "/dev/ttyUSB0", "Location of DMX Controller.")
  dmxGroup := flag.Int("lights", 0, "How many Addressable dmx lights")
  dmxMode := flag.Bool("dmx", false, "Control DMX Lights.")
  dmxRed := flag.String("red", "1,2", "Lights that are red.")
  dmxGreen := flag.String("green", "3,4", "Lights that are green.")
  dmxBlue := flag.String("blue", "5,6", "Lights that are blue.")

	flag.Parse()
  dmxDevice := string(*dmxDev)
	configName := strings.Split(*configFile, ".")

	fmt.Println(configName)
  dmx, e := dmx.NewDMXConnection(dmxDevice)
  if e != nil {
    log.Fatal(e)
  }

  println(dmxGroup, dmxMode, dmxRed, dmxGreen, dmxBlue, dmx)

	viper.SetConfigName(configName[0])
	viper.SetConfigType("toml")
	viper.SetConfigFile(*configFile)

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


	if *listCFG {
		lightc := 0
		for i := 0; i < len(light)+5; i++ {
			light, _ := bridge.GetLight(lightc)

			log.Infof("Light ID: %v, Light Name: %v", light.ID, light.Name)
			lightc += 1

		}
		fmt.Println("Pick the lights you want to add to the group then add them to your config.toml under bridge in LightGroup.")

	} else if *groupSet {
		log.Infof("Creating group using lights: %v", cfg.Bridge.LightGroup)
		group, err := bridge.CreateGroup(huego.Group{
			Name:       *groupName,
			Lights:     cfg.Bridge.LightGroup,
			ID:         *groupID,
		})
		if err != nil {
			log.Warn(err)
		}
		log.Printf("Group %v created", group)
		fmt.Println("Please add the above group ID to your config.toml under bridge GroupNumber")
  } else if *dmxMode {
    hue.FightDMX(dmxDev)
  }else {
		hue.Fighter()
	}
}


