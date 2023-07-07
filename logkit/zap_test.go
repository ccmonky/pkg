package logkit_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/ccmonky/pkg/logkit"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapFields(t *testing.T) {
	buf := new(bytes.Buffer)
	sync := zapcore.AddSync(buf)
	cfg := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), sync, zap.InfoLevel)
	logger := zap.New(core)
	defer logger.Sync()
	logger.Info("json", logkit.ZapJSON("map", map[string]int{
		"one": 1,
		"two": 2,
	}))
	assert.Equalf(t, int64(2), gjson.Get(buf.String(), "map.two").Int(), "map.two")
	buf.Reset()
	r, err := http.NewRequest("POST", "/", nil)
	assert.Nilf(t, err, "new request err == nil")
	logger.Info("request_id", logkit.ZapRequestID(r))
	assert.Equalf(t, "-", gjson.Get(buf.String(), "request_id").String(), "request_id")
	buf.Reset()
	r = r.WithContext(context.WithValue(r.Context(), logkit.RequestIDKey, "abc"))
	logger.Info("request_id", logkit.ZapReqID(r))
	assert.Equalf(t, "abc", gjson.Get(buf.String(), "request_id").String(), "request_id")

	n := 3
	buf.Reset()
	logger.Info("", logkit.ZapAnyN("anyn", 10000, n))
	assert.Equalf(t, int64(10000), gjson.Get(buf.String(), "anyn").Int(), "int")

	buf.Reset()
	logger.Info("", logkit.ZapAnyN("anyn", []int{10000}, n))
	assert.Equalf(t, "[10...", gjson.Get(buf.String(), "anyn").String(), "int-array")

	buf.Reset()
	logger.Info("", logkit.ZapAnyN("anyn", []string{"10000"}, n))
	assert.Equalf(t, "[10...", gjson.Get(buf.String(), "anyn").String(), "string-array")

	buf.Reset()
	logger.Info("", logkit.ZapAnyN("anyn", "10000", n))
	assert.Equalf(t, "100...", gjson.Get(buf.String(), "anyn").String(), "string")

	buf.Reset()
	sp := "10000"
	logger.Info("", logkit.ZapAnyN("anyn", &sp, n))
	assert.Equalf(t, "100...", gjson.Get(buf.String(), "anyn").String(), "string-pointer")

	buf.Reset()
	logger.Info("", logkit.ZapAnyN("anyn", []byte("10000"), n))
	assert.Equalf(t, "100...", gjson.Get(buf.String(), "anyn").String(), "bytes")
}
