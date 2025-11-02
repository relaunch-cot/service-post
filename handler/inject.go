package handler

import "github.com/relaunch-cot/service-post/repositories"

type Handlers struct {
	Post IPostHandler
}

func (h *Handlers) Inject(repositories *repositories.Repositories) {
	h.Post = NewPostHandler(repositories)
}
