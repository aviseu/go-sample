package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/aviseu/go-sample/internal/app/infrastructure"
	"github.com/aviseu/go-sample/internal/app/infrastructure/aggregators"
	"github.com/google/uuid"
)

type Repository interface {
	Find(ctx context.Context, id uuid.UUID) (*aggregators.Task, error)
	Save(ctx context.Context, task *aggregators.Task) error
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, title string) (*aggregators.Task, error) {
	if title == "" {
		return nil, ErrTitleIsRequired
	}

	task := newTask(uuid.New(), title, false).toAggregator()

	if err := s.r.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

func (s *Service) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	task, err := s.r.Find(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrTaskNotFound) {
			return fmt.Errorf("%w: %s", ErrTaskNotFound, id)
		}
		return fmt.Errorf("failed to find task: %w", err)
	}

	d := newFromAggregator(task)
	d.markCompleted()

	if err := s.r.Save(ctx, d.toAggregator()); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}
