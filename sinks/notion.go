package sinks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func (config NotionConfig) Name() string {
	return "Notion"
}

type NotionConfig struct {
	Secret     string `yaml:"secret"`
	DatabaseId string `yaml:"database"`
}
type notionPayload struct {
	Parent     parent                 `json:"parent"`
	Properties map[string]interface{} `json:"properties"`
}
type parent struct {
	DatabaseID string `json:"database_id"`
}

func (config NotionConfig) IsReady() bool {
	if config.DatabaseId == "" || config.Secret == "" {
		return false
	}
	if config.Secret[0:7] != "secret_" {
		return false
	}
	return true
}

func (config NotionConfig) Sink(input map[string]interface{}) error {
	if !config.IsReady() {
		return ErrSinkNotReady
	}

	// Notions is very funny and love games and their API is so clean and I love them
	//
	// If you find yourself lost in this file, first and foremost I'm sorry
	// Also, feel free to increment the counter of hours *wasted* here: 13h
	for k, v := range input {
		switch v.(type) {
		case string:
			textObj := map[string]interface{}{
				"type": "rich_text",
				"rich_text": []map[string]interface{}{
					{
						// "type": "text",
						"text": map[string]interface{}{
							"content": v,
						},
					},
				},
			}
			input[k] = textObj
			break
		// We only parse into int64/float 64
		case int64, float64:
			numberObj := map[string]interface{}{
				"type":   "number",
				"number": v,
			}
			input[k] = numberObj
		default:
			fmt.Println("Got unexpected type", reflect.TypeOf(v))
		}

	}

	payload := notionPayload{
		Parent: parent{
			DatabaseID: config.DatabaseId,
		},
		Properties: input,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Join(errors.New("Error mashaling form into JSON for notion"), err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.notion.com/v1/pages", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errors.Join(errors.New("Error creating request"), err)
	}
	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Secret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-02-22") // Use the latest version supported

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Join(errors.New("Error sending notion request"), err)
	}
	defer resp.Body.Close()
	fmt.Println(string(payloadBytes))

	if resp.StatusCode > 299 {
		respBodyStr, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Join(errors.New("Error reading notion response"), err)
		}
		return errors.New(fmt.Sprintf("Got unexpected code from notion (%d). Response from the server: %s", resp.StatusCode, respBodyStr))
	}

	return nil
}
