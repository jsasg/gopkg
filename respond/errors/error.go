package errors

import (
	stderrors "errors"
)

type RespondError struct {
	s string
}

func New(text string) error {
	return &RespondError{text}
}

func (e *RespondError) Error() string {
	return e.s
}

func Is(err error, target error) bool {
	return stderrors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

func Join(err ...error) error {
	return stderrors.Join(err...)
}

func Wrap(err error, text string) error {
	return Join(err, New(text))
}

func Unwrap(err error) []error {
	if e, ok := err.(interface{ Unwrap() []error }); ok {
		return e.Unwrap()
	}
	if e, ok := err.(interface{ Unwrap() error }); ok {
		return []error{e.Unwrap()}
	}

	return []error{err}
}
