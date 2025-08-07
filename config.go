package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TwilioAccountSID string `json:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken  string `json:"TWILIO_AUTH_TOKEN"`
	SourceNumber     string `json:"TWILIO_SOURCE_NUMBER"`
	TargetNumber     string `json:"TWILIO_TARGET_NUMBER"`
	TgBotApiToken    string `json:"TG_BOT_API_TOKEN"`
}

func LoadConfig() (*Config, error) {

	slog.Info("Loading configuration from environment variables")
	godotenv.Load()

	source := os.Getenv("TWILIO_SOURCE_NUMBER")
	if source == "" {
		return nil, fmt.Errorf("missing the TWILIO_SOURCE_NUMBER environment variable.")
	}

	target := os.Getenv("TWILIO_TARGET_NUMBER")
	if target == "" {
		return nil, fmt.Errorf("missing the TWILIO_TARGET_NUMBER environment variable.")
	}

	tg_token := os.Getenv("TG_BOT_API_TOKEN")
	if tg_token == "" {
		return nil, fmt.Errorf("missing the TG_BOT_API_TOKEN environment variable.")
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	if accountSid == "" {
		return nil, fmt.Errorf("missing the TWILIO_ACCOUNT_SID environment variable.")
	}

	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("missing the TWILIO_AUTH_TOKEN environment variable.")
	}

	config := &Config{
		TwilioAccountSID: accountSid,
		TwilioAuthToken:  authToken,
		SourceNumber:     source,
		TargetNumber:     target,
		TgBotApiToken:    tg_token,
	}
	return config, nil
}
