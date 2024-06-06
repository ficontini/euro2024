package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/ficontini/euro2024/types"
)

type WebSocketHandler struct {
	client   *apigatewaymanagementapi.Client
	producer Producer
}

func NewWebSocketHandler(client *apigatewaymanagementapi.Client, producer Producer) *WebSocketHandler {
	return &WebSocketHandler{
		client:   client,
		producer: producer,
	}
}
func (h *WebSocketHandler) HandleMessage(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := event.RequestContext.ConnectionID
	responseMessage := "received"

	var response Response
	if err := json.Unmarshal([]byte(event.Body), &response); err != nil {
		log.Fatal(err)
	}
	if err := h.producer.ProduceData(ctx, response.Matches); err != nil {
		log.Fatal(err)
	}
	_, err := h.client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte(responseMessage),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

type Response struct {
	Matches []*types.Match `json:"matches"`
}
