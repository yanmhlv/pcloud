package pcloud

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	ErrAuthRequired  ErrorCode = 1000
	ErrInvalidName   ErrorCode = 2001
	ErrAccessDenied  ErrorCode = 2003
	ErrAlreadyExists ErrorCode = 2004
	ErrNotFound      ErrorCode = 2005
)

type apiError interface {
	Err() error
}

// Error represents a pCloud API error response.
// Result holds the numeric code; zero means success.
type Error struct {
	Result  ErrorCode `json:"result"`
	Message string    `json:"error"`
}

// Err returns nil when Result is 0, otherwise returns the Error itself.
func (e *Error) Err() error {
	if e.Result == 0 {
		return nil
	}
	return e
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("pcloud error %d", e.Result)
	}
	return fmt.Sprintf("pcloud error %d: %s", e.Result, e.Message)
}

func IsErrorCode(err error, code ErrorCode) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Result == code
	}
	return false
}
