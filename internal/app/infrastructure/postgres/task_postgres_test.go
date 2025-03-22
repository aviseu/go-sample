package postgres_test

import (
	"context"
	"github.com/aviseu/go-sample/internal/app/domain"
	"github.com/aviseu/go-sample/internal/app/infrastructure/postgres"
	"github.com/aviseu/go-sample/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestTaskRepository(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TaskRepositorySuite))
}

type TaskRepositorySuite struct {
	testutils.PostgresSuite
}

func (suite *TaskRepositorySuite) TestAllSuccess() {
	// Prepare
	id3 := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", id3.String(), "task 3", false)
	suite.NoError(err)
	id1 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", id1.String(), "task 1", false)
	suite.NoError(err)
	id2 := uuid.New()
	_, err = suite.DB.Exec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", id2.String(), "task 2", false)
	suite.NoError(err)
	r := postgres.NewTaskRepository(suite.DB)

	// Execute
	tasks, err := r.All(context.Background())

	// Assert
	suite.NoError(err)
	suite.Len(tasks, 3)
	suite.Equal(id1, tasks[0].ID)
	suite.Equal("task 1", tasks[0].Title)
	suite.Equal(id2, tasks[1].ID)
	suite.Equal("task 2", tasks[1].Title)
	suite.Equal(id3, tasks[2].ID)
	suite.Equal("task 3", tasks[2].Title)
}

func (suite *TaskRepositorySuite) TestAllRepositoryFail() {
	// Prepare
	r := postgres.NewTaskRepository(suite.BadDB)

	// Execute
	tasks, err := r.All(context.Background())

	// Assert
	suite.Error(err)
	suite.Nil(tasks)
	suite.ErrorContains(err, "failed to get tasks")
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *TaskRepositorySuite) TestFindSuccess() {
	// Prepare
	id := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", id.String(), "task 1", false)
	suite.NoError(err)
	r := postgres.NewTaskRepository(suite.DB)

	// Execute
	task, err := r.Find(context.Background(), id)

	// Assert
	suite.NoError(err)
	suite.Equal(id, task.ID)
	suite.Equal("task 1", task.Title)
	suite.False(task.Completed)
}

func (suite *TaskRepositorySuite) TestFindNotFound() {
	// Prepare
	r := postgres.NewTaskRepository(suite.DB)

	// Execute
	task, err := r.Find(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.Nil(task)
	suite.ErrorIs(err, domain.ErrTaskNotFound)
}

func (suite *TaskRepositorySuite) TestFindRepositoryFail() {
	// Prepare
	r := postgres.NewTaskRepository(suite.BadDB)

	// Execute
	task, err := r.Find(context.Background(), uuid.New())

	// Assert
	suite.Error(err)
	suite.Nil(task)
	suite.ErrorContains(err, "failed to find task")
	suite.ErrorContains(err, "sql: database is closed")
}

func (suite *TaskRepositorySuite) TestSaveNewSuccess() {
	// Prepare
	id := uuid.New()
	task := &domain.Task{ID: id, Title: "task 1"}
	r := postgres.NewTaskRepository(suite.DB)

	// Execute
	err := r.Save(context.Background(), task)

	// Assert result
	suite.NoError(err)
	suite.Equal(id, task.ID)
	suite.Equal("task 1", task.Title)
	suite.False(task.Completed)

	// Assert state
	var dbTasks []*domain.Task
	err = suite.DB.Select(&dbTasks, "SELECT * FROM tasks")
	suite.NoError(err)
	suite.Len(dbTasks, 1)
	suite.Equal(id, dbTasks[0].ID)
	suite.Equal("task 1", dbTasks[0].Title)
	suite.False(dbTasks[0].Completed)
}

func (suite *TaskRepositorySuite) TestSaveExistingSuccess() {
	// Prepare
	id := uuid.New()
	_, err := suite.DB.Exec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", id.String(), "task 1", false)
	suite.NoError(err)

	task := &domain.Task{ID: id, Title: "task 1 updated", Completed: true}
	r := postgres.NewTaskRepository(suite.DB)

	// Execute
	err = r.Save(context.Background(), task)

	// Assert result
	suite.NoError(err)
	suite.Equal(id, task.ID)
	suite.Equal("task 1 updated", task.Title)
	suite.True(task.Completed)

	// Assert state
	var dbTasks []*domain.Task
	err = suite.DB.Select(&dbTasks, "SELECT * FROM tasks")
	suite.NoError(err)
	suite.Len(dbTasks, 1)
	suite.Equal(id, dbTasks[0].ID)
	suite.Equal("task 1 updated", dbTasks[0].Title)
	suite.True(dbTasks[0].Completed)
}

func (suite *TaskRepositorySuite) TestSaveRepositoryFail() {
	// Prepare
	r := postgres.NewTaskRepository(suite.BadDB)

	// Execute
	err := r.Save(context.Background(), &domain.Task{})

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to save task")
	suite.ErrorContains(err, "sql: database is closed")
}
