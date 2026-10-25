package mappers

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func channel(id string) *discordgo.Channel { return &discordgo.Channel{ID: id} }

func configServers() fakeServers {
	return fakeServers{servers: []entities.Server{testServer()}}
}

func TestMapConfigurationRequests_CopyFields(t *testing.T) {
	get := MapConfigurationGetRequest("guild", "author", discordgo.French)
	assert.Equal(t, amqp.RabbitMQMessage_CONFIGURATION_GET_REQUEST, get.Type)
	assert.Equal(t, amqp.Language_FR, get.Language)
	assert.Equal(t, "author", get.UserID)
	require.NotNil(t, get.ConfigurationGetRequest)
	assert.Equal(t, "guild", get.ConfigurationGetRequest.GuildId)

	server := MapConfigurationServerRequest("guild", "chan", "srv", "author", discordgo.French)
	assert.Equal(t, amqp.RabbitMQMessage_CONFIGURATION_SET_SERVER_REQUEST, server.Type)
	assert.Equal(t, "author", server.UserID)
	require.NotNil(t, server.ConfigurationSetServerRequest)
	assert.Equal(t, "guild", server.ConfigurationSetServerRequest.GuildId)
	assert.Equal(t, "chan", server.ConfigurationSetServerRequest.ChannelId)
	assert.Equal(t, "srv", server.ConfigurationSetServerRequest.ServerId)
}

func TestMapConfigToEmbed_ListsGuildServerAndChannels(t *testing.T) {
	serverService := configServers()
	guild := constants.GuildConfig{
		Name:     "Les Fous",
		Icon:     "guild.png",
		ServerID: "srv",
		ServerChannels: []constants.ServerChannel{
			{Channel: channel("c1"), ServerID: "srv"},
			{Channel: channel("c2"), ServerID: "ghost-srv"},
		},
	}

	embed := MapConfigToEmbed(guild, fakeEmojis{}, serverService, amqp.Language_EN)

	assert.Equal(t, "Les Fous", embed.Title)
	assert.Equal(t, "guild.png", embed.Thumbnail.URL)
	assert.Equal(t, constants.Color, embed.Color)
	assertDefaultFooter(t, embed.Footer)
	assertRendered(t, embed.Description)
	assert.Contains(t, embed.Description, "<Server:srv>")
	assert.Contains(t, embed.Description, "**Imagiro**")

	require.Len(t, embed.Fields, 1)
	serverField := embed.Fields[0]
	assert.Contains(t, serverField.Name, "<misc:dofusGame>")
	assertInOrder(t, serverField.Value,
		"<#c1> <Server:srv> **Imagiro**",
		"<#c2> <Server:ghost-srv> **ghost-srv**")
	assert.False(t, serverField.Inline)
}

func TestMapConfigToEmbed_NoServerNorChannels(t *testing.T) {
	serverService := configServers()

	embed := MapConfigToEmbed(constants.GuildConfig{Name: "Empty"}, fakeEmojis{},
		serverService, amqp.Language_FR)

	assertRendered(t, embed.Description)
	assert.Contains(t, embed.Description, constants.GetGame().Name)
	assert.NotContains(t, embed.Description, "<Server:")
	for _, field := range embed.Fields {
		assertRendered(t, field.Value)
		assert.NotContains(t, field.Value, "<#")
	}
}

func TestMapConfigToEmbed_UnknownGuildServerFallsBackToItsID(t *testing.T) {
	serverService := configServers()

	embed := MapConfigToEmbed(constants.GuildConfig{ServerID: "ghost"}, fakeEmojis{},
		serverService, amqp.Language_EN)

	assert.Contains(t, embed.Description, "<Server:ghost> **ghost**")
}
