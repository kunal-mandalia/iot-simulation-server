# iot-simulation-server
A Golang server to simulate IoT data

## Install

* git clone https://github.com/kunal-mandalia/iot-simulation-server.git
* cd ./iot-simulation-server
* go run main.go

## Test

* go test -v ./... (all packages)

## Build

### Docker image
* docker build -t iot-simulation-server .
* docker run --rm -it -p 8080:8080 iot-simulation-server
