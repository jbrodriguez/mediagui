package lib

import (
	"io"
	"net/http"
)

func RestGet(url string, agent string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"
	req.Header.Set("User-Agent", agent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer Close(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
