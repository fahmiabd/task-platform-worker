package worker

import (
	"encoding/json"
	"log"

	"github.com/fahmiabd/task-platform-worker/internal/task"
)

func HandleMessage(data []byte) error {
	var message task.Message

	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	log.Printf(
		"task_id=%s type=%s",
		message.TaskID,
		message.Type,
	)

	return nil
}
