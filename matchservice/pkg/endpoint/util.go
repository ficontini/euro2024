package matchendpoint

import (
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/ficontini/euro2024/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newProtoMatchFromMatch(m Match) *proto.Match {
	return &proto.Match{
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
		Round: string(m.Round),
	}
}
func NewProtoMatchResponseFromMatchResponse(resp MatchResponse) *proto.MatchResponse {
	var matches []*proto.Match
	for _, m := range resp.Matches {
		matches = append(matches, newProtoMatchFromMatch(m))
	}
	return &proto.MatchResponse{
		Matches: matches,
	}
}
func NewProtoWinnerResponseFromWinnerResponse(resp WinnerResponse) *proto.WinnerResponse {
	return &proto.WinnerResponse{
		Team:  resp.Team,
		Final: newProtoMatchFromMatch(resp.Final),
	}
}
func newMatchFromProtoMatch(m *proto.Match) Match {
	return Match{
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
		Round: types.RoundStatus(m.Round),
	}
}
func NewMatchResponseFromProtoMatchResponse(resp *proto.MatchResponse) MatchResponse {
	var matches []Match
	for _, m := range resp.Matches {
		matches = append(matches, newMatchFromProtoMatch(m))
	}
	return MatchResponse{
		Matches: matches,
	}
}
func NewWinnerResponseFromProtoWinnerResponse(resp *proto.WinnerResponse) WinnerResponse {
	return WinnerResponse{
		Team:  resp.Team,
		Final: newMatchFromProtoMatch(resp.Final),
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
			m.Status,
			m.Round))
	}
	return matches
}
func NewMatchFromWinnerResponse(res WinnerResponse) *types.Match {
	return &types.Match{
		Date:     res.Final.Date,
		Location: types.NewLocation(res.Final.Location.City, res.Final.Location.Stadium),
		Home:     types.NewMatchTeam(res.Final.Home.Team, res.Final.Home.Goals),
		Away:     types.NewMatchTeam(res.Final.Away.Team, res.Final.Away.Goals),
		Status:   res.Final.Status,
		Round:    res.Final.Round,
		Winner:   res.Team,
	}
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
			Round: match.Round,
		})
	}
	return MatchResponse{
		Matches: response,
	}
}
func makeWinnerResponse(match *types.Match) WinnerResponse {
	return WinnerResponse{
		Team: match.Winner,
		Final: Match{
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
			Round: match.Round,
		},
	}
}
