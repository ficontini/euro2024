package transport

import (
	"context"

	matchendpoint "github.com/ficontini/euro2024/matchservice/pkg/endpoint"
	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type grpcServer struct {
	getUpcoming grpctransport.Handler
	getLive     grpctransport.Handler
	getByTeam   grpctransport.Handler
	getWinner   grpctransport.Handler
	proto.UnimplementedMatchesServer
}

func (s *grpcServer) GetUpcoming(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error) {
	_, rep, err := s.getUpcoming.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.MatchResponse), nil
}
func (s *grpcServer) GetLive(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error) {
	_, rep, err := s.getLive.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.MatchResponse), nil
}
func (s *grpcServer) GetByTeam(ctx context.Context, req *proto.TeamRequest) (*proto.MatchResponse, error) {
	_, rep, err := s.getByTeam.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.MatchResponse), nil
}
func (s *grpcServer) GetEuroWinner(ctx context.Context, req *proto.WinnerRequest) (*proto.WinnerResponse, error) {
	_, rep, err := s.getWinner.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.WinnerResponse), nil
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
			options...,
		),
		getByTeam: grpctransport.NewServer(
			endpoints.GetMatchesByTeamEndpoint,
			decodeGRPCTeamReq,
			encodeGRPCResp,
			options...,
		),
		getWinner: grpctransport.NewServer(
			endpoints.GetEuroWinnerEndpoint,
			decodeGRPCReq,
			encodeGRPCWinnerResp,
			options...,
		),
	}
}
func NewGRPCClient(conn *grpc.ClientConn) service.Service {
	var (
		options             = []grpctransport.ClientOption{}
		getUpcomingEndpoint endpoint.Endpoint
		getLiveEndpoint     endpoint.Endpoint
		getByTeamEndpoint   endpoint.Endpoint
		getWinnerEndpoint   endpoint.Endpoint
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
		getByTeamEndpoint = grpctransport.NewClient(
			conn,
			"Matches",
			"GetByTeam",
			encodeGRPTeamRequest,
			decodeGRPResponse,
			proto.MatchResponse{},
			options...,
		).Endpoint()
		getWinnerEndpoint = grpctransport.NewClient(
			conn,
			"Matches",
			"GetEuroWinner",
			encodeGRCPEmptyRequest,
			decodeGRPCWinnerResponse,
			proto.WinnerResponse{},
			options...,
		).Endpoint()

	}
	return matchendpoint.Set{
		GetUpcomingMatchesEndpoint: getUpcomingEndpoint,
		GetLiveMatchesEndpoint:     getLiveEndpoint,
		GetMatchesByTeamEndpoint:   getByTeamEndpoint,
		GetEuroWinnerEndpoint:      getWinnerEndpoint,
	}

}
func decodeGRPCReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return struct{}{}, nil
}
func decodeGRPCTeamReq(_ context.Context, grpcRep interface{}) (interface{}, error) {
	req := grpcRep.(*proto.TeamRequest)
	return &matchendpoint.TeamRequest{Team: req.Name}, nil
}
func encodeGRPCResp(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(matchendpoint.MatchResponse)
	return matchendpoint.NewProtoMatchResponseFromMatchResponse(response), nil

}
func encodeGRPCWinnerResp(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(matchendpoint.WinnerResponse)
	return matchendpoint.NewProtoWinnerResponseFromWinnerResponse(response), nil
}

func encodeGRPRequest(context.Context, interface{}) (request interface{}, err error) {
	return &proto.MatchRequest{}, nil
}
func encodeGRCPEmptyRequest(context.Context, interface{}) (request interface{}, err error) {
	return &proto.WinnerRequest{}, nil
}
func encodeGRPTeamRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*matchendpoint.TeamRequest)
	return &proto.TeamRequest{
		Name: request.Team,
	}, nil
}
func decodeGRPResponse(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*proto.MatchResponse)
	return matchendpoint.NewMatchResponseFromProtoMatchResponse(response), nil
}
func decodeGRPCWinnerResponse(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*proto.WinnerResponse)
	return matchendpoint.NewWinnerResponseFromProtoWinnerResponse(response), nil
}
