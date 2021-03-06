package main

import (
	"time"
    "errors"

	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2/bson"
)

var slateTime = time.Hour

type mUser struct {
	OID           bson.ObjectId `bson:"_id"`
	ID            int64         `bson:"id"`
	Name          string        `bson:"name"`
	TwitterToken  string        `bson:"twitter_token"`
	TwitterSecret string        `bson:"twitter_secret"`
	UpdatedAt     time.Time     `bson:"updated_at"`
	CreatedAt     time.Time     `bson:"created_at"`
}

func findOrCreateUserFromToken(token, secret string) (*mUser, error) {
	var user *mUser
	cUser.Find(bson.M{"twitter_token": token}).One(user)
	if user == nil {
		return createUserFromToken(token, secret)
	}
	if user.TwitterSecret != secret {
        return nil, errors.New("Glitched twitter token.")
	}
    return user, nil
}

func createUserFromToken(token, secret string) (*mUser, error) {
    twitter := anaconda.NewTwitterApi(token, secret)
    self, err := twitter.GetSelf(nil)
    if err != nil {
        return nil, err
    }
    now := time.Now()
    user := &mUser{
        OID: bson.NewObjectId(),
        ID: self.Id,
        Name: self.ScreenName,
        TwitterToken: token,
        TwitterSecret: secret,
        UpdatedAt: now,
        CreatedAt: now,
    }
    return user, cUser.Insert(user)
}

func (u *mUser)updateDataIfStale() error {
    if u.UpdatedAt.Add(slateTime).Before(time.Now()) {
        return nil
    }
    return u.updateData()
}

func (u *mUser)updateData() error {
    self, err := u.getTwitterAPI().GetSelf(nil)
    if err != nil {
        return err
    }
    return cUser.UpdateId(u.OID, bson.M{
        "name": self.ScreenName,
        "updated_at": time.Now(),
    })    
}

func (u *mUser)getTwitterAPI() *anaconda.TwitterApi {
    return anaconda.NewTwitterApi(u.TwitterToken, u.TwitterSecret)
}
