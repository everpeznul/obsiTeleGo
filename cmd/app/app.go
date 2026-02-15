package app

import (
	"log/slog"
	"obsiTeleGo/internal/logger"
	"os"
)

type App struct {
	Logger *logger.Logger
	Log    *slog.Logger
}

func New() App {
	base := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger := initLogger(base)
	log := initAppLog(base)

	return App{
		Logger: logger,
		Log:    log,
	}
}

func initLogger(base *slog.Logger) *logger.Logger {
	return logger.New(base)
}

func initAppLog(base *slog.Logger) *slog.Logger {
	return base.With("logger", "app")
}
