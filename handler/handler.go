package handler

import "github.com/relaunch-cot/service-post/repositories"

type resource struct {
	repositories *repositories.Repositories
}

type IPostHandler interface {
}

func NewPostHandler(repositories *repositories.Repositories) IPostHandler {
	return &resource{
		repositories: repositories,
	}
}
