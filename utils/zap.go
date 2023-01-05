package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const (
	// RequestIDName defines the name of request ID in log
	RequestIDName = "request_id"
)

// ZapJSON 使用json序列化对象的zao.Field
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

// ZapRequestID 获取http请求ID的zao.Field
func ZapRequestID(r *http.Request) zap.Field {
	return zap.String(RequestIDName, GetReqID(r.Context()))
}
