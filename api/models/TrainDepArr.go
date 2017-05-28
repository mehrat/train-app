package models

type TrainDepArr struct {
	Train   string `form:"Train" json:"Train"`
	Departure   int `form:"Departure" json:"Departure"`
	Arrival   int `form:"Arrival" json:"Arrival"`
}
