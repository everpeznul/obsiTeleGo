package logger

import (
	"log/slog"
)

type Logger struct {
	HttpServer *slog.Logger
	BotHandler *slog.Logger
	Domain     *slog.Logger
	Obsidian   *slog.Logger
	Rabbit     *slog.Logger
	Repo       *slog.Logger
}

func New(base *slog.Logger) *Logger {
	return &Logger{
		HttpServer: base.With("logger", "httpServer"),
		BotHandler: base.With("logger", "botHandler"),
		Domain:     base.With("logger", "domain"),
		Obsidian:   base.With("logger", "obsidian"),
		Rabbit:     base.With("logger", "rabbit"),
		Repo:       base.With("logger", "repo"),
	}
}
