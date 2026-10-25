package discord

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/utils/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testutils.LoadI18n()
	m.Run()
}

func TestGetUserID_PrefersGuildMemberOverDirectUser(t *testing.T) {
	cases := []struct {
		name        string
		interaction *discordgo.Interaction
		expected    string
	}{
		{"nil interaction", nil, ""},
		{"guild member", &discordgo.Interaction{
			Member: &discordgo.Member{User: &discordgo.User{ID: "member"}},
			User:   &discordgo.User{ID: "user"},
		}, "member"},
		{"direct message user", &discordgo.Interaction{User: &discordgo.User{ID: "user"}}, "user"},
		{"member without user falls back to user", &discordgo.Interaction{
			Member: &discordgo.Member{}, User: &discordgo.User{ID: "user"},
		}, "user"},
		{"nobody", &discordgo.Interaction{Member: &discordgo.Member{}}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, GetUserID(tc.interaction))
		})
	}
}

func TestGetMemberNickNames_NicknameOverridesUsername(t *testing.T) {
	s, fake := testutils.NewDiscordSession(t)
	fake.Respond(testutils.DiscordResponse{Method: http.MethodGet, PathContains: "/guilds/g1/members",
		Status: http.StatusOK, Body: `[
			{"nick": "Kay", "user": {"id": "u1", "username": "kaysoro"}},
			{"nick": "", "user": {"id": "u2", "username": "elycann"}},
			{"user": {"id": "u3", "username": "colibry"}}
		]`})

	names, err := GetMemberNickNames(s, "g1")

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"u1": "Kay", "u2": "elycann", "u3": "colibry"}, names)
	req := fake.Last(t)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/api/v"+discordgo.APIVersion+"/guilds/g1/members", req.Path)
}

func TestGetMemberNickNames_RequestsUpToMemberListLimit(t *testing.T) {
	s, fake := testutils.NewDiscordSession(t)
	var query string
	s.Client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		query = req.URL.RawQuery
		return fake.RoundTrip(req)
	})
	fake.Respond(testutils.DiscordResponse{PathContains: "/members", Status: http.StatusOK, Body: `[]`})

	names, err := GetMemberNickNames(s, "g1")

	require.NoError(t, err)
	assert.Empty(t, names)
	assert.Equal(t, "limit="+strconv.Itoa(constants.MemberListLimit), query)
}

func TestGetMemberNickNames_DiscordErrorIsReturned(t *testing.T) {
	s, fake := testutils.NewDiscordSession(t)
	fake.Respond(testutils.DiscordResponse{PathContains: "/members", Status: http.StatusForbidden,
		Body: `{"code": 50013, "message": "Missing Permissions"}`})

	names, err := GetMemberNickNames(s, "g1")

	var restErr *discordgo.RESTError
	require.ErrorAs(t, err, &restErr)
	assert.Equal(t, http.StatusForbidden, restErr.Response.StatusCode)
	assert.Nil(t, names)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
