package orders

import (
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetOrders() ([]entities.Order, error) {
	var orders []entities.Order
	response := repo.db.GetDB().
		Model(&entities.Order{}).
		Where("game = ?", constants.GetGame().AMQPGame).
		Preload("Labels").
		Find(&orders)
	return orders, response.Error
}
