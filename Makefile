#!/usr/bin/make -f

########################################
### Setup flags

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
# VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
VERSION :=v0.1.0
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true

export GO111MODULE = on

########################################
### Process build tags

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

########################################
### Process linker flags

ldflags = -X github.com/link-chain/link/version.Version=$(VERSION) \
		  -X github.com/link-chain/link/version.Commit=$(COMMIT) \
		  -X "github.com/link-chain/link/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'


########################################
### Lint
lint: golangci-lint
	golangci-lint run
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

########################################
### Build

all: install lint test-unit

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/linkd ./cmd/linkd
	go build -mod=readonly $(BUILD_FLAGS) -o build/linkcli ./cmd/linkcli

build-docker:
	docker build -t line/link .

build-swagger-docs: statik
	statik -src=client/lcd/swagger-ui -dest=client/lcd -f -m

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/linkd
	go install $(BUILD_FLAGS) ./cmd/linkcli

install-debug: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/linkdebug

clean:
	rm -rf  build/

########################################
### Tools & dependencies

get-tools:
	go get github.com/rakyll/statik
	go get -u github.com/client9/misspell/cmd/misspell
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get github.com/cosmos/tools/cmd/runsim@v1.0.0

golangci-lint:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

statik:
	@go get github.com/rakyll/statik

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify


########################################
### Testing

test: test-all

test-all: test-unit-all test-integration-all

test-integration-all: test-integration test-integration-multi-node

test-unit-all: test-unit test-unit-race test-unit-cover

test-unit:
	@go test -mod=readonly  ./...

test-unit-race:
	@go test -mod=readonly -race  ./...

# `coverage.txt` is used in CircleCi config for the coverage report so if someone updates one, please updates the other too
test-unit-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic ./...

test-integration: build
	@go test -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_test -v

test-integration-multi-node: build-docker
	@go test -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_multi_node_test -v


########################################
### Local TestNet using docker-compose

# Run a 4-node testnet locally
testnet-start:
	$(MAKE) -C  $(CURDIR)/networks/local testnet-start

# Stop testnet
testnet-stop:
	$(MAKE) -C  $(CURDIR)/networks/local testnet-stop

testnet-test:
	$(MAKE) -C  $(CURDIR)/networks/local testnet-test

run-swagger-server:
	linkcli rest-server --trust-node=true

########################################
### Simulation

# include simulations
include sims.mk


.PHONY: all install install-debug go-mod-cache clean build\
    test test-all test-integration-all test-unit-all \
    test-unit test-unit-race test-unit-cover \
    test-integration test-integration-multi-node
