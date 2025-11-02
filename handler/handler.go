package handler

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/post"
	"github.com/relaunch-cot/service-post/repositories"
)

type resource struct {
	repositories *repositories.Repositories
}

type IPostHandler interface {
	CreatePost(ctx *context.Context, in *pb.CreatePostRequest) error
}

func (r *resource) CreatePost(ctx *context.Context, in *pb.CreatePostRequest) error {
	postId := uuid.New().String()
	err := r.repositories.Mysql.CreatePost(ctx, in.UserId, postId, in.Title, in.Content, in.Type, in.UrlImagePost)
	if err != nil {
		return err
	}

	return nil
}

func NewPostHandler(repositories *repositories.Repositories) IPostHandler {
	return &resource{
		repositories: repositories,
	}
}
