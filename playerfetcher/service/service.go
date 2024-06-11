package service

import (
	"github.com/ficontini/euro2024/types"
)

type Service interface {
	FetchData() ([]*types.Player, error)
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
func (svc *basicService) FetchData() ([]*types.Player, error) {
	resp, err := svc.fetcher.FetchData()
	if err != nil {
		return nil, err
	}
	return svc.processor.ProcessData(resp), nil
}
