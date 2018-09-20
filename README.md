# iot-simulation-server
A Golang server to simulate IoT data

## Install

* git clone https://github.com/kunal-mandalia/iot-simulation-server.git
* cd ./iot-simulation-server
* go run main.go

## Simulations

* Status: `curl -X GET http://localhost:8080/status`
* Start: `curl -X POST http://localhost:8080/start`
* Stop: `curl -X POST http://localhost:8080/stop`

Simulation options are provided by optional query params. E.g. setting all available options:
`curl -X POST http://localhost:8080/start?sensor=Gyroscope&measurement=Rotation&frequency=1000&mean=180&stddev=10`


## Test

* go test -v ./... (all packages)

## Build

### Docker image
* docker build -t iot-simulation-server .
* docker run --rm -it -p 8080:8080 iot-simulation-server
