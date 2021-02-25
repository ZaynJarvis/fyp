package storage

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	session *mgo.Session
	images  *mgo.Collection
}

type Images struct {
	ID        string      `json:"id" bson:"id"`
	Result    interface{} `json:"result" bson:"result"`
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

func (m *MongoDB) Data(id string, res []byte) error {
	var result map[string]interface{}
	if err := json.Unmarshal(res, &result); err != nil {
		logrus.Error("unmarshal result failed, res: ", string(res))
	}
	if err := m.images.Insert(&Images{
		ID:        id,
		Result:    result,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}
	return nil
}
