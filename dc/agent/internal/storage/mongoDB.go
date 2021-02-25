package storage

import (
	"time"

	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	session *mgo.Session
	images  *mgo.Collection
}

type Images struct {
	ID        string      `json:"id" bson:"id"`
	Result    interface{} `json:"result" bson:"result"`
	Text      string      `json:"text" bson:"text"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
}

func newMongoDB(addr string) (*MongoDB, error) {
	session, err := mgo.Dial(addr)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return &MongoDB{
		session: session,
		images:  session.DB("dc").C("storage"),
	}, nil
}

func (m *MongoDB) Close() {
	m.session.Close()
}

func (m *MongoDB) Doc(id string, res interface{}) error {
	mres := res.(map[string]interface{})
	text := mres["text"]
	delete(mres, "text")
	if err := m.images.Insert(&Images{
		ID:        id,
		Result:    res,
		Text:      text.(string),
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}
	return nil
}
