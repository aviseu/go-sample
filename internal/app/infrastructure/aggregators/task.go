package aggregators

import (
	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Completed bool      `db:"completed" json:"completed"`
}
