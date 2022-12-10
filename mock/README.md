# mock

mock工具集合。

## Matcher

Matcher根据请求匹配要使用的ResponseMocker，定义如下：

```go
// Matcher 根据请求匹配要使用的ResponseMocker
type Matcher interface {
    Eigenkey(r *http.Request) (string, error)
    Match(r *http.Request) (ResponseMocker, error)
}
```

## ResponseMocker

ResponseMocker是一个Mock Response生成器接口，其定义如下：

```go
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
}
```

本工具库内置如下ResponseMocker实现：

- TransparentResponseMocker: 透明Mocker，即从源服务获取真实响应作为Mock，默认如果不指定则使用此Mocker
- ResponseMockerFromURL: 请求URL的获取整个响应作为Mock，通常用于Mock平台

```json
{
    "response_mocker": "ResponseMockerFromURL",
    "response_from_url": "http://mock.alibaba-inc.com/ws/test/xxx?abc=123"
}
```

- ResponseMockerBuilder: Response Mock构造器，通过指定状态码、头和Body生成Mock响应，通常用于静态mock

```json
{
    "response_mocker": "ResponseMockerBuilder",
    "status_code": 200,
    "header": {
        "X-Test": ["abc"],
        "X-Test2": ["def", "123"]
    },
    "body": "this is a test"
}
```

```json
{
    "response_mocker": "ResponseMockerBuilder",
    "status_code": 200,
    "header": {
        "X-Test": ["abc"],
        "X-Test2": ["def", "123"]
    },
    "body_from_url": "http://oss.alibaba-inc.com/ws/test/xxx/body?abc=123"
}
```
