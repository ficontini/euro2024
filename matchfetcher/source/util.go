package util

import (
	"github.com/ficontini/euro2024/types"
)

const (
	FINAL     = "Final"
	SEMIFINAL = "Semi-finals"
	QUARTER   = "Quarter-finals"
	ROUNDOF16 = "Round of 16"
)

func GetRound(round string) types.RoundStatus {
	switch round {
	case FINAL:
		return types.FINAL
	case SEMIFINAL:
		return types.SEMIFINAL
	case QUARTER:
		return types.QUARTER
	case ROUNDOF16:
		return types.ROUNDOF16
	default:
		return types.GROUPS
	}
}
