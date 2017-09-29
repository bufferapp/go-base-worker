package mongodb

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

// MongoType all things Mongo
type MongoType struct {
	Session *mgo.Session
}

// Initialize mongo connection
func Initialize(mongoConnectionString string) MongoType {
	sess, err := mgo.Dial(mongoConnectionString)
	if err != nil {
		fmt.Println("No mongo connection : (")
		panic(err)
	}
	defer sess.Close()
	sess.SetMode(mgo.Secondary, true)
	s := MongoType{
		Session: sess,
	}
	return s
}
