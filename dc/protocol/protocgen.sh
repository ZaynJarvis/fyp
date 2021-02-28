#!/bin/sh
# need to have protoc installed on local machine
cd ./protocol
rm -rf ../api
mkdir -p ../api
protoc --go_out=../api --go_opt=paths=source_relative \
    --go-grpc_out=../api --go-grpc_opt=paths=source_relative ./collection.proto
python3 -m grpc_tools.protoc -I. --python_out=../demo/backend \
    --grpc_python_out=../demo/backend ./collection.proto
cd -
