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
	Email      string         `json:"email,omitempty"`
	UserID     string         `json:"userId,omitempty"`
	DataFields map[string]any `json:"dataFields,omitempty"`
}

// APIError represents an error returned from the Iterable API
type APIError struct {
	Code    string         `json:"code"`
	Message string         `json:"msg"`
	Params  map[string]any `json:"params,omitempty"`
}

// List represents an Iterable list
type List struct {
	CreatedAt   int64  `json:"createdAt"`
	Description string `json:"description,omitempty"`
	ID          int    `json:"id"`
	ListType    string
	Name        string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Iterable API error: %s - %s", e.Code, e.Message)
}

// newRequest creates a new HTTP request to the Iterable API
func (c *Client) newRequest(method, path string, body any) (*http.Request, error) {
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

// doPlain sends an HTTP request and returns the response body without parsing.
func (c *Client) doPlain(req *http.Request) (*[]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return &bodyBytes, nil
}

// do sends an API request and returns the response
func (c *Client) do(req *http.Request, v any) error {
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

// MergeUsers merges two Iterable users
func (c *Client) MergeUsers(src, dst string) (*APIError, error) {
	path := "users/merge"
	body := map[string]string{
		"destinationEmail": dst,
		"sourceEmail":      src,
	}
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	var response APIError
	err = c.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
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

// UserUpdateRequest represents a request to update a user's profile in Iterable
type UserUpdateRequest struct {
	Email              string         `json:"email,omitempty"`
	UserID             string         `json:"userId,omitempty"`
	DataFields         map[string]any `json:"dataFields,omitempty"`
	CreateNewFields    bool           `json:"createNewFields,omitempty"`
	MergeNestedObjects bool           `json:"mergeNestedObjects,omitempty"`
	PreferUserId       bool           `json:"preferUserId,omitempty"`
}

// UpdateUser updates an Iterable user's profile
func (c *Client) UpdateUser(user UserUpdateRequest) error {
	req, err := c.newRequest("POST", "users/update", user)
	if err != nil {
		return err
	}

	var response APIError
	err = c.do(req, &response)
	if err != nil {
		return err
	}

	if response.Code != "Success" {
		return fmt.Errorf("failed to update user: %v", response)
	}

	return nil
}

// GetLists retrieves the lists associated with the Iterable account
func (c *Client) GetLists() (*[]List, error) {
	req, err := c.newRequest("GET", "lists", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Lists []List `json:"lists"`
	}
	err = c.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.Lists, nil
}

// GetListUsers retrieves the users in a specific Iterable list
func (c *Client) GetListUsers(listId string, preferUserId bool) (*[]byte, error) {
	query := url.Values{}
	query.Set("preferUserId", fmt.Sprintf("%t", preferUserId))
	query.Set("listId", listId)
	path := fmt.Sprintf("lists/getUsers?%s", query.Encode())
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doPlain(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Campaign represents an Iterable campaign
type Campaign struct {
	CampaignState       string   `json:"campaignState"`
	CreatedAt           int      `json:"createdAt"`
	CreatedByUserId     string   `json:"createdByUserId"`
	EndedAt             int      `json:"endedAt,omitempty"`
	ID                  int      `json:"id"`
	Labels              []string `json:"labels,omitempty"`
	ListIds             []int    `json:"listIds,omitempty"`
	MessageMedium       string   `json:"messageMedium"`
	Name                string   `json:"name"`
	RecurringCampaignId int      `json:"recurringCampaignId,omitempty"`
	SendSize            int      `json:"sendSize,omitempty"`
	StartAt             int      `json:"startAt,omitempty"` // milliseconds
	SuppressionListIds  []int    `json:"suppressionListIds,omitempty"`
	TemplateId          int      `json:"templateId"`
	Type                string   `json:"type"`
	UpdatedAt           int      `json:"updatedAt"`
	UpdatedByUserId     string   `json:"updatedByUserId,omitempty"`
	WorkflowId          int      `json:"workflowId,omitempty"`
}

// GetCampaigns retrieves the campaigns associated with the Iterable account
func (c *Client) GetCampaigns() (*[]Campaign, error) {
	req, err := c.newRequest("GET", "campaigns", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Campaigns []Campaign `json:"campaigns"`
	}
	err = c.do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.Campaigns, nil
}

// DeleteUser removes a user from Iterable by their email address
func (c *Client) DeleteUser(email string) error {
	path := fmt.Sprintf("users/%s", email)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	var response APIError
	err = c.do(req, &response)
	if err != nil {
		return err
	}

	if response.Code != "Success" {
		return fmt.Errorf("failed to delete user: %v", response)
	}

	return nil
}
