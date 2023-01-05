package utils

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// AfterWithSkew 在时间偏差为d的情况下判断时间a > b
func AfterWithSkew(a, b time.Time, d time.Duration) bool {
	return a.Sub(b) > -d
}

// AfterEqualWithSkew 在时间偏差为d的情况下判断时间a >= b
func AfterEqualWithSkew(a, b time.Time, d time.Duration) bool {
	return a.Sub(b) >= -d
}

// BeforeWithSkew 在时间偏差为d的情况下判断时间a < b
func BeforeWithSkew(a, b time.Time, d time.Duration) bool {
	return a.Sub(b) < d
}

// BeforeEqualWithSkew 在时间偏差为d的情况下判断时间a <= b
func BeforeEqualWithSkew(a, b time.Time, d time.Duration) bool {
	return a.Sub(b) <= d
}

// Duration 定义可json化的time.Duration
type Duration struct {
	time.Duration
}

// MarshalJSON 序列化Duration
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON 反序列化Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
