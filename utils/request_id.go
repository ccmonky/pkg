package utils

import (
	"context"
)

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds th unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

// GetReqID 用于从context中提取请求ID
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return "-"
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return "-"
}
