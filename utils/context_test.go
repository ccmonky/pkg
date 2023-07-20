package utils_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccmonky/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestWithoutCancel(t *testing.T) {
	// test WithValue
	key := struct{}{}
	ctx := context.WithValue(context.Background(), key, "value")
	ctx = utils.WithoutCancel(ctx)
	assert.Equal(t, ctx.Value(key), "value")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
		io.WriteString(w, "ok")
	}))
	defer ts.Close()

	// test WithCancel
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	testCancel(t, ts, ctx)

	// test WithoutCancel
	ctx = utils.WithoutCancel(ctx)
	testWithoutCancel(t, ts, ctx)

	// test WithCancel again
	ctx, cancel = context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	testCancel(t, ts, ctx)
}

func testCancel(t *testing.T, ts *httptest.Server, ctx context.Context) {
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)
	client := &http.Client{}
	rp, err := client.Do(req)
	assert.Nil(t, rp)
	assert.NotNil(t, err)
}

func testWithoutCancel(t *testing.T, ts *httptest.Server, ctx context.Context) {
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)
	client := &http.Client{}
	rp, err := client.Do(req)
	assert.Nil(t, err)
	defer rp.Body.Close()
	data, err := ioutil.ReadAll(rp.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(data), "ok")
}
