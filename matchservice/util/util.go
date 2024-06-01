package util

import (
	matchendpoint "github.com/ficontini/euro2024/matchservice/pkg/endpoint"
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/ficontini/euro2024/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewProtoMatchResponseFromMatchResponse(resp matchendpoint.MatchResponse) *proto.MatchResponse {
	var matches []*proto.Match
	for _, m := range resp.Matches {
		matches = append(matches,
			&proto.Match{
				Home: &proto.Team{
					Team:  m.Home.Team,
					Goals: int64(m.Home.Goals),
				},
				Away: &proto.Team{
					Team:  m.Away.Team,
					Goals: int64(m.Away.Goals),
				},
				Date:   timestamppb.New(m.Date),
				Status: string(m.Status),
				Location: &proto.Location{
					City:    m.Location.City,
					Stadium: m.Location.Stadium,
				},
			})
	}
	return &proto.MatchResponse{
		Matches: matches,
	}
}
func NewMatchResponseFromProtoMatchResponse(resp *proto.MatchResponse) matchendpoint.MatchResponse {
	var matches []matchendpoint.Match
	for _, m := range resp.Matches {
		matches = append(matches,
			matchendpoint.Match{
				Home: matchendpoint.Team{
					Team:  m.Home.Team,
					Goals: int(m.Home.Goals),
				},
				Away: matchendpoint.Team{
					Team:  m.Away.Team,
					Goals: int(m.Away.Goals),
				},
				Date:   m.Date.AsTime(),
				Status: types.MatchStatus(m.Status),
				Location: matchendpoint.Location{
					City:    m.Location.City,
					Stadium: m.Location.Stadium,
				},
			},
		)
	}
	return matchendpoint.MatchResponse{
		Matches: matches,
	}
}
