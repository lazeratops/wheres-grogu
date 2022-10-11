package main

import (
	"bytes"
	"encoding/json"
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
	webhookEndpointKey = "SLACK_GROGU_WEBHOOK_URL"
	slackMessageText   = "Din Djarin has been awarded!"
)

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// TODO: do some kind of auth handling, make sure only
	// permitted person can use this webhook.

	// Make the actual HTTP request to Slack
	webhookEndpoint := os.Getenv(webhookEndpointKey)
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
