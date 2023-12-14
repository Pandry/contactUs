package captchas

import (
	"log"

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
	resp, err := ts.Verify(response, ip)
	if err != nil {
		log.Println("Error validating captcha", err)
		return false
	}
	if resp.Success {
		return true
	}
	return false
}
