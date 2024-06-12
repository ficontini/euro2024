package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Store struct {
	User UserStorer
	Auth AuthStorer
}

func New() (*Store, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &Store{
		User: NewDynamoDBUserStore(client),
		Auth: NewDynamoDBAuthStore(client),
	}, nil
}
