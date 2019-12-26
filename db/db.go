package db

import (
	"gopkg.in/mgo.v2"
)

const (
	DBNAME = "marathon_test"
	URI    = "mongodb://localhost:27017"
)

func GetDB() (*mgo.Database, *mgo.Session, error) {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	db := session.DB("marathon_test")
	return db, session.Clone(), nil
}
