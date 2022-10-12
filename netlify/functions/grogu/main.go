package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"net/http"
	"os"
)

// The POST request data which will be sent to Slack
type message struct {
	Text string `json:"text"`
}

const (
	keyWebhookEndpint     = "SLACK_GROGU_WEBHOOK_URL"
	keySlackSigningSecret = "SLACK_SIGNING_SECRET"
	slackMessageText      = "Din Djarin has been awarded!"
)

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	webhookEndpoint := os.Getenv(keyWebhookEndpint)
	slackSigningSecret := os.Getenv(keySlackSigningSecret)
	if webhookEndpoint == "" || slackSigningSecret == "" {
		return nil, errors.New("function not configured; missing critical env variables")
	}
	fmt.Println("request body", request.Body, request.Path)
	if err := validateSignature(request.Headers, request.Body, slackSigningSecret); err != nil {
		return nil, err
	}

	// Make the actual HTTP request to Slack
	msg := message{Text: slackMessageText}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	req, err := http.NewRequest("POST", webhookEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to Slack hook: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send Slack message: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Slack response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected Slack status code: %d (body: %s)", res.StatusCode, resBody)
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func validateSignature(headers map[string]string, body string, signingSecret string) error {
	var slackSig, reqTimestamp string

	const errMsg = "unauthorized request"
	fmt.Println("headers", headers)

	if sig, ok := headers["x-slack-signature"]; ok {
		slackSig = sig
	} else {
		return fmt.Errorf("%s: missing signature", errMsg)
	}
	// Check if the key had an empty value just in case
	if slackSig == "" {
		return fmt.Errorf("%s: invalid signature", errMsg)
	}

	if timestamp, ok := headers["x-slack-request-timestamp"]; ok {
		reqTimestamp = timestamp
	} else {
		return fmt.Errorf("%s: missing request timestamp", errMsg)
	}

	version := getVersion(reqTimestamp, body)
	fmt.Println("versionNumber:", version)
	sig := computeSig(version, signingSecret)
	if hmac.Equal([]byte(sig), []byte(slackSig)) {
		return nil
	}
	return errors.New("failed to validate given auth signature")
}

func getVersion(timestamp string, body string) string {
	return fmt.Sprintf("v0:%s:%s", timestamp, body)
}

func computeSig(version, signingSecret string) string {
	sig := hmac.New(sha256.New, []byte(signingSecret))
	sig.Write([]byte(version))
	return fmt.Sprintf("v0=%s", base64.StdEncoding.EncodeToString(sig.Sum(nil)))
}
