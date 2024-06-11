package main

import (
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	matchGrpcAddrEnvVar  = "GRPC_LISTENER"
	playerGrpcAddrEnvVar = "PLAYER_GRPC_LISTENER"
)

type Config struct {
	MatchServiceConn  *grpc.ClientConn
	PlayerServiceConn *grpc.ClientConn
}

type OptionFunc func(*Config) error

func defaultConfig() (*Config, error) {
	var (
		matchGRPCAddr  = os.Getenv(matchGrpcAddrEnvVar)
		playerGRPCAddr = os.Getenv(playerGrpcAddrEnvVar)
	)
	if matchGRPCAddr == "" {
		return nil, fmt.Errorf("%s env var not set", matchGrpcAddrEnvVar)
	}
	if playerGRPCAddr == "" {
		return nil, fmt.Errorf("%s env var not set", playerGrpcAddrEnvVar)
	}
	matchServiceConn, err := newGRPCConnection(matchGRPCAddr)
	if err != nil {
		return nil, err
	}

	playerServiceConn, err := newGRPCConnection(playerGRPCAddr)
	if err != nil {
		return nil, err
	}
	return &Config{
		MatchServiceConn:  matchServiceConn,
		PlayerServiceConn: playerServiceConn,
	}, nil
}

func NewConfig(options ...OptionFunc) (*Config, error) {
	cfg, err := defaultConfig()
	if err != nil {
		return nil, err
	}
	for _, fn := range options {
		if err := fn(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func WithMatchServiceConn(conn *grpc.ClientConn) OptionFunc {
	return func(c *Config) error {
		if conn == nil {
			return fmt.Errorf("error establishing the match service connection")
		}
		c.MatchServiceConn = conn
		return nil
	}
}

func WithPlayerServiceConn(conn *grpc.ClientConn) OptionFunc {
	return func(c *Config) error {
		if conn == nil {
			return fmt.Errorf("error establishing the player service connection")
		}
		c.PlayerServiceConn = conn
		return nil
	}
}

func (c *Config) Close() {
	if c.MatchServiceConn != nil {
		c.MatchServiceConn.Close()
	}
	if c.PlayerServiceConn != nil {
		c.PlayerServiceConn.Close()
	}
}

func newGRPCConnection(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
