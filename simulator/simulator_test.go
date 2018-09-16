package simulator

import (
	"log"
	"testing"
	"time"
)

func TestSetProperties(t *testing.T) {
	instance := GetInstance()
	SetProperties("Thermometer", "Temperature")
	if instance.Sensor.Name != "Thermometer" || instance.Sensor.Measurement.Name != "Temperature" {
		t.Error("Expected Sensor: Thermometer, Measurement: Temperature, got ", instance.Sensor.Name, instance.Sensor.Measurement.Name)
	}
}

func TestSimulation(t *testing.T) {
	GetInstance()
	SetProperties("Thermometer", "Temperature")

	simulationCount := 0
	countSimulations := func(simulation *Simulation) {
		log.Print(simulation)
		simulationCount++
	}
	go Start(countSimulations)
	time.Sleep(6 * time.Second)
	if simulationCount != 3 {
		t.Error("Expected 3 simulations to have occured in 6 seconds, got", simulationCount)
	}
}
