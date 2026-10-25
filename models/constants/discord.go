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

func GetIntents() discordgo.Intent {
	return discordgo.IntentGuildMembers |
		discordgo.IntentGuilds |
		discordgo.IntentGuildMessages |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildWebhooks
}
