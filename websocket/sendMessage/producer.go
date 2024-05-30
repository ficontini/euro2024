package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/euro2024/types"
)

type Producer interface {
	ProduceData(context.Context, []*types.Match) error
}

type SQSProducer struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSProducer(queueURL string, client *sqs.Client) Producer {
	return &SQSProducer{
		client:   client,
		queueURL: queueURL,
	}
}
func (svc *SQSProducer) ProduceData(ctx context.Context, matches []*types.Match) error {
	b, err := json.Marshal(matches)
	if err != nil {
		return err
	}
	_, err = svc.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(b)),
		QueueUrl:    &svc.queueURL,
	})
	return err
}
