package api

import (
	"errors"
	"github.com/aviseu/go-sample/internal/errs"
)

var ErrInvalidID = errs.NewValidationError(errors.New("invalid ID"))
