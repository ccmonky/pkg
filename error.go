package pkg

import "github.com/pkg/errors"

// PanicIfError panic if err != nil
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicfIfError(err error, format string, args ...interface{}) {
	if err != nil {
		panic(errors.WithMessagef(err, format, args...))
	}
}
