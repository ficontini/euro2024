package api

import (
	"errors"
	"strings"
)

func validateTeamParameter(value string) (string, error) {
	if len(value) == 0 {
		return "", errors.New("empty value")
	}
	if len(value) == 1 {
		return "", errors.New("invalid value")
	}
	var (
		builder strings.Builder
		first   = strings.ToUpper(string(value[0]))
		last    = strings.ToLower(string(value[1:]))
	)
	builder.WriteString(first)
	builder.WriteString(last)
	return builder.String(), nil
}
