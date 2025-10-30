package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func FetchAndEncodeImage(url string, httpClient http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "wheel-of-wonder (+https://github.com/kylecain/wheel-of-wonder)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(imageBytes)
	return fmt.Sprintf("data:%s;base64,%s", resp.Header.Get("Content-Type"), encoded), nil
}
