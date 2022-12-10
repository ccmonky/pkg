package baggage

import (
	"errors"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

// CanonicalKey 规范化Key，ParamKey与此一致
// 1. 全部使用中划线，用户输入的下划线自动转化为中划线(如`mozi_tid`->`x-tproxy-user-mozi-tid`)
// 2. 全小写，内部保存的key也是全小写，为了查询变量方便；
//
// NOTE: 使用中划线的原因是为了单一化规则，便于记忆
func CanonicalKey(k string) string {
	return strings.ToLower(strings.ReplaceAll(k, "_", "-"))
}

// CanonicalHeaderKey 规范化HeaderKey
// 1. 全部使用中划线，用户输入的下划线自动转化为中划线(如`mozi_tid`->`X-Tproxy-User-Mozi-Tid`)
// 2. Header使用texproto.CanonicalMIMEHeaderKey规范化
//
// NOTE: 使用中划线的原因是默认nginx的underscores_in_headers默认为off，即不透传包含下划线的头，导致穿透困难！
func CanonicalHeaderKey(k string) string {
	return textproto.CanonicalMIMEHeaderKey(strings.ReplaceAll(k, "_", "-"))
}

// New 新建Baggage
func New(domains ...string) *Baggage {
	domainPrefix := strings.Join(domains, "-") + "-"
	user := &Baggage{
		HeaderPrefix: CanonicalHeaderKey(domainPrefix),
		ParamPrefix:  CanonicalKey(domainPrefix),
	}
	return user
}

// Option 用于添加附加属性，也可以直接使用WithAttr指定key添加
// 使用场景主要是：某些用户域内置属性但是不是必须属性，可定义Option，初始化时传入，使用者无需记住属性名
type Option func(*Baggage)

// Baggage 类似于opentracing的Baggage，用于跨越网络边界传递信息
type Baggage struct {
	Info         map[string]string
	HeaderPrefix string
	ParamPrefix  string
	Errs         []error
}

// WithHeaderPrefix 设定HeaderPrefix
// NOTE: 初始化阶段使用，否则会提取不到数据
func (u *Baggage) WithHeaderPrefix(v string) *Baggage {
	u.HeaderPrefix = CanonicalHeaderKey(v)
	return u
}

// WithParamPrefix 设定ParamPrefix
// NOTE: 初始化阶段使用，否则会提取不到数据
func (u *Baggage) WithParamPrefix(v string) *Baggage {
	u.ParamPrefix = CanonicalKey(v)
	return u
}

// WithAttr 设定AuthBackend
func (u *Baggage) WithAttr(k, v string) *Baggage {
	if u.Info == nil {
		u.Info = make(map[string]string)
	}
	u.Info[CanonicalKey(k)] = v
	return u
}

// WithOption 使用Option设定属性
func (u *Baggage) WithOption(opt Option) *Baggage {
	opt(u)
	return u
}

// WithInfo 设定Info
func (u *Baggage) WithInfo(info map[string]string) *Baggage {
	if info == nil {
		u.Info = nil
	}
	for k, v := range info {
		u.Info[CanonicalKey(k)] = v
	}
	return u
}

// ExtendInfo 追加info包含的字段，同名会覆盖！
func (u *Baggage) ExtendInfo(info map[string]string) *Baggage {
	for k, v := range info {
		u.Info[CanonicalKey(k)] = v
	}
	return u
}

// WithError 追加error到Errs
func (u *Baggage) WithError(v error) *Baggage {
	u.Errs = append(u.Errs, v)
	return u
}

// Headers 将用户信息转换为HTTP头附加的前缀
func (u Baggage) Headers() url.Values {
	vs := url.Values{}
	for k, v := range u.Info {
		vs.Set(CanonicalHeaderKey(u.HeaderPrefix+k), v)
	}
	return vs
}

// InjectHeaders 将用户信息转换为头注入到http.Request上
func (u Baggage) InjectHeaders(r *http.Request) {
	if r == nil {
		u.Errs = append(u.Errs, errors.New("SetHeaders: request is nil"))
		return
	}
	headers := u.Headers()
	for k := range headers {
		r.Header.Set(k, headers.Get(k))
	}
}

// Params 将用户信息转换为HTTP请求参数附加的前缀
func (u Baggage) Params() url.Values {
	vs := url.Values{}
	for k, v := range u.Info {
		vs.Set(CanonicalKey(u.ParamPrefix+k), v)
	}
	return vs
}

// InjectParams 将用户信息转换为参数注入到http.Request上
func (u Baggage) InjectParams(r *http.Request) {
	if r == nil {
		u.Errs = append(u.Errs, errors.New("SetParams: request is nil"))
		return
	}
	vs := r.URL.Query()
	params := u.Params()
	for k := range params {
		v := params.Get(k)
		vs.Set(k, v)
		if r.Form != nil {
			r.Form.Set(k, v)
		}
	}
	r.URL.RawQuery = vs.Encode()
}

// Extract 从请求提取用户信息
func (u *Baggage) Extract(r *http.Request, domains ...string) *Baggage {
	if r == nil {
		u.Errs = append(u.Errs, errors.New("Extract: request is nil"))
		return u
	}
	if r.Form == nil {
		_ = r.FormValue("") // NOTE: 解析Form
	}
	if u.Info == nil {
		u.Info = make(map[string]string)
	}
	for k := range r.Header {
		if strings.HasPrefix(CanonicalHeaderKey(k), u.HeaderPrefix) {
			u.Info[CanonicalKey(k[len(u.HeaderPrefix):])] = r.Header.Get(k)
		}
	}
	for k := range r.Form {
		ck := CanonicalKey(k)
		if strings.HasPrefix(ck, u.ParamPrefix) {
			u.Info[ck[len(u.ParamPrefix):]] = r.FormValue(k)
		}
	}
	return u
}

// Attr 返回用户信息属性
func (u Baggage) Attr(name string) Value {
	return Value(u.Info[CanonicalKey(name)])
}

// Value 代表一个可转换值
type Value string

// Int 返回Value代表的int值
func (v Value) Int() (int, error) {
	return strconv.Atoi(string(v))
}

// MustInt 返回Value代表的int值, 如果转换失败则panic
func (v Value) MustInt() int {
	i, err := strconv.Atoi(string(v))
	if err != nil {
		panic(err)
	}
	return i
}

// Int64 返回Value代表的int64值
func (v Value) Int64() (int64, error) {
	return strconv.ParseInt(string(v), 10, 64)
}

// MustInt64 返回Value代表的int64值
func (v Value) MustInt64() int64 {
	i64, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		panic(err)
	}
	return i64
}

// Float64 返回Value代表的float64值
func (v Value) Float64() (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}

// MustFloat64 返回Value代表的float64值
func (v Value) MustFloat64() float64 {
	f64, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		panic(err)
	}
	return f64
}

// String 返回Value代表的String值
func (v Value) String() string {
	return string(v)
}

// Bytes 返回Value代表的[]byte值
func (v Value) Bytes() []byte {
	return []byte(v)
}

// Bool 返回Value代表的bool值, 属性为["", "false", "False"]之一时返回false，其他返回true
func (v Value) Bool() bool {
	s := string(v)
	if s == "False" || s == "false" || s == "" {
		return false
	}
	return true
}
