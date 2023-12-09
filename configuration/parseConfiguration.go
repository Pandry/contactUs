package configuration

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// TODO: Read file
func ParseFile(path string) (Configuration, error) {
	filesContent, err := ioutil.ReadFile(path)
	if err != nil {
		return Configuration{}, err
	}
	return ParseBytes(filesContent)
}

func ParseBytes(config []byte) (Configuration, error) {
	c := GetDefaultConfig()
	err := yaml.Unmarshal(config, &c)
	if err != nil {
		panic(err)
	}

	// TODO: Check is sink exists

	Config = c
	// for name, sink := range c.Sinks {
	// }
	return Config, nil
}

func GetDefaultConfig() Configuration {
	return Configuration{
		Sinks:          make(map[string]Sink),
		Forms:          make(map[string]Form),
		TrustProxyes:   false,
		ClientIPHeader: "X-Forwarded-For",
	}
}
