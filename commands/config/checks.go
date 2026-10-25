package config

import (
	"context"

	"github.com/bwmarrin/discordgo"
	contract "github.com/kaellybot/kaelly-commands"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/utils/discord"
	"github.com/kaellybot/kaelly-retro-discord/utils/middlewares"
)

// checkServerChannelID binds the channel a thread belongs to: a thread has no
// configuration of its own, and that is how the binding is read back.
func (command *Command) checkServerChannelID(ctx context.Context, s *discordgo.Session,
	i *discordgo.InteractionCreate, next middlewares.NextFunc) error {
	return command.checkChannelID(ctx, s, i, func(channelCtx context.Context) error {
		channelID, ok := channelCtx.Value(constants.ContextKeyChannel).(string)
		if !ok {
			return next(channelCtx)
		}

		return next(context.WithValue(channelCtx, constants.ContextKeyChannel,
			discord.ResolveParentChannelID(s, channelID)))
	})
}

func (command *Command) checkChannelID(ctx context.Context, s *discordgo.Session,
	i *discordgo.InteractionCreate, next middlewares.NextFunc) error {
	data := i.ApplicationCommandData()
	for _, subCommand := range data.Options {
		for _, option := range subCommand.Options {
			if option.Name == contract.ConfigChannelOptionName {
				return next(context.WithValue(ctx, constants.ContextKeyChannel, option.ChannelValue(s).ID))
			}
		}
	}

	return next(ctx)
}
