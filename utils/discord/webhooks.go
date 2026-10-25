package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	i18n "github.com/kaysoro/discordgo-i18n"
)

func GetPaginationButtons(page, pages int, crafter CraftPageCustomID,
	lg discordgo.Locale, emojiService emojis.Service) []discordgo.MessageComponent {
	lastPage := pages - 1
	if pages == 0 || constants.DefaultPage == lastPage {
		return []discordgo.MessageComponent{}
	}

	// Stale buttons may carry a page outside of the current range.
	page = min(max(page, constants.DefaultPage), lastPage)

	previousPage := max(page-1, constants.DefaultPage)

	nextPage := min(page+1, lastPage)

	buttons := make([]discordgo.MessageComponent, 0)
	if previousPage > constants.DefaultPage {
		buttons = append(buttons, discordgo.Button{
			CustomID: crafter(previousPage),
			Label:    i18n.Get(lg, "default.page.previous"),
			Style:    discordgo.SecondaryButton,
			Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDPrevious),
		})
	} else {
		buttons = append(buttons, discordgo.Button{
			CustomID: crafter(constants.DefaultPage),
			Label:    i18n.Get(lg, "default.page.first"),
			Style:    discordgo.SecondaryButton,
			Disabled: page <= constants.DefaultPage,
			Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDFirst),
		})
	}

	if nextPage < lastPage {
		buttons = append(buttons, discordgo.Button{
			CustomID: crafter(nextPage),
			Label:    i18n.Get(lg, "default.page.next"),
			Style:    discordgo.SecondaryButton,
			Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDNext),
		})
	} else {
		buttons = append(buttons, discordgo.Button{
			CustomID: crafter(lastPage),
			Label:    i18n.Get(lg, "default.page.last"),
			Style:    discordgo.SecondaryButton,
			Disabled: page >= lastPage,
			Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDLast),
		})
	}

	return buttons
}

func BuildDefaultFooter(lg discordgo.Locale) *discordgo.MessageEmbedFooter {
	return &discordgo.MessageEmbedFooter{
		Text: i18n.Get(lg, "default.footer", i18n.Vars{
			"name":      constants.Name,
			"version":   constants.Version,
			"changelog": i18n.Get(lg, "default.changelog"),
		}),
		IconURL: constants.AvatarIcon,
	}
}
