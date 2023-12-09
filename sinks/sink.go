package sinks

import (
	"bytes"
	"text/template"
)

/*
To add:
Slack
Zulip
Notion
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
