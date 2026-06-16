package entity

import "time"

type Task struct {
	ID           string
	Type         string
	Payload      string
	Status       string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	StartedAt    *time.Time
	CompletedAt  *time.Time
	ErrorMessage *string
	RetryCount   int
}
