# RDMA Go Examples

## Go bindings for `rdma/rsocket.h`

### Prerequisites

- RDMA devices (Soft-RoCE is also available).
- RDMA-Core library installed. (Run `apt install rdma-core` on Debian-based distros)

### Components

- `rsocket.go`: C APIs provided by `rdma-core`, serves as wrapper for `ibverbs` which provides socket-like APIs.
- `rsocket.Listener`: Implements `net.Listener` and serves as connection handler for incoming clients.
- `rsocket.Conn`: Implements `net.Conn` and serves as connection handle for RDMA qps.

## Simple Ping-Pong

This example implements a ping-pong session between client and server:

```bash
# Start server
go run ping-pong/server/main.go
# Start client
go run ping-pong/client/main.go
```

## gRPC Example

This example uses `rsocket.Conn` and `rsocket.Listener` as client and server implementation for gRPC:

```bash
# Gen code from .proto file
bash ./grpc/proto/protobuf_gen.sh
# Start server
go run grpc/server/main.go
# Start client
go run grpc/client/main.go
```
