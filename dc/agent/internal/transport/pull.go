package transport

import (
	"github.com/zaynjarvis/fyp/dc/api"
)

type PullModel struct {

}

func (p PullModel) Start() {
	panic("implement me")
}

func (p PullModel) Stop() {
	panic("implement me")
}

func (p PullModel) SendNotification(event api.CollectionEvent) {
	panic("implement me")
}

func (p PullModel) ReceiveConfig() api.CollectionConfig {
	panic("implement me")
}
