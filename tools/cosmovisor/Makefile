#!/usr/bin/make -f

COSMOVISOR_TAG=v1.4.0

DEBUG=true

all: test

test: test-single

test-single: test-single-manual-download test-single-auto-download

test-single-manual-download:
	COSMOVISOR_TAG=$(COSMOVISOR_TAG) DEBUG=$(DEBUG) sh test_single.sh

test-single-auto-download:
	COSMOVISOR_TAG=$(COSMOVISOR_TAG) DEBUG=$(DEBUG) DAEMON_ALLOW_DOWNLOAD_BINARIES=true sh test_single.sh

.PHONY: test-single-manual-download test-single-auto-download
