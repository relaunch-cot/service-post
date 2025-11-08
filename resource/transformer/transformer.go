package transformer

import (
	"encoding/json"

	libModels "github.com/relaunch-cot/lib-relaunch-cot/models"
	pbBaseModels "github.com/relaunch-cot/lib-relaunch-cot/proto/base_models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetPostToBaseModels(post *libModels.Post) (*pbBaseModels.Post, error) {
	b, err := json.Marshal(post)
	if err != nil {
		return nil, status.Error(codes.Internal, "error marshalling post. Details: "+err.Error())
	}

	var pbPost pbBaseModels.Post
	err = json.Unmarshal(b, &pbPost)
	if err != nil {
		return nil, status.Error(codes.Internal, "error unmarshalling post. Details: "+err.Error())
	}

	return &pbPost, nil
}

func GetAllPostsToBaseModels(posts []*libModels.Post) ([]*pbBaseModels.Post, error) {
	var pbPosts []*pbBaseModels.Post
	for _, post := range posts {
		pbPost, err := GetPostToBaseModels(post)
		if err != nil {
			return nil, err
		}
		pbPosts = append(pbPosts, pbPost)
	}
	return pbPosts, nil
}

func GetAllPostsFromUserToBaseModels(posts []*libModels.Post) ([]*pbBaseModels.Post, error) {
	var pbPosts []*pbBaseModels.Post
	for _, post := range posts {
		pbPost, err := GetPostToBaseModels(post)
		if err != nil {
			return nil, err
		}
		pbPosts = append(pbPosts, pbPost)
	}
	return pbPosts, nil
}

func GetAllLikesFromPostOrCommentToBaseModels(postLikes *libModels.PostLikes) (*pbBaseModels.PostLikes, error) {
	var pbPostLikes *pbBaseModels.PostLikes
	b, err := json.Marshal(postLikes)
	if err != nil {
		return nil, status.Error(codes.Internal, "error marshalling post likes. Details: "+err.Error())
	}

	err = json.Unmarshal(b, &pbPostLikes)
	if err != nil {
		return nil, status.Error(codes.Internal, "error unmarshalling post likes. Details: "+err.Error())
	}

	return pbPostLikes, nil
}

func GetAllCommentsFromPostToBaseModels(postComments *libModels.PostComments) (*pbBaseModels.PostComments, error) {
	var pbPostComments *pbBaseModels.PostComments
	b, err := json.Marshal(postComments)
	if err != nil {
		return nil, status.Error(codes.Internal, "error marshalling post comments. Details: "+err.Error())
	}

	err = json.Unmarshal(b, &pbPostComments)
	if err != nil {
		return nil, status.Error(codes.Internal, "error unmarshalling post comments. Details: "+err.Error())
	}

	return pbPostComments, nil
}
