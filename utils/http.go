package utils

import (
	"bufio"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func TryRead(r *http.Request) error {
	_, err := TryReadByte(r)
	return err
}

func TryReadByte(r *http.Request) (b byte, err error) {
	if r == nil {
		return 0, errors.New("nil http request")
	}
	onebuf := make([]byte, 1)
	n, err := r.Body.Read(onebuf)
	if n > 0 {
		b = onebuf[0]
		r.Body = io.NopCloser(NewAheadReader(onebuf, r.Body))
	}
	return
}

func NewAheadReader(ahead []byte, reader io.Reader) io.Reader {
	return &AheadReader{
		ahead:  ahead,
		Reader: reader,
	}
}

type AheadReader struct {
	ahead []byte
	io.Reader
}

func (ar *AheadReader) Read(p []byte) (n int, err error) {
	var i int
	if len(ar.ahead) > 0 {
		i = copy(p, ar.ahead)
		ar.ahead = nil
	}
	n, err = ar.Reader.Read(p[i:])
	return n + i, err
}

func TryBufRead(r *http.Request) error {
	_, err := TryBufReadByte(r)
	return err
}

// TryGetFirstByteOfRequestBody try to read first byte of http request body, and then fill back to body
func TryBufReadByte(r *http.Request) (b byte, err error) {
	if r == nil {
		return 0, errors.New("nil http request")
	}
	br := bufio.NewReader(r.Body)
	b, err = br.ReadByte()
	if err != nil {
		return
	}
	err = br.UnreadByte()
	if err != nil {
		return
	}
	r.Body = io.NopCloser(br)
	return
}

var methods = map[string]struct{}{
	http.MethodGet:     struct{}{},
	http.MethodHead:    struct{}{},
	http.MethodPost:    struct{}{},
	http.MethodPut:     struct{}{},
	http.MethodPatch:   struct{}{},
	http.MethodDelete:  struct{}{},
	http.MethodConnect: struct{}{},
	http.MethodOptions: struct{}{},
	http.MethodTrace:   struct{}{},
}

// IsHTTPMethodValid test if method is valid http method
func IsHTTPMethodValid(method string) bool {
	if _, ok := methods[strings.ToUpper(method)]; ok {
		return true
	}
	return false
}

var (
	_ io.Reader = (*AheadReader)(nil)
)
