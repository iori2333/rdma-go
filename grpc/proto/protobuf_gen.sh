#!/bin/bash

protoc --go_out=. --go-grpc_out=. ./grpc/proto/message.proto
