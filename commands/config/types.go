package config

import (
	"github.com/kaellybot/kaelly-retro-discord/commands"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	"github.com/kaellybot/kaelly-retro-discord/services/feeds"
	"github.com/kaellybot/kaelly-retro-discord/services/guilds"
	"github.com/kaellybot/kaelly-retro-discord/services/servers"
	"github.com/kaellybot/kaelly-retro-discord/services/twitters"
	"github.com/kaellybot/kaelly-retro-discord/utils/requests"
)

const (
	defaultIconSize = "128"
)

type Command struct {
	commands.AbstractCommand
	emojiService   emojis.Service
	feedService    feeds.Service
	guildService   guilds.Service
	serverService  servers.Service
	twitterService twitters.Service
	requestManager requests.RequestManager
	handlers       commands.DiscordHandlers
}
