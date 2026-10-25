package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
)

func GetUserID(i *discordgo.Interaction) string {
	if i == nil {
		return ""
	}

	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}

	if i.User != nil {
		return i.User.ID
	}

	return ""
}

func GetMemberNickNames(s *discordgo.Session, guildID string) (map[string]any, error) {
	members, err := s.GuildMembers(guildID, "", constants.MemberListLimit)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]any)
	for _, member := range members {
		username := member.Nick
		if len(username) == 0 {
			username = member.User.Username
		}
		properties[member.User.ID] = username
	}

	return properties, nil
}
