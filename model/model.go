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
	GetKeys() map[string]string
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

type IDataSource interface {
	GetDataSources() ([]AssetEvent, map[string]string)
	GetOEMForAsset() []string
	GetKeys() map[string]string
}

type Service struct {
	datasource IDataSource
}

func NewService(d IDataSource) Service {
	return Service{datasource: d}
}

func (s Service) GetKeys() map[string]string {
	return s.datasource.GetKeys()
}
