package domain

import "github.com/google/uuid"

type Task struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Completed bool      `json:"completed" db:"completed"`
}

func NewTask(id uuid.UUID, title string, completed bool) *Task {
	return &Task{
		ID:        id,
		Title:     title,
		Completed: completed,
	}
}

func (t *Task) markCompleted() {
	t.Completed = true
}
