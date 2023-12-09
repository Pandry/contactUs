package configuration_test

import (
	"contactUs/configuration"
	"contactUs/sinks"
	"reflect"
	"testing"
)

type confParseTest struct {
	config   string
	expected configuration.Configuration
}

// TODO: actually check output
func TestConfParsing(t *testing.T) {
	defaultConfig := configuration.GetDefaultConfig()
	testConfigs := []confParseTest{
		{
			`sinks:
  test:
    webhook:
      url: https://example.com
      method: POST
      headers:
        foo: bar
`, defaultConfig{
				Sinks: map[string]configuration.Sink{
					"test": {
						Webhook: sinks.WebhookConfig{
							URL:     "https://example.com",
							Method:  "POST",
							Headers: map[string]string{"foo": "bar"}}}}}},
		// 		{
		// 			`
		// forms:
		// www-contact:
		// redirect: https://example.com
		// inputs:
		//   - name
		//   - surname
		//   - weight
		// captcha:
		//   enabled: false
		// `, configuration.Configuration{Sinks: make(map[string]configuration.Sink)}},
		// 		{
		// 			`sinks:
		// test:
		// webhook:
		//   url: https://example.com
		//   method: POST
		//   headers:
		//     foo: bar
		// forms:
		// www-contact:
		// redirect: https://example.com
		// inputs:
		//   - name
		//   - surname
		//   - weight
		// captcha:
		//   enabled: false
		// sinks:
		//   - test
		// `, configuration.Configuration{Sinks: make(map[string]configuration.Sink)}},
		// 		{"forms: {}", configuration.Configuration{}},
	}

	for i, conf := range testConfigs {
		t.Log("Running test", i+1)
		confStr := conf.config
		parsedConf, err := configuration.ParseBytes([]byte(confStr))
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if !reflect.DeepEqual(parsedConf, conf.expected) {
			t.Error("Config struct not matching", parsedConf, conf.expected)
			t.Fail()
		}
	}
}
