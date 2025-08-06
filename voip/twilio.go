package voip

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func MakeCall(source, target string) error {

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if accountSid == "" || authToken == "" {
		log.Panic("Please set the TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN environment variables.")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateCallParams{}
	params.SetTo(target)
	params.SetFrom(source)
	params.SetUrl("http://demo.twilio.com/docs/voice.xml")

	response, err := client.Api.CreateCall(params)
	if err != nil {
		return fmt.Errorf("call to %s: %w", target, err)
	} else {
		slog.Info("Call created successfully", "sid", *response.Sid, "status", *response.Status)
	}

	return nil
}
