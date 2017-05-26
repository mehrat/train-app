package models

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Train struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Number   string        `form:"Number" json:"Number" bson:"Number"`
	Name     string        `form:"Name" json:"Name"  bson:"Name"`
	Schedule map[string]time.Time    `form:"Schedule" json:"Schedule" bson:"Schedule"`
}
