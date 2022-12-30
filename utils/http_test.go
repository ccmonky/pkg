package pkg_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/ccmonky/pkg"
	"github.com/stretchr/testify/assert"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/ccmonky/pkg
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkGetTypeId-12    	1000000000	         0.3010 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMakeZeroSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = make([]byte, 0, 0)
	}
}

// goos: darwin
// goarch: amd64
// pkg: github.com/ccmonky/pkg
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkNewRequest-12    	 2527928	       461.7 ns/op	     552 B/op	       7 allocs/op
func BenchmarkNewRequest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("12345"))
	}
}

// goos: darwin
// goarch: amd64
// pkg: github.com/ccmonky/pkg
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkReadBody-12    	139470220	         8.719 ns/op	       0 B/op	       0 allocs/op
func BenchmarkReadBody(b *testing.B) {
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("12345"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		zerobuf := make([]byte, 0, 0)
		_, _ = r.Body.Read(zerobuf)
	}
}

func TestReadBody(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("12345"))
	assert.Nilf(t, err, "new request err")
	zerobuf := make([]byte, 0, 0)
	n, err := r.Body.Read(zerobuf)
	assert.Equalf(t, 0, n, "read bytes count")
	assert.Nilf(t, err, "read body err")
}

// goos: darwin
// goarch: amd64
// pkg: github.com/ccmonky/pkg
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkReadBodyEOF-12    	149360746	         7.734 ns/op	       0 B/op	       0 allocs/op
func BenchmarkReadBodyEOF(b *testing.B) {
	r, _ := http.NewRequest(http.MethodPost, "/", ioutil.NopCloser(strings.NewReader("12345")))
	_, _ = io.ReadAll(r.Body)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		zerobuf := make([]byte, 0, 0)
		_, _ = r.Body.Read(zerobuf)
	}
}

func TestReadBodyEOF(t *testing.T) {
	for _, body := range []io.ReadCloser{
		io.NopCloser(strings.NewReader("12345")),
		io.NopCloser(bytes.NewReader([]byte("12345"))),
	} {
		r, err := http.NewRequest(http.MethodPost, "/", body)
		assert.Nilf(t, err, "new request err")
		data, err := io.ReadAll(r.Body)
		assert.Nilf(t, err, "read all err")
		assert.Equalf(t, "12345", string(data), "req body")
		zerobuf := make([]byte, 0, 0)
		n, err := r.Body.Read(zerobuf)
		assert.Equalf(t, 0, n, "read bytes count, %v", err)
		assert.Equalf(t, io.EOF, err, "read body err, count %d, %v", n, err)
	}
}

func TestReadBodyEOFWithBuffer(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString("12345"))
	r, err := http.NewRequest(http.MethodPost, "/", body)
	assert.Nilf(t, err, "new request err")
	data, err := io.ReadAll(r.Body)
	assert.Nilf(t, err, "read all err")
	assert.Equalf(t, "12345", string(data), "req body")
	zerobuf := make([]byte, 1) // NOTE: 使用make([]byte, 0, 0)判断会失败！
	n, err := r.Body.Read(zerobuf)
	assert.Equalf(t, io.EOF, err, "read body err, count %d, %v", n, err)
}

func TestTryReadByte(t *testing.T) {
	b, err := pkg.TryBufReadByte(nil)
	assert.Equalf(t, uint8(0), b, "first byte")
	assert.NotNilf(t, err, "nil request")

	b, err = pkg.TryReadByte(nil)
	assert.Equalf(t, uint8(0), b, "first byte")
	assert.NotNilf(t, err, "nil request")
}

func TestBufio(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString("12345"))
	r, err := http.NewRequest(http.MethodPost, "/", body)
	assert.Nilf(t, err, "new request err")
	data, err := io.ReadAll(r.Body)
	assert.Nilf(t, err, "read all err")
	assert.Equalf(t, "12345", string(data), "req body")
	err = pkg.TryRead(r)
	assert.Equalf(t, io.EOF, err, "read body err %v", err)

	body = io.NopCloser(bytes.NewBufferString("12345"))
	r, err = http.NewRequest(http.MethodPost, "/", body)
	assert.Nilf(t, err, "new request err")
	buf := make([]byte, 2)
	n, err := r.Body.Read(buf)
	assert.Nilf(t, err, "read two err")
	assert.Equalf(t, 2, n, "read two")
	err = pkg.TryRead(r)
	assert.Nilf(t, err, "read third byte err")
	data, err = io.ReadAll(r.Body)
	assert.Nilf(t, err, "read all err")
	assert.Equalf(t, "345", string(data), "req body left")

	err = pkg.TryRead(r)
	assert.Equalf(t, io.EOF, err, "read body err %v", err)
}

// goos: darwin
// goarch: amd64
// pkg: github.com/ccmonky/pkg
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkTryBufRead-12    	 1108556	      1159 ns/op	    4704 B/op	       8 allocs/op
func BenchmarkTryBufRead(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r, _ := http.NewRequest(http.MethodPost, "/", ioutil.NopCloser(strings.NewReader("12345")))
		pkg.TryBufRead(r)
	}
}

// goos: darwin
// goarch: amd64
// pkg: github.com/ccmonky/pkg
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkTryRead-12    	 2460674	       467.2 ns/op	     561 B/op	       8 allocs/op
func BenchmarkTryRead(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r, _ := http.NewRequest(http.MethodPost, "/", ioutil.NopCloser(strings.NewReader("12345")))
		pkg.TryRead(r)
	}
}
