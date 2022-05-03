#!/bin/sh
set -e

go run ./v5config/gen > v5-config.schema.json
go run ./v4config/gen > v4-config.schema.json
