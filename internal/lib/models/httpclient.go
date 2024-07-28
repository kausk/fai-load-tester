package models

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(timeoutSeconds time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeoutSeconds * time.Second,
		},
	}
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// TODO: POST requests with headers and headers
