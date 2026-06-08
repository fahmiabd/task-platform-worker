package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) MarkProcessing(
	ctx context.Context,
	taskID string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE tasks
		SET status = 'processing'
		WHERE id = $1
		`,
		taskID,
	)

	return err
}

func (r *TaskRepository) MarkCompleted(
	ctx context.Context,
	taskID string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE tasks
		SET status = 'completed'
		WHERE id = $1
		`,
		taskID,
	)

	return err
}
