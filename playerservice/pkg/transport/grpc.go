package transport

import (
	"context"

	playerendpoint "github.com/ficontini/euro2024/playerservice/pkg/endpoint"
	"github.com/ficontini/euro2024/playerservice/pkg/service"
	"github.com/ficontini/euro2024/playerservice/proto"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type grpcServer struct {
	getPlayersByTeam grpctransport.Handler
	proto.UnimplementedPlayersServer
}

func (s *grpcServer) GetByTeam(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	_, rep, err := s.getPlayersByTeam.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.Response), nil
}
func NewGRPCServer(endpoints playerendpoint.Set) proto.PlayersServer {
	options := []grpctransport.ServerOption{}
	return &grpcServer{
		getPlayersByTeam: grpctransport.NewServer(
			endpoints.GetPlayersByTeamEndpoint,
			decodeGRPCReq,
			encodeGRPCResp,
			options...,
		),
	}
}
func NewGRPCClient(conn *grpc.ClientConn) service.Service {
	var (
		options                  = []grpctransport.ClientOption{}
		getPlayersByTeamEndpoint endpoint.Endpoint
	)
	{
		getPlayersByTeamEndpoint = grpctransport.NewClient(
			conn,
			"Players",
			"GetByTeam",
			encodeGRPCRequest,
			decodeGRPCResponse,
			proto.Response{},
			options...,
		).Endpoint()
	}
	return playerendpoint.Set{
		GetPlayersByTeamEndpoint: getPlayersByTeamEndpoint,
	}

}
func decodeGRPCReq(_ context.Context, grpcRep interface{}) (interface{}, error) {
	req := grpcRep.(*proto.Request)
	return playerendpoint.Request{Team: req.Name}, nil
}
func encodeGRPCResp(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(playerendpoint.Response)

	var players []*proto.Player
	for _, r := range response.Players {
		players = append(players, newProtoPlayer(r))
	}
	return &proto.Response{
		Players: players,
	}, nil
}

func newProtoPlayer(player playerendpoint.Player) *proto.Player {
	return &proto.Player{
		FirstName:    player.FirstName,
		LastName:     player.LastName,
		Age:          int64(player.Age),
		Position:     player.Position,
		Goals:        int64(player.Goals),
		Assists:      int64(player.Assists),
		PassAccuracy: int64(player.PassAccuracy),
		YellowCards:  int64(player.YellowCards),
		RedCards:     int64(player.RedCards),
	}
}
func encodeGRPCRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(playerendpoint.Request)
	return &proto.Request{
		Name: request.Team,
	}, nil
}
func decodeGRPCResponse(_ context.Context, resp interface{}) (interface{}, error) {
	protoResp := resp.(*proto.Response)

	var players []playerendpoint.Player
	for _, p := range protoResp.Players {
		players = append(players, newPlayer(p))
	}
	return playerendpoint.Response{Players: players}, nil
}

func newPlayer(player *proto.Player) playerendpoint.Player {
	return playerendpoint.Player{
		FirstName:    player.FirstName,
		LastName:     player.LastName,
		Age:          int(player.Age),
		Position:     player.Position,
		Goals:        int(player.Goals),
		Assists:      int(player.Assists),
		PassAccuracy: int(player.PassAccuracy),
		YellowCards:  int(player.YellowCards),
		RedCards:     int(player.RedCards),
	}
}
