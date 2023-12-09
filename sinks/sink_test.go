package sinks

import (
	"log"
	"testing"
)

type templateTestScaffold struct {
	name     string
	template string
	input    map[string]string
	output   string
}

func TestTemplate(t *testing.T) {
	// Triples of values (template, input, output)
	values := []templateTestScaffold{
		{"Dummy", `{{ .foo }}`, map[string]string{"foo": "bar"}, "bar"},
		{"Empty value", `{{ .foo }}`, map[string]string{"foo": ""}, ""},
		{"Empty key, default", `{{ .foo }}{{ or .bar "" }}`, map[string]string{"foo": "baz"}, "baz"},
	}
	for i, v := range values {
		t.Log("Running test ", v.name, "(", i+1, ")")
		if res := evaluateTemplate(v.template, v.input); v.output != res {
			log.Fatalln("Test", v.name, "(", i+1, ")", "failed: Got: \""+res+"\", expected: \""+v.output+"\"")
			t.Fail()
		}
	}
}
