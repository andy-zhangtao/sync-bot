package github

import "context"

type Helper struct {
	token string
	repo  string
	ctx   context.Context
}

func NewHelper(token, repo string, ctx context.Context) *Helper {
	return &Helper{
		token: token,
		repo:  repo,
		ctx:   ctx,
	}
}
