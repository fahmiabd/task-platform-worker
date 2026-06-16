package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

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
		return err
	}

	log.Printf(
		"task_id=%s type=%s",
		message.TaskID,
		message.Type,
	)

	task, err := taskRepo.FindByID(ctx, message.TaskID)
	if err != nil {
		log.Printf(
			"task_id=%s error on FindByID",
			message.TaskID,
		)
		return ErrAck
	}

	if task.RetryCount >= 3 {
		err = taskRepo.MarkFailed(ctx, message.TaskID)
		if err != nil {
			taskRepo.IncrementRetryCount(ctx, message.TaskID)
			return ErrRetry
		}

		return ErrAck
	}

	// test retry
	// if task.RetryCount < 3 {
	// 	log.Printf(
	// 		"task_id=%s failed, retrying count=%d",
	// 		message.TaskID,
	// 		task.RetryCount,
	// 	)
	// 	taskRepo.IncrementRetryCount(ctx, message.TaskID)
	// 	return ErrRetry
	// }

	if err := taskRepo.MarkProcessing(
		ctx,
		message.TaskID,
	); err != nil {
		taskRepo.IncrementRetryCount(ctx, message.TaskID)
		return ErrRetry
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
		taskRepo.IncrementRetryCount(ctx, message.TaskID)
		return ErrRetry
	}

	log.Printf(
		"task_id=%s completed",
		message.TaskID,
	)

	return nil
}
