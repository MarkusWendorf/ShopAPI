#!/usr/bin/env sh

# Only run setup on first start
if [ -e setupCompleted ]; then
    ./shopApi
else
    ./shopApi --setup-only
    touch setupCompleted
    ./shopApi
fi

