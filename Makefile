STACK_NAME :=demo
match_fetcher:
	@go build -o bin/match_fetcher ./match_fetcher
	@./bin/match_fetcher
player_fetcher:
	@go build -o bin/player_fetcher ./player_fetcher
	@./bin/player_fetcher
player_storer: delete deploy
	@go build -o bin/player_storer ./player_storer
	@./bin/player_storer
receiver:
	@go build -o bin/match_receiver ./match_receiver
	@./bin/match_receiver
match:
	@go build -o bin/match ./matchservice
	@./bin/match
player:
	@go build -o bin/player ./playerservice
	@./bin/player
gateway: 
	@go build -o bin/gateway ./gateway
	@./bin/gateway
test: 
	@go test -v ./... --count=1
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative playerservice/proto/*.proto
deploy:
	@aws cloudformation deploy --template-file cloudformation/template.yaml --stack-name $(STACK_NAME)
	@aws cloudformation wait stack-create-complete --stack-name $(STACK_NAME)
delete:
	@aws cloudformation delete-stack --stack-name $(STACK_NAME)
	@aws cloudformation wait stack-delete-complete --stack-name $(STACK_NAME)
.PHONY: match_fetcher receiver match gateway player_fetcher player_storer player
