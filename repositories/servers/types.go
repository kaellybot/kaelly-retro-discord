package servers

import (
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
)

type Repository interface {
	GetServers() ([]entities.Server, error)
}

type Impl struct {
	db databases.MySQLConnection
}
