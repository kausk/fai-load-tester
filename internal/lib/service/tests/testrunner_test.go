package tests

import (
	"testing"

	"loadtester/internal/lib/models"
	"loadtester/internal/lib/service"

	"github.com/stretchr/testify/assert"
)

func TestTestRunnerMultiPhase(t *testing.T) {
	tcOneIteration := models.TestCaseGetRequest{
		URL:           "http://google.com",
		NumIterations: 2,
	}
	tp := models.TestPlan{
		Name: "QPS Test",
		Phases: []models.TestPhase{
			{
				Name:            "50 QPS",
				NumVirtualUsers: 25,
				TestCase:        tcOneIteration,
				DurationSeconds: 4,
			},
			{
				Name:            "100 QPS",
				NumVirtualUsers: 50,
				TestCase:        tcOneIteration,
				DurationSeconds: 4,
			},
		},
		HTTPTimeoutMilliseconds: 500,
	}
	server := NewMockHTTPServer(false, false)
	tr := service.NewTestRunnerWithHTTPClientFactory(tp, func() service.HTTPClientInterface {
		return NewMockHTTPClient(server, 50)
	})
	tr.Start()
	_ = <-tr.Done
	firstPhase := tr.MetricsByPhase[0]
	assert.Equal(t, 25*4, firstPhase.NumVUsCreated) // 25 VUs, 2 iterations, 3 seconds
	assert.Equal(t, 0, firstPhase.NumVUsFailed)
	// TODO: finer-grained testing of metrics calculations
	assert.NotZero(t, firstPhase.AvgVUDuration)
	assert.NotZero(t, firstPhase.P50VUDuration)
	assert.NotZero(t, firstPhase.P95VUDuration)
	assert.NotZero(t, firstPhase.P99VUDuration)
	assert.NotZero(t, firstPhase.MaxVUDuration)
	// coarse-grained testing of overall metrics
	assert.Equal(t, (25*4)+(50*4), tr.OverallMetrics.NumVUsCreated) // 25 VUs, 2 iterations, 3 seconds + 50 VUs, 2 iterations, 3 seconds
	assert.Equal(t, 0, tr.OverallMetrics.NumVUsFailed)
	assert.NotZero(t, tr.OverallMetrics.AvgVUDuration)
	assert.NotZero(t, tr.OverallMetrics.MaxVUDuration)
}

func TestTestRunnerErrors(t *testing.T) {

}

func TestTestRunnerTimeouts(t *testing.T) {

}
