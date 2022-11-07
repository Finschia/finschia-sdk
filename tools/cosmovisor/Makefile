#!/usr/bin/make -f

COSMOVISOR_TAG=v1.4.0

all: test-single

test-single:
	COSMOVISOR_TAG=$(COSMOVISOR_TAG) sh test_single.sh

.PHONY: all test-single
