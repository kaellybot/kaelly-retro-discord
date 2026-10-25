package job

import (
	"github.com/kaellybot/kaelly-retro-discord/commands"
	"github.com/kaellybot/kaelly-retro-discord/services/books"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	"github.com/kaellybot/kaelly-retro-discord/services/guilds"
	"github.com/kaellybot/kaelly-retro-discord/services/servers"
	"github.com/kaellybot/kaelly-retro-discord/utils/requests"
)

const (
	userProperty = "user"
)

type Command struct {
	commands.AbstractCommand
	bookService    books.Service
	emojiService   emojis.Service
	guildService   guilds.Service
	serverService  servers.Service
	requestManager requests.RequestManager
	slashHandlers  commands.DiscordHandlers
	userHandlers   commands.DiscordHandlers
}
