package datasource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Go-routine-4595/how-oem-view/model"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

type EventHubConfig struct {
	Connection   string `yaml:"connection"`
	EventHubName string `yaml:"EventHubName"`
}

type FCTSDataModel struct {
	SiteCode    string       `json:"site_code"`
	SensorId    string       `json:"sensor_id"`
	DataSource  string       `json:"data_source"`
	TimeStamp   int64        `json:"time_stamp"`
	Value       string       `json:"value"`
	Uom         string       `json:"uom"`
	Quality     string       `json:"quality"`
	Annotations []Annotation `json:"annotations"`
}

type Annotation struct {
	Properties []map[string]interface{} `json:"properties"`
}

func NewDataSourceEH(connectionString string, container string, eh string) *DataSource {
	return start(connectionString, container, eh)
}

func start(con string, cont string, eh string) *DataSource {
	var (
		ds  DataSource
		sig chan os.Signal
	)
	ds = DataSource{
		dataLock: sync.RWMutex{},
		Keys:     make(map[string]string),
		Alarms:   make([]string, 0),
	}

	sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	_, processorCancel := context.WithCancel(context.TODO())

	// create a consumer client using a connection string to the namespace and the event hub
	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(con, eh, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		panic(err)
	}

	go func() {
		<-sig
		consumerClient.Close(context.TODO())
		processorCancel()
		os.Exit(0)
	}()

	//defer consumerClient.Close(context.TODO())

	eventHubProps, err := consumerClient.GetEventHubProperties(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Partition: ", eventHubProps.PartitionIDs)

	for _, partitionID := range eventHubProps.PartitionIDs {
		//partitionID := "1"
		fmt.Printf("Partition ID: %s\n", partitionID)
		partitionClient, err := consumerClient.NewPartitionClient(partitionID, nil)

		if err != nil {
			panic(err)
		}
		_, err = consumerClient.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
			Prefetch: 1,
		})

		if err != nil {
			panic(err)
		}
		go ds.processEvents(partitionClient, partitionID)
	}

	return &ds
}

func (ds *DataSource) processEvents(partitionClient *azeventhubs.PartitionClient, partionId string) error {
	var (
		k       map[string]string
		alarms  []string
		eventsM model.Events
		fevent  FCTSDataModel
	)
	k = make(map[string]string)

	defer partitionClient.Close(context.TODO())
	for {
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Fatal().Err(err).Msg("failed to receive events")
		}

		log.Debug().Str("partitionId", partionId).Int("processing_events", len(events)).Msg("processing events")

		for _, event := range events {
			log.Trace().Str("event_body", string(event.Body)).Msg("received event")
			err = json.Unmarshal(event.Body, &fevent)
			if err != nil {
				log.Error().Err(errors.Join(err, fmt.Errorf(" processEvents unmarshall FCTSDataModel error data: %q", string(event.Body)))).Msg("dataprocess")
				continue
			}

			err = json.Unmarshal([]byte(fevent.Value), &eventsM)
			if err != nil {
				log.Error().Err(errors.Join(err, fmt.Errorf(" processEvents unmarshall Events error data: %q", fevent.Value))).Msg("dataprocess")
				continue
			}
			if err == nil {

				ds.dataLock.Lock()

				ds.Events = append(ds.Events, eventsM.AssetEvents...)
				k = getKey(eventsM, k)
				for _, v := range k {
					ds.Keys[v] = k[v]
				}
				//ds.Keys = k
				alarms = getAlarms(ds.Events)
				ds.Alarms = append(ds.Alarms, alarms...)
				//ds.Alarms = alarms

				log.Debug().Int("event_size", len(ds.Events)).Msg("dataprocess")

				ds.dataLock.Unlock()
			} else {
				log.Error().Err(errors.Join(err, fmt.Errorf(" stdinReader unmarshall error data: %q", event.Body))).Msg("dataprocess")
			}
		}
	}
}

func closePartitionResources(partitionClient *azeventhubs.ProcessorPartitionClient) {
	defer partitionClient.Close(context.TODO())
}
