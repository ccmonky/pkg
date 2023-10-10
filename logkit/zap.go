package logkit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ccmonky/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RequestIDName field name of request id in log
var RequestIDName = "request_id"

var RequestIDKey = utils.RequestIDKey

var (
	// ZapReqID short name for ZapRequestID
	ZapReqID = ZapRequestID

	// GetReqID short name for GetRequestID
	GetReqID = utils.GetRequestID
)

// ZapRequestID `zap.Field` to record request id which will get from `http.Request`
func ZapRequestID(r *http.Request) zap.Field {
	return zap.String(RequestIDName, GetReqID(r.Context()))
}

// ZapJSON zap.Field will marshal obj as json to record
func ZapJSON(key string, obj interface{}) zap.Field {
	return zap.Reflect(key, &objectJsonMarshaler{obj: obj})
}

type objectJsonMarshaler struct {
	obj interface{}
}

func (j *objectJsonMarshaler) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(j.obj)
	if err != nil {
		return nil, fmt.Errorf("json marshaling failed: %w", err)
	}
	return bytes, nil
}

// ZapAnyN like zap.Any, but limit to n bytes for non-primitive type values
func ZapAnyN(key string, value interface{}, n int) zap.Field {
	switch val := value.(type) {
	case zapcore.ObjectMarshaler, zapcore.ArrayMarshaler, []bool, []complex128, []complex64, []float64, []float32,
		[]int, []int64, []int32, []int16, []int8, []string, []uint, []uint64, []uint32, []uint16, []uintptr, []time.Time,
		[]time.Duration, []error:
		s := fmt.Sprintf("%v", value)
		if len(s) > n {
			return zap.String(key, s[:n]+"...")
		}
		return zap.Any(key, value)
	case string:
		if len(val) > n {
			return zap.String(key, val[:n]+"...")
		}
		return zap.String(key, val)
	case *string:
		if val == nil {
			return zap.String(key, "<nil>")
		}
		if len(*val) > n {
			return zap.String(key, (*val)[:n]+"...")
		}
		return zap.String(key, *val)
	case []byte:
		if val == nil {
			return zap.ByteString(key, []byte("<nil>"))
		}
		if len(val) > n {
			return zap.ByteString(key, append(val[:n], []byte("...")...))
		}
		return zap.ByteString(key, val)
	default:
		return zap.Any(key, value)
	}
}
