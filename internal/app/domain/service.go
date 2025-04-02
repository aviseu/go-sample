package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Repository interface {
	Find(ctx context.Context, id uuid.UUID) (*Task, error)
	Save(ctx context.Context, task *Task) error
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, title string) (*Task, error) {
	if title == "" {
		return nil, ErrTitleIsRequired
	}

	task := &Task{ID: uuid.New(), Title: title}

	if err := s.r.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

func (s *Service) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	task, err := s.r.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	task.Completed = true
	if err := s.r.Save(ctx, task); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}
