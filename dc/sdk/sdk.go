package sdk

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/api"
	"google.golang.org/grpc"
)

type (
	CollectionSDK struct {
		agentAddr     string
		collectionCfg *api.CollectionConfig
		imgCh         chan *api.ImageReport
		tags          []Tag
	}
	Tag struct {
		K string
		V string
	}
	With func(*CollectionSDK)
)

func WithTag(K string, V string) With {
	return func(sdk *CollectionSDK) {
		sdk.tags = append(sdk.tags, Tag{K, V})
	}
}

func WithSession(uid string) With {
	return WithTag("session", uid)
}

func New(agentAddr, service string, ws ...With) *CollectionSDK {
	sdk := &CollectionSDK{
		agentAddr: agentAddr,
		imgCh:     make(chan *api.ImageReport),
		tags:      []Tag{{K: "service", V: service}},
	}
	for _, w := range ws {
		w(sdk)
	}
	return sdk
}

// connect to a local agent, if an agent doesn't exist, collection is not possible
// same model as Consul
func (c *CollectionSDK) Connect(ctx context.Context) error {
	conn, err := grpc.Dial(c.agentAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	clt := api.NewLocalClient(conn)
	errCh := make(chan error)
	go func() {
		stream, err := clt.Image(ctx)
		if err != nil {
			errCh <- err
			return
		}
		for img := range c.imgCh {
			logrus.Info("here")
			if err := stream.Send(img); err != nil {
				logrus.Error("stream send image has error: ", err)
			}
		}
	}()
	return <-errCh
}

// based on collection config, determine whether the image should be stored
func (c *CollectionSDK) Image(img []byte, result interface{}, tags ...Tag) {
	data, err := json.Marshal(result)
	if err != nil {
		logrus.Debug("marshal failed, err: ", err)
		return
	}
	ts := make([]*api.Tag, len(c.tags)+len(tags))
	for _, t := range append(c.tags, tags...) {
		ts = append(ts, &api.Tag{
			Key: t.K,
			Val: t.V,
		})
	}
	c.imgCh <- &api.ImageReport{Img: img, Result: data, Tags: ts}
}
