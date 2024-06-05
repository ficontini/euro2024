macth_fetcher:
	@go build -o bin/match_fetcher ./match_fetcher
	@./bin/match_fetcher
team_fetcher:
	@go build -o bin/team_fetcher ./team_fetcher
	@./bin/team_fetcher
receiver:
	@go build -o bin/match_receiver ./match_receiver
	@./bin/match_receiver
match:
	@go build -o bin/match ./matchservice
	@./bin/match
gateway: 
	@go build -o bin/gateway ./gateway
	@./bin/gateway
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative matchservice/proto/*.proto

.PHONY: match_fetcher receiver match gateway team_fetcher
