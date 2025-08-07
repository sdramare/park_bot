package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"park_bot/voip"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("Starting app...")
	config, err := LoadConfig()

	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(config.TgBotApiToken)
	if err != nil {
		slog.Error("Could not create Telegram bot", "error", err)
		os.Exit(1)
	}

	voipClient := voip.NewVoipClient(config.TwilioAccountSID, config.TwilioAuthToken)

	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "url", r.URL.String())

		update, err := bot.HandleUpdate(r)
		if err != nil {
			slog.Error("Failed to handle update", "error", err)
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		if update.Message != nil {
			slog.Info(
				"Received message",
				"chat_id", update.Message.Chat.ID,
				"user_name", update.Message.Chat.UserName,
				"text", update.Message.Text)

			if update.Message.Text == "/open_underground_parking" {
				err := voipClient.MakeCall(config.SourceNumber, config.TargetNumber)
				if err != nil {
					slog.Error("Failed to open", "error", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Failed to make call: %s", err.Error()))
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Open successfully.")
					bot.Send(msg)
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	adapter := httpadapter.NewV2(root)

	slog.Info("Server starting")

	lambda.Start(adapter.ProxyWithContext)
}
