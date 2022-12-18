package github

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (h *Helper) UpdateBranch(sha string) error {
	url := fmt.Sprintf(UpdateBranch, h.repo)

	bData, err := json.Marshal(struct {
		Sha string `json:"sha"`
	}{
		Sha: sha,
	})
	if err != nil {
		return err
	}

	data, err := h.postRequest(url, strings.NewReader(string(bData)))
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
