package tests

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"loadtester/pkg/service"
)

type MockHTTPClient struct {
	server          *MockHTTPServer
	clientLatencyMS int
}

type MockHTTPServer struct {
	SuccessfulQueries int
	FailedQueries     int
	mu                sync.Mutex // lock for updating above counts since MockHTTPClients are run concurrently
	// TODO: currently not used in unit tests, but intended to simulate error conditions
	returnErrors   bool
	returnTimeouts bool
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
	m.server.mu.Lock()
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
	m.server.SuccessfulQueries++
	m.server.mu.Unlock()
	return &http.Response{
		StatusCode: 200,
	}, nil
}

var _ service.HTTPClientInterface = &MockHTTPClient{}
