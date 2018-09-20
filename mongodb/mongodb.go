package mongodb

import (
	"github.com/globalsign/mgo"
)

// Dial opens a mongo connection with an URL on secondary.
func Dial(url string) (sess *mgo.Session, err error) {
	if sess, err = mgo.Dial(url); err != nil {
		return sess, err
	}
	sess.SetMode(mgo.Secondary, true)
	return sess, nil
}
