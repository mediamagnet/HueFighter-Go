package config

type Configuration struct {
	Twitch TwitchConfiguration
	IRC    IRCConfiguration
	Bridge BridgeConfiguration
}
