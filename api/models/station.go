package models

import (
	"labix.org/v2/mgo/bson"
)

type Station struct {
	Id   bson.ObjectId `json:"id" bson:"_id"`
	Code string        `form:"Code" json:"Code"`
	Name string        `form:"Name" json:"Name"`
}
