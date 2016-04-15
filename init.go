package main

import (
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

var db *mgo.Database
var cUser *mgo.Collection

func init() {
	//Database
	session, err := mgo.Dial(os.Getenv("OPENSHIFT_MONGODB_DB_URL"))
	if err != nil {
		glog.Fatal(err)
	}
    db = session.DB("booltter")
    cUser = db.C("User")
    err = cUser.EnsureIndex(mgo.Index{
        Unique: true,
        Key: []string{"id", "name", "twitter_token"},
    }) 
    if err != nil {
        glog.Fatal(err)
    }

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
}
