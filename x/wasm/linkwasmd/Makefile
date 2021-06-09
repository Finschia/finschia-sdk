#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# ascribe tag only if on a release/ branch, otherwise pick branch name and concatenate commit hash
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
ifeq (, $(findstring release/,$(BRANCH)))
  VERSION = $(BRANCH)-$(COMMIT)
endif

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/line/lfb-sdk | sed  's/ /\@/g')
OST_VERSION := $(shell go list -m github.com/line/ostracon | sed 's:.* ::') # grab everything after the space in "github.com/line/ostracon v0.34.7"
DOCKER := $(shell which docker)
BUILDDIR ?= $(CURDIR)/build
TEST_DOCKER_REPO=jackzampolin/linktest

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (rocksdb,$(findstring rocksdb,$(LBM_BUILD_OPTIONS)))
  build_tags += gcc rocksdb
else ifeq (cleveldb,$(findstring cleveldb,$(LBM_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
else
  build_tags += goleveldb
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/line/lfb-sdk/version.Name=link \
		  -X github.com/line/lfb-sdk/version.AppName=linkwasmd \
		  -X github.com/line/lfb-sdk/version.Version=$(VERSION) \
		  -X github.com/line/lfb-sdk/version.Commit=$(COMMIT) \
		  -X github.com/line/lfb-sdk/x/wasm/linkwasmd/app.ProposalsEnabled=true \
		  -X "github.com/line/lfb-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/line/ostracon/version.TMCoreSemVer=$(OST_VERSION)

ifeq (cleveldb,$(findstring cleveldb,$(LBM_BUILD_OPTIONS)))
  ldflags += -X github.com/line/lfb-sdk/types.DBBackend=cleveldb
endif
ifeq (,$(findstring nostrip,$(LBM_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(LBM_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

#$(info $$BUILD_FLAGS is [$(BUILD_FLAGS)])

# The below include contains the tools target.
include ../../../contrib/devtools/Makefile

###############################################################################
###                              Build                                      ###
###############################################################################

all: install lint test

build: BUILD_ARGS=-o $(BUILDDIR)/

build: go.sum $(BUILDDIR)/
	go build -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./cmd/linkwasmd

install: go.sum $(BUILDDIR)/
	go install $(BUILD_FLAGS) $(BUILD_ARGS) ./cmd/linkwasmd

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-reproducible: go.sum
	$(DOCKER) rm latest-build || true
	$(DOCKER) run --volume=$(CURDIR):/sources:ro \
        --env TARGET_PLATFORMS='linux/amd64 darwin/amd64 linux/arm64 windows/amd64' \
        --env APP=lfb \
        --env VERSION=$(VERSION) \
        --env COMMIT=$(COMMIT) \
        --env LEDGER_ENABLED=$(LEDGER_ENABLED) \
        --name latest-build cosmossdk/rbuilder:latest
	$(DOCKER) cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

clean:
	rm -rf $(BUILDDIR)/ artifacts/

distclean: clean
	rm -rf vendor/


########################################
### Testing


test: test-unit
test-all: test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 5 TestAppImportExport

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 10 TestFullAppSimulation
  

.PHONY: all build-linux install \
	go-mod-cache clean build \
	test test-all test-cover test-unit test-race \
