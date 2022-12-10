package eigenkey_test

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"gitlab.alibaba-inc.com/t3/pkg/eigenkey"
)

func TestURLValues(t *testing.T) {
	u := url.Values{}
	u.Add("a", "1")
	u.Add("a", "2")
	u.Add("b", "4")
	c := u["c"]
	if len(c) != 0 {
		t.Fatal("should ==")
	}
}

func TestHash(t *testing.T) {
	key := "http://yunfei.liu:123456@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2#dummy:HTTP/1.1:d=hello&d=world&c=hi&c="
	var h string
	h = eigenkey.MD5(key)
	if h != "ab582e91ea63f409a47825f388d0422f" {
		t.Fatalf("should ==, got %s", h)
	}
	h = eigenkey.SHA1(key)
	if h != "b48f81edcebb6c467b1bc4665b2674ec3bb535a8" {
		t.Fatalf("should ==, got %s", h)
	}
	h = eigenkey.SHA256(key)
	if h != "98d8616a4d4a22fdeaa8853de6336269d013b45fb9498ce21dd17c0d61aa3b9d" {
		t.Fatalf("should ==, got %s", h)
	}
}

func TestDefaultHTTPEigenkeyFunc(t *testing.T) {
	info := &eigenkey.HTTPRequestInfo{}
	var key string
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Path = "/ws/authcar/jwks"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "/ws/authcar/jwks" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.UseArguments = []string{"b", "a"}
	info.Arguments = url.Values{
		"a": []string{"1", "2"},
		"b": []string{"5"},
	}
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "/ws/authcar/jwks?b=5&a=1&a=2" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Scheme = "http"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "http:///ws/authcar/jwks?b=5&a=1&a=2" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Host = "authcar.amap.test"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "http://authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Username = "yunfei.liu"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "http://yunfei.liu@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Password = "123456"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "http://yunfei.liu:123456@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Fragment = "dummy"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "http://yunfei.liu:123456@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2#dummy" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Proto = "HTTP/1.1"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "http://yunfei.liu:123456@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2#dummy:HTTP/1.1" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.UseHeaders = []string{"d", "c"}
	info.Headers = url.Values{
		"c": []string{"hi", ""},
		"d": []string{"hello", "world"},
	}
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "http://yunfei.liu:123456@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2#dummy:HTTP/1.1:d=hello&d=world&c=hi&c=" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	info.Method = "POST"
	key = eigenkey.DefaultHTTPEigenkeyFunc("", info)
	if key != "POST:http://yunfei.liu:123456@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2#dummy:HTTP/1.1:d=hello&d=world&c=hi&c=" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	key = eigenkey.DefaultHTTPEigenkeyFunc("tproxy", info)
	if key != "tproxy:POST:http://yunfei.liu:123456@authcar.amap.test/ws/authcar/jwks?b=5&a=1&a=2#dummy:HTTP/1.1:d=hello&d=world&c=hi&c=" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	key = eigenkey.DefaultHTTPEigenkeyFunc("tproxy", info, []eigenkey.KeyPostFunc{eigenkey.MD5}...)
	if key != "d935b6bbc01a5153fc87f4d93cce38f9" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	key = eigenkey.DefaultHTTPEigenkeyFunc("tproxy", info, []eigenkey.KeyPostFunc{eigenkey.SHA1}...)
	if key != "9e7185abd3e0c01bb1a254a2052e50bb6242d083" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	key = eigenkey.DefaultHTTPEigenkeyFunc("tproxy", info, []eigenkey.KeyPostFunc{eigenkey.SHA256}...)
	if key != "b5597568965500652ed7c82dc798c11107de031355e78bfe197eec34e7733f13" {
		t.Fatalf("shoudl ==, got %s", key)
	}
	key = eigenkey.DefaultHTTPEigenkeyFunc("tproxy", info, []eigenkey.KeyPostFunc{eigenkey.Prefix64}...)
	if key != "tproxy:POST:http://yunfei.liu:123456@authcar.amap.test/ws/authca" {
		t.Fatalf("shoudl ==, got %s", key)
	}
}

