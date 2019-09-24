#!/usr/bin/make -f

GOFILE = $(wildcard *.go)
GONAME = $(shell basename `pwd`)

BUILD_FLAGS = -mod=readonly -tags $(GONAME)

all: build

build:
	go build $(BUILD_FLAGS) -o build/$(GONAME) $(GOFILE)

install:
	go install $(BUILD_FLAGS) $(GOFILE)

clean:
	rm -rf build

.PHONY: all build install clean
