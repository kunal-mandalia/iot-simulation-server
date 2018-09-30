package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"./timeseriesdb"
)

func SetUp() {
	_, err := timeseriesdb.ClearTestData()
	if err != nil {
		panic(err)
	}
	go Main()
}

func TestSimulationStatusHandler(t *testing.T) {
	SetUp()
	resp, err := http.Get("http://localhost:8080/status")
	if err != nil {
		t.Error("Get /status should not error, got ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	status := string(body)
	if !strings.Contains(status, "Simulation is not running") {
		t.Error("Expected no simulations to be running by default but got, ", status)
	}
}

func TestStartSimulationHandler(t *testing.T) {
	SetUp()
	_, err := http.Get("http://localhost:8080/start?mode=test&frequency=200")
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
	time.Sleep(2 * time.Second)

	// Stop the simulation
	_, err = http.Get("http://localhost:8080/stop")
	if err != nil {
		t.Error("Error stopping simulation ", err)
	}
	count := timeseriesdb.CountRecordsSince("test_simulation", 200)
	if count <= 0 {
		t.Error("Expected to create at least one datapoint but found 0")
	}
}
