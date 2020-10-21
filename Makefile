.PHONY: lint vendor

export GO111MODULE=on

lint:
	golangci-lint run

vendor:
	go mod vendor