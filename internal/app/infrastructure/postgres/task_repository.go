package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aviseu/go-sample/internal/app/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sort"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) All(ctx context.Context) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.SelectContext(ctx, &tasks, "SELECT * FROM tasks ORDER BY title")
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Title < tasks[j].Title
	})

	return tasks, nil
}

func (r *TaskRepository) Find(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	err := r.db.GetContext(ctx, &task, "SELECT * FROM tasks WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	return &task, nil
}

func (r *TaskRepository) Save(ctx context.Context, task *domain.Task) error {
	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO tasks (id, title, completed)
		VALUES (:id, :title, :completed)
		ON CONFLICT (id) DO UPDATE SET title = :title, completed = :completed
		RETURNING id`,
		task,
	)
	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}
