package utils

import (
	"bytes"
	"text/template"
	"time"
)

// ExecuteTemplate used to test template and funcs
func ExecuteTemplate(text string, data interface{}, funcMap template.FuncMap) (result string, err error) {
	var tmpl *template.Template
	tmpl, err = template.New("test_" + time.Now().Format("20060102030405006")).Funcs(funcMap).Parse(text)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return
	}
	return buf.String(), nil
}
