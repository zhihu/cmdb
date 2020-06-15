#!/bin/sh
set -e
protoc -I/usr/local/include -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  ./pkg/api/v1/*.proto

protoc -I/usr/local/include -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
 ./pkg/api/v1/*.proto

protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:. \
  ./pkg/api/v1/*.proto

go-bindata -pkg v1 -nometadata -prefix pkg/api/v1/ -o ./pkg/api/v1/service.swagger.go ./pkg/api/v1/service.swagger.json

cd ./cmd/server/ && wire && cd ../../