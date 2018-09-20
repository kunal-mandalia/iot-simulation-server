package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSimulationStatusHandler(t *testing.T) {
	go Main()
	resp, err := http.Get("http://localhost:8080/status")
	if err != nil {
		t.Error("Get /status should not error, got ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Print(string(body))
	status := string(body)
	if !strings.Contains(status, "Simulation is not running") {
		t.Error("Expected no simulations to be running by default but got, ", status)
	}
}

func TestStartSimulationHandler(t *testing.T) {
	// Start server and begin simulating
	// Q: Can go routines be run for a particular lifespan e.g. 30seconds?
	// as the thread will continue in the background after the function returns
	go Main()
	_, err := http.Get("http://localhost:8080/start?mode=test&freqency=1")
	if err != nil {
		t.Error("/start should not error, got ", err)
	}
	// Expect simulation to have started
	resp, err := http.Get("http://localhost:8080/status")
	if err != nil {
		t.Error("Get /status should not error, got ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	status := string(body)
	if strings.Contains(status, "Simulation is not running") {
		t.Error("Expected simulation to be running but got, ", status)
	}

	time.Sleep(10 * time.Second)

	// Stop the simulation
	_, err = http.Get("http://localhost:8080/stop")
	if err != nil {
		t.Error("Error stopping simulation ", err)
	}
	// TODO assert against count of data points between start and stop
}
