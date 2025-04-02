package testutils

import (
	"context"
	"github.com/aviseu/go-sample/internal/app/infrastructure"
	"github.com/aviseu/go-sample/internal/app/infrastructure/aggregators"
	"github.com/google/uuid"
	"sort"
)

type TaskRepositoryOptional func(*TaskRepository)

func TaskRepositoryWithError(err error) TaskRepositoryOptional {
	return func(r *TaskRepository) {
		r.err = err
	}
}

func TaskRepositoryWithTask(t *aggregators.Task) TaskRepositoryOptional {
	return func(r *TaskRepository) {
		r.Records[t.ID] = t
	}
}

type TaskRepository struct {
	Records map[uuid.UUID]*aggregators.Task
	err     error
}

func NewTaskRepository(opts ...TaskRepositoryOptional) *TaskRepository {
	r := &TaskRepository{Records: make(map[uuid.UUID]*aggregators.Task)}
	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *TaskRepository) All(_ context.Context) ([]*aggregators.Task, error) {
	if r.err != nil {
		return nil, r.err
	}

	tasks := make([]*aggregators.Task, 0, len(r.Records))
	for _, task := range r.Records {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Title < tasks[j].Title
	})

	return tasks, nil
}

func (r *TaskRepository) Find(_ context.Context, id uuid.UUID) (*aggregators.Task, error) {
	if r.err != nil {
		return nil, r.err
	}

	task, ok := r.Records[id]
	if !ok {
		return nil, infrastructure.ErrTaskNotFound
	}

	return task, nil
}

func (r *TaskRepository) Save(_ context.Context, task *aggregators.Task) error {
	if r.err != nil {
		return r.err
	}

	r.Records[task.ID] = task

	return nil
}
