package main

import (
	"log"

	"github.com/fahmiabd/task-platform-worker/internal/worker"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("tasks.>", func(msg *nats.Msg) {
		if err := worker.HandleMessage(msg.Data); err != nil {
			log.Printf("failed to handle message: %v", err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("worker started")

	select {}
}
