package transport

import "github.com/zaynjarvis/fyp/dc/api"

type CollectionService interface {
	Start()
	Stop()
	SendNotification(*api.CollectionEvent)
	RecvConfig() <-chan *api.CollectionConfig
}

func New(port string, info *api.AgentInfo, push bool) CollectionService {
	if push {
		return newPushModel(port, info)
	}
	panic("not implemented")
}
