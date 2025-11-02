package methods

import (
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/post"
	"google.golang.org/grpc"
)

func RegisterGrpcServices(s *grpc.Server) {
	pb.RegisterPostServiceServer(s, resource.Server.Post)
}
