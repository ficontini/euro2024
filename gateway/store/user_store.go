package store

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ficontini/euro2024/types"
)

const table = "users"

type UserStorer interface {
	Insert(context.Context, *types.User) error
}

type dynamoDBUserStore struct {
	client *dynamodb.Client
	table  *string
}

func NewDynamoDBUserStore(client *dynamodb.Client) UserStorer {
	return &dynamoDBUserStore{
		client: client,
		table:  aws.String(table),
	}
}
func (s *dynamoDBUserStore) Insert(ctx context.Context, user *types.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           s.table,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(email)"),
	})
	var ae *dynamodbtypes.ConditionalCheckFailedException
	if errors.As(err, &ae) {
		return ErrEmailAlreadyInUse
	}
	return err
}

var ErrEmailAlreadyInUse = errors.New("email already in use")
