package transport

import (
	"context"

	matchendpoint "github.com/ficontini/euro2024/matchservice/endpoint"
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/ficontini/euro2024/matchservice/service"
	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type grpcServer struct {
	getUpcoming grpctransport.Handler
	proto.UnimplementedMatchesServer
}

func (s *grpcServer) GetUpcoming(ctx context.Context, req *proto.UpcomingRequest) (*proto.UpcomingResponse, error) {
	_, rep, err := s.getUpcoming.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.UpcomingResponse), nil
}

func NewGRPCServer(endpoints matchendpoint.Set) proto.MatchesServer {
	options := []grpctransport.ServerOption{}
	return &grpcServer{
		getUpcoming: grpctransport.NewServer(
			endpoints.GetLiveMatchesEndpoint,
			decodeGRPCReq,
			encodeGRPCResp,
			options...,
		),
	}
}
func NewGRPCClient(conn *grpc.ClientConn) service.Service {
	var (
		options     = []grpctransport.ClientOption{}
		getEndpoint endpoint.Endpoint
	)
	{
		getEndpoint = grpctransport.NewClient(
			conn,
			"proto.GetUpcoming",
			"Matches",
			encodeGRPRequest,
			decodeGRPResponse,
			proto.UpcomingResponse{},
			options...,
		).Endpoint()
	}
	return matchendpoint.Set{
		GetUpcomingMatchesEndpoint: getEndpoint,
	}

}
func decodeGRPCReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return struct{}{}, nil
}
func encodeGRPCResp(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*proto.UpcomingResponse)
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
	response := matchendpoint.MatchResponse{
		Matches: matches,
	}
	return response, nil
}

func encodeGRPRequest(context.Context, interface{}) (request interface{}, err error) {
	return nil, nil
}

func decodeGRPResponse(context.Context, interface{}) (response interface{}, err error) {
	return nil, nil
}
