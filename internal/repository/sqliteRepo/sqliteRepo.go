package sqliteRepo

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type SQLiteRepo struct {
	db  *sql.DB
	Log *slog.Logger
}

func New(db *sql.DB, log *slog.Logger) *SQLiteRepo {
	sqLiteRepo := SQLiteRepo{
		db:  db,
		Log: log,
	}

	if err := sqLiteRepo.initSchema(); err != nil {
		sqLiteRepo.Log.Error("Init Schema Error", "error", err)
	}

	return &sqLiteRepo
}

func (s *SQLiteRepo) initSchema() error {
	query := `
    CREATE TABLE IF NOT EXISTS threads (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        threadID INTEGER NOT NULL UNIQUE,
        name TEXT NOT NULL
    )`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("create table error: %w", err)
	}

	return nil
}

func (s *SQLiteRepo) GetThreadName(ctx context.Context, threadID int64) (string, error) {
	var name string
	err := s.db.QueryRowContext(ctx, "SELECT name FROM threads WHERE threadID = ?", threadID).
		Scan(&name)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("thread %d not found", threadID)
	}
	if err != nil {
		return "", fmt.Errorf("query thread error: %w", err)
	}

	return name, nil
}
func (s *SQLiteRepo) NewThread(ctx context.Context, threadID int64, name string) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO threads (threadID, name) VALUES (?, ?)",
		threadID, name,
	)

	if err != nil {
		return fmt.Errorf("exec thread error: %w", err)
	}

	return nil
}
