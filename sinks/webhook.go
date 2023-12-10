package sinks

import (
	"net/http"
	"strings"
)

type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

func (config WebhookConfig) Name() string {
	return "Webhook"
}

func (config WebhookConfig) IsReady() bool {
	return config.checkConfig()
}
func (config WebhookConfig) checkConfig() bool {
	if config.URL == "" {
		return false
	}
	// TODO: Check URL
	method := config.Method
	if method == "" {
		config.Method = "GET"
	}
	if method != "POST" && method != "GET" && method != "PUT" && method != "DELETE" && method != "PATCH" {
		return false
	}
	return true
}
func (config WebhookConfig) Sink(input map[string]interface{}) error {
	if !config.IsReady() {
		return ErrSinkNotReady
	}
	url := evaluateTemplate(config.URL, input)
	method := config.Method
	if method == "" {
		method = "GET"
	}
	body := strings.NewReader(evaluateTemplate(config.Body, input))

	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	for key, value := range config.Headers {
		req.Header.Add(key, value)
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
