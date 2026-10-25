package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	contract "github.com/kaellybot/kaelly-commands"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	"github.com/kaellybot/kaelly-retro-discord/services/guilds"
	"github.com/kaellybot/kaelly-retro-discord/services/servers"
	"github.com/kaellybot/kaelly-retro-discord/utils/requests"
	"github.com/kaellybot/kaelly-retro-discord/utils/testutils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	guildID        = "guild-1"
	currentChannel = "chan-current"
	otherChannel   = "chan-other"
	authorID       = "author-1"
)

func TestMain(m *testing.M) {
	testutils.LoadI18n()
	m.Run()
}

// ---- service doubles ----

type findCall struct {
	name   string
	locale discordgo.Locale
	limit  int
}

type fakeServers struct {
	servers.Service
	found []entities.Server
	byID  map[string]entities.Server
	calls []findCall
}

func (f *fakeServers) FindServers(name string, locale discordgo.Locale, limit int) []entities.Server {
	f.calls = append(f.calls, findCall{name, locale, limit})
	return f.found
}

func (f *fakeServers) GetServer(id string) (entities.Server, bool) {
	server, found := f.byID[id]
	return server, found
}

type fakeGuilds struct {
	guilds.Service
}

type fakeEmojis struct {
	emojis.Service
}

func (fakeEmojis) GetMiscStringEmoji(emojiID constants.EmojiMiscID) string {
	return ":" + string(emojiID) + ":"
}

func (fakeEmojis) GetEntityStringEmoji(entityID string, _ constants.EmojiEntity) string {
	return ":" + entityID + ":"
}

type requestCall struct {
	routingKey string
	message    *amqp.RabbitMQMessage
	callback   requests.RequestCallback
	properties []map[string]any
}

type fakeRequestManager struct {
	requests.RequestManager
	err   error
	calls []requestCall
}

func (f *fakeRequestManager) Request(_ *discordgo.Session, _ *discordgo.InteractionCreate,
	routingKey string, message *amqp.RabbitMQMessage, callback requests.RequestCallback,
	properties ...map[string]any) error {
	f.calls = append(f.calls, requestCall{routingKey, message, callback, properties})
	return f.err
}

type fixture struct {
	cmd     *Command
	servers *fakeServers
	manager *fakeRequestManager
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	viper.Set(constants.ConfigID, "config-discord-id")
	t.Cleanup(func() { viper.Set(constants.ConfigID, nil) })
	f := &fixture{
		servers: &fakeServers{},
		manager: &fakeRequestManager{},
	}
	f.cmd = New(fakeEmojis{}, fakeGuilds{}, f.servers, f.manager)
	return f
}

// ---- entities ----

func server(id, frLabel, enLabel string) entities.Server {
	var labels []entities.ServerLabel
	if frLabel != "" {
		labels = append(labels, entities.ServerLabel{Locale: amqp.Language_FR, ServerID: id, Label: frLabel})
	}
	if enLabel != "" {
		labels = append(labels, entities.ServerLabel{Locale: amqp.Language_EN, ServerID: id, Label: enLabel})
	}
	return entities.Server{ID: id, Labels: labels}
}

// ---- interactions ----

func interaction(interactionType discordgo.InteractionType,
	data discordgo.InteractionData) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		ID: "interaction-id", AppID: "app", Token: "tok", Type: interactionType,
		GuildID: guildID, ChannelID: currentChannel,
		Locale: discordgo.French, Data: data,
		Member: &discordgo.Member{User: &discordgo.User{ID: authorID}},
	}}
}

func configCommand(subCommand string,
	options ...*discordgo.ApplicationCommandInteractionDataOption) *discordgo.InteractionCreate {
	return interaction(discordgo.InteractionApplicationCommand, discordgo.ApplicationCommandInteractionData{
		Name: contract.ConfigCommandName,
		Options: []*discordgo.ApplicationCommandInteractionDataOption{{
			Name: subCommand, Type: discordgo.ApplicationCommandOptionSubCommand, Options: options,
		}},
	})
}

func serverAutocomplete(
	options ...*discordgo.ApplicationCommandInteractionDataOption) *discordgo.InteractionCreate {
	event := configCommand(contract.ConfigServerSubCommandName, options...)
	event.Type = discordgo.InteractionApplicationCommandAutocomplete
	return event
}

