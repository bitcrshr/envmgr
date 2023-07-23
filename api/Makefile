SHELL := /bin/bash

.PHONY: proto
proto:
	cd proto/src && \
	protoc --go_out=../compiled/go --go_opt=paths=source_relative \
		--go-grpc_out=../compiled/go --go-grpc_opt=paths=source_relative \
		--ts_out ../compiled/ts \
		environment_manager.proto

.PHONY: run
run:
	go run main.go