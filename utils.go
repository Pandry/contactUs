package main

import (
	"strconv"
	"strings"
)

func stringMapToInterface(input map[string]string) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if strings.ToLower(strings.TrimSpace(v)) == "true" {
			output[k] = true
			continue
		}
		if strings.ToLower(strings.TrimSpace(v)) == "false" {
			output[k] = false
			continue
		}
		if val, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			output[k] = val
			continue
		}
		if val, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			output[k] = val
			continue
		}
		output[k] = v
	}
	return output
}
