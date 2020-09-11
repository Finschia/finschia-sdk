#!/usr/bin/make -f

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

########################################
### Lint
lint: golangci-lint
	golangci-lint run
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

golangci-lint:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

########################################
### Testing

test: test-all

test-all: test-unit-all

test-unit-all: test-unit test-unit-race test-unit-cover

test-unit:
	@go test -mod=readonly -p 4  ./...

test-unit-race:
	@go test -mod=readonly -p 4 -race  ./...

# `coverage.txt` is used in CircleCi config for the coverage report so if someone updates one, please updates the other too
test-unit-cover:
	@go test -mod=readonly -p 4 -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: all clean \
    test test-all test-unit-all \
    test-unit test-unit-race test-unit-cover
