package timeseriesdb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
			log.Print("Error initialising db: ", db, err)
		}
	}
}

// WriteSingleDatapoint to timeseries db
func WriteSingleDatapoint(req WriteSingleDatapointRequest) {
	url := "http://localhost:8086/write?db=" + req.Db + "&precision=s"
	data := req.Measurement + " " + req.Field + "=" + strconv.FormatFloat(req.Value, 'f', -1, 64) + " " + strconv.FormatInt(time.Now().Unix(), 10)
	res, err := http.Post(url, "application/octet-stream", bytes.NewBufferString(data))
	if err == nil {
		log.Print(res)
	}
}

// CountRecordsSince returns the number of datapoints since a number of seconds ago
func CountRecordsSince(db string, seconds int) int {
	url := "http://localhost:8086/query"
	data := "db=&epoch=ms&q=SELECT count(*) AS \"RecordCount\" FROM \"" + db + "\".\"autogen\".\"Gyroscope\" WHERE time > now() - " + strconv.Itoa(seconds) + "s"
	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	// POST /query?db=&epoch=ms&q=SELECT+count%28%2A%29+FROM+%22simulation%22.%22autogen%22.%22Gyroscope%22+WHERE+time+%3E+now%28%29+-+10m&rp
	// curl -X POST http://localhost:8086/query -d "db=&epoch=ms&q=SELECT+count%28%2A%29+FROM+%22simulation%22.%22autogen%22.%22Gyroscope%22+WHERE+time+%3E+now%28%29+-+10s&rp"
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	results := string(body)
	count := numberOfRecords(results)
	return count
}

func numberOfRecords(results string) int {
	type SeriesRow struct {
		Name    string      `json:"name"`
		Columns []string    `json:"columns"`
		Values  [][]float64 `json:"values"`
	}
	type ResultRow struct {
		StatementID int         `json:"statement_id"`
		Series      []SeriesRow `json:"series"`
	}
	type Results struct {
		Results []ResultRow `json:"results"`
	}
	res := Results{}
	json.Unmarshal([]byte(results), &res)
	count := res.Results[0].Series[0].Values[0][1]
	return int(count)
}

// ClearTestData from test db for integration testing
func ClearTestData() (resp *http.Response, err error) {
	url := "http://localhost:8086/query"
	data := "db=test_simulation&epoch=ms&q=DROP SERIES FROM /.*/"
	resp, err = http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	return resp, err
}
