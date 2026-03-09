package httpServer

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Log      *slog.Logger
	LogLevel *slog.LevelVar
	Router   *gin.Engine
	srv      *http.Server
}

func New(log *slog.Logger, logLevel *slog.LevelVar) *HttpServer {
	r := gin.Default()

	s := &HttpServer{
		Log:      log,
		LogLevel: logLevel,
		Router:   r,
	}

	s.setupRoutes()
	return s
}

func (s *HttpServer) setupRoutes() {
	s.Router.Use(corsMiddleware())

	configPath := os.Getenv("CONFIG_PATH")

	s.Router.GET("/config.json", func(c *gin.Context) {
		s.Log.Info("Send Config Request", "ip", c.ClientIP())
		c.File(configPath)
	})

	s.Router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	s.Router.PUT("/admin/loglevel", func(c *gin.Context) {
		newLevel := c.Query("level")
		s.Log.Info("Change LogLevel Request", "level", newLevel, "ip", c.ClientIP())

		switch newLevel {
		case "debug":
			s.LogLevel.Set(slog.LevelDebug)
		case "info":
			s.LogLevel.Set(slog.LevelInfo)
		case "warn":
			s.LogLevel.Set(slog.LevelWarn)
		case "error":
			s.LogLevel.Set(slog.LevelError)
		default:
			s.Log.Warn("Change LogLevel Bad Level Warn")
			c.String(400, "unknown level")
			return
		}

		s.Log.Info("Change LogLevel Successfully", "new_level", newLevel)
		c.String(200, "Log level changed to %s", newLevel)
	})
}

func (s *HttpServer) Run(addr string) error {
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	s.Log.Info("Starting Gin Server", "addr", addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	s.Log.Info("Shutting down Gin Server")
	return s.srv.Shutdown(ctx)
}

func corsMiddleware() gin.HandlerFunc {
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
