package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"net/http"
	"time"
)

type StatePoint struct {
	Timestamp     time.Time
	NewState      string
	PreviousState string
}

type EquipementState struct {
	EquipmentID string
	State       string
	SateList    []StatePoint
}

func NewEquipmentSate(key string) *EquipementState {
	return &EquipementState{
		EquipmentID: key,
		State:       inActive,
		SateList:    make([]StatePoint, 0),
	}
}

func (e *EquipementState) BuildState(events []AssetEvent, key string) {
	for _, elem := range events {
		if elem.AssetName == key {
			if len(e.SateList) == 0 {
				item := StatePoint{
					NewState:      elem.EventStatus,
					PreviousState: e.State,
					Timestamp:     elem.Timestamp,
				}
				e.State = item.NewState
				e.SateList = append(e.SateList, item)
			} else {
				if e.State != elem.EventStatus {
					item := StatePoint{
						NewState:      elem.EventStatus,
						PreviousState: e.State,
						Timestamp:     elem.Timestamp,
					}
					e.State = item.NewState
					e.SateList = append(e.SateList, item)
				}
			}
		}
	}
}

func (e EquipementState) GetState() string {
	return e.State
}

func (e EquipementState) GetSateList() []StatePoint {
	return e.SateList
}

func (e EquipementState) PrintPointSate() {
	fmt.Println()
	fmt.Println("Equipment State:", e.EquipmentID)
	for _, elem := range e.SateList {
		fmt.Println(" time : ", elem.Timestamp, "  state: ", elem.NewState)
	}
}

func isStateCahnged(s1 StatePoint, s2 StatePoint) bool {
	if s1.NewState == s2.NewState {
		return false
	}
	return true
}

func StateGraph(res http.ResponseWriter, req *http.Request) {
	events, _ := openFile("UAS_external_event_queue.jsonl")

	key := req.URL.Query().Get("key")

	if key == "" {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(fmt.Sprintf("key not found: %s", key)))
		return
	}

	equipmentState := NewEquipmentSate(key)

	equipmentState.BuildState(events, key)
	stateList := equipmentState.GetSateList()

	//equipmentState.PrintPointSate()

	p := plot.New()
	p.Title.Text = "Equipment: " + key + " State"
	p.X.Label.Text = "Timestamp"
	p.Y.Label.Text = "State"

	//pts := make(plotter.XYs, len(stateList))
	pts := make(plotter.XYs, 0)

	for i, state := range stateList {
		var t_pts plotter.XY

		if i != 0 {
			//fmt.Println("P0(t): ", stateList[i-1].Timestamp.Unix(), "P1(t): ", state.Timestamp.Unix(), " delta: ", state.Timestamp.Unix()-stateList[i-1].Timestamp.Unix()-1)
			if state.Timestamp.Unix()-stateList[i-1].Timestamp.Unix()-1 > 0 {
				t_pts.X = float64(state.Timestamp.Unix() - 1)
				if stateList[i-1].NewState == inActive {
					t_pts.Y = float64(0)
				} else {
					t_pts.Y = float64(1)
				}
				//fmt.Println("point added: ", int64(t_pts.X), "   ", int64(t_pts.Y))
				pts = append(pts, t_pts)
			}
		}
		t_pts.X = float64(state.Timestamp.Unix())
		if state.NewState == inActive {
			t_pts.Y = float64(0)
		} else {
			t_pts.Y = float64(1)
		}
		//fmt.Println("point added: ", int64(t_pts.X), "   ", int64(t_pts.Y))
		pts = append(pts, t_pts)

	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		processError(err)
		return
	}
	p.Add(line)

	/*
		if err := p.Save(4*vg.Inch, 4*vg.Inch, "state.png"); err != nil {
			processError(err)
			return
		}

		http.ServeFile(res, req, "state.png")
	*/
	res.Header().Set("Content-Type", "image/png")
	iow, err := p.WriterTo(10*vg.Inch, 6*vg.Inch, "png")
	if err != nil {
		fmt.Println(err)
	}
	iow.WriteTo(res)
}
