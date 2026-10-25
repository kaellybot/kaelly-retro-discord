package emojis

import (
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
)

type Repository interface {
	GetEmojis() ([]entities.Emoji, error)
}

type Impl struct {
	db databases.MySQLConnection
}
