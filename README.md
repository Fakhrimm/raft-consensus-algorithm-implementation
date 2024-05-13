## What is this?
An implementation of the Raft consensus algorithm.

## How to set up?
1. Install [Go 1.22.3](https://go.dev/dl/) and update your system environment variables accordingly
2. Install [protoc 27.0-rc2](https://grpc.io/docs/protoc-installation/#install-pre-compiled-binaries-any-os)
3. Install the protocol compiler plugins for Go using the following commands:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```