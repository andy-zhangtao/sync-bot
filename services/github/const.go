package github

const (
	NewCommit     = "https://api.github.com/repos/%s/git/commits"
	GetLastCommit = "https://api.github.com/repos/%s/git/refs/heads/main"
	GetRepoTree   = "https://api.github.com/repos/%s/git/commits/%s"
	NewRepoTree   = "https://api.github.com/repos/%s/git/trees"
	UpdateBranch  = "https://api.github.com/repos/%s/git/refs/heads/main"
)
