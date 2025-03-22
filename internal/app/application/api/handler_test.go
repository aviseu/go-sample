package api_test

import (
	"errors"
	"github.com/aviseu/go-sample/internal/app/application"
	"github.com/aviseu/go-sample/internal/app/domain"
	"github.com/aviseu/go-sample/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) TestAllSuccess() {
	// Prepare
	id1 := uuid.New()
	id2 := uuid.New()
	r := testutils.NewTaskRepository(
		testutils.TaskRepositoryWithTask(&domain.Task{ID: id1, Title: "task 1", Completed: false}),
		testutils.TaskRepositoryWithTask(&domain.Task{ID: id2, Title: "task 2", Completed: true}),
	)
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodGet, "/api/tasks", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert result
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"tasks":[{"id":"`+id1.String()+`","title":"task 1","completed":false},{"id":"`+id2.String()+`","title":"task 2","completed":true}]}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestAllRepositoryFail() {
	// Prepare
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithError(errors.New("boom!")))
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodGet, "/api/tasks", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert result
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal(`{"message":"Internal Server Error"}`+"\n", rr.Body.String())

	// Assert log
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(logs[0], `"level":"ERROR"`)
	suite.Contains(logs[0], `"msg":"boom!"`)
}

func (suite *HandlerSuite) TestCreateSuccess() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodPost, "/api/tasks", strings.NewReader(`{"title":"task 1"}`))
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state
	suite.Len(r.Records, 1)
	var task *domain.Task
	for _, t := range r.Records {
		task = t
	}
	suite.NotNil(task)
	suite.NotEmpty(task.ID)
	suite.Equal("task 1", task.Title)
	suite.False(task.Completed)

	// Assert result
	suite.Equal(oghttp.StatusCreated, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+task.ID.String()+`","title":"task 1","completed":false}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestCreateInvalidRequest() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodPost, "/api/tasks", strings.NewReader(`{"invalid":"task 1"}`))
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state
	suite.Empty(r.Records)

	// Assert result
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"message":"title is required"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestCreateRepositoryFail() {
	// Prepare
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithError(errors.New("boom!")))
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodPost, "/api/tasks", strings.NewReader(`{"title":"task 1"}`))
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state
	suite.Empty(r.Records)

	// Assert result
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal(`{"message":"Internal Server Error"}`+"\n", rr.Body.String())

	// Assert log
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(logs[0], `"level":"ERROR"`)
	suite.Contains(logs[0], `"msg":"failed to save task: boom!"`)
}

func (suite *HandlerSuite) TestMarkCompletedSuccess() {
	// Prepare
	id := uuid.New()
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithTask(&domain.Task{ID: id, Title: "task 1", Completed: false}))
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodPut, "/api/tasks/"+id.String()+"/complete", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state
	suite.True(r.Records[id].Completed)

	// Assert result
	suite.Equal(oghttp.StatusNoContent, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestMarkCompletedInvalidID() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodPut, "/api/tasks/invalid/complete", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state
	suite.Empty(r.Records)

	// Assert result
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"message":"invalid task ID"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestMarkCompletedNotFound() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodPut, "/api/tasks/"+uuid.New().String()+"/complete", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state
	suite.Empty(r.Records)

	// Assert result
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"message":"failed to find task: task not found"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestMarkCompletedRepositoryFail() {
	// Prepare
	id := uuid.New()
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithError(errors.New("boom!")))
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodPut, "/api/tasks/"+id.String()+"/complete", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert state
	suite.Empty(r.Records)

	// Assert result
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal(`{"message":"Internal Server Error"}`+"\n", rr.Body.String())

	// Assert log
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(lbuf.String(), `"level":"ERROR"`)
	suite.Contains(lbuf.String(), `"msg":"failed to find task: boom!"`)
}

func (suite *HandlerSuite) TestFindSuccess() {
	// Prepare
	id := uuid.New()
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithTask(&domain.Task{ID: id, Title: "task 1", Completed: false}))
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodGet, "/api/tasks/"+id.String(), nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert result
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+id.String()+`","title":"task 1","completed":false}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestFindInvalidID() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodGet, "/api/tasks/invalid", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert result
	suite.Equal(oghttp.StatusBadRequest, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"message":"invalid task ID"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestFindNotFound() {
	// Prepare
	r := testutils.NewTaskRepository()
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodGet, "/api/tasks/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert result
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"message":"task not found"}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *HandlerSuite) TestFindRepositoryFail() {
	// Prepare
	r := testutils.NewTaskRepository(testutils.TaskRepositoryWithError(errors.New("boom!")))
	s := domain.NewService(r)
	lbuf, log := testutils.NewLogger()
	h := application.APIHandler(log, s)

	req := httptest.NewRequest(oghttp.MethodGet, "/api/tasks/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert result
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal(`{"message":"Internal Server Error"}`+"\n", rr.Body.String())

	// Assert log
	logs := testutils.LogLines(lbuf)
	suite.Len(logs, 1)
	suite.Contains(lbuf.String(), `"level":"ERROR"`)
	suite.Contains(lbuf.String(), `"msg":"boom!"`)
}
