package constants

import (
	"github.com/spf13/viper"
)

const (
	AboutRoutingKey                = "requests.about"
	AlignRequestRoutingKey         = "requests.books"
	ConfigurationRequestRoutingKey = "requests.configs"
	HelpRoutingKey                 = "requests.help"
	JobRequestRoutingKey           = "requests.books"

	GuildNewsRoutingKey = "news.guilds"
)

func GetRabbitMQClientID() string {
	return Name + "-" + viper.GetString(DiscordShardID)
}
