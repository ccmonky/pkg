# eigenkey

计算特征键。

## HTTPRequestEigenkeyExtractor

根据配置从http请求提取特征键。

```go
url := "https://yfliu:123@xxxstore.xxx.com/ws/xxxsdk/login?a=1&z=2&b=3#fragment"
r, _ := http.NewRequest("GET", url, nil)
g := &eigenkey.HTTPRequestEigenkeyExtractor{}
_ = g.Provision()
fmt.Println(g.Eigenkey(r)) // "/ws/xxxsdk/login"
g.RequestExtractor.UseArguments = []string{"a", "b"} // NOTE: not thread safe
fmt.Println(g.Eigenkey(r)) // "/ws/xxxsdk/login?a=1&b=3"
```
