package constants

import (
	"github.com/bwmarrin/discordgo"
)

const (
	MaxBookRowPerEmbed = 15

	MemberListLimit = 1000

	DefaultPage = 0

	MaxChoices = 25
)

const (
	DiscordCodeNotFound        = 0     // 404: Not Found.
	DiscordCodeTooManyWebhooks = 30007 // Maximum number of webhooks reached (15).
)

func GetIntents() discordgo.Intent {
	return discordgo.IntentGuildMembers |
		discordgo.IntentGuilds |
		discordgo.IntentGuildMessages |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildWebhooks
}
