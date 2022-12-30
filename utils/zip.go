package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// UnzipFirst used to unzip first file in a request body with zip archived
func UnzipFirst(payload []byte) (content []byte, err error) {
	var z *zip.Reader
	z, err = zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if err != nil {
		return nil, errors.Wrap(err, "new zip reader error")
	}
	if len(z.File) == 0 { // need?
		return nil, errors.Errorf("no valid zip files found")
	}
	r, err := z.File[0].Open()
	if err != nil {
		return nil, errors.Wrap(err, "open zip file error")
	}
	content, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "zip readall error")
	}
	r.Close()

	return content, nil
}

// Unzip used to unzip all files in a request body with zip archived
func Unzip(payload []byte) (map[string][]byte, error) {
	contents := make(map[string][]byte, 1)
	var err error
	var z *zip.Reader
	z, err = zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if err != nil {
		return nil, errors.Wrap(err, "new zip reader error")
	}
	for _, f := range z.File {
		r, err := f.Open()
		if err != nil {
			return nil, errors.Wrap(err, "open zip file error")
		}
		content, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, errors.Wrap(err, "zip readall error")
		}
		contents[f.Name] = content
		r.Close()
	}

	return contents, nil
}

// Zip compress the payload and return zip bytes, file name pattern is `file\d+`
func Zip(files ...[]byte) ([]byte, error) {
	if len(files) == 0 {
		return nil, errors.New("empty input to zip")
	}
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for i, file := range files {
		f, err := w.Create(fmt.Sprintf("file%d", i+1))
		if err != nil {
			return nil, errors.Wrap(err, "create zip writer error")
		}
		n, err := f.Write(file)
		if err != nil {
			return nil, errors.Wrap(err, "write error on zip writer")
		}
		if n != len(file) {
			return nil, errors.Wrapf(err, "incomplete write error for file %d", i)
		}
	}

	err := w.Close()
	if err != nil {
		return nil, errors.Wrap(err, "close zip writer error")
	}

	return buf.Bytes(), nil
}
