package sinks

import (
	"bytes"
	"text/template"
)

type ISink interface {
	Sink(input map[string]string) error
	IsReady() bool
	Name() string
}

/*
To add:
Slack
Zulip
Email
sqlite
PostgreSQL
*/

func evaluateTemplate(t string, v any) string {
	var buf bytes.Buffer
	if err := template.Must(template.New("template").Parse(t)).Execute(&buf, v); err == nil {
		return buf.String()
	}
	return ""
}
