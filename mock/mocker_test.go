package mock_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccmonky/pkg/mock"
)

func TestUnmarshalResponseMocker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Demo", "demo")
		io.WriteString(w, "hello")
	}))
	defer ts.Close()

	var cases = []struct {
		data     []byte
		id       string
		body     string
		hdrKey   string
		hdrValue string
	}{
		{
			[]byte(""),
			"",
			"",
			"",
			"",
		},
		{
			[]byte(fmt.Sprintf(`{
				"response_mocker": "ResponseMockerFromURL",
				"response_from_url": "%s",
				"options": {
					"latency": "3ms"
				}
			}`, ts.URL)),
			"ResponseMockerFromURL",
			"hello",
			"X-Demo",
			"demo",
		},
		{
			[]byte(`{
				"response_mocker": "ResponseMockerBuilder",
				"status_code": 200,
				"header": {
					"X-Test": ["abc"]
				},
				"body": "this is a test",
				"options": {
					"latency": "3ms"
				}
			}`),
			"ResponseMockerBuilder",
			"this is a test",
			"X-Test",
			"abc",
		},
		{
			[]byte(fmt.Sprintf(`{
				"response_mocker": "ResponseMockerBuilder",
				"status_code": 200,
				"header": {
					"X-Test": ["123"]
				},
				"body_from_url": "%s",
				"options": {
					"latency": "3ms"
				}
			}`, ts.URL)),
			"ResponseMockerBuilder",
			"hello",
			"X-Test",
			"123",
		},
	}

	dummyRq, err := http.NewRequest("GET", "http://www.example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range cases {
		mocker, err := mock.UnmarshalResponseMocker(tc.data)
		if err != nil {
			t.Fatal(err)
		}
		if mocker.ID() != tc.id {
			t.Fatal("should ==")
		}
		if !mocker.IsTransparent() {
			rp, err := mocker.Mock(dummyRq)
			if err != nil {
				t.Fatal(err)
			}
			body, err := ioutil.ReadAll(rp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != tc.body {
				t.Fatalf("%s should == %s", string(body), tc.body)
			}
			if rp.Header.Get(tc.hdrKey) != tc.hdrValue {
				t.Fatal("should ==")
			}
		}
		if mocker.Extension() != nil && mocker.Extension().Latency.Duration > 0 {
			time.Sleep(mocker.Extension().Latency.Duration)
		}
	}
}
