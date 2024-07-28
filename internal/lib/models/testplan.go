package models

// TestPlan contains a sequence of TestPhases
type TestPlan struct {
	Name           string
	Phases         []TestPhase
	SpawnVUsEvenly bool // if true, the test runner will spawn the virtual users evenly across each second duration of the test phase
}

// TestPhase represents one phase of a TestPlan.
// Contains a TestCase and a number of virtual users executing that TestCase
type TestPhase struct {
	Name            string
	NumVirtualUsers int // number of virtual users executing the TestCase that will be spawned per second
	TestPlan        TestPlan
}

// TestCase represents a single test case to be executed by a virtual user
// In the future, a test case can excecute a user-defined closure to determine success/failure
// For now, we will just check that the status code is 2xx
type TestCase struct {
	URL             string
	DurationSeconds int
	Method          string // only GET supported for now
	NumIterations   int    // number of times to repeat the request
}
