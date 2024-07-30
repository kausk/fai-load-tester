package main

import (
	"fmt"
	"time"

	"loadtester/internal/lib/models"
	"loadtester/internal/lib/service"
)

func main() {
	loadTest()
}

func loadTest() {
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
	fmt.Println(tr.Results())
}
