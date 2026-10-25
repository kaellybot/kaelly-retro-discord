package constants

import "github.com/rs/zerolog"

const (
	LogCommand         = "command"
	LogCorrelationID   = "correlationID"
	LogGame            = "game"
	LogEvent           = "event"
	LogCommandOption   = "option"
	LogGuildCount      = "guildCount"
	LogGuildID         = "guildID"
	LogChannelID       = "channelID"
	LogEntity          = "entity"
	LogEntityCount     = "entityCount"
	LogEmojiType       = "emojiType"
	LogInteractionType = "interactionType"
	LogShard           = "shard"
	LogFileName        = "fileName"
	LogLocale          = "locale"
	LogPanic           = "panic"

	LogLevelFallback = zerolog.InfoLevel
)
