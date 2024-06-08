package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/euro2024/types"
)

const (
	playersTable = "players"
)

type Storer interface {
	InsertPlayer(context.Context, *types.Player) error
}

type DynamoDBStore struct {
	client *dynamodb.Client
	table  *string
}

func NewDynamoDBStore(client *dynamodb.Client) Storer {
	return &DynamoDBStore{
		client: client,
		table:  aws.String(playersTable),
	}
}
func (s *DynamoDBStore) InsertPlayer(ctx context.Context, player *types.Player) error {
	item, err := attributevalue.MarshalMap(player)
	if err != nil {
		return err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: s.table,
		Item:      item,
	})
	return err
}
