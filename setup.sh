#!/usr/bin/env bash
mkdir -p /data/db
mongod --bind_ip "$(hostname -i)" &
cd /go/src/ShopAPI
go get ./...
go build
./shopApi --ip "$(hostname -i)" --setup-only
./shopApi --ip "$(hostname -i)"
