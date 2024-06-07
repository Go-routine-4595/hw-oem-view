package model

import (
	"bytes"
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	active   = "Active"
	inActive = "InActive"
)

func (s Service) DrawOEMBarPlot(val string) ([]byte, error) {

	var (
		valuesA plotter.Values
		valuesI plotter.Values
		res     *bytes.Buffer
	)

	res = new(bytes.Buffer)

	events, keys := s.datasource.GetDataSources()
	p := plot.New()

	p.Title.Text = "Bar chart OEM alarm activated/deactivate"
	p.Y.Label.Text = "value"

	for _, k := range keys {
		valuesA = append(valuesA, CreateBarsValues(events, k, active))
		valuesI = append(valuesI, CreateBarsValues(events, k, inActive))
	}
	w := vg.Points(20)
	// active
	barsA, err := plotter.NewBarChart(valuesA, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -w

	// inactive
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

	//res.Header().Set("Content-Type", "image/png")
	iow, err := p.WriterTo(10*vg.Inch, 6*vg.Inch, "png")
	if err != nil {
		fmt.Println(err)
	}
	iow.WriteTo(res)

	return res.Bytes(), nil

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
