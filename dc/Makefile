.PHONY: protoc

protoc:
	@./protocol/protocgen.sh > /dev/null

minio:
	@./start_mio.sh

mongo:
	docker run -p 27017:27017 mongo

build: buildc builda builds buildf builddc

buildc:
	docker build -f Dockerfile.cloud -t cloud .

builda:
	docker build -f Dockerfile.agent -t agent .

builds:
	docker build -t websvr ./webserver

buildd:
	docker build -t demo ./demo

builddc:
	docker build -t cat-ui ./demo-cat
	docker build -t demo-cat ./demo-cat/backend

buildf:
	docker build -t fr ./face_recognition

buildui:
	docker build -t webui ./browser

# demo ml server
dml:
	./demo/backend/start.sh

# demo web server (mongo adaptor) and frontend
dweb:
	node webserver/index.js

cfge:
	curl --header "Content-Type: application/json" --request POST \
		 --data '{"service":"zayn","version":"123", \
		 "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017", \
		 "rules": [{"field":"nose_tip", "op": 1, "operand":"", "sample_rate": 1}]}' \
		 http://localhost:8900/services/zayn

cfgne:
	curl --header "Content-Type: application/json" --request POST \
		 --data '{"service":"zayn","version":"123", \
		 "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017", \
		 "rules": [{"field":"nose_tip", "op": 2, "operand":"", "sample_rate": 1}]}' \
		 http://localhost:8900/services/zayn

cfglt:
	curl --header "Content-Type: application/json" --request POST \
		 --data '{"service":"zayn","version":"123", \
		 "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017", \
		 "rules": [{"field":"off_center", "op": 5, "operand":"0.3", "sample_rate": 1}]}' \
		 http://localhost:8900/services/zayn

cfggt:
	curl --header "Content-Type: application/json" --request POST \
		 --data '{"service":"zayn","version":"123", \
		 "object_storage_path": "minio:9000", "document_storage_path": "mongo:27017", \
		 "rules": [{"field":"off_center", "op": 7, "operand":"0.3", "sample_rate": 1}]}' \
		 http://localhost:8900/services/zayn
