#!/bin/bash

# This creates a docker containers each time this runs
# They are deleted immediately after (--rm)

docker run \
    --rm \
    -v ./examples:/app/examples \
    -v ./benchmark-results:/app/benchmark-results \
    grits \
    "$@"