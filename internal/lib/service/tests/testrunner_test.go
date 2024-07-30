package tests

import (
	"fmt"
	"testing"
	"time"

	"loadtester/internal/lib/models"
	"loadtester/internal/lib/service"
)

func TestTestRunner50QPS(t *testing.T) {
	tcOneIteration := models.TestCaseGetRequest{
		URL:           "http://google.com",
		NumIterations: 1,
	}
	tp := models.TestPlan{
		Name: "QPS Test",
		Phases: []models.TestPhase{
			{
				Name:            "50 QPS",
				NumVirtualUsers: 50,
				TestCase:        tcOneIteration,
				DurationSeconds: 1,
			},
		},
		HTTPTimeoutMilliseconds: 5000,
	}
	server := NewMockHTTPServer(false, false)
	tr := service.NewTestRunnerWithHTTPClientFactory(tp, func() service.HTTPClientInterface {
		return NewMockHTTPClient(server)
	})
	tr.Start()
	time.Sleep(5 * time.Second)
	fmt.Println(tr.Results())
}

func TestTestRunnerErrors(t *testing.T) {

}

func TestTestRunnerTimeouts(t *testing.T) {

}

func TestTestRunnerMultiplePhases(t *testing.T) {

}

func TestTestRunnerMultipleVirtualWorkers(t *testing.T) {
}

func TestTestRunnerExtendedDuration(t *testing.T) {

}
