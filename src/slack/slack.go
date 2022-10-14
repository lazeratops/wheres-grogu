package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	keyWebhookEndpoint    = "SLACK_GROGU_WEBHOOK_URL"
	keySlackSigningSecret = "SLACK_SIGNING_SECRET"
)

// message represents the POST request data which will be sent to Slack
type message struct {
	Text string `json:"text"`
}

// PostMessage validates the request using the given validation params
// and sends the given message string to the Slack webhook endpoint
func PostMessage(text string, vp ValidationParams) error {
	webhookEndpoint := os.Getenv(keyWebhookEndpoint)
	slackSigningSecret := os.Getenv(keySlackSigningSecret)
	if webhookEndpoint == "" || slackSigningSecret == "" {
		return errors.New("function not configured; missing critical env variables")
	}
	if err := validateSignature(vp.Headers, vp.Body, slackSigningSecret); err != nil {
		return err
	}
	return postMessage(text, webhookEndpoint)
}

func postMessage(text, endpoint string) error {
	msg := message{Text: text}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create POST request to Slack hook: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read Slack response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected Slack status code: %d (body: %s)", res.StatusCode, resBody)
	}
	return nil
}
