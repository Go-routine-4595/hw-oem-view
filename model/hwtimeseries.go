package model

import (
	"bytes"
	"github.com/wcharczuk/go-chart/v2"
	"time"
)

type TimeSerie struct {
	XValues []time.Time
	YValues []float64
}

func (s Service) DrawOEMChartSeries(val string) ([]byte, error) {
	/*
	   This is basically the other timeseries example, except we switch to hour intervals and specify a different formatter from default for the xaxis tick labels.
	*/

	var (
		graph  chart.Chart
		series []chart.Series
		res    *bytes.Buffer
	)

	res = new(bytes.Buffer)

	events, keys := s.datasource.GetDataSources()

	if val != "" {
		timeSeries := createTimeSeries(events, val)
		series = append(series, BuildGraphTimeSerie(timeSeries, val))
	} else {
		for i := range keys {
			timeSeries := createTimeSeries(events, i)
			series = append(series, BuildGraphTimeSerie(timeSeries, i))
		}
	}

	graph = chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeHourValueFormatter,
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	graph.Height = 1000
	graph.Width = 1600

	//res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
	return res.Bytes(), nil
}

func BuildGraphTimeSerie(timeSeries TimeSerie, name string) chart.TimeSeries {

	serie := chart.TimeSeries{
		Name:    name,
		XValues: timeSeries.XValues,
		YValues: timeSeries.YValues,
	}

	return serie
}

func (s Service) DrawOEMChartBar(val string) ([]byte, error) {
	/*
	   This is basically the other timeseries example, except we switch to hour intervals and specify a different formatter from default for the xaxis tick labels.
	*/

	var (
		graph  chart.BarChart
		bars   []chart.Value
		title  string
		events []AssetEvent
		keys   map[string]string
		res    *bytes.Buffer
	)

	res = new(bytes.Buffer)

	events, keys = s.datasource.GetDataSources()

	if val != "" {
		title = val
		bars = append(bars, createBars(events, val)...)
	} else {
		title = "all data"
		for i := range keys {
			bars = append(bars, createBars(events, i)...)
		}
	}

	graph = chart.BarChart{
		Title: title,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		UseBaseValue: true,
		BaseValue:    0.0,
		Height:       1000,
		BarWidth:     3000,
		Bars:         bars,
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: barMax(bars),
			},
		},
	}

	graph.Render(chart.PNG, res)
	return res.Bytes(), nil
}

func barMax(slice []chart.Value) float64 {
	m := slice[0].Value
	for _, v := range slice {
		if v.Value > m {
			m = v.Value
		}
	}
	return float64(m)
}

func createTimeSeries(events []AssetEvent, filter string) TimeSerie {
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

func createBars(events []AssetEvent, filter string) []chart.Value {
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
