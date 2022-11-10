#!/usr/bin/make -f

COSMOVISOR_TAG=v1.4.0

DEBUG=true

install:
	go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@$(COSMOVISOR_TAG)

install-simd:
	go install ./../../...

.PHONY: install install-simd

test: test-single

test-single: test-single-manual-download test-single-auto-download

test-single-manual-download: install install-simd
	COSMOVISOR_TAG=$(COSMOVISOR_TAG) DEBUG=$(DEBUG) sh test_single.sh

test-single-auto-download: install install-simd
	COSMOVISOR_TAG=$(COSMOVISOR_TAG) DEBUG=$(DEBUG) DAEMON_ALLOW_DOWNLOAD_BINARIES=true sh test_single.sh
