package service

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"gonum.org/v1/gonum/stat"

	"loadtester/internal/lib/models"
)

type TestRunner struct {
	StartTime time.Time
	TestPlan  models.TestPlan
	status    models.Status
	// EndTime   time.Time // needed for time duration calculations
	MetricsByPhase []models.VirtualUserAggregatedMetrics
	OverallMetrics *models.TestRunAggregatedMetrics
	clientFactory  func() HTTPClientInterface
	Done           chan bool // TODO: implement and adopt in unit tests
}

func NewTestRunner(tp models.TestPlan) *TestRunner {
	return &TestRunner{
		StartTime: time.Now(),
		TestPlan:  tp,
		status:    models.WaitingForStart,
		clientFactory: func() HTTPClientInterface {
			return NewHTTPClient(tp.HTTPTimeoutMilliseconds)
		},
		Done: make(chan bool, 1),
	}
}

func NewTestRunnerWithHTTPClientFactory(tp models.TestPlan, factory func() HTTPClientInterface) *TestRunner {
	return &TestRunner{
		StartTime:     time.Now(),
		TestPlan:      tp,
		status:        models.WaitingForStart,
		clientFactory: factory,
		Done:          make(chan bool, 1),
	}
}

func (t *TestRunner) Start() {
	t.status = models.Running
	for _, phase := range t.TestPlan.Phases {
		fmt.Println("Starting phase: ", phase.Name)
		metrics := t.executePhase(phase)
		t.MetricsByPhase = append(t.MetricsByPhase, metrics)
	}
	t.OverallMetrics = calculateOverallMetrics(t.MetricsByPhase)
	t.status = models.Finished
	t.Done <- true
}

func (t *TestRunner) Status() models.Status {
	return t.status
}

func (t TestRunner) Results() *models.TestRunAggregatedMetrics {
	return t.OverallMetrics
}

func (t *TestRunner) executePhase(tp models.TestPhase) models.VirtualUserAggregatedMetrics {
	var wg sync.WaitGroup
	var vus []*VirtualUser
	if t.TestPlan.SpawnVUsEvenly {
		vus = t.spawnVUsEvenly(tp, &wg)
	} else {
		vus = t.spawnVUsAtStartOfSecond(tp, &wg)
	}
	wg.Wait() // wait for any remaining virtual users to finish
	metrics := calculatePhaseMetrics(vus)
	printPhaseMetrics(tp.Name, metrics)
	return metrics
}

func (t *TestRunner) spawnVUsAtStartOfSecond(tp models.TestPhase, wg *sync.WaitGroup) []*VirtualUser {
	ticker := time.NewTicker(1 * time.Second)
	done := createTestDurationChannel(tp.DurationSeconds)
	var vus []*VirtualUser
	for {
		select {
		case <-done:
			ticker.Stop()
			fmt.Printf("Test Phase %s is done\n", tp.Name)
			return vus
		case <-ticker.C:
			for i := 0; i < tp.NumVirtualUsers; i++ {
				wg.Add(1)
				vu := t.spawnVirtualUser(tp.TestCase, wg)
				vus = append(vus, vu)
			}
		}
	}
	return vus
}

func (t *TestRunner) spawnVUsEvenly(tp models.TestPhase, wg *sync.WaitGroup) []*VirtualUser {
	ticker := time.NewTicker(time.Duration(1/tp.NumVirtualUsers) * time.Second)
	done := createTestDurationChannel(tp.DurationSeconds)
	var vus []*VirtualUser
	for {
		select {
		case <-done:
			ticker.Stop()
			fmt.Printf("Test Phase %s is done\n", tp.Name)
			return vus
		case <-ticker.C:
			wg.Add(1)
			vu := t.spawnVirtualUser(tp.TestCase, wg)
			vus = append(vus, vu)
		}
	}
	return vus
}

func (t *TestRunner) spawnVirtualUser(tc models.TestCaseGetRequest, wg *sync.WaitGroup) *VirtualUser {
	client := t.clientFactory() // in order to simulate a real user, we should create a new client for each virtual user
	vu := NewVirtualUser(client)
	go func() {
		defer wg.Done()
		vu.ExecuteTestCase(tc)
	}()
	return vu
}

func calculateOverallMetrics(phaseMetrics []models.VirtualUserAggregatedMetrics) *models.TestRunAggregatedMetrics {
	overallMetrics := &models.TestRunAggregatedMetrics{}
	durations := make([]float64, 0, len(phaseMetrics))
	numVUsSucceeded := make([]float64, 0, len(phaseMetrics))
	for _, metrics := range phaseMetrics {
		overallMetrics.NumVUsCreated += metrics.NumVUsCreated
		overallMetrics.NumVUsSucceeded += metrics.NumVUsSucceeded
		overallMetrics.NumVUsFailed += metrics.NumVUsFailed
		if overallMetrics.MaxVUDuration < metrics.MaxVUDuration {
			overallMetrics.MaxVUDuration = metrics.MaxVUDuration
		}
		numVUsSucceeded = append(numVUsSucceeded, float64(metrics.NumVUsSucceeded))
		durations = append(durations, float64(metrics.AvgVUDuration))
	}
	overallMetrics.AvgVUDuration = int(stat.Mean(durations, numVUsSucceeded))
	return overallMetrics
}

func calculatePhaseMetrics(virtualUsers []*VirtualUser) models.VirtualUserAggregatedMetrics {
	totalNumVUs := len(virtualUsers)
	successfulVUs := 0
	failedVUs := 0
	executionDurations := make([]float64, 0, len(virtualUsers))
	for _, vu := range virtualUsers {
		if vu.status == models.CompletedWithSuccess {
			successfulVUs++
		}
		if vu.status == models.CompletedWithFailure {
			failedVUs++
		}
		executionDurations = append(executionDurations, float64(vu.duration))
	}
	// Sort the durations slice
	sort.Slice(executionDurations, func(i, j int) bool {
		return executionDurations[i] < executionDurations[j]
	})
	return models.VirtualUserAggregatedMetrics{
		NumVUsCreated:   totalNumVUs,
		NumVUsSucceeded: successfulVUs,
		NumVUsFailed:    failedVUs,
		AvgVUDuration:   int(stat.Mean(executionDurations, nil)),
		P50VUDuration:   int(stat.Quantile(0.50, stat.Empirical, executionDurations, nil)),
		P95VUDuration:   int(stat.Quantile(0.95, stat.Empirical, executionDurations, nil)),
		P99VUDuration:   int(stat.Quantile(0.99, stat.Empirical, executionDurations, nil)),
		MaxVUDuration:   int(executionDurations[len(executionDurations)-1]),
	}
}

func createTestDurationChannel(duration time.Duration) chan bool {
	done := make(chan bool) // will signal when the test is done
	go func() {
		// Wait for the test duration to complete
		time.Sleep(duration * time.Second)
		// Signal that the test is done
		done <- true
	}()
	return done
}

func printPhaseMetrics(name string, metrics models.VirtualUserAggregatedMetrics) {
	fmt.Printf("Phase %s metrics: %s\n", name, metrics.String())
}
