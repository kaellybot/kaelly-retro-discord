package emojis

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	repository "github.com/kaellybot/kaelly-retro-discord/repositories/emojis"
)

type Service interface {
	GetMiscStringEmoji(emojiID constants.EmojiMiscID) string
	GetMiscEmoji(emojiID constants.EmojiMiscID) *discordgo.ComponentEmoji
	GetEntityStringEmoji(entityID string, entityType constants.EmojiEntity) string
	GetEntityEmoji(entityID string, entityType constants.EmojiEntity) *discordgo.ComponentEmoji
}

type Impl struct {
	emojiStore map[constants.EmojiType]emojiStore
	repository repository.Repository
}

type emojiStore map[string]entities.Emoji
