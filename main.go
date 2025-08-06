package main

import (
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

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_API_TOKEN"))
	if err != nil {
		log.Fatal("create tgbotapi ", err)
	}

	source := os.Getenv("TWILIO_SOURCE_NUMBER")
	if source == "" {
		log.Fatal("Please set the TWILIO_SOURCE_NUMBER environment variable.")
	}

	target := os.Getenv("TWILIO_TARGET_NUMBER")
	if target == "" {
		log.Fatal("Please set the TWILIO_TARGET_NUMBER environment variable.")
	}

	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "url", r.URL.String())

		updates := bot.ListenForWebhookRespReqFormat(w, r)
		for update := range updates {

			if update.Message != nil {
				slog.Info("Received message", "chat_id", update.Message.Chat.ID, "user_name", update.Message.Chat.UserName, "text", update.Message.Text)

				if update.Message.Text == "/open" {
					err := voip.MakeCall(source, target)
					if err != nil {
						slog.Error("Failed to make call", "error", err)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to make call: "+err.Error())
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Call initiated successfully.")
						bot.Send(msg)
					}
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
