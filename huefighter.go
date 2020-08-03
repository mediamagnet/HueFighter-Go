package main

import (
	"Huefighter-go/config"
	"Huefighter-go/hue"
	"flag"
	"fmt"
	"github.com/amimof/huego"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {

	listCFG := flag.Bool("list", false, "List detected Lights and exit")
	groupSet := flag.Bool("new", false, "Setup a group of lights")
	groupName := flag.String("name", "", "Name of the light group")
	groupID := flag.Int("id", 0, "ID of the Light group")

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
	flag.Parse()

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
	} else {
		hue.Fighter()
	}
}


