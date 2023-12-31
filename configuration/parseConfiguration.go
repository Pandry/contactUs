package configuration

import (
	"contactUs/captchas"
	"contactUs/sinks"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// TODO: Read file
func ParseFile(path string) (*Configuration, error) {
	filesContent, err := ioutil.ReadFile(path)
	if err != nil {
		return &Configuration{}, err
	}
	return ParseBytes(filesContent)
}

func ParseBytes(config []byte) (*Configuration, error) {
	c := GetDefaultConfig()
	err := yaml.Unmarshal(config, &c)
	if err != nil {
		panic(err)
	}

	// Load sinks configurations
	for sinkName, sink := range c.Sinks {
		fmt.Println("Loading sink", sinkName)
		sink.ActiveSinksList = make([]sinks.ISink, 0)
		if sink.Webhook.IsReady() {
			sink.ActiveSinksList = append(sink.ActiveSinksList, sink.Webhook)
			fmt.Println("Loaded subsink", sink.Webhook.Name(), "for sink", sinkName)
		}
		if sink.Notion.IsReady() {
			sink.ActiveSinksList = append(sink.ActiveSinksList, sink.Notion)
			fmt.Println("Loaded subsink", sink.Notion.Name(), "for sink", sinkName)
		}
		fmt.Println("Sink", sinkName, "loaded")
	}

	// Loads captcha configurations
	for formName, form := range c.Forms {
		if form.Captcha.Enabled {
			fmt.Println("Loading captcha ", form.Captcha.Provider, " for form ", formName)
			switch form.Captcha.Provider {
			case "turnstile":
				form.Captcha.Captcha = captchas.TurnstileConfig{
					Key: form.Captcha.Secret,
				}
			default:
				panic("Could not find captcha provider")
			}
		}
	}
	Config = &c
	return Config, nil
}

func GetDefaultConfig() Configuration {
	return Configuration{
		Sinks:        make(map[string]*Sink),
		Forms:        make(map[string]*Form),
		TrustProxyes: false,
		// ClientIPHeader: "X-Forwarded-For",
	}
}
