package constants

import "github.com/bwmarrin/discordgo"

type GuildConfig struct {
	Name           string
	Icon           string
	ServerID       string
	ServerChannels []ServerChannel
}

type ServerChannel struct {
	Channel  *discordgo.Channel
	ServerID string
}
