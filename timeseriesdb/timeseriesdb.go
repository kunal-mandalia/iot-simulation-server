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

// TODO create simulation, test_simulation db if they don't exist
