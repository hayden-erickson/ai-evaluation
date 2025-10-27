package main

import (
    "context"
    "fmt"
    "log"

    twilio "github.com/twilio/twilio-go"
    openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioNotifier struct {
    client   *twilio.RestClient
    from     string
    dryRun   bool
}

func newTwilioNotifier(accountSID, authToken, from string, dryRun bool) *TwilioNotifier {
    client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSID, Password: authToken})
    return &TwilioNotifier{client: client, from: from, dryRun: dryRun}
}

func (t *TwilioNotifier) SendSMS(ctx context.Context, to, body string) error {
    if t.dryRun {
        log.Printf("[dry-run] sms to %s: %s", to, body)
        return nil
    }
    params := &openapi.CreateMessageParams{}
    params.SetTo(to)
    params.SetFrom(t.from)
    params.SetBody(body)
    _, err := t.client.Api.CreateMessage(params)
    if err != nil {
        return fmt.Errorf("twilio send error: %w", err)
    }
    return nil
}


