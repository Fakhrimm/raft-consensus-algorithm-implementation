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
4. Load the powershell script to simplify executions:
   ```bash
   . ./scripts.ps1
    ```
5. Run program using the scripts, check the params inside `scripts.ps1`
   - Tidy: runs go mod tidy on controller and node
   - Proto: sets up proto library from /proto to /controller and /node
   - WifiIP: returns the current ip address of your wifi adapter
   - Server: runs a single server on a separate cmd
   - Servers: runs a batch of server and a controller on current terminal (**Note: Will overwrite Hostfile**)
