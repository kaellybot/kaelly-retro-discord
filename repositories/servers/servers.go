package servers

import (
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetServers() ([]entities.Server, error) {
	var servers []entities.Server
	response := repo.db.GetDB().
		Model(&entities.Server{}).
		Where("game = ?", constants.GetGame().AMQPGame).
		Preload("Labels").
		Find(&servers)
	return servers, response.Error
}
