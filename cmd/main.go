package main

import (
	"context"
	"log/slog"
	"obsiTeleGo/cmd/app"
	"obsiTeleGo/config"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	a := app.New(&app.Options{
		Repo: os.Getenv("REPO"),
	})
	defer a.DBClose()
	defer a.Rabbit.Conn.Close()
	defer a.Rabbit.Ch.Close()

	configPath := os.Getenv("CONFIG_PATH")
	if err := config.InitConfig(configPath); err != nil {
		a.Log.Error("Failed to initialize config", "error", err)
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(a.BotHandler.Handle),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/init_thread", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		a.BotHandler.InitThreadHandler(ctx, b, update)
	})

	go b.Start(ctx)

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/config/config.json", func(c *gin.Context) {
		c.File(configPath)
	})

	r.PUT("/admin/loglevel", func(c *gin.Context) {
		levelStr := c.Query("level")

		switch levelStr {
		case "debug":
			app.LogLevel.Set(slog.LevelDebug)
		case "info":
			app.LogLevel.Set(slog.LevelInfo)
		case "warn":
			app.LogLevel.Set(slog.LevelWarn)
		case "error":
			app.LogLevel.Set(slog.LevelError)
		default:
			c.String(400, "unknown level")
			return
		}

		c.String(200, "Log level changed to %s", levelStr)
	})
	r.Run(":8080")
}
