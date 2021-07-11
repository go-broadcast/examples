# Redis example

This example shows how to scale out [broadcast](https://github.com/go-broadcast/broadcast) using [go-broadcast/redis](https://github.com/go-broadcast/redis). 

## Running the example

1. clone the repository
```bash
git clone https://github.com/go-broadcast/examples.git
```
2. start a Redis instance
```bash
docker run --name test-redis -d -p 6379:6379 redis
```
3. start one instance of the application
```bash
cd ./examples/cmd/redis
go run . -port 5200
```
4. start a second instance of the application
```bash
go run . -port 5300
```
5. open a browser and navigate to localhost:5200
