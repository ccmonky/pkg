package errors_test

import (
	"errors"
	"testing"

	pkgerrors "github.com/ccmonky/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type errorWrapper struct {
	error
	msg string
}

func (w errorWrapper) Unwrap() error {
	return w.error
}

func (w errorWrapper) Error() string {
	return w.msg + ": " + w.error.Error()
}

func TestIsError(t *testing.T) {
	type NotFoundError pkgerrors.Error
	type AlreadyExistsError pkgerrors.Error
	errNfe := NotFoundError{pkgerrors.New("xxx not found")}
	errAee := AlreadyExistsError{pkgerrors.New("xxx already exists")}
	assert.Equalf(t, errNfe.Error(), "xxx not found", "not found error")
	assert.Equalf(t, errAee.Error(), "xxx already exists", "already exists error")
	assert.Truef(t, pkgerrors.Is[NotFoundError](errNfe), "errNfe is notfound")
	assert.Truef(t, !pkgerrors.Is[AlreadyExistsError](errNfe), "errNfe is not already existsnotfound")
	assert.Truef(t, pkgerrors.Is[AlreadyExistsError](errAee), "errAee is already")
	assert.Truef(t, !pkgerrors.Is[NotFoundError](errAee), "errAee is not notfound")
	errNfe2 := NotFoundError{errorWrapper{
		error: errors.New("hahaha"),
		msg:   "wrapper",
	}}
	assert.Truef(t, pkgerrors.Is[NotFoundError](errNfe2), "errNfe2 is notfound")
	assert.Truef(t, !pkgerrors.Is[AlreadyExistsError](errNfe2), "errNfe2 is not already existsnotfound")
	errNfeW := errorWrapper{
		error: errNfe,
		msg:   "wrapper",
	}
	errAeeW := errorWrapper{
		error: errAee,
		msg:   "wrapper",
	}
	assert.Equalf(t, errNfeW.Error(), "wrapper: xxx not found", "not found error")
	assert.Equalf(t, errAeeW.Error(), "wrapper: xxx already exists", "already exists error")
	assert.Truef(t, pkgerrors.Is[NotFoundError](errNfeW), "errNfeW is notfound")
	assert.Truef(t, !pkgerrors.Is[AlreadyExistsError](errNfeW), "errNfeW is not already")
	assert.Truef(t, pkgerrors.Is[AlreadyExistsError](errAeeW), "errAeeW is already")
	assert.Truef(t, !pkgerrors.Is[NotFoundError](errAeeW), "errAeeW is not already")
	type XXXError pkgerrors.Error
	errXXX := XXXError{pkgerrors.New("xxx error")}
	errNfXXX := NotFoundError{pkgerrors.New(errXXX)}
	assert.Truef(t, pkgerrors.Is[NotFoundError](errNfXXX), "errnfxxx is notfound error")
	assert.Truef(t, pkgerrors.Is[XXXError](errNfXXX), "errnfxxx is also xxx error")
	assert.Truef(t, !pkgerrors.Is[AlreadyExistsError](errNfXXX), "errnfxxx is not already exists error")
}
