package infrastructure

import (
	"errors"
	"github.com/aviseu/go-sample/internal/errs"
)

var ErrTaskNotFound = errs.NewValidationError(errors.New("task not found"))
