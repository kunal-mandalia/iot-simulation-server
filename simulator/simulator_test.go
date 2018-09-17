package simulator

import (
	"log"
	"testing"
	"time"
)

func TestSetProperties(t *testing.T) {
	instance := GetInstance()
	sensorName := "Thermometer"
	measurementName := "Temperature"
	SetProperties(Properties{SensorName: sensorName, MeasurementName: measurementName, FrequencyOfReading: 0, DesiredMean: 180, DesiredStdDev: 10})

	if instance.Sensor.Name != "Thermometer" || instance.Sensor.Measurement.Name != "Temperature" {
		t.Error("Expected Sensor: Thermometer, Measurement: Temperature, got ", instance.Sensor.Name, instance.Sensor.Measurement.Name)
	}
}

func TestSimulation(t *testing.T) {
	GetInstance()
	SetProperties(Properties{SensorName: "Thermometer", MeasurementName: "Temperature", FrequencyOfReading: 2000, DesiredMean: 18, DesiredStdDev: 2})

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
