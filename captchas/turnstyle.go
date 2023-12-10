package captchas

import (
	turnstile "github.com/meyskens/go-turnstile"
)

type TurnstyleConfig struct {
	Key string `yaml:"key"`
}

func (config TurnstyleConfig) GetResponseFieldName() string {
	return "cf-turnstile-response"
}

func (config TurnstyleConfig) VerifyCaptcha(response string, ip string) bool {
	ts := turnstile.New(config.Key)
	//Get IP from RemoteAddr
	// ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if resp, err := ts.Verify(response, ip); err == nil && resp.Success {
		return true
	}
	return false
}
