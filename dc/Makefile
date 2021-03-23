.PHONY: protoc

protoc:
	@./protocol/protocgen.sh > /dev/null

minio:
	@./start_mio.sh

mongo:
	docker run -p 27017:27017 mongo

build: buildc builda builds buildd

buildc:
	docker build -f Dockerfile.cloud -t cloud .

builda:
	docker build -f Dockerfile.agent -t agent .

builds:
	docker build -t websvr ./webserver

buildd:
	docker build -t demo ./demo

buildui:
	docker build -t webui ./browser

# demo ml server
dml:
	./demo/backend/start.sh

# demo web server (mongo adaptor) and frontend
dweb:
	node webserver/index.js

cfge:
	curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"service":"abc","version":"123", "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017",  "rules": [{"field":"nose_tip", "op": 1, "operand":"", "sample_rate": 1}]}' \
  http://localhost:8900


cfgne:
	curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"service":"abc","version":"123", "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017",  "rules": [{"field":"nose_tip", "op": 2, "operand":"", "sample_rate": 1}]}' \
  http://localhost:8900


cfglt:
	curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"service":"abc","version":"123", "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017",  "rules": [{"field":"off_center", "op": 5, "operand":"0.3", "sample_rate": 1}]}' \
  http://localhost:8900

cfggt:
	curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"service":"abc","version":"123", "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017",  "rules": [{"field":"off_center", "op": 7, "operand":"0.3", "sample_rate": 1}]}' \
  http://localhost:8900

