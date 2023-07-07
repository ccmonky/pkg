package utils_test

import (
	"mime"
	"testing"
)

func TestParseMediaType(t *testing.T) {
	// `()<>@,;:\"/[]?=`
	_, _, err := mime.ParseMediaType("application@x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}
}
