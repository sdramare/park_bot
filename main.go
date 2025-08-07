package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"park_bot/voip"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("Starting app...")
	godotenv.Load()

	source := os.Getenv("TWILIO_SOURCE_NUMBER")
	if source == "" {
		log.Fatal("Please set the TWILIO_SOURCE_NUMBER environment variable.")
	}

	target := os.Getenv("TWILIO_TARGET_NUMBER")
	if target == "" {
		log.Fatal("Please set the TWILIO_TARGET_NUMBER environment variable.")
	}

	tg_token := os.Getenv("TG_BOT_API_TOKEN")
	if tg_token == "" {
		log.Fatal("Please set the TG_BOT_API_TOKEN environment variable.")
	}

	bot, err := tgbotapi.NewBotAPI(tg_token)
	if err != nil {
		log.Fatal("Could not create Telegram bot: ", err)
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	if accountSid == "" {
		log.Fatal("Please set the TWILIO_ACCOUNT_SID environment variable.")
	}

	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if authToken == "" {
		log.Fatal("Please set the TWILIO_AUTH_TOKEN environment variable.")
	}

	voipClient := voip.NewVoipClient(accountSid, authToken)

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
				err := voipClient.MakeCall(source, target)
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
