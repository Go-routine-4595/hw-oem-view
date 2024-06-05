package main

import (
	"fmt"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"net/http"

	"gonum.org/v1/plot"
)

const (
	active   = "Active"
	inActive = "InActive"
)

func mytest(res http.ResponseWriter, req *http.Request) {
	groupA := plotter.Values{20, 35, 30, 35, 27}
	groupB := plotter.Values{25, 32, 34, 20, 25}
	groupC := plotter.Values{12, 28, 15, 21, 8}

	p := plot.New()

	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w := vg.Points(20)

	barsA, err := plotter.NewBarChart(groupA, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -w

	barsB, err := plotter.NewBarChart(groupB, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)

	barsC, err := plotter.NewBarChart(groupC, w)
	if err != nil {
		panic(err)
	}
	barsC.LineStyle.Width = vg.Length(0)
	barsC.Color = plotutil.Color(2)
	barsC.Offset = w

	p.Add(barsA, barsB, barsC)
	p.Legend.Add("Group A", barsA)
	p.Legend.Add("Group B", barsB)
	p.Legend.Add("Group C", barsC)
	p.Legend.Top = true
	p.NominalX("One", "Two", "Three", "Four", "Five")
	res.Header().Set("Content-Type", "image/png")
	iow, err := p.WriterTo(10*vg.Inch, 6*vg.Inch, "png")
	if err != nil {
		fmt.Println(err)
	}
	iow.WriteTo(res)
}
func drawOEMBarPlot(res http.ResponseWriter, req *http.Request) {

	var (
		valuesA plotter.Values
		valuesI plotter.Values
	)

	events, keys := openFile("UAS_external_event_queue.jsonl")
	p := plot.New()

	p.Title.Text = "Bar chart OEM alarm activated/deactivate"
	p.Y.Label.Text = "value"

	for _, k := range keys {
		valuesA = append(valuesA, CreateBarsValues(events, k, active))
		valuesI = append(valuesI, CreateBarsValues(events, k, inActive))
	}
	w := vg.Points(20)
	// active
	//barsA, err := plotter.NewBarChart(group[i*2], w)
	barsA, err := plotter.NewBarChart(valuesA, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -w

	// inactive
	//barsI, err := plotter.NewBarChart(group[i*2+1], w)
	barsI, err := plotter.NewBarChart(valuesI, w)
	if err != nil {
		panic(err)
	}
	barsI.LineStyle.Width = vg.Length(0)
	barsI.Color = plotutil.Color(1)

	p.Add(barsA, barsI)
	p.Legend.Add("I", barsI)
	p.Legend.Add("A", barsA)
	p.Legend.Top = true
	//p.NominalX("BAGT134", "BAGT215", "BAGT126", "BAGT209", "BAGT220", "BAGT219", "BAGT201", "BAGT211", "BAGS012", "BAGT204", "BAGT133", "BAGT205", "BAGT137", "BAGT218", "BAGT203", "BAGS022", "BAGT207")

	res.Header().Set("Content-Type", "image/png")
	iow, err := p.WriterTo(10*vg.Inch, 6*vg.Inch, "png")
	if err != nil {
		fmt.Println(err)
	}
	iow.WriteTo(res)

}

func CreateBarsValues(events []AssetEvent, filter string, status string) float64 {
	var (
		barValue float64
	)

	for _, event := range events {
		if event.AssetName == filter {
			if event.EventStatus == status {
				barValue += 1
			}
		}
	}
	return barValue
}
