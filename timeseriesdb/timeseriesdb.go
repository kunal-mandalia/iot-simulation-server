package timeseriesdb

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"
)

// WriteSingleDatapointRequest payload to specify db, measurement, field, and value
type WriteSingleDatapointRequest struct {
	Db          string
	Measurement string
	Field       string
	Value       float64
}

// Initialise dbs
func Initialise() {
	log.Print("Initialising dbs...")
	// Create dbs for simulations if they don't already exist
	dbs := []string{"simulation", "test_simulation"}
	for _, db := range dbs {
		url := "http://localhost:8086/query"
		data := "q=CREATE DATABASE \"" + db + "\""
		_, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
		if err != nil {
			log.Print(err)
		}
	}
}

// WriteSingleDatapoint to timeseries db
func WriteSingleDatapoint(req WriteSingleDatapointRequest) {
	url := "http://localhost:8086/write?db=" + req.Db + "&precision=s"
	ts := time.Now().UTC().String()
	log.Print(ts)
	data := req.Measurement + " " + req.Field + "=" + strconv.FormatFloat(req.Value, 'f', -1, 64) + " " + strconv.FormatInt(time.Now().Unix(), 10)
	res, err := http.Post(url, "application/octet-stream", bytes.NewBufferString(data))
	if err == nil {
		log.Print(res)
	}
}

// GetRecordsBetween two timestamps
func GetRecordsBetween(t1 time.Time, t2 time.Time) {
	t := time.Now()
	log.Print(t)
}
