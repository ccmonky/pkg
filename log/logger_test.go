package log_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/ccmonky/pkg/log"
)

func TestStdLogger(t *testing.T) {
	var cases = []struct {
		level   log.Level
		out     io.Writer
		lines   int
		keyword string
	}{
		{
			log.DebugLevel,
			new(bytes.Buffer),
			3,
			"[logger -1]",
		},
		{
			log.InfoLevel,
			new(bytes.Buffer),
			2,
			"[logger 0]",
		},
		{
			log.ErrorLevel,
			new(bytes.Buffer),
			1,
			"[logger 2]",
		},
	}
	for _, tc := range cases {
		logger := log.NewStdLogger(
			log.WithLevel(tc.level),
			log.WithOut(tc.out),
		)
		logger.Debug("msg 1", "logger", tc.level)
		logger.Info("msg 2", "logger", tc.level)
		logger.Error("msg 3", "logger", tc.level)
		lines := strings.Split(strings.TrimSpace(tc.out.(*bytes.Buffer).String()), "\n")
		if len(lines) != tc.lines {
			t.Log(lines)
			t.Fatalf("level %v, should == %d, got %d", tc.level, tc.lines, len(lines))
		}
		for _, line := range lines {
			if !strings.Contains(line, tc.keyword) {
				t.Fatal("should contain")
			}
		}
	}
}
