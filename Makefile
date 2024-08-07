STACK_NAME :=demo
TEMPLATE := cloudformation/template.yaml
match_fetcher:
	@go build -o bin/match_fetcher ./matchfetcher
	@./bin/match_fetcher
player_fetcher:
	@go build -o bin/player_fetcher ./playerfetcher
	@./bin/player_fetcher
player_storer: delete deploy
	@go build -o bin/player_storer ./playerstorer
	@./bin/player_storer
match:
	@go build -o bin/match ./matchservice/cmd
	@./bin/match
player:
	@go build -o bin/player ./playerservice/cmd
	@./bin/player
gateway: 
	@go build -o bin/gateway ./gateway
	@./bin/gateway
test: 
	@go test -v ./... --count=1
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative matchservice/proto/*.proto
deploy:
	@aws cloudformation deploy --template-file $(TEMPLATE) --stack-name $(STACK_NAME)
	@aws cloudformation wait stack-create-complete --stack-name $(STACK_NAME)
delete:
	@aws cloudformation delete-stack --stack-name $(STACK_NAME)
	@aws cloudformation wait stack-delete-complete --stack-name $(STACK_NAME)
.PHONY: match_fetcher match gateway player_fetcher player_storer player
