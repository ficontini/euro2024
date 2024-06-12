package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ficontini/euro2024/types"
)

const (
	authTable = "auths"
)

type AuthStorer interface {
	Insert(context.Context, *types.Auth) error
	Get(context.Context, *types.AuthFilter) (*types.Auth, error)
	Delete(context.Context, *types.AuthFilter) error
}

type dynamoDBAuthStore struct {
	client *dynamodb.Client
	table  *string
}

func NewDynamoDBAuthStore(client *dynamodb.Client) AuthStorer {
	return &dynamoDBAuthStore{
		client: client,
		table:  aws.String(authTable),
	}
}

func (s *dynamoDBAuthStore) Insert(ctx context.Context, auth *types.Auth) error {
	item, err := attributevalue.MarshalMap(auth)
	if err != nil {
		return err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: s.table,
		Item:      item,
	})

	return err
}
func (s *dynamoDBAuthStore) Get(ctx context.Context, filter *types.AuthFilter) (*types.Auth, error) {
	key, err := attributevalue.MarshalMap(filter)
	if err != nil {
		return nil, err
	}
	res, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: s.table,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, fmt.Errorf("resource not found")
	}
	var auth *types.Auth
	if err := attributevalue.UnmarshalMap(res.Item, &auth); err != nil {
		return nil, err
	}
	return auth, nil
}
func (s *dynamoDBAuthStore) Delete(ctx context.Context, filter *types.AuthFilter) error {
	key, err := attributevalue.MarshalMap(filter)
	if err != nil {
		return err
	}
	_, err = s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: s.table,
		Key:       key,
	})
	return err
}