func stringOption(name, value string) *discordgo.ApplicationCommandInteractionDataOption {
	return &discordgo.ApplicationCommandInteractionDataOption{
		Name: name, Type: discordgo.ApplicationCommandOptionString, Value: value,
	}
}

func focused(option *discordgo.ApplicationCommandInteractionDataOption,
) *discordgo.ApplicationCommandInteractionDataOption {
	option.Focused = true
	return option
}

func otherChannelOption() *discordgo.ApplicationCommandInteractionDataOption {
	return &discordgo.ApplicationCommandInteractionDataOption{
		Name: contract.ConfigChannelOptionName, Type: discordgo.ApplicationCommandOptionChannel, Value: otherChannel,
	}
}

// addGuild registers the guild and the given channels in the session state.
func addGuild(t *testing.T, s *discordgo.Session, channels ...*discordgo.Channel) {
	t.Helper()
	require.NoError(t, s.State.GuildAdd(&discordgo.Guild{ID: guildID, OwnerID: "someone-else"}))
	for _, channel := range channels {
		channel.GuildID = guildID
		require.NoError(t, s.State.ChannelAdd(channel))
	}
}

func textChannel(id string) *discordgo.Channel {
	return &discordgo.Channel{ID: id, Type: discordgo.ChannelTypeGuildText}
}

// nextRecorder is a middleware "next" capturing the context it was called with.
type nextRecorder struct {
	called bool
	ctx    context.Context
}

func (n *nextRecorder) next(ctx context.Context) error {
	n.called = true
	n.ctx = ctx
	return nil
}

func editContent(t *testing.T, fake *testutils.FakeDiscord) string {
	t.Helper()
	payload := fake.LastEdit(t)
	require.NotNil(t, payload.Content)
	return *payload.Content
}

func assertEditedWebhook(t *testing.T, fake *testutils.FakeDiscord) {
	t.Helper()
	req := fake.Last(t)
	assert.Equal(t, "PATCH", req.Method)
	assert.Equal(t, "/api/v9/webhooks/app/tok/messages/@original", req.Path)
}

// ---- Matches / descriptions ----

func TestMatches_OnlyConfigApplicationCommandsAndAutocompletes(t *testing.T) {
	cmd := newFixture(t).cmd
	other := interaction(discordgo.InteractionApplicationCommand,
		discordgo.ApplicationCommandInteractionData{Name: "job"})
	component := interaction(discordgo.InteractionMessageComponent,
		discordgo.MessageComponentInteractionData{CustomID: "/config"})

	assert.True(t, cmd.Matches(configCommand(contract.ConfigGetSubCommandName)))
	assert.True(t, cmd.Matches(serverAutocomplete()))
	assert.False(t, cmd.Matches(other))
	assert.False(t, cmd.Matches(component), "config has no message components")
	assert.Equal(t, "config", cmd.GetName())
}

func TestGetDescriptions_OneEntryPerSubCommandWithDiscordID(t *testing.T) {
	cmd := newFixture(t).cmd

	descriptions := cmd.GetDescriptions(discordgo.French)

	require.Len(t, descriptions, 2)
	for i, sub := range []string{"get", "server"} {
		assert.Equal(t, "/config "+sub, descriptions[i].Name)
		assert.Equal(t, "</config "+sub+":config-discord-id>", descriptions[i].CommandID)
		assert.Equal(t, testutils.Translations(t, discordgo.French, "config.help.tutorial."+sub)[0],
			descriptions[i].TutorialURL)
		assert.NotEmpty(t, descriptions[i].Description)
		assert.NotContains(t, descriptions[i].Description, "{{", "placeholders are resolved")
	}
	assert.Equal(t, testutils.Translations(t, discordgo.French, "config.help.detailed.get")[0],
		descriptions[0].Description)
}

func TestGetDescriptions_Localized(t *testing.T) {
	cmd := newFixture(t).cmd

	fr := cmd.GetDescriptions(discordgo.French)
	en := cmd.GetDescriptions(discordgo.EnglishGB)

	assert.Equal(t, testutils.Translations(t, discordgo.EnglishGB, "config.help.detailed.server")[0],
		en[1].Description)
	assert.NotEqual(t, fr[1].Description, en[1].Description)
}

// ---- checks ----

