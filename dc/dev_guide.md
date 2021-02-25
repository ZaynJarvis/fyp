# run the cloud

go run cloud/cmd/main.go

# configure the cloud

curl --header "Content-Type: application/json" \
 --request POST \
 --data '{"service":"abc","version":"123", "object_storage_path": "localhost:9000", "document_storage_path": "localhost:27017", "rules": [{"field":"confidence", "op": 5, "operand":"0.9", "sample_rate": 1}]}' \
 http://localhost:8900

# run minIO

./bootstrap.sh

# run mongoDB

docker run -p 27017:27017 mongo

# run agent

go run agent/cmd/main.go 1

# configure example and run SDK

go run sdk/cmd/main.go

# run webserver

node webserver/index.js

# run browser

yarn run dev

# check collected result
