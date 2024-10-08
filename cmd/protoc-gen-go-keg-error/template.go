package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

func Is{{.CamelValue}}(err error) bool {
	e := response.FromError(err)
	return e.GetReason() == {{.Name}}_{{.Value}}.String() && e.GetCode() == {{.Code}} 
}

func {{.CamelValue}}(args ...interface{}) *response.Response {
	 return response.Newf({{.HttpCode}}, {{.Code}}, {{.Name}}_{{.Value}}.String(), "{{.Msg}}", args...)
}

{{- end }}
`

type errorInfo struct {
	Name       string
	Value      string
	HttpCode   int
	CamelValue string
	Code       int
	Msg        string
}
type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
