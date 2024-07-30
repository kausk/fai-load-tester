package tests

import (
	"errors"
	"net/http"
	"time"

	"loadtester/internal/lib/service"
)

type MockHTTPClient struct {
	server          *MockHTTPServer
	clientLatencyMS int
}

type MockHTTPServer struct {
	TotalQueries      int
	SuccessfulQueries int
	FailedQueries     int
	returnErrors      bool
	returnTimeouts    bool
}

func NewMockHTTPServer(returnErrors, returnTimeouts bool) *MockHTTPServer {
	return &MockHTTPServer{
		returnErrors:   returnErrors,
		returnTimeouts: returnTimeouts,
	}
}

func NewMockHTTPClient(mockHTTPServer *MockHTTPServer, clientLatencyMS int) *MockHTTPClient {
	return &MockHTTPClient{
		server:          mockHTTPServer,
		clientLatencyMS: clientLatencyMS,
	}
}

func (m MockHTTPClient) Get(_ string) (*http.Response, error) {
	time.Sleep(time.Duration(m.clientLatencyMS) * time.Millisecond) // simulate latency
	if m.server.returnTimeouts {
		m.server.FailedQueries++
		return nil, errors.New("method timed out")
	}
	if m.server.returnErrors {
		m.server.FailedQueries++
		return &http.Response{
			StatusCode: 404,
		}, nil
	}
	return &http.Response{
		StatusCode: 200,
	}, nil
}

var _ service.HTTPClientInterface = &MockHTTPClient{}
