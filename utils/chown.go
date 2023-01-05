//go:build !linux
// +build !linux

package utils

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
