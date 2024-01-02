#!/usr/bin/env bash
export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I/usr/local/protoc/include -I. --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false health.proto