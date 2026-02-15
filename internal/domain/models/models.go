package models

import "time"

type Thread struct {
	ChatID    int64
	ThreadID  int64
	Name      string
	CreatedAt time.Time
}

type Message struct {
	Thread    string
	Text      string
	Author    string
	Timestamp time.Time
}
