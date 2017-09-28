package mongodb

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

// MongoType all things Mongo
type MongoType struct {
	session *mgo.Session
}

// Initialize mongo connection
func Initialize(mongoConnectionString string) MongoType {
	sess, err := mgo.Dial(mongoConnectionString)
	if err != nil {
		fmt.Println("No mongo connection : (")
		panic(err)
	}

	defer sess.Close()

	s := MongoType{
		session: sess,
	}
	return s

}
