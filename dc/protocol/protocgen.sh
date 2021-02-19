#!/bin/sh
# need to have protoc installed on local machine
cd ./protocol
mkdir ../api
protoc --go_out=../api --go_opt=paths=source_relative \
    --go-grpc_out=../api --go-grpc_opt=paths=source_relative ./config.proto
cd -
