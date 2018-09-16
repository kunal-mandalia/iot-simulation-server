package main

import (
	"log"
	"net/http"
	"os"

	"./simulator"
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
		sensorName := "Gyroscope"
		measurementName := "Rotation"
		logValueChange := func(simulation *simulator.Simulation) {
			log.Print("Simulated value changed:")
			log.Print(simulation)
		}

		if r.URL.Query().Get("sensor") != "" {
			sensorName = r.URL.Query().Get("sensor")
		}
		if r.URL.Query().Get("measurement") != "" {
			measurementName = r.URL.Query().Get("measurement")
		}
		simulator.SetProperties(sensorName, measurementName)
		go simulator.Start(logValueChange)

		message := "Started " + sensorName + " simulation for measurement " + measurementName
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

func statusHandler(w http.ResponseWriter, r *http.Request) {
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
	simulator.GetInstance()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/start", startSimulationHandler)
	http.HandleFunc("/stop", stopSimulationHandler)
	http.HandleFunc("/status", statusHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("Starting IoT simulation on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
