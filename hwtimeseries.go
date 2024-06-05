package main

import (
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"net/http"
)

func drawOEMChartSeries(res http.ResponseWriter, req *http.Request) {
	/*
	   This is basically the other timeseries example, except we switch to hour intervals and specify a different formatter from default for the xaxis tick labels.
	*/

	var graph chart.Chart
	var series []chart.Series

	// Query parameters are returned as a map[string][]string
	params := req.URL.Query()
	// To get the value of a specific key:
	// (assuming the parameter exists)
	val := params.Get("key")

	fmt.Printf("Parsed param: %s \n", val)

	events, keys := openFile("UAS_external_event_queue.jsonl")
	if val != "" {
		timeSeries := CreateTimeSeries(events, val)
		series = append(series, BuildGraphTimeSerie(timeSeries, val))
	} else {
		for i := range keys {
			timeSeries := CreateTimeSeries(events, i)
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

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func BuildGraphTimeSerie(timeSeries TimeSerie, name string) chart.TimeSeries {

	serie := chart.TimeSeries{
		Name:    name,
		XValues: timeSeries.XValues,
		YValues: timeSeries.YValues,
	}

	return serie
}

func drawOEMChartBar(res http.ResponseWriter, req *http.Request) {
	/*
	   This is basically the other timeseries example, except we switch to hour intervals and specify a different formatter from default for the xaxis tick labels.
	*/

	var (
		graph  chart.BarChart
		bars   []chart.Value
		title  string
		events []AssetEvent
		keys   map[string]string
	)

	// Query parameters are returned as a map[string][]string
	params := req.URL.Query()
	// To get the value of a specific key:
	// (assuming the parameter exists)
	val := params.Get("key")

	fmt.Printf("Parsed param: %s \n", val)

	events, keys = openFile("UAS_external_event_queue.jsonl")
	if val != "" {
		title = val
		bars = append(bars, CreateBars(events, val)...)
	} else {
		title = "all data"
		for i := range keys {
			bars = append(bars, CreateBars(events, i)...)
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

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
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
