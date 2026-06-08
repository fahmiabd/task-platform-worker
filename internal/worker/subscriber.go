package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/fahmiabd/task-platform-worker/internal/repository"
	"github.com/fahmiabd/task-platform-worker/internal/task"
)

func HandleMessage(
	ctx context.Context,
	taskRepo *repository.TaskRepository,
	data []byte,
) error {
	var message task.Message

	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	log.Printf(
		"task_id=%s type=%s",
		message.TaskID,
		message.Type,
	)

	if err := taskRepo.MarkProcessing(
		ctx,
		message.TaskID,
	); err != nil {
		return err
	}

	log.Printf(
		"task_id=%s processing",
		message.TaskID,
	)

	time.Sleep(3 * time.Second)

	if err := taskRepo.MarkCompleted(
		ctx,
		message.TaskID,
	); err != nil {
		return err
	}

	log.Printf(
		"task_id=%s completed",
		message.TaskID,
	)

	return nil
}
