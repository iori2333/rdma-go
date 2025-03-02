#!/bin/bash

protoc --go_out=. --go-grpc_out=. ./examples/proto/message.proto
