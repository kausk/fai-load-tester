package tests

import (
	"testing"

	"loadtester/pkg/models"
	"loadtester/pkg/service"

	"github.com/stretchr/testify/assert"
)

func TestTestRunnerMultiPhase(t *testing.T) {
	tcOneIteration := models.TestCaseGetRequest{
		URL:           "http://google.com",
		NumIterations: 2,
	}
	phase1 := models.TestPhase{
		Name:            "50 QPS",
		NumVirtualUsers: 25,
		TestCase:        tcOneIteration,
		DurationSeconds: 4,
	}
	phase2 := models.TestPhase{
		Name:            "100 QPS",
		NumVirtualUsers: 50,
		TestCase:        tcOneIteration,
		DurationSeconds: 4,
	}
	tpSpawnVUsAtStart := models.TestPlan{
		Name: "QPS Test - Spawn VUs at start",
		Phases: []models.TestPhase{
			phase1,
			phase2,
		},
		HTTPTimeoutMilliseconds: 500,
		SpawnVUsEvenly:          false,
	}
	tpSpawnVUsEvenly := models.TestPlan{
		Name: "QPS Test - Spawn VUs at start",
		Phases: []models.TestPhase{
			phase1,
			phase2,
		},
		HTTPTimeoutMilliseconds: 500,
		SpawnVUsEvenly:          true,
	}
	for _, tp := range []models.TestPlan{tpSpawnVUsAtStart, tpSpawnVUsEvenly} {
		server := NewMockHTTPServer(false, false)
		tr := service.NewTestRunnerWithHTTPClientFactory(tp, func() service.HTTPClientInterface {
			return NewMockHTTPClient(server, 50)
		})
		tr.Start()
		_ = <-tr.Done
		firstPhase := tr.MetricsByPhase[0]
		assertPhaseMetrics(t, firstPhase, 25*4, 0)
		assert.Equal(t, 25*4*2+50*4*2, server.SuccessfulQueries) // total number of queries to server in first phase
		assertOverallMetrics(t, tr.OverallMetrics, (25*4)+(50*4), 0)
	}
}

func assertPhaseMetrics(t *testing.T, metrics models.VirtualUserAggregatedMetrics, numVUsCreated, numVUsFailed int) {
	assert.Equal(t, 25*4, metrics.NumVUsCreated) // 25 VUs, 2 iterations, 3 seconds
	assert.Equal(t, 0, metrics.NumVUsFailed)
	// TODO: finer-grained testing of metrics calculations
	assert.NotZero(t, metrics.AvgVUDuration)
	assert.NotZero(t, metrics.P50VUDuration)
	assert.NotZero(t, metrics.P95VUDuration)
	assert.NotZero(t, metrics.P99VUDuration)
	assert.NotZero(t, metrics.MaxVUDuration)
}

func assertOverallMetrics(t *testing.T, metrics models.OverallMetrics, numVUSCreated, numVUSFailed int) {
	// coarse-grained testing of overall metrics
	assert.Equal(t, numVUSCreated, metrics.NumVUsCreated)
	assert.Equal(t, numVUSFailed, metrics.NumVUsFailed)
	assert.NotZero(t, metrics.AvgVUDuration)
	assert.NotZero(t, metrics.MaxVUDuration)
}

// TODO: simulate error conditions
func TestTestRunnerErrors(t *testing.T) {

}

func TestTestRunnerTimeouts(t *testing.T) {

}
