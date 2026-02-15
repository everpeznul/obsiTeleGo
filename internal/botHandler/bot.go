package botHandler

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotHandler struct {
	Log *slog.Logger
	// Queue *rabbit.Queue
	// Repo *repository.Repo
}

func New(log *slog.Logger) *BotHandler {
	return &BotHandler{
		Log: log,
	}
}

func (botHandler *BotHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	botHandler.Log.Debug("New Message", "thread id", update.Message.MessageThreadID, "message text", update.Message.Text)

	/*threadName := repo.getThreadName(ctx, update.Message.Chat.ID, update.MessageThreadID)
	botHandler.Queue.SendMessage(threadName, update.Message.Text)
	*/
}
