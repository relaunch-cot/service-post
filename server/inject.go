package server

import (
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/post"
	"github.com/relaunch-cot/service-post/handler"
)

type Servers struct {
	Post pb.PostServiceServer
}

func (s *Servers) Inject(handler *handler.Handlers) {
	s.Post = NewPostServer(handler)
}
