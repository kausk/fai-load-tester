package models

import "time"

// TestPlan contains a sequence of TestPhases
type TestPlan struct {
	Name                    string
	Phases                  []TestPhase
	SpawnVUsEvenly          bool // if true, the test runner will spawn the virtual users evenly across each second duration of the test phase
	HTTPTimeoutMilliseconds time.Duration
}

// TestPhase represents one phase of a TestPlan.
// Contains a TestCaseGetRequest and a number of virtual users executing that TestCaseGetRequest
type TestPhase struct {
	Name            string
	NumVirtualUsers int // number of virtual users executing the TestCaseGetRequest that will be spawned per second
	TestCase        TestCaseGetRequest
	DurationSeconds time.Duration // duration of the test phase
	Metrics         VirtualUserAggregatedMetrics
}

// TestCaseGetRequest represents a single test case to be executed by a virtual user
// In the future, a test case can execute a user-defined closure to determine success/failure
// In the future, a test case can also execute POST/PUT/DELETE requests, and pass in custom headers and payloads
// For now, we will just check that the status code is 2xx
type TestCaseGetRequest struct {
	URL           string
	NumIterations int // number of times the virtual user will repeat the request
}
