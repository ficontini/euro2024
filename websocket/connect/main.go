package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func connectHandler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Websocket connected: %s\n", request.RequestContext.ConnectionID)
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Connected to Websocket",
	}
	return response, nil
}

func main() {
	lambda.Start(connectHandler)
}
