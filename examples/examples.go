/*
This file provides an example of how to use the load tester library.
It creates a test plan with two test phases and runs the test runner in order to execute the load test.
*/
package main

import (
	"fmt"
	"loadtester/pkg/models"
	"loadtester/pkg/service"
)

func main() {
	loadTestQPS()
}

func loadTestQPS() {
	// Create a test plan with phases and a test case
	tcOneIteration := models.TestCaseGetRequest{
		URL:           "http://google.com",
		NumIterations: 1,
	}
	tcFiveIterations := models.TestCaseGetRequest{
		URL:           "http://google.com",
		NumIterations: 5,
	}
	tp := models.TestPlan{
		Name: "QPS Test",
		Phases: []models.TestPhase{
			{
				Name:            "10 QPS",
				NumVirtualUsers: 10,
				TestCase:        tcOneIteration,
				DurationSeconds: 5,
			},
			{
				Name:            "20 QPS",
				NumVirtualUsers: 4,
				TestCase:        tcFiveIterations,
				DurationSeconds: 5,
			},
		},
	}
	// Create a test runner and start the test
	tr := service.NewTestRunner(tp)
	tr.Start()
	<-tr.Done
	// Print the overall results
	// Per-phase results also available in tr.ResultsByPhase()
	fmt.Println(tr.OverallResults().String())
}

// TODO: add more examples
