package config

import (
	"github.com/kaellybot/kaelly-retro-discord/commands"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	"github.com/kaellybot/kaelly-retro-discord/services/guilds"
	"github.com/kaellybot/kaelly-retro-discord/services/servers"
	"github.com/kaellybot/kaelly-retro-discord/utils/requests"
)

const (
	defaultIconSize = "128"
)

type Command struct {
	commands.AbstractCommand
	emojiService   emojis.Service
	guildService   guilds.Service
	serverService  servers.Service
	requestManager requests.RequestManager
	handlers       commands.DiscordHandlers
}
