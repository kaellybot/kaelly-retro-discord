package discord

import (
	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-retro-discord/commands"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	"github.com/kaellybot/kaelly-retro-discord/services/guilds"
)

const (
	welcomeMessagePermissions = discordgo.PermissionViewChannel |
		discordgo.PermissionSendMessages | discordgo.PermissionEmbedLinks
)

type Service interface {
	Listen() error
	IsConnected() bool
	Shutdown()
}

type Impl struct {
	shardID      int
	session      *discordgo.Session
	commands     []commands.DiscordCommand
	guildService guilds.Service
	emojiService emojis.Service
	broker       amqp.MessageBroker
}
