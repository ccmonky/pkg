package eigenkey

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/ccmonky/typemap"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// HTTPRequestEigenkeyGen 定义根据请求信息生成特征键的函数
type HTTPRequestEigenkeyGen func(namespace string, info *HTTPRequestInfo, postFuncs ...KeyPostFunc) string

// DefaultHTTPEigenkeyFunc 定义默认的HTTP请求特征提取键函数
func DefaultHTTPEigenkeyFunc(ns string, info *HTTPRequestInfo, postFns ...KeyPostFunc) string {
	var parts []string
	if info.RemoteAddr != "" {
		parts = append(parts, info.RemoteAddr)
	}
	if ns != "" {
		parts = append(parts, ns)
	}
	if info.Method != "" {
		parts = append(parts, info.Method)
	}
	var us string
	u := info.URL()
	if u != nil {
		us = u.String()
		if us != "" {
			parts = append(parts, us)
		}
	}
	if info.Proto != "" {
		parts = append(parts, info.Proto)
	}
	hs := info.HeaderString()
	if hs != "" {
		parts = append(parts, hs)
	}
	key := strings.Join(parts, ":")
	for _, fn := range postFns {
		key = fn(key)
	}
	return key
}

// HTTPRequestInfo 根据RequestExtractor抽取得到的关键信息
type HTTPRequestInfo struct {
	Method       string
	Scheme       string
	Host         string
	Path         string
	Fragment     string
	Username     string
	Password     string
	Proto        string
	UseArguments []string
	Arguments    url.Values
	UseHeaders   []string
	Headers      url.Values
	RemoteAddr   string
}

// QueryString 根据UseArguments和Arguments生成RawQuery
func (i HTTPRequestInfo) QueryString() string {
	if i.Arguments == nil {
		return ""
	}
	var buf strings.Builder
	for _, k := range i.UseArguments {
		vs := i.Arguments[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

// HeaderString 根据UseHeaders和Headers生成HeaderString
func (i HTTPRequestInfo) HeaderString() string {
	if i.Headers == nil {
		return ""
	}
	var buf strings.Builder
	for _, k := range i.UseHeaders {
		vs := i.Headers[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

// URL 根据HTTPRequestInfo生成URL
func (i HTTPRequestInfo) URL() *url.URL {
	var user *url.Userinfo
	if i.Username != "" {
		if i.Password != "" {
			user = url.UserPassword(i.Username, i.Password)
		} else {
			user = url.User(i.Username)
		}
	}
	u := &url.URL{
		Scheme:   i.Scheme,
		Host:     i.Host,
		Path:     i.Path,
		Fragment: i.Fragment,
		User:     user,
		RawQuery: i.QueryString(),
	}

	return u
}

// HTTPRequestExtractor 根据http请求抽取RequestInfo
type HTTPRequestExtractor struct {
	UseMethod     bool     `json:"use_method"`
	UseScheme     bool     `json:"use_scheme"`
	UseHost       bool     `json:"use_host"`
	UsePath       bool     `json:"use_path"`
	UseFragment   bool     `json:"use_fragment"`
	UseUsername   bool     `json:"use_username"`
	UsePassword   bool     `json:"use_password"`
	UseProto      bool     `json:"use_proto"`
	UseArguments  []string `json:"use_arguments"`
	UseHeaders    []string `json:"use_headers"`
	UseRemoteAddr bool     `json:"use_remote_addr"`
}

// Extract 抽取HTTP特征
func (e HTTPRequestExtractor) Extract(r *http.Request) (*HTTPRequestInfo, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	info := &HTTPRequestInfo{
		UseArguments: e.UseArguments,
		UseHeaders:   e.UseHeaders,
	}
	if len(info.UseArguments) > 0 {
		info.Arguments = make(url.Values)
	}
	if len(info.UseHeaders) > 0 {
		info.Headers = make(url.Values)
	}
	if e.UseMethod {
		info.Method = r.Method
	}
	if r.URL != nil {
		if e.UseScheme {
			info.Scheme = r.URL.Scheme
		}
		if e.UseHost {
			info.Host = r.URL.Host
		}
		if e.UsePath {
			info.Path = r.URL.Path
		}
		if e.UseFragment {
			info.Fragment = r.URL.Fragment
		}
		if r.URL.User != nil {
			if e.UseUsername {
				info.Username = r.URL.User.Username()
			}
			if e.UsePassword {
				pwd, ok := r.URL.User.Password()
				if ok {
					info.Password = pwd
				}
			}
		}
	}
	if e.UseProto {
		info.Proto = r.Proto
	}
	for _, arg := range info.UseArguments {
		info.Arguments[arg] = r.Form[arg]
	}

	for _, hdr := range info.UseHeaders {
		info.Headers[hdr] = r.Header.Values(hdr)
	}

	if e.UseRemoteAddr {
		info.RemoteAddr = r.RemoteAddr
	}
	return info, nil
}

// HTTPRequestEigenkeyExtractor http请求特征提取器
type HTTPRequestEigenkeyExtractor struct {
	Namespace        string                `json:"namespace"`
	KeyFuncName      string                `json:"key_func_name"`
	KeyPostFuncNames []string              `json:"key_post_func_names"`
	RequestExtractor *HTTPRequestExtractor `json:"request_extractor"`
	CleanPath        bool                  `json:"clean_path"`

	keyFn      HTTPRequestEigenkeyGen
	keyPostFns []KeyPostFunc
}

// Provision 初始化
func (g *HTTPRequestEigenkeyExtractor) Provision() error {
	var err error
	g.keyFn, err = typemap.Get[HTTPRequestEigenkeyGen](context.Background(), g.KeyFuncName)
	if err != nil {
		return err
	}
	if g.keyFn == nil {
		return errors.Errorf("http request eigenkey func %s is nil", g.KeyFuncName)
	}
	for _, postName := range g.KeyPostFuncNames {
		fn, err := typemap.Get[KeyPostFunc](context.Background(), postName)
		if err != nil {
			return err
		}
		if fn != nil {
			g.keyPostFns = append(g.keyPostFns, fn)
		}
	}
	if g.RequestExtractor == nil {
		g.RequestExtractor = &HTTPRequestExtractor{
			UsePath: true,
		}
	}
	return nil
}

// Eigenkey 从给定的请求中提取Eigenkey
func (g HTTPRequestEigenkeyExtractor) Eigenkey(r *http.Request) (string, error) {
	info, err := g.RequestExtractor.Extract(r)
	if err != nil {
		return "", err
	}
	if g.CleanPath && g.RequestExtractor.UsePath {
		info.Path = httprouter.CleanPath(info.Path)
	}
	return g.keyFn(g.Namespace, info, g.keyPostFns...), nil
}
