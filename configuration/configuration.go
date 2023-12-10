package configuration

import (
	"contactUs/captchas"
	"contactUs/sinks"
)

var Config *Configuration

type Configuration struct {
	Sinks          map[string]*Sink `yaml:"sinks"`
	Forms          map[string]Form  `yaml:"forms"`
	TrustProxyes   bool             `yaml:"trustProxies"`
	ClientIPHeader string           `yaml:"clientIPHeader"`
	ListenPort     int              `yaml:"port"`
}

type Sink struct {
	ActiveSinksList []sinks.ISink
	Webhook         sinks.WebhookConfig `yaml:"webhook"`
	Notion          sinks.NotionConfig  `yaml:"notion"`
}

type Form struct {
	Inputs   []string             `yaml:"inputs"`
	Redirect string               `yaml:"redirect"`
	Captcha  CaptchaConfiguration `yaml:"captcha"`
	Sinks    []string             `yaml:"sinks"`
}

type CaptchaConfiguration struct {
	Enabled  bool   `yaml:"enabled"`
	Provider string `yaml:"provider"` // turnstyle, reCaptcha, hCaptcha
	Secret   string `yaml:"secret"`
	Captcha  captchas.ICaptcha
}
