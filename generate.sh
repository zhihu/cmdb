#!/bin/sh
set -e
protoc -I/usr/local/include -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  ./pkg/api/v1/object.proto

protoc -I/usr/local/include -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  ./pkg/api/v1/object.proto

protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:. \
  ./pkg/api/v1/object.proto

echo "package v1"  > ./pkg/api/v1/object.swagger.go  \
    && echo "var Swagger=\`" >> ./pkg/api/v1/object.swagger.go \
	&& cat ./pkg/api/v1/object.swagger.json >> ./pkg/api/v1/object.swagger.go \
	&& echo "\`" >> ./pkg/api/v1/object.swagger.go

cd ./cmd/server/ && wire && cd ../../