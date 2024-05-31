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
			endpoints.GetUpcomingMatchesEndpoint,
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
			"Matches",
			"GetUpcoming",
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
	return &proto.UpcomingResponse{
		Matches: matches,
	}, nil
}

func encodeGRPRequest(context.Context, interface{}) (request interface{}, err error) {
	return &proto.UpcomingRequest{}, nil
}

func decodeGRPResponse(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*proto.UpcomingResponse)
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
