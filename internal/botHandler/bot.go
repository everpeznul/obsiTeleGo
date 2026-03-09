package botHandler

import (
	"context"
	"log/slog"
	"obsiTeleGo/internal/rabbitmq"
	"obsiTeleGo/internal/repository"

	"github.com/go-telegram/bot"
)

type BotHandler struct {
	Log    *slog.Logger
	Repo   repository.Repo
	Rabbit *rabbitmq.Rabbit
}

func New(log *slog.Logger, repo repository.Repo, rabbit *rabbitmq.Rabbit) *BotHandler {
	return &BotHandler{
		Log:    log,
		Repo:   repo,
		Rabbit: rabbit,
	}
}

func (botHandler *BotHandler) Run(ctx context.Context, token string) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.MessageHandle),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		botHandler.Log.Error("Bot Token Error", "error", err)
		return err
	}
	botHandler.Log.Info("New Bot Successfully")

	b.RegisterHandler(bot.HandlerTypeMessageText, "/init_thread", bot.MatchTypePrefix, botHandler.InitThreadHandler)

	botHandler.Log.Info("Bot Start Successfully")
	b.Start(ctx)

	return nil
}
