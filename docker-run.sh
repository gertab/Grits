#!/bin/bash

# This creates a docker containers each time this runs
# Clean containers using: docker rm $(docker ps -a -q  --filter ancestor=grits)

docker run \
    -v ./examples:/app/examples \
    -v ./benchmark-results:/app/benchmark-results \
    grits \
    "$@"