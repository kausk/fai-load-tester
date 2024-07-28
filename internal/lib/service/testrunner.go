package service

import (
	"time"

	"loadtester/internal/lib/models"
)

type TestRunner struct {
	StartTime time.Time
	TestPlan  models.TestPlan
	// EndTime   time.Time // needed for time duration calculations
	// Results   []TestResult
}

func NewTestRunner(tp models.TestPlan) *TestRunner {
	return &TestRunner{
		StartTime: time.Now(),
		TestPlan:  tp,
	}
}

func (t *TestRunner) Start() {
	// initialize test session
	// create workers
	// start workers
	// wait for workers to finish
	// collect results
	// end test session
}

func (t TestRunner) Stop() {

}

func (t TestRunner) Status() {

}

func (t TestRunner) Results() {

}
