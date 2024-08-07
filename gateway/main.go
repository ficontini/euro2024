package main

import (
	"log"
	"os"

	"github.com/ficontini/euro2024/gateway/api"
	"github.com/ficontini/euro2024/gateway/service"
	"github.com/ficontini/euro2024/gateway/store"
	"github.com/ficontini/euro2024/matchservice/pkg/transport"
	playertransport "github.com/ficontini/euro2024/playerservice/pkg/transport"
	"github.com/ficontini/euro2024/util"
	"github.com/gofiber/fiber/v2"
)

const (
	listenAddrEnVar = "GATEWAY_ADDR"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer cfg.Close()

	store, err := store.New()
	if err != nil {
		log.Fatal(err)
	}

	var (
		listenAddr     = os.Getenv(listenAddrEnVar)
		matchService   = transport.NewGRPCClient(cfg.MatchServiceConn)
		playerService  = playertransport.NewGRPCClient(cfg.PlayerServiceConn)
		userService    = service.NewUserService(store)
		authMiddleware = api.NewAuthMiddleware(userService)
		app            = fiber.New(config)
		apiv1          = app.Group("/api/v1/euro")
		users          = apiv1.Group("/user")
		matches        = apiv1.Group("/match", authMiddleware)
		team           = apiv1.Group("/team", authMiddleware)
		userHandler    = api.NewUserHandler(userService)
		matchHandler   = api.NewMatchHandler(matchService)
		playerHandler  = api.NewPlayerHandler(playerService)
	)
	users.Post("/sign-up", userHandler.HandlePostUser)
	users.Post("/sign-in", userHandler.HandleAuthenticate)
	users.Post("/sign-out", authMiddleware, userHandler.HandlePostSignOut)
	apiv1.Get("/winner", authMiddleware, matchHandler.HandleGetWinner)
	matches.Get("/upcoming", matchHandler.HandleGetUpcomingMatches)
	matches.Get("/live", matchHandler.HandleGetLiveMatches)
	team.Get("/:team/matches", matchHandler.HandleGetMatchesByTeam)
	team.Get("/:team/players", playerHandler.HandleGetPlayersByTeam)
	log.Fatal(app.Listen(listenAddr))
}

func init() {
	if err := util.Load(); err != nil {
		log.Fatal(err)
	}
}
