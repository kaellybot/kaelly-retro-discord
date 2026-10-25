package cities

import (
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetCities() ([]entities.City, error) {
	var cities []entities.City
	response := repo.db.GetDB().
		Model(&entities.City{}).
		Where("game = ?", constants.GetGame().AMQPGame).
		Preload("Labels").Find(&cities)
	return cities, response.Error
}
