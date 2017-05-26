package models

type StationTrainTime struct {
	Station string `form:"Station" json:"Station"`
	Train   string `form:"Train" json:"Train"`
}
