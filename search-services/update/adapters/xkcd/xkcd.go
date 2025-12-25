package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"yadro.com/course/update/core"
)

type Client struct {
	log    *slog.Logger
	client http.Client
	url    string
}

func NewClient(url string, timeout time.Duration, log *slog.Logger) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("empty base url specified")
	}
	return &Client{
		client: http.Client{Timeout: timeout},
		log:    log,
		url:    url,
	}, nil
}

func (c Client) Get(ctx context.Context, id int) (core.XKCDInfo, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", c.url, id)
	var info core.XKCDInfo

	// 2. Создаем запрос с Context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.log.Error("failed to create request", "id", id, "error", err)
		return info, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("failed to execute request", "id", id, "error", err)
		return info, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.log.Warn("request failed", "id", id, "status", resp.Status)
		return info, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		c.log.Error("failed to decode json response", "id", id, "error", err)
		return info, fmt.Errorf("failed to decode response: %w", err)
	}

	return info, nil
}

func (c Client) LastID(ctx context.Context) (int, error) {
	url := fmt.Sprintf("%s/info.0.json", c.url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.log.Error("failed to create 'latest' request", "error", err)
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("failed to execute 'latest' request", "error", err)
		return 0, fmt.Errorf("request 'latest' failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.log.Warn("'latest' request failed", "status", resp.Status)
		return 0, fmt.Errorf("request 'latest' failed with status: %s", resp.Status)
	}

	var result struct {
		ID int `json:"num"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.log.Error("failed to decode 'latest' json response", "error", err)
		return 0, fmt.Errorf("failed to decode 'latest' response: %w", err)
	}

	return result.ID, nil
}
