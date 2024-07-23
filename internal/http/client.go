package http

import (
	"bytes"
	"ccrayz/event-scanner/internal/indexer/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout = 30
)

type Client struct {
	baseURL string
	client  *http.Client
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

type ResponseData struct {
	jsonrpc string `json:jsonrpc`
	id      int    `json:id`
	result  Result `json:result`
}

type Result struct {
	TotalConnected int `json:"totalConnected"`
	Peers          []models.PeerInfo
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Duration(defaultTimeout) * time.Second,
		},
	}
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL, "/") {
		c.baseURL += "/"
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	println(req.Method)
	println(req.URL.String())
	println(req.Body)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if err := CheckResponse(resp); err != nil {
		return resp, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		fmt.Errorf("failed to unmarshal as rpcResponse: %w", err)
	}

	// TODO(ccrayz): 이부분을 해결해야함
	// buf := new(bytes.Buffer)
	// teeReader := io.TeeReader(resp.Body, buf)
	// decErr := json.NewDecoder(teeReader).Decode(v)
	// if decErr == io.EOF {
	// 	decErr = nil
	// }
	// if decErr != nil {
	// 	err = fmt.Errorf("%s: %s", decErr.Error(), buf.String())
	// }

	return resp, err
}

func CheckResponse(resp *http.Response) error {
	if c := resp.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("API error with status code %d: %w", resp.StatusCode, err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("API error with status code %d: %s", resp.StatusCode, string(data))
	}

	message := ""
	if value, ok := raw["message"].(string); ok {
		message = value
	} else if value, ok := raw["error"].(string); ok {
		message = value
	}

	return fmt.Errorf("API error with status code %d: %s", resp.StatusCode, message)
}
