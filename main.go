package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"./simulator"
	"./timeseriesdb"
)

const defaultPort = "8080"

func formatMessage(message string) string {
	return "\n    " + message + "\n\n"
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	status := formatMessage("IoT Simulation Server is running")
	w.Write([]byte(status))
}

func startSimulationHandler(w http.ResponseWriter, r *http.Request) {
	simulation := simulator.GetInstance()
	if simulation.IsRunning == true {
		message := "A simulation is already running. Go to /status for details or /stop to stop the current simulation"
		w.WriteHeader(http.StatusLocked)
		w.Write([]byte(formatMessage(message)))
	} else {
		simulationProperties := simulator.Properties{SensorName: "Gyroscope", MeasurementName: "Rotation", FrequencyOfReading: 0, DesiredMean: 180, DesiredStdDev: 10}

		db := "simulation"
		if r.URL.Query().Get("mode") == "test" {
			db = "test_simulation"
		}
		if r.URL.Query().Get("sensor") != "" {
			simulationProperties.SensorName = r.URL.Query().Get("sensor")
		}
		if r.URL.Query().Get("measurement") != "" {
			simulationProperties.MeasurementName = r.URL.Query().Get("measurement")
		}
		if r.URL.Query().Get("frequency") != "" {
			frequency, err := strconv.ParseInt(r.URL.Query().Get("frequency"), 10, 32)
			if err != nil {
				panic(err)
			}
			simulationProperties.FrequencyOfReading = time.Duration(frequency)
		}
		if r.URL.Query().Get("mean") != "" {
			mean, err := strconv.ParseFloat(r.URL.Query().Get("mean"), 64)
			if err != nil {
				panic(err)
			}
			simulationProperties.DesiredMean = mean
		}
		if r.URL.Query().Get("stddev") != "" {
			stddev, err := strconv.ParseFloat(r.URL.Query().Get("stddev"), 64)
			if err != nil {
				panic(err)
			}
			simulationProperties.DesiredStdDev = stddev
		}

		simulator.SetProperties(simulationProperties)
		writeToDb := func(simulation *simulator.Simulation) {
			dp := timeseriesdb.WriteSingleDatapointRequest{Db: db, Measurement: simulation.Sensor.Name, Field: simulation.Sensor.Measurement.Name, Value: simulation.Sensor.Measurement.Value}
			timeseriesdb.WriteSingleDatapoint(dp)
		}

		go simulator.Start(writeToDb)

		message := "Started " + simulationProperties.SensorName + " simulation for measurement " + simulationProperties.MeasurementName
		log.Print(message)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formatMessage(message)))
	}
}

func stopSimulationHandler(w http.ResponseWriter, r *http.Request) {
	simulator.Stop()
	message := "Stopped simulation"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(formatMessage(message)))
}

func simulationStatusHandler(w http.ResponseWriter, r *http.Request) {
	simulation := simulator.GetInstance()
	var message string

	if simulation.IsRunning == true {
		message = "Simulation is running for Sensor: " + simulation.Sensor.Name + " on Measurement: " + simulation.Sensor.Measurement.Name
	} else {
		message = "Simulation is not running"
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(formatMessage(message)))
}

func main() {
	go timeseriesdb.Initialise()
	simulator.GetInstance()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status", simulationStatusHandler)
	// TODO change simulation start/stop to respond to POST
	// requests only
	http.HandleFunc("/start", startSimulationHandler)
	http.HandleFunc("/stop", stopSimulationHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("Starting IoT simulation on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Main is a duplicate of main for running system tests
func Main() {
	main()
}
