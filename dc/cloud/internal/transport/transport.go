package transport

import "github.com/zaynjarvis/fyp/dc/api"

type CollectionService interface {
	Start()
	Stop()
	RecvNotification() <-chan *api.CollectionEvent
	SendConfig(*api.CollectionConfig)
}

func New(port string, push bool) CollectionService {
	if push {
		return newPushModel(port)
	}
	panic("not implemented")
}
