########################################
### Process build tags

ifeq ($(WITH_CLEVELDB),yes)
  CGO_ENABLED=1
  build_tags += cleveldb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))
src := $(wildcard app/*.go) $(wildcard cmd/*.go) $(wildcard types/*.go)

########################################
### Process linker flags

ldflags = -X github.com/line/lbm-sdk/version.Name=link \
		  -X github.com/line/lbm-sdk/version.ServerName=linkd \
		  -X github.com/line/lbm-sdk/version.ClientName=linkcli \
		  -X github.com/line/lbm-sdk/version.Version=$(VERSION) \
		  -X github.com/line/lbm-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/line/lbm-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/line/lbm-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
CLI_TEST_BUILD_FLAGS := -tags "cli_test $(build_tags)"

########################################
### Lint

lint: $(src)
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

########################################
### Build

all: build install lint

build: go.sum $(src)
	CGO_ENABLED=$(CGO_ENABLED) go build -mod=readonly $(BUILD_FLAGS) -o build/linkwasmd ./cmd/linkwasmd
	CGO_ENABLED=$(CGO_ENABLED) go build -mod=readonly $(BUILD_FLAGS) -o build/linkwasmcli ./cmd/linkwasmcli

install: go.sum $(src)
	CGO_ENABLED=$(CGO_ENABLED) go install $(BUILD_FLAGS) ./cmd/linkwasmd
	CGO_ENABLED=$(CGO_ENABLED) go install $(BUILD_FLAGS) ./cmd/linkwasmcli

clean:
	rm -rf  build/

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify


########################################
### Testing

test: test-all

test-all: test-unit-all test-integration

test-unit-all: test-unit

test-unit: $(src)
	@go test -mod=readonly -p 4  ./...

test-integration: build
	@go test -mod=readonly -p 4 `go list ./cli_test/...` $(CLI_TEST_BUILD_FLAGS) -v
