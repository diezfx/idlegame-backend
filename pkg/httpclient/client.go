package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const httpSuccessThreshold = 400

var globalMiddleware = []HttpMiddleware{}

type RequestHandler func(request *http.Request) (*http.Response, error)
type HttpMiddleware func(RequestHandler) RequestHandler

const defaultTimeout = 30 * time.Second

type Config struct {
	Host    string
	Timeout time.Duration
}

type Client struct {
	host        string
	middlewares []HttpMiddleware

	client http.Client
}

func New(cfg Config) *Client {
	timeout := defaultTimeout

	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}

	return &Client{
		client: http.Client{Timeout: timeout},
		host:   cfg.Host,
	}
}

func (c *Client) RegisterMiddleware(m HttpMiddleware) {
	c.middlewares = append(c.middlewares, m)
}

// ensure result only pointes
func (c *Client) Get(ctx context.Context, path string, result any) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) Post(ctx context.Context, path string, requestBody any, result any) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, requestBody, result)
}

func (c *Client) doRequest(ctx context.Context, method string, path string, requestBody any, responseBody any) (*http.Response, error) {

	addr, err := url.JoinPath(c.host, path)
	if err != nil {
		return nil, fmt.Errorf("join path: %w", err)
	}

	var requestBodyJSON []byte
	if requestBody != nil {
		requestBodyJSON, err = json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, addr, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	requestFunc := func(request *http.Request) (*http.Response, error) {
		return c.client.Do(request)
	}
	for _, m := range globalMiddleware {
		requestFunc = m(requestFunc)
	}

	for _, m := range c.middlewares {
		requestFunc = m(requestFunc)
	}

	response, err := requestFunc(request)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	if responseBody != nil && response.StatusCode <= httpSuccessThreshold {
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			return nil, fmt.Errorf("unmarshal response: %w", err)
		}
	}
	return response, nil
}
