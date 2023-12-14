package captchas

import (
	turnstile "github.com/meyskens/go-turnstile"
)

type TurnstileConfig struct {
	Key string `yaml:"key"`
}

func (config TurnstileConfig) GetResponseFieldName() string {
	return "cf-turnstile-response"
}

func (config TurnstileConfig) VerifyCaptcha(response string, ip string) bool {
	ts := turnstile.New(config.Key)

	if resp, err := ts.Verify(response, ip); err == nil && resp.Success {
		return true
	}
	return false
}
