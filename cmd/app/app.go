package app

import (
	"log/slog"
	"obsiTeleGo/internal/botHandler"
	"obsiTeleGo/internal/logger"
	"os"
)

type App struct {
	Logger     *logger.Logger
	Log        *slog.Logger
	BotHandler *botHandler.BotHandler
}

func New() App {
	base := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger := initLogger(base)
	log := initAppLog(base)
	botHandler := initBotHandler(logger.BotHandler)

	return App{
		Logger:     logger,
		Log:        log,
		BotHandler: botHandler,
	}
}

func initLogger(base *slog.Logger) *logger.Logger {
	return logger.New(base)
}

func initAppLog(base *slog.Logger) *slog.Logger {
	return base.With("logger", "app")
}

func initBotHandler(log *slog.Logger) *botHandler.BotHandler {
	return botHandler.New(log)
}
