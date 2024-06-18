package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/euro2024/types"
	"github.com/google/uuid"
)

type Producer interface {
	ProduceData(context.Context, []*types.Player) error
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
func (svc *SQSProducer) ProduceData(ctx context.Context, players []*types.Player) error {
	var (
		i    int
		step = 10
	)
	for i < len(players) {
		j := i + step
		if j > len(players) {
			j = len(players)
		}
		if err := svc.batchRequestEntry(ctx, players[i:j]); err != nil {
			return err
		}
		i += step
	}
	return nil
}

func (svc *SQSProducer) batchRequestEntry(ctx context.Context, players []*types.Player) error {
	var (
		entries []sqstypes.SendMessageBatchRequestEntry
	)
	for i := 0; i < len(players); i++ {
		entry, err := createEntry(players[i])
		if err != nil {
			return err
		}
		entries = append(entries, entry)
	}
	_, err := svc.client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: &svc.queueURL,
	})
	return err
}

func createEntry(player *types.Player) (sqstypes.SendMessageBatchRequestEntry, error) {
	var (
		entry sqstypes.SendMessageBatchRequestEntry
	)
	b, err := json.Marshal(player)
	if err != nil {
		return entry, err
	}
	entry = sqstypes.SendMessageBatchRequestEntry{
		Id:          aws.String(uuid.NewString()),
		MessageBody: aws.String(string(b)),
	}
	return entry, nil
}
