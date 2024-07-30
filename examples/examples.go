/*
This file provides an example of how to use the load tester library.
It creates a test plan with two test phases and runs the test runner in order to execute the load test.
*/
package main

import (
	"fmt"
	"time"

	"loadtester/pkg/models"
	"loadtester/pkg/service"
)

func main() {
	loadTestQPS()
}

func loadTestQPS() {
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
	tr := service.NewTestRunner(tp)
	tr.Start()
	time.Sleep(10 * time.Second)
	fmt.Println(tr.OverallResults().String())
}

// TODO: add more examples
