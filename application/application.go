package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-retro-discord/commands"
	"github.com/kaellybot/kaelly-retro-discord/commands/about"
	"github.com/kaellybot/kaelly-retro-discord/commands/align"
	"github.com/kaellybot/kaelly-retro-discord/commands/config"
	"github.com/kaellybot/kaelly-retro-discord/commands/help"
	"github.com/kaellybot/kaelly-retro-discord/commands/job"
	"github.com/kaellybot/kaelly-retro-discord/models/constants"
	"github.com/kaellybot/kaelly-retro-discord/repositories/cities"
	emojiRepo "github.com/kaellybot/kaelly-retro-discord/repositories/emojis"
	feedRepo "github.com/kaellybot/kaelly-retro-discord/repositories/feeds"
	guildRepo "github.com/kaellybot/kaelly-retro-discord/repositories/guilds"
	"github.com/kaellybot/kaelly-retro-discord/repositories/jobs"
	"github.com/kaellybot/kaelly-retro-discord/repositories/orders"
	serverRepo "github.com/kaellybot/kaelly-retro-discord/repositories/servers"
	twitterRepo "github.com/kaellybot/kaelly-retro-discord/repositories/twitters"
	"github.com/kaellybot/kaelly-retro-discord/services/books"
	"github.com/kaellybot/kaelly-retro-discord/services/discord"
	"github.com/kaellybot/kaelly-retro-discord/services/emojis"
	"github.com/kaellybot/kaelly-retro-discord/services/feeds"
	"github.com/kaellybot/kaelly-retro-discord/services/guilds"
	"github.com/kaellybot/kaelly-retro-discord/services/servers"
	"github.com/kaellybot/kaelly-retro-discord/services/twitters"
	"github.com/kaellybot/kaelly-retro-discord/utils/databases"
	"github.com/kaellybot/kaelly-retro-discord/utils/insights"
	"github.com/kaellybot/kaelly-retro-discord/utils/requests"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//nolint:funlen,nolintlint // Yup, but much clearer like that.
func New() (*Impl, error) {
	if !viper.GetBool(constants.Production) {
		log.Info().Msgf("Development mode enabled")
	}

	// Misc
	// The handler is registered before the manager exists, but no message can be
	// consumed until Run() starts the broker, well after New() has returned.
	var requestsManager *requests.RequestManagerImpl
	broker := amqp.New(
		constants.GetRabbitMQClientID(),
		viper.GetString(constants.RabbitMQAddress),
		amqp.WithBindings(requests.GetBinding()),
		// An answer the broker refuses never reaches the consumer: without this the
		// caller would wait out the whole timeout for a reply already thrown away.
		amqp.WithRejectedHandler(func(ctx amqp.Context, message *amqp.RabbitMQMessage, err error) {
			requestsManager.Reject(ctx, message, err)
		}),
	)
	db := databases.New()
	if errDBRun := db.Run(); errDBRun != nil {
		return nil, errDBRun
	}

	// Repositories
	cityRepo := cities.New(db)
	emojiRepo := emojiRepo.New(db)
	feedRepo := feedRepo.New(db)
	guildRepo := guildRepo.New(db)
	jobRepo := jobs.New(db)
	orderRepo := orders.New(db)
	serverRepo := serverRepo.New(db)
	twitterRepo := twitterRepo.New(db)

	// Services
	serverService, errServ := servers.New(serverRepo)
	if errServ != nil {
		return nil, errServ
	}

	bookService, errBook := books.New(jobRepo, cityRepo, orderRepo)
	if errBook != nil {
		return nil, errBook
	}

	feedService, errFeed := feeds.New(feedRepo)
	if errFeed != nil {
		return nil, errFeed
	}

	twitterService, errTwitter := twitters.New(twitterRepo)
	if errTwitter != nil {
		return nil, errTwitter
	}

	emojiService, errEmoji := emojis.New(emojiRepo)
	if errEmoji != nil {
		return nil, errEmoji
	}

	guildService := guilds.New(guildRepo)
	requestsManager = requests.New(broker)

	commands := make([]commands.DiscordCommand, 0)
	commands = append(commands,
		about.New(broker, emojiService),
		align.New(bookService, guildService, serverService, emojiService, requestsManager),
		config.New(emojiService, feedService, guildService,
			serverService, twitterService, requestsManager),
		help.New(broker, &commands, emojiService),
		job.New(bookService, guildService, serverService, emojiService, requestsManager),
	)

	discordService, errDiscord := discord.New(
		viper.GetString(constants.DiscordToken),
		viper.GetInt(constants.DiscordShardID),
		viper.GetInt(constants.DiscordShardCount),
		commands, guildService, emojiService, broker)
	if errDiscord != nil {
		return nil, errDiscord
	}

	// Insights
	probes := insights.NewProbes(broker.IsConnected, db.IsConnected, discordService.IsConnected)
	prom := insights.NewPrometheusMetrics()

	return &Impl{
		broker:         broker,
		db:             db,
		probes:         probes,
		prom:           prom,
		requestManager: requestsManager,
		discordService: discordService,
	}, nil
}

func (app *Impl) Run() error {
	app.probes.ListenAndServe()
	app.prom.ListenAndServe()

	if err := app.broker.Run(); err != nil {
		return err
	}

	app.requestManager.Listen()

	if err := app.discordService.Listen(); err != nil {
		return err
	}

	return nil
}

func (app *Impl) Shutdown() {
	app.discordService.Shutdown()
	app.broker.Shutdown()
	app.db.Shutdown()
	app.prom.Shutdown()
	app.probes.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
