package controllers

import (
	"github.com/mehrat/train-app/api/models"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"sort"
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
	err := session.Find(nil).All(&trains)

	if err != nil {
		panic(err)
	}

	r.JSON(200, trains)
}

func (tc *TrainController) GetTrains(r render.Render, params martini.Params, t models.StationSearch) {

	trains := []models.Train{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")
	err := session.Find(nil).All(&trains)
	tranArr := make(map[string]float64)
	for _, train := range trains {
		var found bool = false
		var t1 int = -1
		var t2 int = -1

		for stn, tim := range train.Schedule {
			if !found && stn == t.From {
				found = true
				t1 = tim
				for stn1, tim1 := range train.Schedule {
					if found && stn1 == t.To {
						t2 = tim1
						break
					}
				}
				break
			}
		}
		if t2 == -1 {
			found = false
		}
		if found {
			duration := t2 - t1
			tranArr[train.Name] = float64(duration) / 100
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
	var arrivalTime int = -1

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
	err := session.Find(nil).All(&trains)

	if err != nil {
		panic(err)
	}
	var isWeekendtrip = "No"
	var from string = params["from"]
	var to string = params["to"]
	for _, train := range trains {
		var found bool = false
		var t1 int = -1
		var t2 int = -1

		for stn, tim := range train.Schedule {
			if !found && stn == from {
				found = true
				t1 = tim

				for stn1, tim1 := range train.Schedule {
					if found && stn1 == to {
						t2 = tim1
						break
					}
				}
				break
			}
		}
		if t2 == -1 {
			found = false
		}
		if found {
			duration := t2 - t1
			if duration <= 180 {
				isWeekendtrip = "Yes"
				break
			}
		}
	}

	r.JSON(200, isWeekendtrip)
}

func (tc *TrainController) SuggestTrains(r render.Render, params martini.Params, bs models.BestSchedule) {
	trains := []models.Train{}
	session := tc.session.DB(os.Getenv("DB_NAME")).C("train")
	err := session.Find(nil).All(&trains)
	trainDep := make(map[string]int)
	trainArr := make(map[string]int)
	diffMap := make(map[string]int)

	for _, train := range trains {
		var found bool = false
		var t1 int = -1
		var t2 int = -1
		for stn, tim := range train.Schedule {
			if !found && stn == bs.From {
				found = true
				t1 = tim
				for stn1, tim1 := range train.Schedule {
					if found && stn1 == bs.To {
						t2 = tim1
						break
					}
				}
				break
			}
		}
		if t2 == -1 {
			found = false
		}
		if found {
			trainDep[train.Name] = t1
			trainArr[train.Name] = t2
			diffMap[train.Name] = bs.TentativeArrival - t2
		}
	}

	result := []models.TrainDepArr{}
	groupedMap := make(map[int][]string)

	//Group trainNames by diffArrivalTime
	for trainName, diffArrivalTime := range diffMap {
		trainList, ok := groupedMap[absVal(diffArrivalTime)]
		if ok {
			trainList = append(trainList, trainName)
			groupedMap[absVal(diffArrivalTime)] = trainList
		} else {
			trainList := []string{trainName}
			groupedMap[absVal(diffArrivalTime)] = trainList
		}
	}

	var keys []int
	for k := range groupedMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		durationMap := make(map[string]int)
		for _, trainName := range groupedMap[key] {
			durationMap[trainName] = trainArr[trainName] - trainDep[trainName]
		}
		durationPairList := models.PairList{}
		durationPairList = getSortedMap(durationMap)

		for _, pair := range durationPairList {
			var trainName string = pair.Key
			trainobj := models.TrainDepArr{trainName, trainDep[trainName], trainArr[trainName] }
			result = append(result, trainobj)
		}
	}

	if err != nil {
		panic(err)
	}
	r.JSON(200, result)
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

func getSortedMap(originalMap map[string]int) models.PairList {
	sortedMap := make(models.PairList, len(originalMap))
	i := 0
	for k, v := range originalMap {
		sortedMap[i] = models.Pair{k, v}
		i++
	}
	sort.Sort(sortedMap)
	return sortedMap
}

func absVal(integ int) int {
	if (integ < 0) {
		return (integ * -1)
	} else {
		return integ
	}
}

