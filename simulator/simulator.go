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
	Value float64
}

type sensor struct {
	Name        string
	Frequency   time.Duration // frequency of reading e.g. 2000 ms
	Measurement measurement
}

type simulationAttributes struct {
	DesiredStdDev float64
	DesiredMean   float64
}

// Simulation defines the Sensor and state (running or not)
type Simulation struct {
	IsRunning  bool
	Sensor     sensor
	Attributes simulationAttributes
}

// Properties to initialize simulation with
type Properties struct {
	IsRunning          bool
	SensorName         string
	MeasurementName    string
	FrequencyOfReading time.Duration
	DesiredStdDev      float64
	DesiredMean        float64
}

type emitValueChange func(simulation *Simulation)

func (s *sensor) generateMeasurementValue() {
	desiredStdDev := instance.Attributes.DesiredStdDev
	desiredMean := instance.Attributes.DesiredMean
	sample := rand.NormFloat64()*desiredStdDev + desiredMean
	s.Measurement.Value = sample
}

func next(callback emitValueChange) {
	if instance.IsRunning == true {
		instance.Sensor.generateMeasurementValue()
		callback(instance)
		var frequency time.Duration
		if instance.Sensor.Frequency == 0 {
			frequency = time.Duration(rand.Intn(2000))
		} else {
			frequency = instance.Sensor.Frequency
		}
		time.Sleep(frequency * time.Millisecond)
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
func SetProperties(properties Properties) {
	instance.Sensor.Name = properties.SensorName
	instance.Sensor.Measurement.Name = properties.MeasurementName
	instance.Sensor.Frequency = properties.FrequencyOfReading
	instance.Attributes.DesiredMean = properties.DesiredMean
	instance.Attributes.DesiredStdDev = properties.DesiredStdDev
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
