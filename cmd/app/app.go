package app

import (
	"log/slog"
	"obsiTeleGo/config"
	"obsiTeleGo/internal/botHandler"
	"obsiTeleGo/internal/httpServer"
	"obsiTeleGo/internal/logger"
	"obsiTeleGo/internal/rabbitmq"
	"obsiTeleGo/internal/repository"
	"os"
)

type database interface {
	Close() error
}

var LogLevel = new(slog.LevelVar)

type App struct {
	Logger     *logger.Logger
	Log        *slog.Logger
	db         database
	Repo       repository.Repo
	Rabbit     *rabbitmq.Rabbit
	BotHandler *botHandler.BotHandler
	HttpServer *httpServer.HttpServer
}

type Options struct {
	Repo string
}

func New(opt *Options) App {
	LogLevel.Set(slog.LevelDebug)

	base := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: LogLevel}))

	logger := initLogger(base)
	log := initAppLog(base)
	log.Info("Init Loggers Successfully")

	configPath := os.Getenv("CONFIG_PATH")
	if err := config.InitConfig(configPath); err != nil {
		log.Error("Init Config Error", "error", err)
		panic(err)
	}
	log.Info("Init Config Successfully")

	repo, db, err := initRepo(logger.Repo)

	if err != nil {

		log.Error("Init Repo Error", "error", err)
		panic(err)
	}
	log.Info("Init Repo Successfully")

	rabbit, err := initRabbitMQ(logger.Rabbit)
	if err != nil {
		log.Error("Init Rabbit Error", "error", err)
		panic(err)
	}
	log.Info("Init RabbitMQ Successfully")

	botHandler := initBotHandler(logger.BotHandler, repo, rabbit)
	log.Info("Init BotHandler Successfully")

	httpServer := initHttpServer(logger.HttpServer)

	return App{
		Logger:     logger,
		Log:        log,
		db:         db,
		Repo:       repo,
		BotHandler: botHandler,
		Rabbit:     rabbit,
		HttpServer: httpServer,
	}
}

func initLogger(base *slog.Logger) *logger.Logger {
	return logger.New(base)
}

func initAppLog(base *slog.Logger) *slog.Logger {
	return base.With("logger", "app")
}

func (a *App) DBClose() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func initBotHandler(log *slog.Logger, repo repository.Repo, rabbit *rabbitmq.Rabbit) *botHandler.BotHandler {
	return botHandler.New(log, repo, rabbit)
}

func initRabbitMQ(log *slog.Logger) (*rabbitmq.Rabbit, error) {

	rabbit := rabbitmq.New(log)

	err := rabbit.Ch.ExchangeDeclare(
		"tg_router",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}
	return rabbit, nil
}

func initHttpServer(log *slog.Logger) *httpServer.HttpServer {

	return httpServer.New(log, LogLevel)
}
