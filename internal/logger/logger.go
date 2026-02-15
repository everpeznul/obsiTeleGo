package logger

import (
	"log/slog"
)

type Logger struct {
	Bot      *slog.Logger
	Domain   *slog.Logger
	Obsidian *slog.Logger
	Rabbit   *slog.Logger
	Repo     *slog.Logger
}

func New(base *slog.Logger) *Logger {
	return &Logger{
		Bot:      base.With("logger", "bot"),
		Domain:   base.With("logger", "domain"),
		Obsidian: base.With("logger", "obsidian"),
		Rabbit:   base.With("logger", "rabbit"),
		Repo:     base.With("logger", "repo"),
	}
}
