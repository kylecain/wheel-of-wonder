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
	logger *slog.Logger
}

func NewMovie(client *http.Client, logger *slog.Logger) *Movie {
	return &Movie{
		client: client,
		logger: logger.With(slog.String("component", "service.movie")),
	}
}

func (s *Movie) doGetJSON(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "wheel-of-wonder (+https://github.com/kylecain/wheel-of-wonder)")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Debug("non-200 response from api", slog.String("url", url), slog.String("status", resp.Status), slog.String("body", string(body)))
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
		return nil, fmt.Errorf("get failed: %w", err)
	}

	var searchData wikipediaSearchResponse
	if err := json.Unmarshal(searchBody, &searchData); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	if len(searchData.Query.Search) == 0 {
		s.logger.Warn("no search results found", slog.String("query", query))
		return nil, fmt.Errorf("no search results found for %s", title)
	}

	bestTitle := searchData.Query.Search[0].Title
	bestTitleURL := strings.ReplaceAll(bestTitle, " ", "_")

	summaryURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", url.PathEscape(bestTitleURL))

	summaryBody, err := s.doGetJSON(summaryURL)
	if err != nil {
		return nil, fmt.Errorf("get failed: %w", err)
	}

	var summaryData wikipediaSummaryResponse
	if err := json.Unmarshal(summaryBody, &summaryData); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
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
		return "", fmt.Errorf("failed to make request: %w", err)
	}

	req.Header.Set("User-Agent", "wheel-of-wonder (+https://github.com/kylecain/wheel-of-wonder)")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(imageBytes)
	return fmt.Sprintf("data:%s;base64,%s", resp.Header.Get("Content-Type"), encoded), nil
}
