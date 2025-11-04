package handler

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/post"
	"github.com/relaunch-cot/service-post/repositories"
	"github.com/relaunch-cot/service-post/resource/transformer"
)

type resource struct {
	repositories *repositories.Repositories
}

type IPostHandler interface {
	CreatePost(ctx *context.Context, in *pb.CreatePostRequest) error
	GetPost(ctx *context.Context, in *pb.GetPostRequest) (*pb.GetPostResponse, error)
	GetAllPosts(ctx *context.Context) (*pb.GetAllPostsResponse, error)
	GetAllPostsFromUser(ctx *context.Context, in *pb.GetAllPostsFromUserRequest) (*pb.GetAllPostsFromUserResponse, error)
	UpdatePost(ctx *context.Context, in *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error)
	DeletePost(ctx *context.Context, in *pb.DeletePostRequest) error
	UpdateLikesFromPost(ctx *context.Context, in *pb.UpdateLikesFromPostRequest) (*pb.UpdateLikesFromPostResponse, error)
}

func (r *resource) CreatePost(ctx *context.Context, in *pb.CreatePostRequest) error {
	postId := uuid.New().String()
	err := r.repositories.Mysql.CreatePost(ctx, in.UserId, postId, in.Title, in.Content, in.Type, in.UrlImagePost)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) GetPost(ctx *context.Context, in *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	response, err := r.repositories.Mysql.GetPost(ctx, in.PostId)
	if err != nil {
		return nil, err
	}

	baseModelsPost, err := transformer.GetPostToBaseModels(response)
	if err != nil {
		return nil, err
	}

	getPostResponse := &pb.GetPostResponse{
		Post: baseModelsPost,
	}

	return getPostResponse, nil
}

func (r *resource) GetAllPosts(ctx *context.Context) (*pb.GetAllPostsResponse, error) {
	response, err := r.repositories.Mysql.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	baseModelsPosts, err := transformer.GetAllPostsToBaseModels(response)
	if err != nil {
		return nil, err
	}

	getAllPostsResponse := &pb.GetAllPostsResponse{
		Posts: baseModelsPosts,
	}

	return getAllPostsResponse, nil
}

func (r *resource) GetAllPostsFromUser(ctx *context.Context, in *pb.GetAllPostsFromUserRequest) (*pb.GetAllPostsFromUserResponse, error) {
	response, err := r.repositories.Mysql.GetAllPostsFromUser(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	baseModelsPosts, err := transformer.GetAllPostsFromUserToBaseModels(response)
	if err != nil {
		return nil, err
	}

	getAllPostsFromUserResponse := &pb.GetAllPostsFromUserResponse{
		Posts: baseModelsPosts,
	}

	return getAllPostsFromUserResponse, nil
}

func (r *resource) UpdatePost(ctx *context.Context, in *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	err := r.repositories.Mysql.UpdatePost(ctx, in.PostId, in.UserId, in.Title, in.Content, in.UrlImagePost)
	if err != nil {
		return nil, err
	}

	post, err := r.repositories.Mysql.GetPost(ctx, in.PostId)
	if err != nil {
		return nil, err
	}

	baseModelsPost, err := transformer.GetPostToBaseModels(post)
	if err != nil {
		return nil, err
	}

	updatePostResponse := &pb.UpdatePostResponse{
		Post: baseModelsPost,
	}

	return updatePostResponse, nil
}

func (r *resource) DeletePost(ctx *context.Context, in *pb.DeletePostRequest) error {
	err := r.repositories.Mysql.DeletePost(ctx, in.PostId, in.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) UpdateLikesFromPost(ctx *context.Context, in *pb.UpdateLikesFromPostRequest) (*pb.UpdateLikesFromPostResponse, error) {
	err := r.repositories.Mysql.UpdateLikesFromPost(ctx, in.PostId, in.UserId, in.Liked)
	if err != nil {
		return nil, err
	}

	likesFromPost, err := r.repositories.Mysql.GetLikesFromPost(ctx, in.PostId)
	if err != nil {
		return nil, err
	}

	baseModelsLikesFromPost, err := transformer.GetLikesFromPostToBaseModels(likesFromPost)
	if err != nil {
		return nil, err
	}

	updateLikesFromPostResponse := &pb.UpdateLikesFromPostResponse{
		LikesFromPost: baseModelsLikesFromPost,
	}

	return updateLikesFromPostResponse, nil
}

func NewPostHandler(repositories *repositories.Repositories) IPostHandler {
	return &resource{
		repositories: repositories,
	}
}
