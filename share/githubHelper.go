package share

import (
	"context"
	"sync-bot/services/github"
	"sync-bot/types"
)

var gitHubHelper *github.Helper

func NewGithubHelper(c types.Config) {
	gitHubHelper = github.NewHelper(c.GHelper.Token, c.GHelper.Repo, context.Background())
}

func GithubHelper() *github.Helper {
	return gitHubHelper
}
