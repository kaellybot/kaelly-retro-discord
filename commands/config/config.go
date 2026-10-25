package config

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	contract "github.com/kaellybot/kaelly-commands"
	"github.com/kaellybot/kaelly-retro-discord/commands"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	"github.com/kaellybot/kaelly-retro-discord/services/guilds"
	"github.com/kaellybot/kaelly-retro-discord/services/servers"
	"github.com/kaellybot/kaelly-retro-discord/utils/checks"
	"github.com/kaellybot/kaelly-retro-discord/utils/middlewares"
	"github.com/kaellybot/kaelly-retro-discord/utils/requests"
	di18n "github.com/kaysoro/discordgo-i18n"
	"github.com/spf13/viper"
)

//nolint:exhaustive // only the useful interaction types are handled.
func New(emojiService emojis.Service, guildService guilds.Service, serverService servers.Service,
	requestManager requests.RequestManager) *Command {
	cmd := Command{
		DiscordID:      viper.GetString(constants.ConfigID),
		emojiService:   emojiService,
		guildService:   guildService,
		serverService:  serverService,
		requestManager: requestManager,
	}

	checkServer := checks.CheckServer(contract.ConfigServerOptionName, cmd.serverService)

	subCommandHandlers := cmd.HandleSubCommands(commands.SubCommandHandlers{
		contract.ConfigGetSubCommandName: middlewares.
			Use(cmd.getRequest),
		contract.ConfigServerSubCommandName: middlewares.
			Use(checkServer, cmd.checkServerChannelID, cmd.serverRequest),
	})

	cmd.handlers = commands.DiscordHandlers{
		discordgo.InteractionApplicationCommand:             subCommandHandlers,
		discordgo.InteractionApplicationCommandAutocomplete: cmd.autocomplete,
	}

	return &cmd
}

func (command *Command) GetName() string {
	return contract.ConfigCommandName
}

func (command *Command) GetDescriptions(lg discordgo.Locale) []commands.Description {
	return []commands.Description{
		{
			Name:        fmt.Sprintf("/%v get", contract.ConfigCommandName),
			CommandID:   fmt.Sprintf("</%v get:%v>", contract.ConfigCommandName, command.DiscordID),
			Description: di18n.Get(lg, fmt.Sprintf("%v.help.detailed.get", contract.ConfigCommandName)),
			TutorialURL: di18n.Get(lg, fmt.Sprintf("%v.help.tutorial.get", contract.ConfigCommandName)),
		},
		{
			Name:        fmt.Sprintf("/%v server", contract.ConfigCommandName),
			CommandID:   fmt.Sprintf("</%v server:%v>", contract.ConfigCommandName, command.DiscordID),
			Description: di18n.Get(lg, fmt.Sprintf("%v.help.detailed.server", contract.ConfigCommandName)),
			TutorialURL: di18n.Get(lg, fmt.Sprintf("%v.help.tutorial.server", contract.ConfigCommandName)),
		},
	}
}

func (command *Command) Matches(i *discordgo.InteractionCreate) bool {
	return commands.IsApplicationCommand(i) &&
		i.ApplicationCommandData().Name == command.GetName()
}

func (command *Command) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return command.CallHandler(s, i, command.handlers)
}

func getServerOptions(ctx context.Context) (entities.Server, string, error) {
	server, ok := ctx.Value(constants.ContextKeyServer).(entities.Server)
	if !ok {
		return entities.Server{}, "",
			fmt.Errorf("cannot cast %v as entities.Server", ctx.Value(constants.ContextKeyServer))
	}

	channelID := ""
	if ctx.Value(constants.ContextKeyChannel) != nil {
		channelID, ok = ctx.Value(constants.ContextKeyChannel).(string)
		if !ok {
			return entities.Server{}, "",
				fmt.Errorf("cannot cast %v as string", ctx.Value(constants.ContextKeyChannel))
		}
	}

	return server, channelID, nil
}
