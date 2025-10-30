package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type wikipediaSearchResponse struct {
	Query struct {
		Search []struct {
			Title string `json:"title"`
		} `json:"search"`
	} `json:"query"`
}

type wikipediaSummaryResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Extract     string `json:"extract"`
	Thumbnail   struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
	ContentUrls struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
}

type Movie struct {
	client *http.Client
}

func NewMovie(client *http.Client) *Movie {
	return &Movie{client: client}
}

func (s *Movie) doGetJSON(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "wheel-of-wonder (+https://github.com/kylecain/wheel-of-wonder)")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		slog.Error("HTTP request failed", "url", url, "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "url", url, "error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Non-200 response from API", "url", url, "status", resp.Status, "body", string(body))
		return nil, fmt.Errorf("non-200 response: %s", resp.Status)
	}

	return body, nil
}

func (s *Movie) FetchMovie(title string) (*model.MovieInfo, error) {
	query := strings.TrimSpace(title + " film")

	searchURL := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch=%s&format=json",
		url.QueryEscape(query),
	)

	searchBody, err := s.doGetJSON(searchURL)
	if err != nil {
		slog.Error("Failed to read search response", "error", err)
		return nil, err
	}

	var searchData wikipediaSearchResponse
	if err := json.Unmarshal(searchBody, &searchData); err != nil {
		slog.Error("Failed to parse search JSON", "error", err)
		return nil, err
	}

	if len(searchData.Query.Search) == 0 {
		slog.Warn("No Wikipedia results found", "query", query)
		return nil, fmt.Errorf("no Wikipedia results found for %q", title)
	}

	bestTitle := searchData.Query.Search[0].Title
	bestTitleURL := strings.ReplaceAll(bestTitle, " ", "_")

	summaryURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", url.PathEscape(bestTitleURL))

	summaryBody, err := s.doGetJSON(summaryURL)
	if err != nil {
		slog.Error("Failed to read summary response", "error", err)
		return nil, err
	}

	var summaryData wikipediaSummaryResponse
	if err := json.Unmarshal(summaryBody, &summaryData); err != nil {
		slog.Error("Failed to parse summary JSON", "error", err)
		return nil, err
	}

	return &model.MovieInfo{
		Title:       summaryData.Title,
		Description: summaryData.Extract,
		ImageURL:    summaryData.Thumbnail.Source,
		ContentURL:  summaryData.ContentUrls.Desktop.Page,
	}, nil
}

func (s *Movie) FetchImageAndEncode(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "wheel-of-wonder (+https://github.com/kylecain/wheel-of-wonder)")

	resp, err := s.client.Do(req)
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
