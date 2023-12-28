#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/foundation/expected_keepers.go -package testutil -destination x/foundation/testutil/expected_keepers_mocks.go
