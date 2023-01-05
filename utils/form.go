package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

// ErrNotImplemented 未实现错误
var ErrNotImplemented = errors.New("not implemented")

// EncodeFormToBody 将r.Form回填到r.Body，一般用于反代，目前仅支持PostForm，Multipart未实现
func EncodeFormToBody(r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	ct, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}
	if r.MultipartForm != nil && (ct == "multipart/form-data" || ct == "multipart/mixed") {
		return encodeMultipartFormToBodyWithContentType(r, &ct)
	}
	return encodePostFormToBodyWithContentType(r, &ct)
}

// EncodePostFormToBody 判断是否执行了request.ParseForm，如果执行了那么执行恢复request.Body，一般用于反向代理
func EncodePostFormToBody(r *http.Request) error {
	return encodePostFormToBodyWithContentType(r, nil)
}

func encodePostFormToBodyWithContentType(r *http.Request, contentType *string) error {
	if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
		// NOTE：r.ParseForm只会对Post、Put和Patch方法执行parsePostForm方法，如果不是这三种方法，可以不必考虑！
		return nil
	}
	ct, err := parseContentType(r, contentType)
	if err != nil {
		return err
	}
	if ct != "application/x-www-form-urlencoded" {
		// NOTE: ParseForm中的parsePostForm方法不会执行，因此也不必恢复Body！
		return nil
	}
	if r.PostForm == nil && r.Body != nil && !IsBodyDrained(r) {
		// NOTE: 如果PostForm为nil并且读取Body不是EOF或ErrBodyReadAfterClose那么认为未执行过ParseForm，直接返回！
		// case 1: 未执行ParseForm！
		// case 2: 执行了ParseForm，但把PostForm设为nil了，同时把Body回填了！
		return nil
	}
	if len(r.PostForm) == 0 {
		// 原则：基于PostForm现状进行Body回填！
		// case 1: r.PostForm不是nil但是是空的，原始r.Body为空
		// case 2: r.PostForm不是nil但是是空的，原始r.Body不为空，但是PostForm被清空！
		// case 3: r.PostForm == nil && r.Body == nil
		// case 4: r.PostForm == nil && r.Body != nil, 原始body为空！
		r.ContentLength = int64(0)
		r.Header.Set("Content-Length", strconv.Itoa(0))
		return nil
	}
	if r.Body == nil {
		// NOTE: 如果len(PostForm)>0，但是Body设为nil了，直接回填！
		return encodePostFormToBody(r)
	}
	// NOTE: 下面处理len(r.PostForm)>0&&r.Body!=nil的场景
	err = TryRead(r)
	if err != nil {
		if err == io.EOF || err == http.ErrBodyReadAfterClose {
			err := r.Body.Close() // NOTE: 重复关闭不会报错
			if err != nil {
				return err
			}
			err = encodePostFormToBody(r)
			if err != nil {
				return err
			}
		} else {
			return err // FIXME: 什么场景？
		}
	}
	// NOTE：既然读取未出现EOF或ErrBodyReadAfterClose，那么认为之前的步骤已经回填过Body了，此处仅返回nil
	return nil
}

// IsBodyDrained 判断请求的Body是否已读取，如果返回io.EOF或http.ErrBodyReadAfterClose认为已读取，另外，如果Body为nil，也返回true
func IsBodyDrained(sr *http.Request) bool {
	if sr == nil {
		panic("nil server request")
	}
	if sr.Body == nil {
		return true // FIXME: 对于server request不会是nil，如果改成nil了，那么认为r.Body已消耗？
	}
	err := TryRead(sr)
	if err == io.EOF || err == http.ErrBodyReadAfterClose {
		return true
	}
	// FIXME: err为非nil的其他错误改如何认定？
	return false
}

func parseContentType(r *http.Request, contentType *string) (string, error) {
	var ct string
	if contentType != nil {
		ct = *contentType
	} else {
		ct = r.Header.Get("Content-Type")
	}
	// RFC 7231, section 3.1.1.5 - empty type
	//   MAY be treated as application/octet-stream
	if ct == "" {
		ct = "application/octet-stream"
	}
	var err error
	ct, _, err = mime.ParseMediaType(ct)
	return ct, err
}

func encodePostFormToBody(r *http.Request) error {
	if r == nil {
		return errors.New("request is nil")
	}
	if r.PostForm == nil {
		return errors.New("request post form is nil")
	}
	body := r.PostForm.Encode()
	if r.ContentLength > 0 && r.ContentLength != int64(len(body)) {
		// 场景：
		// 1. 改写Content-Length, 比如用户执行过r.PostForm.Del("xxx");
		// 2. MultipartForm也会填充PostForm，如何处理？
		r.ContentLength = int64(len(body))
		r.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}
	r.Body = ioutil.NopCloser(strings.NewReader(body))
	return nil
}

// EncodeMultipartFormToBody 未实现，最好不要用于反代！
func EncodeMultipartFormToBody(r *http.Request) error {
	return encodeMultipartFormToBodyWithContentType(r, nil)
}

func encodeMultipartFormToBodyWithContentType(r *http.Request, contentType *string) error {
	ct, err := parseContentType(r, contentType)
	if err != nil {
		return err
	}
	if ct != "multipart/form-data" && ct != "multipart/mixed" {
		// NOTE: ParseMultipartForm方法不会执行，因此也不必恢复Body！
		return nil
	}
	return ErrNotImplemented
}
