package controllers

import (
	"github.com/mehrat/train-app/api/models"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

type (
	TrainController struct {
		session *mgo.Session
	}
)

func NewTrainController(s *mgo.Session) *TrainController {
	return &TrainController{s}
}

func (tc *TrainController) GetTrainSchedule(r render.Render, params martini.Params) {
	result := models.Train{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")
	err := session.Find(bson.M{"Number":params["Number"]}).Select(bson.M{"Schedule": ""}).One(&result)

	if err != nil {
		panic(err)
	}
	r.JSON(200, result.Schedule)
}

func (tc *TrainController) GetAllTrains(r render.Render) {
	trains := []models.Train{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")
	err := session.Find(nil).Limit(100).All(&trains)

	if err != nil {
		panic(err)
	}

	r.JSON(200, trains)
}

func (tc *TrainController) GetTrains(r render.Render, params martini.Params, t models.StationSearch) {

	trains := []models.Train{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")
	err := session.Find(nil).Limit(100).All(&trains)
	tranArr := make(map[string]int64)
	for _, train := range trains {
		var found bool = false
		var t1 int64
		var t2 int64

		for stn, tim := range train.Schedule {
			if !found && stn == t.From {
				found = true
				t1 = tim
			}
			if found && stn == t.To {
				t2 = tim
			}
		}
		if found {
			duration := t2 - t1
			tranArr[train.Name] = duration / 100
		}
	}
	if err != nil {
		panic(err)
	}
	r.JSON(200, tranArr)
}

func (tc *TrainController) GetTrainReachTime(r render.Render, params martini.Params, st models.StationTrainTime) {

	result := models.Train{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")
	err := session.Find(bson.M{"Number":st.Train}).Select(bson.M{"Schedule": ""}).One(&result)
	if err != nil {
		panic(err)
	}
	var arrivalTime int64 = -1

	for stn, tim := range result.Schedule {
		if stn == st.Station {
			arrivalTime = tim
		}
	}

	if arrivalTime == -1 {
		r.Error(404)
		return
	}
	r.JSON(200, arrivalTime)
}

func (tc *TrainController) IsWeekendTrip(r render.Render, params martini.Params) {
	trains := []models.Train{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")
	err := session.Find(nil).Limit(100).All(&trains)

	if err != nil {
		panic(err)
	}
	var isWeekendtrip = "No"
	var from string = params["from"]
	var to string = params["to"]
	for _, train := range trains {
		var found bool = false
		var t1 int64
		var t2 int64

		for stn, tim := range train.Schedule {
			if !found && stn == from {
				found = true
				t1 = tim
			}
			if found && stn == to {
				t2 = tim
			}
		}
		if found {
			duration := t2 - t1
			if (duration / 100) <= 180 {
				isWeekendtrip = "Yes"
				break
			}
		}
	}

	r.JSON(200, isWeekendtrip)
}

func (tc *TrainController) PostTrain(train models.Train, r render.Render) {
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")

	train.Id = bson.NewObjectId()
	train.Number = train.Number
	train.Name = train.Name
	train.Schedule = train.Schedule
	session.Insert(train)

	r.JSON(201, train)
}

