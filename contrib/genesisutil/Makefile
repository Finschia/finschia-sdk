#!/usr/bin/make -f

all: install

DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	go build -o $(DIR)/build/genesis_config $(DIR)/cmd/genesis_config

install:
	go install $(DIR)/cmd/genesis_config

clean:
	rm -rf $(DIR)/build

.PHONY: all build install clean
