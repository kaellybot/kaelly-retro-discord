package translators

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/stretchr/testify/assert"
)

type label struct {
	locale amqp.Language
	value  string
}

func server(id string, labels ...label) entities.Server {
	server := entities.Server{ID: id}
	for _, l := range labels {
		server.Labels = append(server.Labels, entities.ServerLabel{ServerID: id, Locale: l.locale, Label: l.value})
	}
	return server
}

func TestGetEntityLabel_FallbackChain(t *testing.T) {
	frDjaul := label{amqp.Language_FR, "Djaul FR"}
	enDjaul := label{amqp.Language_EN, "Djaul EN"}
	cases := []struct {
		name     string
		entity   entities.LabelledEntity
		locale   discordgo.Locale
		expected string
	}{
		{"locale label wins", server("djaul", frDjaul, enDjaul), discordgo.French, "Djaul FR"},
		{"US English maps to EN labels", server("djaul", frDjaul, enDjaul), discordgo.EnglishUS, "Djaul EN"},
		{"missing locale falls back to default locale", server("djaul", enDjaul), discordgo.German, "Djaul EN"},
		{"unsupported locale falls back to default locale", server("djaul", frDjaul, enDjaul),
			discordgo.Japanese, "Djaul EN"},
		{"no locale nor default label falls back to ID", server("djaul", frDjaul), discordgo.German, "djaul"},
		{"no labels at all falls back to ID", server("djaul"), discordgo.French, "djaul"},
		{"an empty label is still a label", server("djaul", label{amqp.Language_FR, ""}, enDjaul),
			discordgo.French, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, GetEntityLabel(tc.entity, tc.locale))
		})
	}
}

func TestGetLabels_PreserveOrderPerEntityKind(t *testing.T) {
	fr := amqp.Language_FR
	assert.Equal(t, []string{"Bonta", "brakmar"}, GetCitiesLabels([]entities.City{
		{ID: "bonta", Labels: []entities.CityLabel{{Locale: fr, Label: "Bonta"}}},
		{ID: "brakmar"},
	}, discordgo.French))
	assert.Equal(t, []string{"Esprit"}, GetOrdersLabels([]entities.Order{
		{ID: "spirit", Labels: []entities.OrderLabel{{Locale: fr, Label: "Esprit"}}},
	}, discordgo.French))
	assert.Equal(t, []string{"Alchimiste", "Paysan"}, GetJobsLabels([]entities.Job{
		{ID: "alchemist", Labels: []entities.JobLabel{{Locale: fr, Label: "Alchimiste"}}},
		{ID: "farmer", Labels: []entities.JobLabel{{Locale: fr, Label: "Paysan"}}},
	}, discordgo.French))
	assert.Equal(t, []string{"Djaul", "Imagiro"}, GetServersLabels([]entities.Server{
		server("djaul", label{fr, "Djaul"}),
		server("imagiro", label{fr, "Imagiro"}),
	}, discordgo.French))
}

func TestGetLabels_EmptyInputGivesEmptyNonNilSlice(t *testing.T) {
	for name, labels := range map[string][]string{
		"cities":  GetCitiesLabels(nil, discordgo.French),
		"orders":  GetOrdersLabels(nil, discordgo.French),
		"jobs":    GetJobsLabels(nil, discordgo.French),
		"servers": GetServersLabels(nil, discordgo.French),
	} {
		assert.NotNil(t, labels, name)
		assert.Empty(t, labels, name)
	}
}
