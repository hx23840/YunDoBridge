package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	OpenAIAPIKey        string `json:"openai_api_key"`
	OpenAIModelEndpoint string `json:"openai_model_endpoint"`
	CallsBaseURL        string `json:"calls_base_url"`
	CallsAppID          string `json:"calls_app_id"`
	CallsAppToken       string `json:"calls_app_token"`
	Port                string `json:"port"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	return &config, nil
}
