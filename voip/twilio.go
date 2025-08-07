package voip

import (
	"fmt"
	"log/slog"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Voip struct {
	client *twilio.RestClient
}

func NewVoipClient(accountSid, authToken string) *Voip {
	slog.Info("Initializing Twilio client", "account_sid", accountSid)
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &Voip{
		client,
	}
}

func (v *Voip) MakeCall(source, target string) error {

	if source == "" || target == "" {
		return fmt.Errorf("source and target numbers must be provided")
	}

	params := &twilioApi.CreateCallParams{}
	params.SetTo(target)
	params.SetFrom(source)
	params.SetUrl("http://demo.twilio.com/docs/voice.xml")

	response, err := v.client.Api.CreateCall(params)
	if err != nil {
		return fmt.Errorf("call to %s: %w", target, err)
	} else {
		slog.Info("Call created successfully", "sid", *response.Sid, "status", *response.Status)
	}

	return nil
}
