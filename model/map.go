package model

import (
	"bytes"
	"errors"
	"fmt"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"image/color"
	"math"
)

const (
	latitude  = "GPSLatitudeofEquipment"
	longitude = "GPSLongitudeofEquipment"
)

type Position struct {
	lat float64
	lng float64
}

func (s Service) DrawMap(val string) ([]byte, error) {
	var (
		pos []Position
		res *bytes.Buffer
	)

	res = new(bytes.Buffer)
	if val == "" {
		fmt.Fprintln(res, "<html><body>")
		fmt.Fprintf(res, "<h1>missing equipement key</h1>")
		return nil, errors.New("No key found")
	}

	pos = s.getPoints(val)
	if pos == nil {
		return nil, errors.New("no value found")
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
	}
	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}
	imgCtx := gg.NewContextForImage(img)
	//res.Header().Set("Content-Type", "image/png")
	if err := imgCtx.EncodePNG(res); err != nil {
		panic(err)
	}
	return res.Bytes(), nil
}

func (s Service) getPoints(val string) []Position {
	var (
		events []AssetEvent
		pos    []Position
		//keys   map[string]string
	)
	events, _ = s.datasource.GetDataSources()

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
