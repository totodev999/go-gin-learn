package utils

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode  int
	MessageCode MessageCode
	Message     string
	Detail      string
	Err         error
}

func (e *APIError) Error() string {
	base := fmt.Sprintf("%s: %s", e.Message, e.Detail)
	if e.Err != nil {
		return fmt.Sprintf("%s [%v]", base, e.Err)
	}
	return base
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func NewBadRequestError(detail string, err error) *APIError {
	return &APIError{
		StatusCode:  http.StatusBadRequest,
		MessageCode: BadRequest,
		Message:     Messages[BadRequest],
		Detail:      detail,
		Err:         err,
	}
}

func NewNotFoundError(detail string, err error) *APIError {
	return &APIError{
		StatusCode:  http.StatusNotFound,
		MessageCode: NotFound,
		Message:     Messages[NotFound],
		Detail:      detail,
		Err:         err,
	}
}

func NewUnauthorized(detail string, err error) *APIError {
	return &APIError{
		StatusCode:  http.StatusUnauthorized,
		MessageCode: UnAuthorized,
		Message:     Messages[UnAuthorized],
		Detail:      detail,
		Err:         err,
	}
}

func NewDBError(detail string, err error) *APIError {
	return &APIError{
		StatusCode:  http.StatusInternalServerError,
		MessageCode: DBError,
		Message:     Messages[DBError],
		Detail:      detail,
		Err:         err,
	}
}

func NewDuplicateKeyError(detail string, err error) *APIError {
	return &APIError{
		StatusCode:  http.StatusConflict,
		MessageCode: DuplicateKeyError,
		Message:     Messages[DuplicateKeyError],
		Detail:      detail,
		Err:         err,
	}
}

func NewUnknownError(detail string, err error) *APIError {
	return &APIError{
		StatusCode:  http.StatusInternalServerError,
		MessageCode: UnknownError,
		Message:     Messages[UnknownError],
		Detail:      detail,
		Err:         err,
	}
}
