package main

import (
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func pad(message string) string {
	return "\n" + message + "\n\n"
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	status := pad("IoT Simulation Server is running")
	w.Write([]byte(status))
}

func main() {
	http.HandleFunc("/", rootHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("Starting IoT simulation on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
