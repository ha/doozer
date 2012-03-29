#!/bin/sh
set -e

go clean . ./...
rm -rf cmd/doozer/version.go 2>&1 >/dev/null
