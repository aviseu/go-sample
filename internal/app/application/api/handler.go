package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aviseu/go-sample/internal/app/domain"
	"github.com/aviseu/go-sample/internal/app/infrastructure"
	"github.com/aviseu/go-sample/internal/app/infrastructure/aggregators"
	"github.com/aviseu/go-sample/internal/errs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Repository interface {
	All(ctx context.Context) ([]*aggregators.Task, error)
	Find(ctx context.Context, id uuid.UUID) (*aggregators.Task, error)
}

type Handler struct {
	log *slog.Logger
	s   *domain.Service
	r   Repository
}

func NewHandler(log *slog.Logger, s *domain.Service, r Repository) *Handler {
	return &Handler{
		log: log,
		s:   s,
		r:   r,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.All)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Find)
	r.Put("/{id}/complete", h.MarkCompleted)

	return r
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.r.All(r.Context())
	if err != nil {
		h.handleError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := NewTaskListResponse(tasks)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.handleError(err, w)
		return
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req RequestTaskCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleFail(ErrInvalidID, http.StatusBadRequest, w)
		return
	}

	task, err := h.s.Create(r.Context(), req.Title)
	if err != nil {
		if errs.IsValidationError(err) {
			h.handleFail(err, http.StatusBadRequest, w)
			return
		}

		h.handleError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.handleError(err, w)
		return
	}
}

func (h *Handler) MarkCompleted(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.handleFail(errors.New("invalid task ID"), http.StatusBadRequest, w)
		return
	}

	if err := h.s.MarkCompleted(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			h.handleFail(err, http.StatusNotFound, w)
			return
		}

		h.handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.handleFail(errors.New("invalid task ID"), http.StatusBadRequest, w)
		return
	}

	task, err := h.r.Find(r.Context(), id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrTaskNotFound) {
			h.handleFail(err, http.StatusNotFound, w)
			return
		}

		h.handleError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.handleError(err, w)
		return
	}
}

func (h *Handler) handleError(err error, w http.ResponseWriter) {
	h.log.Error(err.Error())
	h.handleFail(errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError, w)
}

func (h *Handler) handleFail(err error, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := NewErrorResponse(err)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
