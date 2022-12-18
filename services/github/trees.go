package github

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (h *Helper) GetTreeWithSha(sha string) (tree GithubRepoTree, err error) {
	url := fmt.Sprintf(GetRepoTree, h.repo, sha)

	data, err := h.getRequest(url)
	if err != nil {
		return tree, err
	}

	err = json.Unmarshal(data, &tree)
	return tree, err
}

func (h *Helper) NewTreeWithContent(prevTreeSha string, trees []RepoTree) (resp NewRepoTreeResponse, err error) {
	url := fmt.Sprintf(NewRepoTree, h.repo)
	payload := NewRepoTreeRequest{
		BaseTree: prevTreeSha,
		Tree:     trees,
	}

	str, _ := json.Marshal(payload)
	data, err := h.postRequest(url, strings.NewReader(string(str)))
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	return resp, err
}
