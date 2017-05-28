package models

type BestSchedule struct {
	From string `form:"From" json:"From"`
	To   string `form:"To" json:"To"`
	TentativeArrival int `form:"TentativeArrival" json:"TentativeArrival"`
}
