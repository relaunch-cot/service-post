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
