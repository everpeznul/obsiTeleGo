package botHandler

import (
	"context"
	"obsiTeleGo/config"
	"os"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (botHandler *BotHandler) MessageHandle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {

		botHandler.Log.Debug("MessageHandle Ignoring Non Message Update")
		return
	}

	threadID := update.Message.MessageThreadID
	text := update.Message.Text
	botHandler.Log.Debug("MessageHandle New Message", "thread id", threadID, "message text", len(text))

	threadName, err := botHandler.Repo.GetThreadName(ctx, int64(threadID))
	if err != nil {

		botHandler.Log.Error("MessageHandle Get Thread Name Error", "error", err)
		return
	}
	botHandler.Log.Debug("MessageHandle Get Thread Name Successfully", "thread", threadName)

	key := "tg." + threadName
	err = botHandler.Rabbit.Ch.PublishWithContext(
		ctx,
		"tg_router",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(text),
		},
	)
	if err != nil {
		botHandler.Log.Error("MessageHandle Push Message Queue Error", "error", err)
		return
	}
	botHandler.Log.Debug("MessageHandle Push Message Queue Successfully", "queue", key)

	botHandler.Log.Info("MessageHandle Successfully")
}

func (botHandler *BotHandler) InitThreadHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	parts := strings.Fields(update.Message.Text)
	threadName := strings.Join(parts[1:], " ")
	botHandler.Log.Info("Init New Thread", "thread id", update.Message.MessageThreadID, "thread name", threadName)

	err := botHandler.Repo.NewThread(ctx, int64(update.Message.MessageThreadID), threadName)
	if err != nil {

		botHandler.Log.Error("Init Thread Repo Error", "error", err)
		return
	}
	botHandler.Log.Debug("Init Thread Repo Successfully")

	newQueue, _ := botHandler.Rabbit.Ch.QueueDeclare(threadName, true, false, false, false, nil)
	err = botHandler.Rabbit.Ch.QueueBind(newQueue.Name, "tg."+threadName, "tg_router", false, nil)
	if err != nil {

		botHandler.Log.Error("Init Thread Queue Error", "error", err)
		return
	}
	botHandler.Log.Debug("Init Thread Queue Successfully")

	configPath := os.Getenv("CONFIG_PATH")
	err = config.AddTopic(configPath, threadName)
	if err != nil {

		botHandler.Log.Error("Update Config Error", "error", err)
		return
	}
	botHandler.Log.Debug("Init Thread Update Config Successfully")

	botHandler.Log.Info("Init Thread Successfully")
}