func TestCheckChannelID_ExplicitChannelOptionWins(t *testing.T) {
	f := newFixture(t)
	s, _ := testutils.NewDiscordSession(t)
	addGuild(t, s, textChannel(otherChannel))
	next := &nextRecorder{}

	require.NoError(t, f.cmd.checkChannelID(context.Background(), s,
		configCommand(contract.ConfigServerSubCommandName,
			stringOption(contract.ConfigServerOptionName, "Djaul"), otherChannelOption()), next.next))

	require.True(t, next.called)
	assert.Equal(t, otherChannel, next.ctx.Value(constants.ContextKeyChannel))
}

func TestCheckChannelID_ServerSubCommandWithoutChannelTargetsWholeGuild(t *testing.T) {
	f := newFixture(t)
	s, _ := testutils.NewDiscordSession(t)
	next := &nextRecorder{}

	require.NoError(t, f.cmd.checkChannelID(context.Background(), s,
		configCommand(contract.ConfigServerSubCommandName,
			stringOption(contract.ConfigServerOptionName, "Djaul")), next.next))

	require.True(t, next.called)
	assert.Nil(t, next.ctx.Value(constants.ContextKeyChannel),
		"no channel means the server is bound to the guild")
}

// ---- autocomplete ----

type autocompletePayload struct {
	Type discordgo.InteractionResponseType `json:"type"`
	Data struct {
		Choices []*discordgo.ApplicationCommandOptionChoice `json:"choices"`
	} `json:"data"`
}

func autocompleteAnswer(t *testing.T, fake *testutils.FakeDiscord) autocompletePayload {
	t.Helper()
	req := fake.Last(t)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "/api/v9/interactions/interaction-id/tok/callback", req.Path)
	var payload autocompletePayload
	req.Decode(t, &payload)
	assert.Equal(t, discordgo.InteractionApplicationCommandAutocompleteResult, payload.Type)
	return payload
}

func choiceNames(choices []*discordgo.ApplicationCommandOptionChoice) []string {
	names := make([]string, 0, len(choices))
	for _, choice := range choices {
		names = append(names, choice.Name)
		if choice.Value != choice.Name {
			names = append(names, fmt.Sprintf("value mismatch: %v", choice.Value))
		}
	}
	return names
}

func TestAutocomplete_ServersAreLocalizedAndSortedByLabel(t *testing.T) {
	f := newFixture(t)
	f.servers.found = []entities.Server{
		server("s1", "Tylézia", "Tylezia"),
		server("s2", "", "Draconiros"), // no FR label: falls back to the default locale
		server("s3", "Imagiro", "Imagiro EN"),
	}
	s, fake := testutils.NewDiscordSession(t)

	require.NoError(t, f.cmd.Handle(s, serverAutocomplete(
		focused(stringOption(contract.ConfigServerOptionName, "i")))))

	payload := autocompleteAnswer(t, fake)
	assert.Equal(t, []string{"Draconiros", "Imagiro", "Tylézia"}, choiceNames(payload.Data.Choices))
	assert.Equal(t, []findCall{{"i", discordgo.French, constants.MaxChoices}}, f.servers.calls)
}

func TestAutocomplete_UnknownOrNoFocusedOptionAnswersNoChoices(t *testing.T) {
	cases := map[string]*discordgo.InteractionCreate{
		"unknown focused option": serverAutocomplete(
			focused(stringOption("unknown", "a"))),
		"nothing focused": serverAutocomplete(
			stringOption(contract.ConfigServerOptionName, "a")),
	}
	for name, event := range cases {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			s, fake := testutils.NewDiscordSession(t)

			require.NoError(t, f.cmd.Handle(s, event))

			assert.Empty(t, autocompleteAnswer(t, fake).Data.Choices)
			assert.Empty(t, f.servers.calls)
		})
	}
}

func TestAutocomplete_DiscordFailureIsSwallowed(t *testing.T) {
	f := newFixture(t)
	s, fake := testutils.NewDiscordSession(t)
	fake.Respond(testutils.DiscordResponse{PathContains: "/callback", Status: 400, Body: `{}`})

	assert.NoError(t, f.cmd.Handle(s, serverAutocomplete(
		focused(stringOption(contract.ConfigServerOptionName, "a")))),
		"a refused autocomplete answer is logged, not reported back")
}
