export CONFIG_PATH=config/local.yaml

.PHONY: server
server:
	go run ./cmd/server

.PHONY: client
client:
	go run ./cmd/client

