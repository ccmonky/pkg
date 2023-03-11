package logkit

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

// Options defines the options for TimedRotatingLogger
type Options struct {
	CronSpec string
	Cron     *cron.Cron
}

// Option defines the option func
type Option func(*Options)

// WithCronSpec set cron expression, refer to `https://www.godoc.org/github.com/robfig/cron`
func WithCronSpec(spec string) Option {
	return func(opts *Options) {
		opts.CronSpec = spec
	}
}

// WithCron 设定cron调度器
func WithCron(c *cron.Cron) Option {
	return func(opts *Options) {
		opts.Cron = c
	}
}

// NewTimedRotatingLogger creates new TimedRotatingLogger
func NewTimedRotatingLogger(logger *Logger, opts ...Option) *TimedRotatingLogger {
	options := &Options{
		CronSpec: "0 0 0 * * *", // midnight
	}
	for _, o := range opts {
		o(options)
	}
	if options.Cron == nil {
		options.Cron = cron.New()
	}

	tl := &TimedRotatingLogger{
		Logger:   logger,
		cronSpec: options.CronSpec,
		cron:     options.Cron,
	}
	err := tl.cron.AddFunc(tl.cronSpec, func() { tl.Rotate() })
	if err != nil {
		panic(fmt.Errorf("bad cron expreesion: %v", err))
	}
	tl.cron.Start() // NOTE：重复启动是安全的

	return tl
}

// TimedRotatingLogger rotates log according to cron expression
type TimedRotatingLogger struct {
	*Logger

	cronSpec string
	cron     *cron.Cron
}

// Close close the logger and stop the scheduler
func (l *TimedRotatingLogger) Close() error {
	l.cron.Stop()
	return l.Logger.Close()
}

// BackupName defines the backup name for gaode
func BackupName(name string, local bool) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	t := currentTime().Add(-time.Hour) // 备份时间应该是该小时的起始时间
	if !local {
		t = t.UTC()
	}

	timestamp := t.Format("2006-01-02-15")
	return filepath.Join(dir, fmt.Sprintf("%s.%s", filename, timestamp))
}

// TimeFromName extract time from backup name for gaode
func TimeFromName(filename, prefix, ext string) (time.Time, error) {
	if strings.HasSuffix(prefix, "-") {
		prefix = prefix[:len(prefix)-1]
	}
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	ts := filepath.Ext(filename)
	if strings.HasPrefix(ts, ".") {
		ts = ts[1:]
	}

	return time.ParseInLocation("2006-01-02-15", ts, time.FixedZone("UTC+8", 8*60*60))
}
