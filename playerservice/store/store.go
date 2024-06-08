package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ficontini/euro2024/types"
)

const (
	playersTable = "players"
	keyGSI       = "team"
	gsi          = "teamGSI"
)

type Storer interface {
	GetPlayersByTeam(context.Context, string) ([]*types.Player, error)
}

type DynamoDBStore struct {
	client *dynamodb.Client
	table  *string
	index  *string
}

func NewDynamoDBStore(client *dynamodb.Client) Storer {
	return &DynamoDBStore{
		client: client,
		table:  aws.String(playersTable),
		index:  aws.String(gsi),
	}
}
func (s *DynamoDBStore) GetPlayersByTeam(ctx context.Context, team string) ([]*types.Player, error) {
	var (
		keyCond = expression.Key("team").Equal(expression.Value(team))
		players []*types.Player
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}
	res, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 s.table,
		IndexName:                 s.index,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}
	if err := attributevalue.UnmarshalListOfMaps(res.Items, &players); err != nil {
		return nil, err
	}

	return players, nil
}
