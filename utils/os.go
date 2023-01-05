package utils

import (
	"os"

	"github.com/pkg/errors"
)

func AtomicWriteFile(file, content string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return errors.WithMessagef(err, "file %s already exists", file)
		} else {
			return errors.WithMessagef(err, "create file %s failed", file)
		}
	} else {
		defer f.Close()
		_, err := f.WriteString(content)
		if err != nil {
			return errors.WithMessagef(err, "write file %s with content %s", file, content)
		}
	}
	return nil
}
