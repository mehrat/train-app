package main

import (
	"github.com/mehrat/train-app/api/controllers"
	"github.com/mehrat/train-app/api/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func main() {

	m := martini.Classic()
	m.Map(models.Database())
	m.Use(render.Renderer())

	tc := controllers.NewTrainController(models.Database())
	sc := controllers.NewStationController(models.Database())

	m.Get("/trains/all", binding.Bind(models.Train{}), tc.GetAllTrains)
	m.Get("/stations/all", binding.Bind(models.Station{}), sc.GetAllStations)
	m.Post("/train/add", binding.Bind(models.Train{}), tc.PostTrain)
	m.Post("/station/add", binding.Bind(models.Station{}), sc.PostStation)

	m.Get("/train/:Number/schedule", tc.GetTrainSchedule)
	m.Get("/trains", tc.GetTrains)

	m.Run()
}