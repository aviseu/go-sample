package domain_test

import (
	"context"
	"errors"
	"github.com/aviseu/go-sample/internal/app/domain"
	"github.com/aviseu/go-sample/internal/app/infrastructure/aggregators"
	"github.com/aviseu/go-sample/internal/errs"
	"github.com/aviseu/go-sample/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestService(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) TestCreateSuccess() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)

	// Execute
	task, err := s.Create(context.Background(), "task 1")

	// Assert result
	suite.NoError(err)
	suite.NotNil(task)
	suite.NotEmpty(task.ID)
	suite.Equal("task 1", task.Title)
	suite.False(task.Completed)

	// Assert state
	suite.Len(r.Records, 1)
	t, ok := r.Records[task.ID]
	suite.True(ok)
	suite.Equal(task.ID, t.ID)
	suite.Equal("task 1", t.Title)
	suite.False(t.Completed)
}

func (suite *ServiceSuite) TestCreateTitleIsRequiredFail() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)

	// Execute
	task, err := s.Create(context.Background(), "")

	// Assert
	suite.Error(err)
	suite.Nil(task)
	suite.ErrorIs(err, domain.ErrTitleIsRequired)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) TestCreateRepositoryFail() {
	// Prepare
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithError(errors.New("boom!")))
	s := domain.NewService(r)

	// Execute
	task, err := s.Create(context.Background(), "task 1")

	// Assert
	suite.Error(err)
	suite.Nil(task)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}

func (suite *ServiceSuite) TestMarkCompletedSuccess() {
	// Prepare
	id := uuid.New()
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithTask(&aggregators.Task{ID: id, Title: "task 1", Completed: false}))
	s := domain.NewService(r)

	// Execute
	err := s.MarkCompleted(context.Background(), id)

	// Assert
	suite.NoError(err)

	// Assert state
	t, ok := r.Records[id]
	suite.True(ok)
	suite.True(t.Completed)
}

func (suite *ServiceSuite) TestMarkCompletedNotFoundFail() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)

	// Execute
	err := s.MarkCompleted(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, domain.ErrTaskNotFound)
	suite.True(errs.IsValidationError(err))
}

func (suite *ServiceSuite) TestMarkCompletedRepositoryFail() {
	// Prepare
	id := uuid.New()
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithError(errors.New("boom!")))
	s := domain.NewService(r)

	// Execute
	err := s.MarkCompleted(context.Background(), id)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "boom!")
	suite.False(errs.IsValidationError(err))
}
