package constants

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	ConfigFileName = ".env"

	// Discord Bot Token.
	DiscordToken = "DISCORD_TOKEN"

	// Shard ID. More on https://discord.com/developers/docs/topics/gateway#sharding.
	DiscordShardID = "DISCORD_SHARD_ID"

	// Total number of shards used to run the entire application.
	DiscordShardCount = "DISCORD_SHARD_COUNT"

	// MySQL URL with the following format: HOST:PORT.
	MySQLURL = "MYSQL_URL"

	// MySQL user.
	MySQLUser = "MYSQL_USER"

	// MySQL password.
	MySQLPassword = "MYSQL_PASSWORD"

	// MySQL database name.
	MySQLDatabase = "MYSQL_DATABASE"

	// RabbitMQ address.
	RabbitMQAddress = "RABBITMQ_ADDRESS"

	// Seconds to wait for a backend answer before telling the caller to retry.
	// Must stay well below Discord's 15-minute interaction token lifetime.
	RequestTimeout = "REQUEST_TIMEOUT"

	// About command ID.
	AboutID = "ABOUT_ID"

	// Align command ID.
	AlignID = "ALIGN_ID"

	// Config command ID.
	ConfigID = "CONFIG_ID"
	// Help command ID.
	HelpID = "HELP_ID"

	// Job command ID.
	JobID = "JOB_ID"

	// Probe port.
	ProbePort = "PROBE_PORT"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	// Default values.
	defaultDiscordToken      = ""
	defaultDiscordShardID    = 0
	defaultRequestTimeout    = 60
	defaultDiscordShardCount = 1
	defaultMySQLURL          = "localhost:3306"
	defaultMySQLUser         = ""
	defaultMySQLPassword     = ""
	defaultMySQLDatabase     = "kaellybot"
	defaultRabbitMQAddress   = "amqp://localhost:5672"
	defaultAboutID           = ""
	defaultAlignID           = ""
	defaultConfigID          = ""
	defaultHelpID            = ""
	defaultJobID             = ""
	defaultProbePort         = 9090
	defaultMetricPort        = 2112
	defaultLogLevel          = zerolog.InfoLevel
	defaultProduction        = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		DiscordToken:      defaultDiscordToken,
		DiscordShardID:    defaultDiscordShardID,
		RequestTimeout:    defaultRequestTimeout,
		DiscordShardCount: defaultDiscordShardCount,
		MySQLURL:          defaultMySQLURL,
		MySQLUser:         defaultMySQLUser,
		MySQLPassword:     defaultMySQLPassword,
		MySQLDatabase:     defaultMySQLDatabase,
		RabbitMQAddress:   defaultRabbitMQAddress,
		AboutID:           defaultAboutID,
		AlignID:           defaultAlignID,
		ConfigID:          defaultConfigID,
		HelpID:            defaultHelpID,
		JobID:             defaultJobID,
		ProbePort:         defaultProbePort,
		MetricPort:        defaultMetricPort,
		LogLevel:          defaultLogLevel.String(),
		Production:        defaultProduction,
	}
}

// GetRequestTimeout returns how long to wait for a backend answer before telling
// the caller to retry. A missing or non-positive value would expire every request
// the instant it is made, so it falls back on the default.
func GetRequestTimeout() time.Duration {
	seconds := viper.GetInt(RequestTimeout)
	if seconds <= 0 {
		seconds = defaultRequestTimeout
	}

	return time.Duration(seconds) * time.Second
}
