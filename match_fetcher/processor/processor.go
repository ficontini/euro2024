package processor

import "github.com/ficontini/euro2024/types"

type Processor interface {
	ProcessData(any) ([]*types.Match, error)
}
