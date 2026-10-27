package twitters

import (
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
)

type Repository interface {
	GetTwitterAccounts() ([]entities.TwitterAccount, error)
}

type Impl struct {
	db databases.MySQLConnection
}
