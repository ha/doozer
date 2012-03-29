#!/bin/sh
set -e

ver=$(./version.sh)
printf 'package main\n\nconst version = `%s`\n' "$ver" > cmd/doozer/version.go
go get -d -v . ./...
go install -v . ./...
