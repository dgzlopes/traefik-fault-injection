.PHONY: lint vendor test

export GO111MODULE=on

lint:
	golangci-lint run

vendor:
	go mod vendor

test:
	go test -v faultinjection_test.go faultinjection.go