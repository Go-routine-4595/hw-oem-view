package datasource

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Go-routine-4595/how-oem-view/model"
	"os"
)

type DataSource struct {
	Events []model.AssetEvent
	Keys   map[string]string
}

func NewDataSource(file string) *DataSource {
	var (
		e []model.AssetEvent
		k map[string]string
	)
	e, k = openFile(file)
	return &DataSource{
		Events: e,
		Keys:   k,
	}
}

func openFile(s string) ([]model.AssetEvent, map[string]string) {
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

	var AssetEvents []model.AssetEvent
	var keys map[string]string

	AssetEvents = make([]model.AssetEvent, 0)
	keys = make(map[string]string)

	for scanner.Scan() {
		var events model.Events

		if len(scanner.Bytes()) != 0 {
			err = json.Unmarshal(scanner.Bytes(), &events)
			//fmt.Println(string(scanner.Bytes()))
			//DisplayTime(events)
			keys = getKey(events, keys)
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

// Get Keys public method
func (d DataSource) GetKeys() map[string]string {
	return d.Keys
}

// Return the data
func (d DataSource) GetDataSources() ([]model.AssetEvent, map[string]string) {
	return d.Events, d.Keys
}

// All asset name form the list
func getKey(events model.Events, list map[string]string) map[string]string {

	for _, i := range events.AssetEvents {
		if _, ok := list[i.AssetName]; !ok {
			list[i.AssetName] = i.AssetName
		}
	}
	return list
}

// Debug stuff
func DisplayTime(events model.Events) {
	for _, i := range events.AssetEvents {
		fmt.Printf("%-20s %-20s %-20s\n", i.Timestamp, i.EventStatus, i.AssetName)
	}
}

// Debug stuff
func DisplayName(events model.Events) {
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
