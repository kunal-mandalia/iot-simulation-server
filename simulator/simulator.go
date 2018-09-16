package simulator

import (
	"math/rand"
	"sync"
	"time"
)

var instance *Simulation
var once sync.Once

type measurement struct {
	Name  string
	Value int
}

type sensor struct {
	Name        string
	Measurement measurement
}

// Simulation defines the Sensor and state (running or not)
type Simulation struct {
	IsRunning bool
	Sensor    sensor
}

type emitValueChange func(simulation *Simulation)

func (s *sensor) generateMeasurementValue() {
	value := rand.Intn(100)
	s.Measurement.Value = value
}

func next(callback emitValueChange) {
	if instance.IsRunning == true {
		instance.Sensor.generateMeasurementValue()
		callback(instance)
		time.Sleep(2 * time.Second)
		next(callback)
	}
}

// GetInstance http://marcio.io/2015/07/singleton-pattern-in-go/
func GetInstance() *Simulation {
	once.Do(func() {
		instance = &Simulation{}
	})
	return instance
}

// SetProperties to set the simulation name and measurement name
func SetProperties(simulationName string, measurementName string) {
	instance.Sensor.Name = simulationName
	instance.Sensor.Measurement.Name = measurementName
}

// Start simulation
func Start(callback emitValueChange) {
	instance.IsRunning = true
	next(callback)
}

// Stop the simulation
func Stop() {
	instance.IsRunning = false
}
