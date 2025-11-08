package handler

import (
	"context"

	"github.com/google/uuid"
	libModels "github.com/relaunch-cot/lib-relaunch-cot/models"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/post"
	"github.com/relaunch-cot/service-post/repositories"
	"github.com/relaunch-cot/service-post/resource/transformer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	GetAllLikesFromPost(ctx *context.Context, in *pb.GetAllLikesFromPostRequest) (*pb.GetAllLikesFromPostResponse, error)
	UpdateLikesFromPostOrComment(ctx *context.Context, in *pb.UpdateLikesFromPostOrCommentRequest) (*pb.UpdateLikesFromPostOrCommentResponse, error)
	CreateCommentOrReply(ctx *context.Context, in *pb.CreateCommentOrReplyRequest) (*pb.CreateCommentOrReplyResponse, error)
	DeleteCommentOrReply(ctx *context.Context, in *pb.DeleteCommentOrReplyRequest) (*pb.DeleteCommentOrReplyResponse, error)
	GetAllCommentsFromPost(ctx *context.Context, in *pb.GetAllCommentsFromPostRequest) (*pb.GetAllCommentsFromPostResponse, error)
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

func (r *resource) GetAllLikesFromPost(ctx *context.Context, in *pb.GetAllLikesFromPostRequest) (*pb.GetAllLikesFromPostResponse, error) {
	allLikesFromPost, err := r.repositories.Mysql.GetAllLikesFromPost(ctx, in.PostId, in.UserId)
	if err != nil {
		return nil, err
	}

	baseModelsAllLikesFromPost, err := transformer.GetAllLikesFromPostOrCommentToBaseModels(allLikesFromPost)
	if err != nil {
		return nil, err
	}

	getLikesFromPostResponse := &pb.GetAllLikesFromPostResponse{
		LikesFromPost: baseModelsAllLikesFromPost,
	}

	return getLikesFromPostResponse, nil
}

func (r *resource) UpdateLikesFromPostOrComment(ctx *context.Context, in *pb.UpdateLikesFromPostOrCommentRequest) (*pb.UpdateLikesFromPostOrCommentResponse, error) {
	err := r.repositories.Mysql.UpdateLikesFromPostOrComments(ctx, in.PostId, in.UserId, in.Type)
	if err != nil {
		return nil, err
	}

	likesFromPostOrComment := new(libModels.PostLikes)
	if in.Type == "post" {
		likesFromPostOrComment, err = r.repositories.Mysql.GetAllLikesFromPost(ctx, in.PostId, in.UserId)
	} else if in.Type == "comment" {
		likesFromPostOrComment, err = r.repositories.Mysql.GetAllLikesFromComment(ctx, in.CommentId, in.UserId)
	}
	if err != nil {
		return nil, err
	}

	baseModelsLikesFromPostOrComment, err := transformer.GetAllLikesFromPostOrCommentToBaseModels(likesFromPostOrComment)
	if err != nil {
		return nil, err
	}

	updateLikesFromPostOrCommentResponse := &pb.UpdateLikesFromPostOrCommentResponse{
		LikesFromPost: baseModelsLikesFromPostOrComment,
	}

	return updateLikesFromPostOrCommentResponse, nil
}

func (r *resource) CreateCommentOrReply(ctx *context.Context, in *pb.CreateCommentOrReplyRequest) (*pb.CreateCommentOrReplyResponse, error) {
	var err error

	if in.Type == "comment" {
		commentId := uuid.New().String()
		err = r.repositories.Mysql.CreateComment(ctx, in.PostId, commentId, in.UserId, in.Content)
	} else if in.Type == "reply" {
		replyId := uuid.New().String()
		err = r.repositories.Mysql.CreateReply(ctx, in.ParentCommentId, replyId, in.UserId, in.Content)
	} else {
		return nil, status.Error(codes.InvalidArgument, "invalid comment type")
	}

	if err != nil {
		return nil, err
	}

	comments, err := r.repositories.Mysql.GetAllCommentsFromPost(ctx, in.PostId, in.UserId)
	if err != nil {
		return nil, err
	}

	baseModelsComment, err := transformer.GetAllCommentsFromPostToBaseModels(comments)
	if err != nil {
		return nil, err
	}

	addCommentToPostResponse := &pb.CreateCommentOrReplyResponse{
		CommentsFromPost: baseModelsComment,
	}

	return addCommentToPostResponse, nil
}

func (r *resource) DeleteCommentOrReply(ctx *context.Context, in *pb.DeleteCommentOrReplyRequest) (*pb.DeleteCommentOrReplyResponse, error) {
	var err error
	var postId *string

	if in.Type == "comment" {
		postId, err = r.repositories.Mysql.DeleteComment(ctx, in.CommentId, in.UserId)
	} else if in.Type == "reply" {
		postId, err = r.repositories.Mysql.DeleteReply(ctx, in.ReplyId, in.UserId)
	} else {
		return nil, status.Error(codes.InvalidArgument, "invalid comment type")
	}
	if err != nil {
		return nil, err
	}

	comments, err := r.repositories.Mysql.GetAllCommentsFromPost(ctx, *postId, in.UserId)
	if err != nil {
		return nil, err
	}

	baseModelsComment, err := transformer.GetAllCommentsFromPostToBaseModels(comments)
	if err != nil {
		return nil, err
	}

	removeCommentFromPostResponse := &pb.DeleteCommentOrReplyResponse{
		CommentsFromPost: baseModelsComment,
	}

	return removeCommentFromPostResponse, nil
}

func (r *resource) GetAllCommentsFromPost(ctx *context.Context, in *pb.GetAllCommentsFromPostRequest) (*pb.GetAllCommentsFromPostResponse, error) {
	comments, err := r.repositories.Mysql.GetAllCommentsFromPost(ctx, in.PostId, in.UserId)
	if err != nil {
		return nil, err
	}

	baseModelsComment, err := transformer.GetAllCommentsFromPostToBaseModels(comments)
	if err != nil {
		return nil, err
	}

	getAllCommentsFromPostResponse := &pb.GetAllCommentsFromPostResponse{
		CommentsFromPost: baseModelsComment,
	}

	return getAllCommentsFromPostResponse, nil
}

func NewPostHandler(repositories *repositories.Repositories) IPostHandler {
	return &resource{
		repositories: repositories,
	}
}
