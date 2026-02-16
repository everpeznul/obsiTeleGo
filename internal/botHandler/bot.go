package botHandler

import (
	"context"
	"log/slog"
	"obsiTeleGo/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotHandler struct {
	Log *slog.Logger
	// Queue *rabbit.Queue
	Repo repository.Repo
}

func New(log *slog.Logger, repo repository.Repo) *BotHandler {
	return &BotHandler{
		Log:  log,
		Repo: repo,
	}
}

func (botHandler *BotHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	botHandler.Log.Debug("New Message", "thread id", update.Message.MessageThreadID, "message text", update.Message.Text)

	threadName, err := botHandler.Repo.GetThreadName(ctx, int64(update.Message.MessageThreadID))
	if err != nil {
		botHandler.Log.Error("Get Thread Name Error", "error", err)
		return
	}

	botHandler.Log.Debug("Get Thread Name", "thread", threadName)
	/*
		botHandler.Queue.SendMessage(threadName, update.Message.Text)
	*/
}
