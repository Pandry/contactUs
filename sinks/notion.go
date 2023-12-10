package sinks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (config NotionConfig) Name() string {
	return "Notion"
}

type NotionConfig struct {
	Secret     string                 `yaml:"secret"`
	DatabaseId string                 `yaml:"database"`
	Entries    map[string]NotionEntry `yaml:"parameters"`
}

type NotionEntry struct {
	NotionColumnName string `yaml:"notionColumnName"`
	NotionColumnType string `yaml:"notionColumnType"`
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

func (config NotionConfig) Sink(input map[string]string) error {
	if !config.IsReady() {
		return ErrSinkNotReady
	}

	// Notions is very funny and love games and their API is so clean and I love them
	//
	// If you find yourself lost in this file, first and foremost I'm sorry
	// Also, feel free to increment the counter of hours *wasted* here: 15h
	notionRowProperties := make(map[string]interface{})

	for k, v := range input {
		if nK, ok := config.Entries[k]; ok {
			// log.Println("Parsing", k, "With value", v, "as", nK.NotionColumnName, "of type", nK.NotionColumnType)
			propertyMap := make(map[string]interface{})
			switch nK.NotionColumnType {
			case "text", "rich_text":
				propertyMap["type"] = "rich_text"
				propertyMap["rich_text"] = []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": v,
						},
					},
				}
			case "title":
				propertyMap["type"] = nK.NotionColumnType
				propertyMap[nK.NotionColumnType] = []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": v,
						},
					},
				}
			case "number":
				parsedNumber, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					log.Println("Notion sink: Parsing failed for float type.", err)
					continue
				}
				propertyMap["type"] = nK.NotionColumnType
				propertyMap[nK.NotionColumnType] = parsedNumber

			default:
				log.Println("Notion sink: Given type (", nK.NotionColumnType, ") was unhandled. Ignoring.")
				continue

			}
			notionRowProperties[nK.NotionColumnName] = propertyMap
		} else {
			log.Println("Notion sink: ignoring field", k, "as it was not found in the notion configuration")
		}
	}

	payload := notionPayload{
		Parent: parent{
			DatabaseID: config.DatabaseId,
		},
		Properties: notionRowProperties,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Join(errors.New("Error mashaling form into JSON for notion"), err)
	}

	// log.Println(string(payloadBytes))

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

	if resp.StatusCode > 299 {
		respBodyStr, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Join(errors.New("Error reading notion response"), err)
		}
		return errors.New(fmt.Sprintf("Got unexpected code from notion (%d). Sent payload: %s. Response from the server: %s", resp.StatusCode, string(payloadBytes), respBodyStr))
	}

	return nil
}
