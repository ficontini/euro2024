fetcher:
	@go build -o bin/match_fetcher ./match_fetcher
	@./bin/match_fetcher
receiver:
	@go build -o bin/match_receiver ./match_receiver
	@./bin/match_receiver
match:
	@go build -o bin/match ./matchservice
	@./bin/match
.PHONY: fetcher receiver match