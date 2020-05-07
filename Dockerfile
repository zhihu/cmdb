FROM golang:1.14.2-alpine3.11 as BUILD

# install all dependencies software
RUN apk update && apk add tzdata git curl unzip grpc protobuf-dev && \
    rm -rf /var/cache/apk/*

# install protoc-gen-go and protoc-gen-grpc-gateway
RUN go get -v -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
	&& go install -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    && go get -v github.com/golang/protobuf/protoc-gen-go

# change protoc-gen-go include postion
RUN mkdir -p /usr/local/include && mv  /usr/include/google /usr/local/include/google

# install gops and wire
RUN go get -u github.com/google/gops \
	&& go get -u github.com/google/wire/cmd/wire \
	&& mv $GOPATH/bin/gops /bin/gops

WORKDIR /app/cmdb
COPY . /app/cmdb
ENV CGO_ENABLED=0
RUN sh ./generate.sh
RUN go build -ldflags="-X 'main.Version=${CMDB_VERSION}'" -o /app/bin/server ./cmd/server

FROM alpine:3.11
COPY --from=BUILD /bin/gops /bin/gops
COPY --from=BUILD /app/bin/server /bin/
ENTRYPOINT ["/bin/server"]