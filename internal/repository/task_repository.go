package repository

import (
	"context"

	"github.com/fahmiabd/task-platform-worker/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) FindByID(
	ctx context.Context,
	taskID string,
) (*entity.Task, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Task])
	if err != nil {
		return nil, err
	}

	return &task, nil
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

func (r *TaskRepository) MarkFailed(
	ctx context.Context,
	taskID string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE tasks
		SET status = 'failed'
		WHERE id = $1
		`,
		taskID,
	)

	return err
}

func (r *TaskRepository) IncrementRetryCount(
	ctx context.Context,
	taskID string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE tasks
		SET retry_count = retry_count + 1
		WHERE id = $1
		`,
		taskID,
	)

	return err
}
