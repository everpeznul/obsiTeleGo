package repository

import "context"

type Repo interface {
	GetThreadName(ctx context.Context, threadID int64) (string, error)
	NewThread(ctx context.Context, threadID int64, text string) error
}
