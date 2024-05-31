package transport

import (
	"context"

	matchendpoint "github.com/ficontini/euro2024/matchservice/pkg/endpoint"
	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	getUpcoming grpctransport.Handler
	getLive grpctransport.Handler
	proto.UnimplementedMatchesServer
}

func (s *grpcServer) GetUpcoming(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error) {
	_, rep, err := s.getUpcoming.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.MatchResponse), nil
}
func (s *grpcServer) GetLive(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error){
	_, rep; err := s.getLive.ServeGRPC(ctx, req)
	if err != nil{
		return nil, err
	}
	return rep.(*proto.MatchResponse), nil 
}

func NewGRPCServer(endpoints matchendpoint.Set) proto.MatchesServer {
	options := []grpctransport.ServerOption{}
	return &grpcServer{
		getUpcoming: grpctransport.NewServer(
			endpoints.GetUpcomingMatchesEndpoint,
			decodeGRPCReq,
			encodeGRPCResp,
			options...,
		),
		getLive: grpctransport.NewServer(
			endpoints.GetLiveMatchesEndpoint,
			decodeGRPCReq,
			encodeGRPCResp,
			options...
		),
	}
}
func NewGRPCClient(conn *grpc.ClientConn) service.Service {
	var (
		options     = []grpctransport.ClientOption{}
		getUpcomingEndpoint endpoint.Endpoint
		getLiveEndpoint endpoint.Endpoint
	)
	{
		getUpcomingEndpoint = grpctransport.NewClient(
			conn,
			"Matches",
			"GetUpcoming",
			encodeGRPRequest,
			decodeGRPResponse,
			proto.MatchResponse{},
			options...,
		).Endpoint()
		getLiveEndpoint = grpctransport.NewClient(
			conn,
			"Matches",
			"GetLive",
			encodeGRPRequest,
			decodeGRPResponse,
			proto.MatchResponse{},
			options...,
		).Endpoint()
	}
	return matchendpoint.Set{
		GetUpcomingMatchesEndpoint: getUpcomingEndpoint,
		GetLiveMatchesEndpoint: getLiveEndpoint,
	}

}
func decodeGRPCReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return struct{}{}, nil
}
func encodeGRPCResp(_ context.Context, resp interface{}) (interface{}, error) {
	var (
		response = resp.(matchendpoint.MatchResponse)
		matches  []*proto.Match
	)
	for _, m := range response.Matches {
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
	}, nil
}

func encodeGRPRequest(context.Context, interface{}) (request interface{}, err error) {
	return &proto.MatchRequest{}, nil
}

func decodeGRPResponse(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*proto.MatchResponse)
	var matches []matchendpoint.Match
	for _, m := range response.Matches {
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
	}, nil
}
