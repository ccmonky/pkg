package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Cause
}

func New[T ~struct{ Cause }](err any) T {
	var e Error
	switch err := err.(type) {
	// case Cause:
	// 	e = Error{err}
	case error:
		e = Error{wrapper{err}}
	default:
		e = Error{wrapper{fmt.Errorf("%v", err)}}
	}
	return T(e)
}

type Cause interface {
	error
	Unwrap() error
}

type wrapper struct {
	error
}

func (w wrapper) Unwrap() error {
	return w.error
}

func Is[T error](err error) bool {
	if err == nil {
		return false
	}
	var e T
	return errors.As(err, &e)
}

var (
	_ error = (*Error)(nil)
)
