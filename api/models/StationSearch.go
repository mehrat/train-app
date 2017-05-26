package models

type StationSearch struct {
	From string `form:"From" json:"From"`
	To   string `form:"To" json:"To"`
}
