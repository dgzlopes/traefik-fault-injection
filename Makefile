.PHONY: lint vendor test k6

export GO111MODULE=on

lint:
	golangci-lint run

vendor:
	go mod vendor

test:
	go test -v faultinjection_test.go faultinjection.go

smoke_test:
	k6 run -e MY_HOSTNAME=localhost:3456 smoke_test.js