package baggage_test

import (
	"net/http"
	"net/textproto"
	"strings"
	"testing"

	"github.com/ccmonky/pkg/baggage"
)

func TestStringsJoin(t *testing.T) {
	ss := []string{}
	if strings.Join(ss, "_") != "" {
		t.Error("should ==")
	}
	ss = []string{"1"}
	if strings.Join(ss, "_") != "1" {
		t.Error("should ==")
	}
}

func TestUserInfoo(t *testing.T) {
	u := baggage.New("x-tprOxy", "user")
	if u.HeaderPrefix != "X-Tproxy-User-" {
		t.Errorf("should ==, got %s", u.HeaderPrefix)
	}
	if u.ParamPrefix != "x-tproxy-user-" {
		t.Errorf("should ==, got %s", u.HeaderPrefix)
	}
	if u.Info["no"] != "" {
		t.Error("should ==")
	}
	if u.Info != nil {
		t.Error("should ==")
	}
	u.WithAttr("auth_backend", "ssoauth")
	if u.Attr("auth-backend") != "ssoauth" {
		t.Error("should ==")
	}
	if u.Info == nil {
		t.Error("should not nil")
	}
	if u.HeaderPrefix != "X-Tproxy-User-" {
		t.Error("should ==")
	}
	if u.ParamPrefix != "x-tproxy-user-" {
		t.Error("should ==")
	}
	u.WithAttr("tid", "1")
	if u.Attr("tid").MustInt64() != 1 {
		t.Error("should ==")
	}
	u.WithInfo(nil)
	if u.Attr("tid") != "" {
		t.Errorf("should ==, got %s", u.Attr("tid"))
	}
	u.WithAttr("auth_backend", "jwtauth")
	u.WithAttr("tid", "2")
	u.WithAttr("uid", "2")
	u.WithAttr("user_name", "alice")
	headers := u.Headers()
	if headers.Get("X-Tproxy-User-User-Name") != "alice" {
		t.Error("should ==")
	}
	params := u.Params()
	if params.Get("x-tproxy-user-uid") != "2" {
		t.Error("should ==")
	}
	r, err := http.NewRequest("GET", "http://authcar.amap.com/ws/authcar/jwks", nil)
	if err != nil {
		t.Error("should ==")
	}
	u2 := baggage.New("x-tproxy", "user")
	if u2.Extract(r).Attr("tid") != "" {
		t.Errorf("should==, got %s", u.Attr("tid"))
	}
	u.InjectHeaders(r)
	if r.Header.Get("X-Tproxy-User-User-Name") != "alice" {
		t.Error("should ==")
	}
	if r.Header.Get("X-Tproxy-User-Uid") != "2" {
		t.Error("should ==")
	}
	if u2.Extract(r).Attr("uid").MustInt64() != 2 {
		t.Error("should==")
	}
	r.Header = nil
	if r.Header.Get("X-Tproxy-User-Uid") != "" {
		t.Error("should ==")
	}
	u.InjectParams(r)
	if r.FormValue("x-tproxy-user-tid") != "2" {
		t.Errorf("should ==, got %s", r.FormValue("x-tproxy-user-tid"))
	}
	u3 := baggage.New("x-tproxy", "user")
	if u3.Extract(r).Attr("tid").MustInt64() != 2 {
		t.Error("should==")
	}
}

func TestHTTPRequestHeader(t *testing.T) {
	hdr := http.Header{}
	hdr.Set("X-Tproxy_user_mozi_tid", "123")
	if hdr.Get("X-Tproxy_user_mozi_tid") != "123" {
		t.Fatal("should ==")
	}
	if hdr.Get("x-tproxy_user_mozi_tid") != "123" { // 大小写无关的！
		t.Fatal("should ==")
	}
}

func TestCanonicalMIMEHeaderKey(t *testing.T) {
	h := "x-tproxy-user-mozi_uid"
	ch := textproto.CanonicalMIMEHeaderKey(h)
	if ch != "X-Tproxy-User-Mozi_uid" {
		t.Fatal(ch)
	}
}
