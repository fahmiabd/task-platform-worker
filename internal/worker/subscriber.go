package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/fahmiabd/task-platform-worker/internal/repository"
	"github.com/fahmiabd/task-platform-worker/internal/task"
)

var (
	ErrRetry = errors.New("retry")
	ErrAck   = errors.New("ack")
)

func HandleMessage(
	ctx context.Context,
	taskRepo *repository.TaskRepository,
	data []byte,
) error {
	var message task.Message

	if err := json.Unmarshal(data, &message); err != nil {
		log.Printf("invalid message")
		return ErrAck
	}

	log.Printf(
		"task_id=%s type=%s",
		message.TaskID,
		message.Type,
	)

	task, err := taskRepo.FindByID(ctx, message.TaskID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrTaskNotFound):
			log.Printf(
				"task_id=%s not found",
				message.TaskID,
			)
			return ErrAck

		default:
			return ErrRetry
		}
	}

	if task.RetryCount >= 3 {
		if err = taskRepo.MarkFailed(ctx, message.TaskID); err != nil {
			return ErrRetry
		}

		return ErrAck
	}

	if err := taskRepo.MarkProcessing(
		ctx,
		message.TaskID,
	); err != nil {
		IncrementRetryCount(ctx, taskRepo, message.TaskID)
		return ErrRetry
	}

	log.Printf(
		"task_id=%s processing",
		message.TaskID,
	)

	if err := taskRepo.MarkCompleted(
		ctx,
		message.TaskID,
	); err != nil {
		IncrementRetryCount(ctx, taskRepo, message.TaskID)
		return ErrRetry
	}

	log.Printf(
		"task_id=%s completed",
		message.TaskID,
	)

	return nil
}

func IncrementRetryCount(ctx context.Context, taskRepo *repository.TaskRepository, taskID string) {
	if err := taskRepo.IncrementRetryCount(ctx, taskID); err != nil {
		log.Printf(
			"task_id=%s increment retry count failed: %v",
			taskID,
			err,
		)
	}
}
