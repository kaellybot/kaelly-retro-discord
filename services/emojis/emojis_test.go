package emojis

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	repository "github.com/kaellybot/kaelly-retro-discord/repositories/emojis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEmojiRepo struct {
	repository.Repository
	emojis []entities.Emoji
	err    error
}

func (repo *fakeEmojiRepo) GetEmojis() ([]entities.Emoji, error) {
	return repo.emojis, repo.err
}

// custom is a Discord application emoji: rendered as <:name:snowflake>.
func custom(emojiType constants.EmojiType, id string) entities.Emoji {
	return entities.Emoji{
		ID: id, Type: emojiType,
		DiscordName: "d_" + id, Snowflake: "sf_" + id, Name: "n_" + id,
	}
}

// unicode is a plain unicode emoji: only Name is set.
func unicode(emojiType constants.EmojiType, id, name string) entities.Emoji {
	return entities.Emoji{ID: id, Type: emojiType, Name: name}
}

func newService(t *testing.T, emojis ...entities.Emoji) *Impl {
	t.Helper()
	service, err := New(&fakeEmojiRepo{emojis: emojis})
	require.NoError(t, err)
	return service
}

func TestNew_PropagatesRepositoryError(t *testing.T) {
	boom := errors.New("db down")
	service, err := New(&fakeEmojiRepo{err: boom})
	require.ErrorIs(t, err, boom)
	assert.Nil(t, service)
}

func TestMapEmojiString(t *testing.T) {
	cases := map[string]struct {
		emoji    *entities.Emoji
		expected string
	}{
		"nil": {nil, ""},
		"custom emoji": {
			&entities.Emoji{ID: "kama", DiscordName: "kama", Snowflake: "123", Name: "💰"}, "<:kama:123>",
		},
		"no snowflake uses name":  {&entities.Emoji{ID: "kama", DiscordName: "kama", Name: "💰"}, "💰"},
		"no discord name":         {&entities.Emoji{ID: "kama", Snowflake: "123", Name: "💰"}, "💰"},
		"blank id uses name":      {&entities.Emoji{ID: "  ", DiscordName: "kama", Snowflake: "123", Name: "💰"}, "💰"},
		"nothing usable is empty": {&entities.Emoji{ID: "kama"}, ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, mapEmojiString(tc.emoji))
		})
	}
}

func TestMapEmoji_UsesSnowflakeAndNameNotDiscordName(t *testing.T) {
	assert.Nil(t, mapEmoji(nil))
	assert.Equal(t, &discordgo.ComponentEmoji{ID: "123", Name: "💰"},
		mapEmoji(&entities.Emoji{ID: "kama", DiscordName: "kama", Snowflake: "123", Name: "💰"}))
}

func TestMiscEmojis(t *testing.T) {
	empty := newService(t)
	assert.Empty(t, empty.GetMiscStringEmoji(constants.EmojiIDGlobal), "no misc store at all")
	assert.Nil(t, empty.GetMiscEmoji(constants.EmojiIDGlobal))

	service := newService(t,
		custom(constants.EmojiTypeMisc, string(constants.EmojiIDGlobal)),
		unicode(constants.EmojiTypeMisc, string(constants.EmojiIDScroll), "📅"),
		// Same id under another type must not leak into misc lookups.
		custom(constants.EmojiType(constants.EmojiTypeJob), string(constants.EmojiIDGithub)),
	)
	assert.Equal(t, "<:d_global:sf_global>", service.GetMiscStringEmoji(constants.EmojiIDGlobal))
	assert.Equal(t, &discordgo.ComponentEmoji{ID: "sf_global", Name: "n_global"},
		service.GetMiscEmoji(constants.EmojiIDGlobal))
	assert.Equal(t, "📅", service.GetMiscStringEmoji(constants.EmojiIDScroll))
	assert.Equal(t, &discordgo.ComponentEmoji{Name: "📅"}, service.GetMiscEmoji(constants.EmojiIDScroll))

	assert.Empty(t, service.GetMiscStringEmoji(constants.EmojiIDGithub), "misc store exists, id missing")
	assert.Nil(t, service.GetMiscEmoji(constants.EmojiIDGithub))
}

func TestEntityEmojis_AreScopedByEntityType(t *testing.T) {
	empty := newService(t)
	assert.Empty(t, empty.GetEntityStringEmoji("bonta", constants.EmojiTypeCity))
	assert.Nil(t, empty.GetEntityEmoji("bonta", constants.EmojiTypeCity))

	service := newService(t,
		custom(constants.EmojiType(constants.EmojiTypeCity), "bonta"),
		unicode(constants.EmojiType(constants.EmojiTypeJob), "bonta", "🪓"),
	)
	assert.Equal(t, "<:d_bonta:sf_bonta>", service.GetEntityStringEmoji("bonta", constants.EmojiTypeCity))
	assert.Equal(t, &discordgo.ComponentEmoji{ID: "sf_bonta", Name: "n_bonta"},
		service.GetEntityEmoji("bonta", constants.EmojiTypeCity))
	assert.Equal(t, "🪓", service.GetEntityStringEmoji("bonta", constants.EmojiTypeJob),
		"same id, different entity type, different emoji")
	assert.Equal(t, &discordgo.ComponentEmoji{Name: "🪓"}, service.GetEntityEmoji("bonta", constants.EmojiTypeJob))

	assert.Empty(t, service.GetEntityStringEmoji("brakmar", constants.EmojiTypeCity))
	assert.Nil(t, service.GetEntityEmoji("brakmar", constants.EmojiTypeCity))
	assert.Empty(t, service.GetEntityStringEmoji("bonta", constants.EmojiTypeOrder))
	assert.Nil(t, service.GetEntityEmoji("bonta", constants.EmojiTypeOrder))
}
