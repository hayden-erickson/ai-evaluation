package main

import (
    "context"
    "fmt"
    "log"

    twilio "github.com/twilio/twilio-go"
    openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioNotifier struct {
    client *twilio.RestClient
    from   string
    dryRun bool
}

func newTwilioNotifier(accountSID, authToken, from string, dryRun bool) *TwilioNotifier {
    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: accountSID,
        Password: authToken,
    })
    return &TwilioNotifier{client: client, from: from, dryRun: dryRun}
}

func (t *TwilioNotifier) SendSMS(ctx context.Context, to string, body string) error {
    if t.dryRun {
        log.Printf("[dry-run] would send SMS to %s: %s", to, body)
        return nil
    }
    params := &openapi.CreateMessageParams{}
    params.SetTo(to)
    params.SetFrom(t.from)
    params.SetBody(body)

    // Use WithContext if available in this SDK version
    var (
        err error
    )
    if creator, ok := interface{}(t.client.ApiV2010).(interface{ CreateMessageWithContext(context.Context, *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error) }); ok {
        _, err = creator.CreateMessageWithContext(ctx, params)
    } else {
        _, err = t.client.ApiV2010.CreateMessage(params)
    }
    if err != nil {
        return fmt.Errorf("twilio send failed: %w", err)
    }
    return nil
}

package main

import (
	"context"
	"fmt"
	"log"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioNotifier struct {
	client *twilio.RestClient
	from   string
	dryRun bool
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
