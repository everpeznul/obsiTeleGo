package app

import (
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	_ "modernc.org/sqlite"
	"obsiTeleGo/internal/botHandler"
	"obsiTeleGo/internal/logger"
	"obsiTeleGo/internal/repository"
	"obsiTeleGo/internal/repository/sqliteRepo"
	"os"
)

type App struct {
	Logger     *logger.Logger
	Log        *slog.Logger
	BotHandler *botHandler.BotHandler
	db         *sql.DB
	Repo       repository.Repo
}

func New() App {
	base := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger := initLogger(base)
	log := initAppLog(base)
	repo, db, err := initSQLiteRepo(logger.Repo)
	botHandler := initBotHandler(logger.BotHandler, repo)

	if err != nil {

		log.Error("Init SQLiteRepo Error", "error", err)
		panic(err)
	}

	return App{
		Logger:     logger,
		Log:        log,
		BotHandler: botHandler,
		db:         db,
		Repo:       repo,
	}
}

func initLogger(base *slog.Logger) *logger.Logger {
	return logger.New(base)
}

func initAppLog(base *slog.Logger) *slog.Logger {
	return base.With("logger", "app")
}

func initSQLiteRepo(log *slog.Logger) (repository.Repo, *sql.DB, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/mydb.sqlite" // локальная разработка
	}

	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, nil, fmt.Errorf("open conn to sqlite error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping sqlite error: %w", err)
	}

	return sqliteRepo.New(db, log), db, nil
}

func initRedisRepo(log *slog.Logger) (redis.Client, *sql.DB, error) {

}

func (a *App) DBClose() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func initBotHandler(log *slog.Logger, repo repository.Repo) *botHandler.BotHandler {
	return botHandler.New(log, repo)
}
