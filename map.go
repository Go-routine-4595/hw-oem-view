package main

import (
	"fmt"
	"image/color"
	"math"
	"net/http"
	"net/url"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

const (
	latitude  = "GPSLatitudeofEquipment"
	longitude = "GPSLongitudeofEquipment"
)

type Position struct {
	lat float64
	lng float64
}

func mymap(res http.ResponseWriter, req *http.Request) {
	var (
		pos    []Position
		params url.Values
		val    string
	)
	// Query parameters are returned as a map[string][]string
	params = req.URL.Query()
	// To get the value of a specific key:
	// (assuming the parameter exists)
	val = params.Get("key")

	if val == "" {
		val = "BAGT201"
	}

	pos = getPoints(val)
	if pos == nil {
		res.Header().Set("Content-Type", "text/html")
		// Write the HTML opening tags and the title
		fmt.Fprintln(res, "<html><body>")
		fmt.Fprintf(res, "<h1>No value for: %s </h1>", val)
		http.Error(res, "No value found", http.StatusNotFound)
		return
	}

	ctx := sm.NewContext()
	ctx.SetSize(1600, 1200)
	ctx.SetZoom(14)
	for _, i := range pos {
		ctx.AddObject(
			sm.NewMarker(
				s2.LatLngFromDegrees(i.lat, i.lng),
				color.RGBA{0xff, 0, 0, 0xff},
				16.0,
			),
		)
		// Debug
		//fmt.Printf("%f/%f\n", i.lat, i.lng)
	}
	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}
	imgCtx := gg.NewContextForImage(img)
	res.Header().Set("Content-Type", "image/png")
	if err := imgCtx.EncodePNG(res); err != nil {
		panic(err)
	}
}

func getPoints(val string) []Position {
	var (
		events []AssetEvent
		pos    []Position
		//keys   map[string]string
	)
	events, _ = openFile("UAS_external_event_queue.jsonl")

	for _, e := range events {
		var (
			asData []AssociatedData
			point  Position
		)

		if e.AssetName == val {
			asData = e.AssociatedData
			for _, d := range asData {

				if d.Name == latitude {
					point.lat = roundTo(d.Value, 5)
				}
				if d.Name == longitude {
					point.lng = roundTo(d.Value, 5)
				}
			}
			pos = append(pos, point)
		}
	}
	return pos
}

// roundTo rounds num to places decimal places
func roundTo(num float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(num*shift) / shift
}
