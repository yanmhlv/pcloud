package pcloud

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsErrorCodeMatch(t *testing.T) {
	err := &Error{Result: ErrNotFound, Message: "not found"}
	if !IsErrorCode(err, ErrNotFound) {
		t.Fatal("expected match for ErrNotFound")
	}
}

func TestIsErrorCodeNoMatch(t *testing.T) {
	err := &Error{Result: ErrNotFound, Message: "not found"}
	if IsErrorCode(err, ErrAccessDenied) {
		t.Fatal("expected no match for ErrAccessDenied")
	}
}

func TestIsErrorCodeNonPcloud(t *testing.T) {
	err := errors.New("some other error")
	if IsErrorCode(err, ErrNotFound) {
		t.Fatal("expected false for non-pcloud error")
	}
}

func TestIsErrorCodeWrapped(t *testing.T) {
	inner := &Error{Result: ErrNotFound, Message: "not found"}
	wrapped := fmt.Errorf("operation failed: %w", inner)
	if !IsErrorCode(wrapped, ErrNotFound) {
		t.Fatal("expected match for wrapped error")
	}
}

func TestErrorWrappingUnwrap(t *testing.T) {
	inner := &Error{Result: ErrNotFound, Message: "File not found"}
	wrapped := fmt.Errorf("stat file 123: %w", inner)
	if !IsErrorCode(wrapped, ErrNotFound) {
		t.Fatal("IsErrorCode should find wrapped pcloud error")
	}
	var pErr *Error
	if !errors.As(wrapped, &pErr) {
		t.Fatal("errors.As should unwrap to *Error")
	}
	if pErr.Result != ErrNotFound {
		t.Fatalf("want ErrNotFound, got %d", pErr.Result)
	}
}
