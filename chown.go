//go:build !linux
// +build !linux

package pkg

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
