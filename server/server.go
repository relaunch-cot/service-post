package server

import "github.com/relaunch-cot/service-post/handler"

type postResource struct {
	handler *handler.Handlers
}

type IPostServer interface {
}

func NewPostServer(handler *handler.Handlers) IPostServer {
	return &postResource{
		handler: handler,
	}
}
