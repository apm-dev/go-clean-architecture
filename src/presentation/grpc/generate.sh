#!/bin/bash

#protoc pbs/*.proto --go_out=plugins=grpc:.

protoc \
  -I . \
  -I $GOPATH/src/ \
  -I $GOPATH/src/github.com/google/protobuf/src/ \
  --proto_path=. \
  --go_out=plugins=grpc:. \
  --govalidators_out=. \
  pbs/*.proto
