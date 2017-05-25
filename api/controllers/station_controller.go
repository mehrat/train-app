package controllers

import (
"github.com/mehrat/train-app/api/models"
"github.com/martini-contrib/render"
"labix.org/v2/mgo"
"labix.org/v2/mgo/bson"
"os"
)

type (
	StationController struct {
		session *mgo.Session
	}
)

func NewStationController(s *mgo.Session) *StationController {
	return &StationController{s}
}

func (tc *StationController) GetAllStations(r render.Render) {
	station := []models.Station{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("station")
	err := session.Find(nil).Limit(100).All(&station)

	if err != nil {
		panic(err)
	}

	r.JSON(200, station)
}

func (pc *StationController) PostStation(Station models.Station, r render.Render) {
	session := pc.session.DB(os.Getenv("DB_NAME")).C("station")

	Station.Id = bson.NewObjectId()
	Station.Code = Station.Code
	Station.Name = Station.Name
	session.Insert(Station)

	r.JSON(201, Station)
}


