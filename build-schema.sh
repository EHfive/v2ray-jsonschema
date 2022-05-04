#!/bin/sh
set -e

go mod vendor
go run ./v5config/gen > v5-config.schema.json
go run ./v4config/gen > v4-config.schema.json

rm -rf vendor