func TestRequestKeyGenerator(t *testing.T) {
	url := "https://yfliu:123@apistore.amap.com/ws/autosdk/login?a=1&z=2&b=3#fragment"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	g := &eigenkey.HTTPRequestEigenkeyExtractor{}
	err = g.Provision()
	if err != nil {
		t.Fatal(err)
	}

	k, err := g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "/ws/autosdk/login" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseScheme = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "https:///ws/autosdk/login" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseScheme = false
	g.RequestExtractor.UseFragment = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "/ws/autosdk/login#fragment" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseFragment = false
	g.RequestExtractor.UseMethod = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "GET:/ws/autosdk/login" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseMethod = false
	g.RequestExtractor.UseHost = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "//apistore.amap.com/ws/autosdk/login" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseMethod = false
	g.RequestExtractor.UseScheme = true
	g.RequestExtractor.UseHost = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "https://apistore.amap.com/ws/autosdk/login" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseScheme = false
	g.RequestExtractor.UseHost = false
	g.RequestExtractor.UseArguments = []string{"a", "b"}
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "/ws/autosdk/login?a=1&b=3" {
		t.Errorf("should ==, got %s", k)
	}

	r.Header.Set("xxx-abc", "bearer 123")
	r.Header.Set("abc-324", "xff-rt")
	g.RequestExtractor.UseArguments = []string{"z", "a"}
	g.RequestExtractor.UseHeaders = []string{"abc-324", "xxx-abc"}
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "/ws/autosdk/login?z=2&a=1:abc-324=xff-rt&xxx-abc=bearer+123" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseUsername = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "//yfliu@/ws/autosdk/login?z=2&a=1:abc-324=xff-rt&xxx-abc=bearer+123" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UsePassword = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "//yfliu:123@/ws/autosdk/login?z=2&a=1:abc-324=xff-rt&xxx-abc=bearer+123" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseScheme = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "https://yfliu:123@/ws/autosdk/login?z=2&a=1:abc-324=xff-rt&xxx-abc=bearer+123" {
		t.Errorf("should ==, got %s", k)
	}

	g.RequestExtractor.UseFragment = true
	k, err = g.Eigenkey(r)
	if err != nil {
		t.Error(err)
	}
	if k != "https://yfliu:123@/ws/autosdk/login?z=2&a=1#fragment:abc-324=xff-rt&xxx-abc=bearer+123" {
		t.Errorf("should ==, got %s", k)
	}
}

func TestHTTPRequestEigenkeyExtractor(t *testing.T) {
	extractor := eigenkey.HTTPRequestEigenkeyExtractor{
		RequestExtractor: &eigenkey.HTTPRequestExtractor{
			UseMethod:    true,
			UsePath:      true,
			UseArguments: []string{"posta", "postb"},
		},
		CleanPath: true,
	}
	err := extractor.Provision()
	if err != nil {
		t.Fatal(err)
	}

	rq, err := http.NewRequest("POST", "http://localhost/?a=1&b=2", bytes.NewReader([]byte(``)))
	rq.Form = url.Values{"posta": []string{"1"}, "postb":[]string{"2"}}
	if err != nil {
		t.Fatal(err)
	}
	ek, err := extractor.Eigenkey(rq)
	if err != nil {
		t.Fatal(err)
	}

	if ek != "POST:/?posta=1&postb=2" {
		t.Fatal(ek)
	}

	extractor = eigenkey.HTTPRequestEigenkeyExtractor{
		RequestExtractor: &eigenkey.HTTPRequestExtractor{
			UseMethod:    true,
			UsePath:      true,
			UseArguments: []string{"a", "b", "posta"},
		},
		CleanPath: true,
	}
	err = extractor.Provision()
	if err != nil {
		t.Fatal(err)
	}
	rq, err = http.NewRequest("POST", "http://localhost/?a=1&b=2", bytes.NewReader([]byte(`posta=1&postb=2`)))
	rq.Header.Set("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}
	ek, err = extractor.Eigenkey(rq)
	if err != nil {
		t.Fatal(err)
	}

	if ek != "POST:/?a=1&b=2&posta=1" {
		t.Fatal(ek)
	}
}
