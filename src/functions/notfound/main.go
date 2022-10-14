package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lazeratops/wheres-grogu/src/slack"
	"net/http"
)

const (
	slackMessageText = "Oops! Din Djarin has not been found!"
)

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if err := slack.PostMessage(slackMessageText, slack.ValidationParams{
		Headers: request.Headers,
		Body:    request.Body,
	}); err != nil {
		return nil, err
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
