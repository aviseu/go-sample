package api

import (
	"github.com/aviseu/go-sample/internal/app/infrastructure/aggregators"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Message: err.Error()}
}

type TaskListResponse struct {
	Tasks []*aggregators.Task `json:"tasks"`
}

func NewTaskListResponse(tasks []*aggregators.Task) *TaskListResponse {
	if tasks == nil {
		tasks = []*aggregators.Task{}
	}
	return &TaskListResponse{
		Tasks: tasks,
	}
}
