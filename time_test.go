package pkg_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ccmonky/pkg"
)

func TestTime(t *testing.T) {
	timeSkew := 5 * time.Minute
	now, err := time.Parse(time.RFC3339, "2018-12-08T13:00:00+08:00")
	if err != nil {
		t.Fatal(err)
	}
	activationTime, err := time.Parse(time.RFC3339, "2018-12-08T12:59:00+08:00")
	if err != nil {
		t.Fatal(err)
	}
	//valid := now.Sub(activationTime) >= -timeSkew
	valid := pkg.AfterEqualWithSkew(now, activationTime, timeSkew)
	if !valid {
		t.Fatal("now should >= activation time with 5 minutes skew")
	}
	activationTime, err = time.Parse(time.RFC3339, "2018-12-08T13:02:00+08:00")
	if err != nil {
		t.Fatal(err)
	}
	//valid = now.Sub(activationTime) >= -timeSkew
	valid = pkg.AfterEqualWithSkew(now, activationTime, timeSkew)
	if !valid {
		t.Fatal("now should >= activation time with 5 minutes skew")
	}
	activationTime, err = time.Parse(time.RFC3339, "2018-12-08T13:06:00+08:00")
	if err != nil {
		t.Fatal(err)
	}
	//valid = now.Sub(activationTime) >= -timeSkew
	valid = pkg.AfterEqualWithSkew(now, activationTime, timeSkew)
	if valid {
		t.Fatal("now should >= activation time with 5 minutes skew")
	}
	endTime, err := time.Parse(time.RFC3339, "2018-12-08T12:56:00+08:00")
	// valid = now.Sub(endTime) <= timeSkew
	valid = pkg.BeforeEqualWithSkew(now, endTime, timeSkew)
	if !valid {
		t.Fatal("now should <= end time with 5 minuts skew")
	}
	endTime, err = time.Parse(time.RFC3339, "2018-12-08T12:54:00+08:00")
	//valid = now.Sub(endTime) <= timeSkew
	valid = pkg.BeforeEqualWithSkew(now, endTime, timeSkew)
	if valid {
		t.Fatal("now should <= end time with 5 minuts skew")
	}
	endTime, err = time.Parse(time.RFC3339, "2018-12-08T13:01:00+08:00")
	//valid = now.Sub(endTime) <= timeSkew
	valid = pkg.BeforeEqualWithSkew(now, endTime, timeSkew)
	if !valid {
		t.Fatal("now should <= end time with 5 minuts skew")
	}
}

type demo struct {
	Elapsed pkg.Duration `json:"elapsed"`
}

func TestDuration(t *testing.T) {
	msgEnc, err := json.Marshal(&demo{
		Elapsed: pkg.Duration{time.Second * 5},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(msgEnc))
	var msg demo
	if err := json.Unmarshal([]byte(`{"elapsed": "3s"}`), &msg); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", msg)
	if err = json.Unmarshal([]byte(`{"elapsed": "0"}`), &msg); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", msg)
	if err = json.Unmarshal([]byte(`{"elapsed": "-1ns"}`), &msg); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", msg)
}
