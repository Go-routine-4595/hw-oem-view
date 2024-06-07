package model

import (
	"time"
)

type IService interface {
	DrawOEMChartSeries(val string) ([]byte, error)
	DrawOEMChartBar(val string) ([]byte, error)
	DrawOEMBarPlot(val string) ([]byte, error)
	DrawMap(val string) ([]byte, error)
	StateGraph(val string) ([]byte, error)
}

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

type DataSource interface {
	GetDataSources() ([]AssetEvent, map[string]string)
}

type Service struct {
	datasource DataSource
}

func NewService(d DataSource) Service {
	return Service{datasource: d}
}
