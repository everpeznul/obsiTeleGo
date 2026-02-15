package main

import (
	"context"
	"obsiTeleGo/cmd/app"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {

	app := app.New()

	app.Log.Info("Starting obsiTeleGo")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(app.BotHandler.Handle),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {

		app.Log.Error("Error starting bot", "error", err)
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/initThread", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		app.BotHandler.InitThreadHandler(ctx, b, update)
	})

	b.Start(ctx)
}
