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

const (
	userTable        = "users"
	userPartitionKey = "email"
)

type UserStorer interface {
	Insert(context.Context, *types.User) error
	GetByEmail(context.Context, string) (*types.User, error)
}

type dynamoDBUserStore struct {
	client *dynamodb.Client
	table  *string
}

func NewDynamoDBUserStore(client *dynamodb.Client) UserStorer {
	return &dynamoDBUserStore{
		client: client,
		table:  aws.String(userTable),
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
func (s *dynamoDBUserStore) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	key, err := attributevalue.Marshal(email)
	if err != nil {
		return nil, err
	}
	res, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: s.table,
		Key:       map[string]dynamodbtypes.AttributeValue{userPartitionKey: key},
	})
	if err != nil {
		return nil, err
	}
	var user *types.User
	if err := attributevalue.UnmarshalMap(res.Item, &user); err != nil {
		return nil, err
	}
	return user, nil
}

var ErrEmailAlreadyInUse = errors.New("email already in use")
