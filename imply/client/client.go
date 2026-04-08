// Copyright IBM Corp. 2026

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents the HTTP client for interacting with the API.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

func normalizeHost(host string) string {
	hostURL := strings.TrimSpace(host)
	hostURL = strings.Replace(hostURL, ".app.imply.io", ".api.imply.io", 1)
	hostURL = strings.TrimRight(hostURL, "/")
	return hostURL + "/v1"
}

// NewClient creates and returns a new Client.
func NewClient(host, apiKey *string) (*Client, error) {
	if host == nil || strings.TrimSpace(*host) == "" {
		return nil, errors.New("host cannot be nil or empty")
	}
	if apiKey == nil || strings.TrimSpace(*apiKey) == "" {
		return nil, errors.New("apiKey cannot be nil or empty")
	}

	return &Client{
		HostURL:    normalizeHost(*host),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ApiKey:     "Basic " + strings.TrimSpace(*apiKey),
	}, nil
}

// doRequest performs the actual HTTP request to the API.
func (c *Client) doRequest(method, path string, body any) (map[string]any, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.HostURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	if len(respBody) == 0 || resp.StatusCode == http.StatusNoContent {
		return map[string]any{}, nil
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return result, nil
}

func (c *Client) Get(path string) (map[string]any, error) {
	return c.doRequest(http.MethodGet, path, nil)
}
func (c *Client) Post(path string, body any) (map[string]any, error) {
	return c.doRequest(http.MethodPost, path, body)
}
func (c *Client) Put(path string, body any) (map[string]any, error) {
	return c.doRequest(http.MethodPut, path, body)
}
func (c *Client) Delete(path string) error {
	_, err := c.doRequest(http.MethodDelete, path, nil)
	return err
}
func (c *Client) DeleteWithBody(path string, body any) (map[string]any, error) {
	return c.doRequest(http.MethodDelete, path, body)
}
