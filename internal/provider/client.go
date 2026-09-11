package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func (c *Client) endpointURL(path string) string {
	return c.baseURL + "/" + strings.TrimLeft(path, "/")
}

func (c *Client) request(ctx context.Context, method, path string, body []byte) (json.RawMessage, error) {
	requestURL := c.endpointURL(path)
	req, err := http.NewRequestWithContext(ctx, method, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("herodesk api %s %s failed with status %d: %s", method, requestURL, resp.StatusCode, string(respBody))
	}

	if len(respBody) == 0 {
		return json.RawMessage("{}"), nil
	}

	if !json.Valid(respBody) {
		encoded := []byte(fmt.Sprintf("%q", string(respBody)))
		return json.RawMessage(encoded), nil
	}

	return json.RawMessage(respBody), nil
}

func (c *Client) create(ctx context.Context, endpoint string, payload json.RawMessage) (json.RawMessage, error) {
	return c.request(ctx, http.MethodPost, endpoint, payload)
}

func (c *Client) read(ctx context.Context, endpoint, id string) (json.RawMessage, error) {
	return c.request(ctx, http.MethodGet, endpoint+"/"+url.PathEscape(id), nil)
}

func (c *Client) update(ctx context.Context, endpoint, id string, payload json.RawMessage) (json.RawMessage, error) {
	return c.request(ctx, http.MethodPut, endpoint+"/"+url.PathEscape(id), payload)
}

func (c *Client) delete(ctx context.Context, endpoint, id string) error {
	_, err := c.request(ctx, http.MethodDelete, endpoint+"/"+url.PathEscape(id), nil)
	return err
}

func (c *Client) list(ctx context.Context, endpoint string) (json.RawMessage, error) {
	return c.request(ctx, http.MethodGet, endpoint, nil)
}
