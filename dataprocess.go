package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"os"
	"time"
)

type AssociatedData struct {
	Name      string    `json:"Name"`
	Quality   string    `json:"Quality"`
	Timestamp time.Time `json:"Timestamp"`
	Value     float64   `json:"Value"`
}

type AssetEvent struct {
	AssetName      string           `json:"AssetName"`
	AssociatedData []AssociatedData `json:"AssociatedData"`
	EventName      string           `json:"EventName"`
	EventStatus    string           `json:"EventStatus"`
	Timestamp      time.Time        `json:"Timestamp"`
	CreatedUser    string           `json:"CreatedUser"`
}

type Events struct {
	AssetEvents []AssetEvent `json:"Events"`
}

type TimeSerie struct {
	XValues []time.Time
	YValues []float64
}

func openFile(s string) ([]AssetEvent, map[string]string) {
	const maxCapacity = 512 * 1024

	f, err := os.Open(s)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	// Create a new Scanner for the file
	scanner := bufio.NewScanner(f)
	// use a custom function for the line delimiter (here the file is delimited with CR CR (0xOD 0xOD)
	scanner.Split(splitCRCR)
	// increase the buffer size to 2Mb
	buf := []byte{}
	scanner.Buffer(buf, 2048*1024)

	var AssetEvents []AssetEvent
	var keys map[string]string

	AssetEvents = make([]AssetEvent, 0)
	keys = make(map[string]string)

	for scanner.Scan() {
		var events Events

		if len(scanner.Bytes()) != 0 {
			err = json.Unmarshal(scanner.Bytes(), &events)
			//fmt.Println(string(scanner.Bytes()))
			//DisplayTime(events)
			keys = GetKey(events, keys)
			//DisplayName(events)
			if err != nil {
				processError(err)
			}
			AssetEvents = append(AssetEvents, events.AssetEvents...)
		}
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}

	return AssetEvents, keys
}

// All asset name form the list
func GetKey(events Events, list map[string]string) map[string]string {

	for _, i := range events.AssetEvents {
		if _, ok := list[i.AssetName]; !ok {
			list[i.AssetName] = i.AssetName
		}
	}
	return list
}

// Debug stuff
func DisplayTime(events Events) {
	for _, i := range events.AssetEvents {
		fmt.Printf("%-20s %-20s %-20s\n", i.Timestamp, i.EventStatus, i.AssetName)
	}
}

// Debug stuff
func DisplayName(events Events) {
	for _, i := range events.AssetEvents {
		fmt.Printf("%-20s \n", i.AssetName)
	}
}

// splitCRCR is a bufio.SplitFunc for splitting data at \r\r delimiters.
func splitCRCR(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte("\r\r")); i >= 0 {
		// We have a full \r\r-terminated line.
		return i + 2, data[0:i], nil
	}

	// If at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func CreateTimeSeries(events []AssetEvent, filter string) TimeSerie {
	var timeSeries TimeSerie
	var count float64

	for _, event := range events {
		if event.AssetName == filter {
			timeSeries.XValues = append(timeSeries.XValues, event.Timestamp)
			if event.EventStatus == "InActive" {
				count += 1
				timeSeries.YValues = append(timeSeries.YValues, count)
			} else {
				count -= 1
				timeSeries.YValues = append(timeSeries.YValues, count)
			}
		}

	}
	return timeSeries
}

func CreateBars(events []AssetEvent, filter string) []chart.Value {
	var (
		barActive   chart.Value
		barInactive chart.Value
		bars        []chart.Value
	)

	barActive.Label = filter + "A"
	barInactive.Label = filter + "I"
	for _, event := range events {
		if event.AssetName == filter {
			if event.EventStatus == "InActive" {
				barActive.Value += 1
			} else {
				barInactive.Value += 1
			}
		}
	}
	bars = append(bars, barActive)
	bars = append(bars, barInactive)

	return bars
}
