package api

import "github.com/aviseu/go-sample/internal/app/domain"

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Message: err.Error()}
}

type TaskListResponse struct {
	Tasks []*domain.Task `json:"tasks"`
}

func NewTaskListResponse(tasks []*domain.Task) *TaskListResponse {
	if tasks == nil {
		tasks = []*domain.Task{}
	}
	return &TaskListResponse{
		Tasks: tasks,
	}
}
