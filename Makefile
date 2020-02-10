#!/usr/bin/make -f

########################################
### Setup flags

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
BASE_VERSION := $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BASE_VERSION := $(if $(BASE_VERSION), $(BASE_VERSION), v0.0.0)
VERSION := $(BASE_VERSION)-$(shell basename $(shell git symbolic-ref -q HEAD --short))+$(shell date '+%Y%m%d%H%M%S')
VERSION := $(strip $(VERSION))
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

ldflags = -X github.com/line/link/version.Version=$(VERSION) \
		  -X github.com/line/link/version.Commit=$(COMMIT) \
		  -X "github.com/line/link/version.BuildTags=$(build_tags_comma_sep)"

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

build: go.sum build-swagger-docs
	go build -mod=readonly $(BUILD_FLAGS) -o build/linkd ./cmd/linkd
	go build -mod=readonly $(BUILD_FLAGS) -o build/linkcli ./cmd/linkcli

build-contract-test-hook:
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_test_hook ./cmd/contract_test_hook

build-docker:
	docker build -t line/link .

build-swagger-docs: statik versioning-swagger-docs
	statik -src=client/lcd/swagger-ui -dest=client/lcd -f -m

versioning-swagger-docs:
	perl -pi -e 's/version: "v[^\s]+"/version: "$(strip $(BASE_VERSION))"/' client/lcd/swagger-ui/swagger.yaml

install: go.sum build-swagger-docs
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

yq:
	@go get github.com/mikefarah/yq/v2@v2.4.1

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
	@go test -mod=readonly `go list ./cli_test/...` -tags=cli_test -v

test-integration-multi-node: build-docker
	@go test -mod=readonly `go list ./cli_test/...` -tags=cli_multi_node_test -v


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

setup-contract-test-data: build build-contract-test-hook yq
	echo 'Prepare data for the contract tests' ; \
	./contract_test/testdata/prepare_dredd.sh ; \
	./contract_test/testdata/prepare_chain_state.sh

start-link: setup-contract-test-data
	pkill linkd || true
	./build/linkd --home /tmp/contract_test/.linkd start &
	@sleep 5s
	./contract_test/testdata/wait-for-it.sh localhost 26657

setup-transactions: start-link
	@bash ./contract_test/testdata/setup.sh

contract-test: setup-transactions
	@echo "Running LINK LCD for contract tests"
	@bash ./contract_test/testdata/generate_tx_iteratively.sh &
	./contract_test/testdata/run_dredd.sh

run-lcd-contract-test:
	@echo "Running LINK LCD for contract tests"
	lsof -i tcp:1317 | grep -v PID | awk '{print $$2}' | xargs kill || true
	./build/linkcli rest-server --laddr tcp://0.0.0.0:1317 --home /tmp/contract_test/.linkcli --node http://localhost:26657 --chain-id lcd --trust-node || true

dredd-test:
	cp client/lcd/swagger-ui/swagger.yaml /tmp/contract_test/swagger.yaml
	@bash ./contract_test/testdata/replace_symbols.sh --replace_tx_hash
	@bash ./contract_test/testdata/generate_tx_iteratively.sh &
	./contract_test/testdata/wait-for-it.sh localhost 26657
	dredd; pkill -f ./contract_test/testdata/generate_tx_iteratively.sh

stop-dredd-test:
	./contract_test/testdata/stop_dredd_test.sh

########################################
### Simulation

# include simulations
include sims.mk


.PHONY: all install install-debug go-mod-cache clean build\
    test test-all test-integration-all test-unit-all \
    test-unit test-unit-race test-unit-cover \
    test-integration test-integration-multi-node
