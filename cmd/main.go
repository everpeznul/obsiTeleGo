package main

import (
	"context"
	"errors"
	"net/http"
	"obsiTeleGo/cmd/app"
	"os"
	"os/signal"
	"time"
)

func main() {
	a := app.New(&app.Options{
		Repo: os.Getenv("REPO"),
	})

	defer a.DBClose()
	defer a.Rabbit.Conn.Close()
	defer a.Rabbit.Ch.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := a.BotHandler.Run(ctx, os.Getenv("TELEGRAM_BOT_TOKEN")); err != nil {
			a.Log.Error("Failed to start bot", "error", err)
		}
	}()

	go func() {
		a.Log.Info("Start Gin Server on :8080 Successfully")
		if err := a.HttpServer.Run(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Log.Error("Gin Server Shutdown Error", "error", err)
		}
	}()

	// блокировка main
	<-ctx.Done()
	a.Log.Info("Interrupt signal received, Starting graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.HttpServer.Stop(shutdownCtx); err != nil {
		a.Log.Error("Gin Server forced to shutdown", "error", err)
	}

	a.Log.Info("Graceful shutdown completed, Exiting")
}
