package transport

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/ficontini/euro2024/matchservice/service"
	"github.com/ficontini/euro2024/types"
)

type Consumer interface {
	Start(context.Context)
	Stop(context.Context)
}

type sqsConsumer struct {
	client   *sqs.Client
	queueURL string
	service  service.Service
	stopch   chan struct{}
	wg       sync.WaitGroup
}

func NewSQSConsumer(client *sqs.Client, queueURL string, service service.Service) Consumer {
	return &sqsConsumer{
		client:   client,
		queueURL: queueURL,
		service:  service,
		stopch:   make(chan struct{}),
	}
}
func (c *sqsConsumer) Start(ctx context.Context) {
	c.wg.Add(1)
	defer c.wg.Done()
	go func() {
		for {
			select {
			case <-c.stopch:
				return
			default:
				if err := c.consumeMessage(ctx); err != nil {
					log.Println("error consuming msg: ", err)
				}
			}
		}
	}()
}
func (c *sqsConsumer) Stop(ctx context.Context) {
	close(c.stopch)
	c.wg.Wait()
}
func (c *sqsConsumer) consumeMessage(ctx context.Context) error {
	res, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: &c.queueURL,
	})
	if err != nil {
		return err
	}
	for _, msg := range res.Messages {
		var matches []*types.Match
		if err := json.Unmarshal([]byte(*msg.Body), &matches); err != nil {
			log.Println("error unmarshalling msg: ", err)
			continue
		}
		if err := c.processData(ctx, matches); err != nil {
			log.Println("error processing data: ", err)
		}
		_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &c.queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Println("Error deleting message:", err)
		}
	}
	time.Sleep(5 * time.Minute)
	return nil
}

func (c *sqsConsumer) processData(ctx context.Context, matches []*types.Match) error {
	return c.service.ProcessData(ctx, matches)
}
