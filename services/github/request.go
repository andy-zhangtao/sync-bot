package github

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (h *Helper) getRequest(url string) (data []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.token))

	res, err := client.Do(req)
	if err != nil {
		return data, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	return body, nil
}

func (h *Helper) postRequest(url string, payload *strings.Reader) (data []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return data, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.token))

	res, err := client.Do(req)
	if err != nil {
		return data, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	return body, nil
}
