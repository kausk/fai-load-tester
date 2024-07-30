package service

import (
	"net/http"
	"time"
)

type HTTPClientInterface interface {
	Get(url string) (*http.Response, error)
}

type HTTPClient struct {
	client *http.Client
}

var _ HTTPClientInterface = &HTTPClient{}

func NewHTTPClient(timeoutMilliseconds time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeoutMilliseconds * time.Second,
		},
	}
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	resp, err := c.client.Get(url)
	return resp, err
}

// TODO: POST requests with headers and headers
