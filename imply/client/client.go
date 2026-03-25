package client

import (
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

// NewClient creates and returns a new Client.
func NewClient(host, apiKey *string) (*Client, error) {
	if host == nil || *host == "" {
		return nil, errors.New("host cannot be nil or empty")
	}

	if apiKey == nil || *apiKey == "" {
		return nil, errors.New("apiKey cannot be nil or empty")
	}

	hostURL := *host
	// Replace .app.imply.io with .api.imply.io
	hostURL = strings.Replace(hostURL, ".app.imply.io", ".api.imply.io", 1)

	// Ensure host URL ends with a slash
	if !strings.HasSuffix(hostURL, "/") {
		hostURL += "/"
	}

	return &Client{
		HostURL:    hostURL + "v1",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ApiKey:     "Basic " + *apiKey,
	}, nil
}

// doRequest performs the actual HTTP request to the API.
func (c *Client) doRequest(method, path string, body any) (map[string]any, error) {
	// Prepare the request body if necessary
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = strings.NewReader(string(jsonBody))
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, c.HostURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set the headers
	req.Header.Set("Authorization", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Handle non-OK status codes
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Unmarshal the response JSON
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return result, nil
}

// HTTP Methods for API interaction

// Get performs a GET request to the specified path.
func (c *Client) Get(path string) (map[string]any, error) {
	return c.doRequest(http.MethodGet, path, nil)
}

// Post performs a POST request to the specified path with the given body.
func (c *Client) Post(path string, body any) (map[string]any, error) {
	return c.doRequest(http.MethodPost, path, body)
}

// Put performs a PUT request to the specified path with the given body.
func (c *Client) Put(path string, body any) (map[string]any, error) {
	return c.doRequest(http.MethodPut, path, body)
}

// Delete performs a DELETE request to the specified path.
func (c *Client) Delete(path string) error {
	_, err := c.doRequest(http.MethodDelete, path, nil)
	return err
}
