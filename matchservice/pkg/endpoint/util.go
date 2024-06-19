package matchendpoint

import (
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/ficontini/euro2024/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewProtoMatchResponseFromMatchResponse(resp MatchResponse) *proto.MatchResponse {
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
func NewMatchResponseFromProtoMatchResponse(resp *proto.MatchResponse) MatchResponse {
	var matches []Match
	for _, m := range resp.Matches {
		matches = append(matches,
			Match{
				Home: Team{
					Team:  m.Home.Team,
					Goals: int(m.Home.Goals),
				},
				Away: Team{
					Team:  m.Away.Team,
					Goals: int(m.Away.Goals),
				},
				Date:   m.Date.AsTime(),
				Status: types.MatchStatus(m.Status),
				Location: Location{
					City:    m.Location.City,
					Stadium: m.Location.Stadium,
				},
			},
		)
	}
	return MatchResponse{
		Matches: matches,
	}
}

func NewMatchesFromMatchResponse(response MatchResponse) []*types.Match {
	var matches []*types.Match
	for _, m := range response.Matches {
		matches = append(matches, types.NewMatch(
			m.Date,
			types.NewLocation(m.Location.City, m.Location.Stadium),
			types.NewMatchTeam(m.Home.Team, m.Home.Goals),
			types.NewMatchTeam(m.Away.Team, m.Away.Goals),
			m.Status))
	}
	return matches
}

func makeResponse(matches []*types.Match) MatchResponse {
	var response []Match
	for _, match := range matches {
		response = append(response, Match{
			Date: match.Date,
			Location: Location{
				match.Location.City,
				match.Location.Stadium,
			},
			Status: match.Status,
			Home: Team{
				Team:  match.Home.Name,
				Goals: match.Home.Goals,
			},
			Away: Team{
				Team:  match.Away.Name,
				Goals: match.Away.Goals,
			},
		})
	}
	return MatchResponse{
		Matches: response,
	}
}
