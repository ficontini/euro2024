package main

import (
	"context"

	"github.com/ficontini/euro2024/types"
)

type Service interface {
	FetchMatches(context.Context) ([]*types.Match, error)
}

type basicService struct {
	fetcher   Fetcher
	processor Processor
}

func New(fetcher Fetcher, processor Processor) Service {
	return &basicService{
		fetcher:   fetcher,
		processor: processor,
	}
}

func (svc *basicService) FetchMatches(ctx context.Context) ([]*types.Match, error) {
	res, err := svc.fetcher.FetchData()
	if err != nil {
		return nil, err
	}
	return svc.processor.ProcessData(res)
}
