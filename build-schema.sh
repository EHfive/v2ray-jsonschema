#!/bin/sh
set -e

go run ./v5config/gen > v5-config-schema.json
