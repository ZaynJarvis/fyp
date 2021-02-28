.PHONY: protoc

protoc:
	@./protocol/protocgen.sh > /dev/null

minio:
	@./start_mio.sh

mongo:
	docker run -p 27017:27017 mongo

cloud:
	go run cloud/cmd/main.go

agent:
	go run agent/cmd/main.go 1

# demo ml server
dml:
	./demo/backend/start.sh

# demo web server (mongo adaptor) and frontend
dweb:
	node webserver/index.js
	cd demo && yarn run dev

