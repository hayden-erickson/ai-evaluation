package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TwilioClient handles Twilio API communications
type TwilioClient struct {
	AccountSID string
	AuthToken  string
	FromNumber string
	BaseURL    string
	HTTPClient *http.Client
}

// TwilioResponse represents the response from Twilio API
type TwilioResponse struct {
	SID          string `json:"sid"`
	Status       string `json:"status"`
	To           string `json:"to"`
	From         string `json:"from"`
	Body         string `json:"body"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// NewTwilioClient creates a new Twilio client
func NewTwilioClient(accountSID, authToken, fromNumber string) *TwilioClient {
	return &TwilioClient{
		AccountSID: accountSID,
		AuthToken:  authToken,
		FromNumber: fromNumber,
		BaseURL:    fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s", accountSID),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendSMS sends an SMS message via Twilio
func (c *TwilioClient) SendSMS(toNumber, message string) error {
	// Build request URL
	requestURL := fmt.Sprintf("%s/Messages.json", c.BaseURL)

	// Build request body
	data := url.Values{}
	data.Set("To", toNumber)
	data.Set("From", c.FromNumber)
	data.Set("Body", message)

	// Create HTTP request
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.AccountSID, c.AuthToken)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Twilio API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var twilioResp TwilioResponse
	if err := json.Unmarshal(body, &twilioResp); err != nil {
		return fmt.Errorf("failed to parse Twilio response: %w", err)
	}

	// Check for error in response
	if twilioResp.ErrorCode != 0 {
		return fmt.Errorf("Twilio error %d: %s", twilioResp.ErrorCode, twilioResp.ErrorMessage)
	}

	return nil
}

// SendBatchSMS sends multiple SMS messages in batch
func (c *TwilioClient) SendBatchSMS(recipients map[string]string) error {
	var errors []error

	for phoneNumber, message := range recipients {
		if err := c.SendSMS(phoneNumber, message); err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s: %w", phoneNumber, err))
		} else {
			// Add a small delay to avoid rate limiting
			time.Sleep(100 * time.Millisecond)
		}
	}

	if len(errors) > 0 {
		var errorMessages bytes.Buffer
		for _, err := range errors {
			errorMessages.WriteString(err.Error())
			errorMessages.WriteString("; ")
		}
		return fmt.Errorf("batch send errors: %s", errorMessages.String())
	}

	return nil
}
