package domain

import (
	"github.com/aviseu/go-sample/internal/app/infrastructure/aggregators"
	"github.com/google/uuid"
)

type task struct {
	id        uuid.UUID
	title     string
	completed bool
}

func newTask(id uuid.UUID, title string, completed bool) *task {
	return &task{
		id:        id,
		title:     title,
		completed: completed,
	}
}

func (t *task) markCompleted() {
	t.completed = true
}

func newFromAggregator(t *aggregators.Task) *task {
	return &task{
		id:        t.ID,
		title:     t.Title,
		completed: t.Completed,
	}
}

func (t *task) toAggregator() *aggregators.Task {
	return &aggregators.Task{
		ID:        t.id,
		Title:     t.title,
		Completed: t.completed,
	}
}
