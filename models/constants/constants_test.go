package constants

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRabbitMQClientID_SuffixedWithShardID(t *testing.T) {
	previous := viper.Get(DiscordShardID)
	t.Cleanup(func() { viper.Set(DiscordShardID, previous) })

	viper.Set(DiscordShardID, 3)
	assert.Equal(t, "KaellyRetro-3", GetRabbitMQClientID())

	viper.Set(DiscordShardID, nil)
	assert.Equal(t, "KaellyRetro-", GetRabbitMQClientID())
}

func TestGetDefaultConfigValues_CoversEveryConfigKeyWithUsableDefaults(t *testing.T) {
	defaults := GetDefaultConfigValues()

	for _, key := range []string{DiscordToken, DiscordShardID, DiscordShardCount, MySQLURL, MySQLUser,
		MySQLPassword, MySQLDatabase, RabbitMQAddress, AboutID, AlignID, ConfigID, HelpID,
		JobID, ProbePort, MetricPort, LogLevel, Production, RequestTimeout} {
		assert.Contains(t, defaults, key)
	}
	assert.Len(t, defaults, 18)

	assert.Equal(t, 1, defaults[DiscordShardCount], "a single shard by default")
	assert.Equal(t, false, defaults[Production], "never register commands globally by accident")
	assert.NotEqual(t, defaults[ProbePort], defaults[MetricPort], "both servers bind at once")
	assert.Equal(t, 60, defaults[RequestTimeout], "well below Discord's 15-minute token lifetime")

	rawLevel, ok := defaults[LogLevel].(string)
	require.True(t, ok)
	level, err := zerolog.ParseLevel(rawLevel)
	require.NoError(t, err)
	assert.Equal(t, zerolog.InfoLevel, level)
}

func TestGetIntents_NoPrivilegedMessageContent(t *testing.T) {
	intents := GetIntents()

	assert.Zero(t, intents&discordgo.IntentMessageContent, "message content intent was removed on purpose")
	for _, required := range []discordgo.Intent{discordgo.IntentGuildMembers, discordgo.IntentGuilds,
		discordgo.IntentGuildMessages, discordgo.IntentGuildMessageReactions, discordgo.IntentGuildWebhooks} {
		assert.Equal(t, required, intents&required)
	}
}

func TestGetGame_IsDofusRetro(t *testing.T) {
	game := GetGame()
	assert.Equal(t, "DOFUS Retro", game.Name)
	assert.NotEmpty(t, game.Icon)
	assert.EqualValues(t, "DOFUS_RETRO", game.AMQPGame.String())
}

func TestGraphists_HaveNameAndLink(t *testing.T) {
	for _, graphist := range []Graphist{GetGraphistElycann(), GetGraphistColibry()} {
		assert.NotEmpty(t, graphist.Name)
		assert.Regexp(t, `^https://`, graphist.URL)
	}
}

// A missing or non-positive timeout would expire every request the instant it is
// made, leaving every command answering "try again later".
func TestGetRequestTimeout_FallsBackOnUnusableValues(t *testing.T) {
	previous := viper.Get(RequestTimeout)
	t.Cleanup(func() { viper.Set(RequestTimeout, previous) })

	for _, value := range []any{nil, 0, -1} {
		viper.Set(RequestTimeout, value)
		assert.Equal(t, 60*time.Second, GetRequestTimeout())
	}

	viper.Set(RequestTimeout, 5)
	assert.Equal(t, 5*time.Second, GetRequestTimeout())
}
