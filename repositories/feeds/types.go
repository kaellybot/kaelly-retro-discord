package feeds

import (
	"github.com/kaellybot/kaelly-retro-discord/models/entities"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
)

type Repository interface {
	GetFeedTypes() ([]entities.FeedType, error)
}

type Impl struct {
	db databases.MySQLConnection
}
