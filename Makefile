.PHONY: test up lint

SERVER_DIR = server

test:
	cd $(SERVER_DIR) && go test -v ./...

lint:
	cd $(SERVER_DIR) && golangci-lint run

up:
	cd $(SERVER_DIR) && air
