#!/usr/bin/env bash
go mod tidy

echo "Terminating linkcli"
pkill linkcli

echo "Terminating linkd"
pkill -9 linkd

echo "Terminating contract-test"
kill "$(lsof -i tcp:61322 | tail -n 1 | awk '{print $2}')"
