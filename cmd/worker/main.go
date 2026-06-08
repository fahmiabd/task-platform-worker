package main

import (
	"context"
	"log"

	"github.com/fahmiabd/task-platform-worker/internal/repository"
	"github.com/fahmiabd/task-platform-worker/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
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

	_, err = nc.Subscribe("tasks.>", func(msg *nats.Msg) {
		if err := worker.HandleMessage(
			ctx,
			taskRepo,
			msg.Data,
		); err != nil {
			log.Printf("failed to handle message: %v", err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("worker started")

	select {}
}
