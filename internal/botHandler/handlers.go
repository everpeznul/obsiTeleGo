package botHandler

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (botHandler *BotHandler) InitThreadHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.Fields(update.Message.Text)
	threadName := strings.Join(parts[1:], " ")
	botHandler.Log.Info("Init New Thread", "thread id", update.Message.MessageThreadID, "thread name", threadName)

	err := botHandler.Repo.NewThread(ctx, int64(update.Message.MessageThreadID), threadName)
	if err != nil {

		botHandler.Log.Error("Init Thread Error", "error", err)
		return
	}
	botHandler.Log.Info("Init Repo Thread")

	newQueue, _ := botHandler.Rabbit.Ch.QueueDeclare(threadName, true, false, false, false, nil)
	err = botHandler.Rabbit.Ch.QueueBind(newQueue.Name, "tg."+threadName, "tg_router", false, nil)
	if err != nil {

		botHandler.Log.Error("Init Queue Error", "error", err)
		return
	}
	botHandler.Log.Info("Init Queue Thread")
}
