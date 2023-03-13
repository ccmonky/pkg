package logkit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// RequestIDName field name of request id in log
var RequestIDName = "request_id"

var (
	// ZapReqID short name for ZapRequestID
	ZapReqID = ZapRequestID

	// GetReqID short name for GetRequestID
	GetReqID = GetRequestID
)

// ZapRequestID `zap.Field` to record request id which will get from `http.Request`
func ZapRequestID(r *http.Request) zap.Field {
	return zap.String(RequestIDName, GetReqID(r.Context()))
}

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds th unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

// GetRequestID get request id from context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return "-"
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return "-"
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
