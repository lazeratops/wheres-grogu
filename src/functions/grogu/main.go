package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lazeratops/wheres-grogu/src/slack"
	"net/http"
)

const (
	msgFound    = "Din Djarin has been awarded!"
	msgMove     = "Din Djarin is on the move!"
	msgNotFound = "Oops! Din Djarin has not been found!"
)

var resBadParam = &events.APIGatewayProxyResponse{
	StatusCode: http.StatusBadRequest,
	Body:       "command must specify 'found', 'notfound', or 'onthemove'",
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Get "text" param from body
	fmt.Println("query str:", request.QueryStringParameters)
	var cmdParam string
	if t, ok := request.QueryStringParameters["text"]; ok {
		cmdParam = t
	} else {
		return resBadParam, nil
	}

	var slackMsgText string
	switch cmdParam {
	case "found":
		slackMsgText = msgFound
		break
	case "oopsnotfound":
		slackMsgText = msgNotFound
		break
	case "onthemove":
		slackMsgText = msgMove
		break
	default:
		return resBadParam, nil
	}

	if err := slack.PostMessage(slackMsgText, slack.ValidationParams{
		Headers: request.Headers,
		Body:    request.Body,
	}); err != nil {
		return nil, err
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
