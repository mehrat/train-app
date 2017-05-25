package models

import (
	"labix.org/v2/mgo"
	"os"
)

func Database() *mgo.Session {
	session, err := mgo.Dial("trainapp_database_1")

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	session.DB(os.Getenv("DB_NAME"))

	return session
}

