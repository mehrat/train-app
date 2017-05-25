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

func (tc *TrainController) PostTrain(train models.Train, r render.Render) {
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")

	train.Id = bson.NewObjectId()
	train.Number = train.Number
	train.Name = train.Name
	train.Schedule = train.Schedule
	session.Insert(train)

	r.JSON(201, train)
}

