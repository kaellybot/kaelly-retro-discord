package mappers

import (
	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/i18n"
)

func requestBackbone(authorID string, requestType amqp.RabbitMQMessage_Type,
	lg discordgo.Locale) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     requestType,
		Language: i18n.MapDiscordLocale(lg),
		Game:     constants.GetGame().AMQPGame,
		UserID:   authorID,
	}
}
