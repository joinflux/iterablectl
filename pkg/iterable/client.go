package iterable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.iterable.com/api/"
)

// Client represents an API client for Iterable
type Client struct {
	BaseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Iterable API client
func NewClient(apiKey string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		BaseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// User represents an Iterable user
type User struct {
	Email              string                 `json:"email,omitempty"`
	UserID             string                 `json:"userId,omitempty"`
	DataFields         map[string]interface{} `json:"dataFields,omitempty"`
	MergeNestedObjects bool                   `json:"mergeNestedObjects,omitempty"`
}

// APIError represents an error returned from the Iterable API
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"msg"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Iterable API error: %s - %s", e.Code, e.Message)
}

// newRequest creates a new HTTP request to the Iterable API
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)

	return req, nil
}

// do sends an API request and returns the response
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var apiErr APIError
		if err = json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("API request failed with status %d: %v", resp.StatusCode, err)
		}
		return &apiErr
	}

	if v != nil {
		if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

// GetUser retrieves a user by email from Iterable
func (c *Client) GetUser(email string) (*User, error) {
	path := fmt.Sprintf("users/%s", email)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		User User `json:"user"`
	}
	err = c.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.User, nil
}

// UpdateUser updates an Iterable user's profile
func (c *Client) UpdateUser(user User) error {
	req, err := c.newRequest("POST", "users/update", user)
	if err != nil {
		return err
	}

	var response struct {
		Success bool `json:"success"`
	}
	err = c.do(req, &response)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("failed to update user")
	}

	return nil
}
