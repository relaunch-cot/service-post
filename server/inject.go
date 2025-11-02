package server

import "github.com/relaunch-cot/service-post/handler"

type Servers struct {
	Post IPostServer
}

func (s *Servers) Inject(handler *handler.Handlers) {
	s.Post = NewPostServer(handler)
}
