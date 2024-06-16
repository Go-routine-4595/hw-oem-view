package datasource

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Go-routine-4595/how-oem-view/model"

	"github.com/rs/zerolog/log"
)

type DataSource struct {
	Events   []model.AssetEvent
	Keys     map[string]string
	Alarms   []string
	dataLock sync.RWMutex
}

func NewDataSource(file string) *DataSource {
	var (
		e []model.AssetEvent
		k map[string]string
	)

	if file == "" {
		return stdinReader()
	}

	e, k = openFile(file)
	return &DataSource{
		Events:   e,
		Keys:     k,
		Alarms:   getAlarms(e),
		dataLock: sync.RWMutex{},
	}
}

func stdinReader() *DataSource {
	var (
		ds  DataSource
		sig chan os.Signal
	)
	ds.dataLock = sync.RWMutex{}

	sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		os.Exit(0)
	}()

	go func() {
		var (
			reader *bufio.Reader
			err    error
			input  string
			k      map[string]string
			alarms []string
		)

		log.Info().Msg("Starting server")
		reader = bufio.NewReader(os.Stdin)
		k = make(map[string]string)

		for {

			var (
				events model.Events
			)

			input, err = reader.ReadString('\n')
			fmt.Println(input)
			if err != nil {
				if err == io.EOF {
					fmt.Println("EOF reached")
					fmt.Println("Stats")
					ds.dataLock.Lock()
					fmt.Println("Events: ", len(ds.Events))
					fmt.Println("Keys: ", len(ds.Keys))
					fmt.Println("Alarms: ", len(ds.Alarms))
					ds.dataLock.Unlock()
					break
				} else {
					log.Err(errors.Join(err, fmt.Errorf(" stdinReader stdin reading"))).Msg("dataprocess")
				}
			}

			err = json.Unmarshal([]byte(input), &events)
			if err == nil {
				fmt.Println("events:", events)
				ds.dataLock.Lock()

				ds.Events = append(ds.Events, events.AssetEvents...)
				k = getKey(events, k)
				alarms = getAlarms(ds.Events)
				ds.Keys = k
				ds.Alarms = alarms

				ds.dataLock.Unlock()
			} else {
				log.Err(errors.Join(err, fmt.Errorf(" stdinReader unmarshall error data: %q", input))).Msg("dataprocess")
			}

		}
	}()

	return &ds
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

// Get Alarm lit
func (d DataSource) GetOEMForAsset() []string {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()
	return d.Alarms
}

// Get Keys public method
func (d DataSource) GetKeys() map[string]string {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()
	return d.Keys
}

// Return the data
func (d DataSource) GetDataSources() ([]model.AssetEvent, map[string]string) {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()
	return d.Events, d.Keys
}

// Return all OEM Alarms
func getAlarms(events []model.AssetEvent) []string {
	var alarmsList []string

	for _, i := range events {
		if !contains(alarmsList, i.EventName) {
			alarmsList = append(alarmsList, i.EventName)
		}
	}

	return alarmsList
}

// Function to check if a string is in a list of strings
func contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
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
