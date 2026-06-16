package main

import (
	"context"
	"errors"
	"log"

	"github.com/fahmiabd/task-platform-worker/internal/repository"
	"github.com/fahmiabd/task-platform-worker/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5432/task_platform",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskRepo := repository.NewTaskRepository(db)

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := js.Consumer(
		context.Background(),
		"TASKS",
		"worker",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("worker started")

	for {
		msgs, err := consumer.Fetch(1)
		if err != nil {
			continue
		}

		for msg := range msgs.Messages() {
			err := worker.HandleMessage(
				ctx,
				taskRepo,
				msg.Data(),
			)

			switch {
			case err == nil:
				msg.Ack()

			case errors.Is(err, worker.ErrRetry):
				msg.Nak()

			case errors.Is(err, worker.ErrAck):
				msg.Ack()

			default:
				msg.Nak()
			}
		}
	}
}
