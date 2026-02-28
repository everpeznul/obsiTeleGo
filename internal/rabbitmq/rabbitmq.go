package rabbitmq

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Log  *slog.Logger
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func New(log *slog.Logger) *Rabbit {
	conn, _ := amqp.Dial("amqp://guest:guest@obsidian_rabbitmq:5672/")
	ch, _ := conn.Channel()

	return &Rabbit{
		Log:  log,
		Conn: conn,
		Ch:   ch,
	}
}
