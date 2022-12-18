package github

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (h *Helper) GetLastCommit() (commit GithubLastCommit, err error) {
	url := fmt.Sprintf(GetLastCommit, h.repo)

	data, err := h.getRequest(url)
	if err != nil {
		return commit, err
	}

	err = json.Unmarshal(data, &commit)
	return commit, err
}

func (h *Helper) CreateNewCommit(prevSha, currentSha, name, email string) (sha string, err error) {
	c := NewCommitRequest{
		Message: "commit with rest api",
		Parents: []string{
			prevSha,
		},
		Tree: currentSha,
		Author: struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{Name: name, Email: email},
	}

	url := fmt.Sprintf(NewCommit, h.repo)

	cData, err := json.Marshal(c)
	if err != nil {
		return sha, err
	}

	data, err := h.postRequest(url, strings.NewReader(string(cData)))
	if err != nil {
		return sha, err
	}

	var s CommitResponse
	err = json.Unmarshal(data, &s)
	if err != nil {
		return sha, err
	}

	return s.Sha, nil
}
