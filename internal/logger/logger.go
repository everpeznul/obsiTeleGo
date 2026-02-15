package logger

import (
	"log/slog"
)

type Logger struct {
	BotHandler *slog.Logger
	Domain     *slog.Logger
	Obsidian   *slog.Logger
	Rabbit     *slog.Logger
	Repo       *slog.Logger
}

func New(base *slog.Logger) *Logger {
	return &Logger{
		BotHandler: base.With("logger", "botHandler"),
		Domain:     base.With("logger", "domain"),
		Obsidian:   base.With("logger", "obsidian"),
		Rabbit:     base.With("logger", "rabbit"),
		Repo:       base.With("logger", "repo"),
	}
}
