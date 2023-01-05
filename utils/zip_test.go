package utils_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ccmonky/pkg/utils"
)

func TestUnzipFirst(t *testing.T) {
	file, err := os.Open("testdata/log.dat.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	content, err := utils.UnzipFirst(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
}

func TestUnzip(t *testing.T) {
	file, err := os.Open("testdata/log.dat.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	contents, err := utils.Unzip(data)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := contents["28c7a71c-8f28-4bf7-98ac-1414a917c77f"]; !ok {
		t.Fatal("should contains")
	}
	if _, ok := contents["not-exist"]; ok {
		t.Fatal("should not contains")
	}
}

func TestZip(t *testing.T) {
	file1 := []byte("this is a test")
	file2 := []byte("this is test 2")
	zipData, err := utils.Zip(file1, file2)
	if err != nil {
		t.Fatal(err)
	}
	data, err := utils.Unzip(zipData)
	if err != nil {
		t.Fatal(err)
	}
	var content []byte
	var ok bool
	if content, ok = data["file1"]; !ok {
		t.Fatal("should contains")
	}
	if !bytes.Equal(content, file1) {
		t.Fatalf("should ==, got %s", string(content))
	}
	if content, ok = data["file2"]; !ok {
		t.Fatal("should contains")
	}
	if !bytes.Equal(content, file2) {
		t.Fatal("should ==")
	}
	if _, ok = data["not-exist"]; ok {
		t.Fatal("should not contains")
	}
}
