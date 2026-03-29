#!/bin/sh
set -eu

src_dir=$(dirname "${0}")/..

swag init \
    --generalInfo "${src_dir}/internal/routes/router.go" \
    --output "${src_dir}/internal/handlers/docs" \
    --outputTypes go,json \
    --parseDependency --parseInternal --parseDepth 1