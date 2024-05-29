package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ficontini/euro2024/types"
)

type Service interface {
	SendData(context.Context, []*types.Match) error
}

type basicService struct {
	client   *sqs.Client
	queueURL string
}

func New(queueURL string) (Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := sqs.NewFromConfig(cfg)
	return &basicService{
		client:   client,
		queueURL: queueURL,
	}, nil
}
func (svc *basicService) SendData(ctx context.Context, matches []*types.Match) error {
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
