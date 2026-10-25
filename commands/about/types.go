package about

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-retro-discord/commands"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
)

type Command struct {
	commands.AbstractCommand
	broker       amqp.MessageBroker
	emojiService emojis.Service
}
