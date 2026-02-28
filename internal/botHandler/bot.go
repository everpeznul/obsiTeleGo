package botHandler

import (
	"context"
	"log/slog"
	"obsiTeleGo/internal/rabbitmq"
	"obsiTeleGo/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	amqp "github.com/rabbitmq/amqp091-go"
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

func (botHandler *BotHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	threadID := update.Message.MessageThreadID
	text := update.Message.Text

	botHandler.Log.Debug("New Message", "thread id", threadID, "message text", text)
	threadName, err := botHandler.Repo.GetThreadName(ctx, int64(threadID))
	if err != nil {
		botHandler.Log.Error("Get Thread Name Error", "error", err)
		return
	}

	botHandler.Log.Debug("Get Thread Name", "thread", threadName)

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
		botHandler.Log.Error("Push Message Queue Error", "error", err)
		return
	}
	botHandler.Log.Debug("Push Message Queue", "queue", key)
}
