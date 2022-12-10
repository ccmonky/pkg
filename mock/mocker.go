package mock

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/ccmonky/pkg"
	"github.com/ccmonky/typemap"
)

// Matcher 根据请求匹配要使用的ResponseMocker
type Matcher interface {
	// Match 根据传入请求，计算请求特征值，根据特征值匹配ResponseMocker，并返回请求特征值和ResponseMocker返回
	// Note:
	// 1. 如果计算请求特征值出错，那么特征值返回空字符串，ResponseMocker为nil
	// 2. 如果ResponseMocker为nil，表明未匹配，调用方需处理此场景
	Match(*http.Request) (string, ResponseMocker, error)
	Eigenkey(*http.Request) (string, error)
}

// ResponseMocker 定义生成mock response的工具类接口
// Usage：
// 1. 三方扩展mocker需要将ID和自身实例注册到MetaOfResponseMocker资源;
// 2. 然后UnmarshalResponseMocker方法解析json得到实例应用即可，
//    注意，json内需要额外包含`"response_mocker": "ID"`字段
type ResponseMocker interface {
	// ID is the type name of the ResponseMocker. It
	// must be unique and properly namespaced.
	ID() string

	// New returns a pointer to a new, empty
	// instance of the ResponseMocker's type. This
	// method must not have any side-effects,
	// and no other initialization should
	// occur within it.
	New() ResponseMocker

	// IsTransparent 特化Mocker，没有Mock实现，只有一个用途：用于标识从源服务获取Mock响应
	IsTransparent() bool

	// Mock 根据请求得到Mock响应或错误
	Mock(*http.Request) (*http.Response, error)

	// Extension 返回描述公共行为的参数对象Options，包含一些内置可选参数，如latency
	Extension() *Options
}

// Option 描述一些公共行为，如latency
type Options struct {
	Latency pkg.Duration `json:"latency"`
}

func UnmarshalResponseMocker(jsonBytes []byte) (ResponseMocker, error) {
	if len(bytes.TrimSpace(jsonBytes)) == 0 {
		return new(TransparentResponseMocker), nil
	}
	id := gjson.GetBytes(jsonBytes, "response_mocker").String()
	generator, err := typemap.Get[ResponseMocker](context.Background(), id)
	if err != nil {
		return nil, errors.WithMessagef(err, "get response mocker %s failed", id)
	}
	mocker := generator.New()
	err = json.Unmarshal(jsonBytes, mocker)
	if err != nil {
		return nil, errors.WithMessagef(err, "unmarshal response mocker for %s failed", id)
	}
	return mocker, nil
}

// TransparentResponseMocker 透明Mocker，即从源服务获取真实响应作为Mock，默认如果不指定则使用此Mocker
type TransparentResponseMocker struct {
	*Options `json:"options"`
}

func (mr TransparentResponseMocker) ID() string {
	return ""
}

func (mr TransparentResponseMocker) New() ResponseMocker {
	return new(TransparentResponseMocker)
}

func (mr TransparentResponseMocker) IsTransparent() bool {
	return true
}

func (mr TransparentResponseMocker) Mock(*http.Request) (*http.Response, error) {
	panic("not implement -- TransparentResponseMocker should not use Mock method")
}

func (mr TransparentResponseMocker) Extension() *Options {
	return mr.Options
}

// ResponseMockerFromURL 请求URL的获取整个响应作为Mock，通常用于Mock平台
type ResponseMockerFromURL struct {
	*Options        `json:"options"`
	ResponseFromURL string `json:"response_from_url"`
}

func (mr ResponseMockerFromURL) ID() string {
	return "ResponseMockerFromURL"
}

func (mr ResponseMockerFromURL) New() ResponseMocker {
	return new(ResponseMockerFromURL)
}

func (mr ResponseMockerFromURL) IsTransparent() bool {
	return false
}

func (mr ResponseMockerFromURL) Mock(r *http.Request) (*http.Response, error) {
	rp, err := http.Get(mr.ResponseFromURL)
	if err != nil {
		return nil, errors.WithMessagef(err, "get mock response from %s failed", mr.ResponseFromURL)
	}
	return rp, nil
}

func (mr ResponseMockerFromURL) Extension() *Options {
	return mr.Options
}

// ResponseMockerBuilder Response Mock构造器，通过指定状态码、头和Body生成Mock响应，通常用于静态mock
type ResponseMockerBuilder struct {
	*Options    `json:"options"`
	StatusCode  int         `json:"status_code"`
	Header      http.Header `json:"header,omitempty"`
	Body        string      `json:"body,omitempty"`          // NOTE: 与BodyFromURL二选一即可，都存在以Body为主
	BodyFromURL string      `json:"body_from_url,omitempty"` // NOTE: 根据请求URL结果作为响应Body，如OSS场景
}

func (mr ResponseMockerBuilder) ID() string {
	return "ResponseMockerBuilder"
}

func (mr ResponseMockerBuilder) New() ResponseMocker {
	return new(ResponseMockerBuilder)
}

func (mr ResponseMockerBuilder) IsTransparent() bool {
	return false
}

func (mr ResponseMockerBuilder) Mock(r *http.Request) (*http.Response, error) {
	header := mr.Header
	if header == nil {
		header = http.Header{}
	}
	var body io.ReadCloser
	if mr.Body != "" {
		body = NewResponseBodyFromString(mr.Body)
	} else {
		if mr.BodyFromURL != "" {
			rp, err := http.Get(mr.BodyFromURL)
			if err != nil {
				return nil, errors.WithMessagef(err, "get mock response body from %s failed", mr.BodyFromURL)
			}
			rawBody, err := ioutil.ReadAll(rp.Body)
			if err != nil {
				return nil, errors.WithMessagef(err, "read mock response body from %s failed", mr.BodyFromURL)
			}
			body = NewResponseBodyFromBytes(rawBody)
		}
	}
	rp := &http.Response{
		Status:        strconv.Itoa(mr.StatusCode),
		StatusCode:    mr.StatusCode,
		Body:          body,
		Header:        header,
		ContentLength: -1,
	}
	return rp, nil
}

func (mr ResponseMockerBuilder) Extension() *Options {
	return mr.Options
}

func NewResponseBodyFromString(body string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBufferString(body))
}

func NewResponseBodyFromBytes(body []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(body))
}

var (
	_ ResponseMocker = (*ResponseMockerFromURL)(nil)
	_ ResponseMocker = (*ResponseMockerBuilder)(nil)
	_ ResponseMocker = (*TransparentResponseMocker)(nil)
)
