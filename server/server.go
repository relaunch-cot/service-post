package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/post"
	"github.com/relaunch-cot/service-post/handler"
)

type postResource struct {
	handler *handler.Handlers
	pb.UnimplementedPostServiceServer
}

func (r *postResource) CreatePost(ctx context.Context, in *pb.CreatePostRequest) (*empty.Empty, error) {
	err := r.handler.Post.CreatePost(&ctx, in)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *postResource) GetPost(ctx context.Context, in *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	response, err := r.handler.Post.GetPost(&ctx, in)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *postResource) GetAllPosts(ctx context.Context, in *empty.Empty) (*pb.GetAllPostsResponse, error) {
	response, err := r.handler.Post.GetAllPosts(&ctx)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *postResource) GetAllPostsFromUser(ctx context.Context, in *pb.GetAllPostsFromUserRequest) (*pb.GetAllPostsFromUserResponse, error) {
	response, err := r.handler.Post.GetAllPostsFromUser(&ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *postResource) UpdatePost(ctx context.Context, in *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	response, err := r.handler.Post.UpdatePost(&ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *postResource) DeletePost(ctx context.Context, in *pb.DeletePostRequest) (*empty.Empty, error) {
	err := r.handler.Post.DeletePost(&ctx, in)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *postResource) GetAllLikesFromPost(ctx context.Context, in *pb.GetAllLikesFromPostRequest) (*pb.GetAllLikesFromPostResponse, error) {
	response, err := r.handler.Post.GetAllLikesFromPost(&ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *postResource) UpdateLikesFromPostOrComment(ctx context.Context, in *pb.UpdateLikesFromPostOrCommentRequest) (*pb.UpdateLikesFromPostOrCommentResponse, error) {
	response, err := r.handler.Post.UpdateLikesFromPostOrComment(&ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *postResource) CreateCommentOrReply(ctx context.Context, in *pb.CreateCommentOrReplyRequest) (*pb.CreateCommentOrReplyResponse, error) {
	response, err := r.handler.Post.CreateCommentOrReply(&ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *postResource) DeleteCommentOrReply(ctx context.Context, in *pb.DeleteCommentOrReplyRequest) (*pb.DeleteCommentOrReplyResponse, error) {
	response, err := r.handler.Post.DeleteCommentOrReply(&ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *postResource) GetAllCommentsFromPost(ctx context.Context, in *pb.GetAllCommentsFromPostRequest) (*pb.GetAllCommentsFromPostResponse, error) {
	response, err := r.handler.Post.GetAllCommentsFromPost(&ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewPostServer(handler *handler.Handlers) pb.PostServiceServer {
	return &postResource{
		handler: handler,
	}
}
