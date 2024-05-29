package transport

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ficontini/euro2024/matchservice/service"
	"github.com/ficontini/euro2024/types"
)

func NewSQSConsumer(client *sqs.Client, queueURL string, svc service.Service) {
	ctx := context.Background()
	for {
		result, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: &queueURL,
		})
		if err != nil {
			log.Println("Error receiving message:", err)
			continue
		}

		if len(result.Messages) > 0 {
			msg := result.Messages[0]
			var matches []*types.Match
			if err := json.Unmarshal([]byte(*msg.Body), &matches); err != nil {
				log.Println("error unmarshalling msg", err)
			}
			ProcessData(ctx, matches, svc)
			_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Println("Error deleting message:", err)
			}
		}
	}
}

func ProcessData(ctx context.Context, matches []*types.Match, svc service.Service) {
	if err := svc.ProcessData(ctx, matches); err != nil {
		log.Fatal(err)
	}
}
