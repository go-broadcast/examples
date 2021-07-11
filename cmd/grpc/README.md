# gRPC example

This example shows how to integrate [broadcast](https://github.com/go-broadcast/broadcast) with [gRPC](https://grpc.io/) server streams.
The main piece of integration is in [server.go](https://github.com/go-broadcast/examples/blob/main/server.go) file 

## Running the example

1. clone the repository
```bash
git clone https://github.com/go-broadcast/examples.git
```
2. start the application
```bash
cd ./examples/client
npm install -y
npx webpack
cd ../cmd/grpc
go run .
```
3. navigate to localhost:5200
