package botHandler

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (botHandler *BotHandler) InitThreadHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.Fields(update.Message.Text)
	text := strings.Join(parts[1:], " ")
	botHandler.Log.Info("Init New Thread", "thread id", update.Message.MessageThreadID, "thread name", text)

	/*	 botHandler.Repo.NewThread(update.Message.ThreadID, update.Message.Text)
	 */
}
