package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Cause
}

func New(err any) Cause {
	switch err := err.(type) {
	case error:
		return wrapper{err}
	default:
		return wrapper{fmt.Errorf("%v", err)}
	}
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
