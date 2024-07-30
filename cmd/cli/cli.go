/*
This file exposes the CLI interface for the load tester. It takes in the following flags:
- qps: queries per second
- duration: duration of the test in seconds
- url: URL to test
- virtualUsers: number of virtual users. Under 100 QPS, 1 virtual user = 1 query per second. Above 100 QPS, the queries will be distributed evenly across the virtual users.
It then creates a test plan and starts the test runner.

Usage:
go run cmd/cli/cli.go -qps 10 -duration 10 -url http://google.com -virtualUsers 10
*/
package main

import (
	"flag"
	"fmt"
	"time"

	"loadtester/pkg/models"
	"loadtester/pkg/service"
)

func main() {
	qps := flag.Int("qps", 10, "Queries per second")
	duration := flag.Int("duration", 10, "Duration of the test in seconds")
	url := flag.String("url", "http://google.com", "URL to test")
	virtualUsers := flag.Int("virtualUsers", 10, "Number of virtual users. Under 100 QPS, 1 virtual user = 1 query per second. Above 100 QPS, the queries will be distributed evenly across the virtual users. Note that the actual QPS might be slightly under what is specified in -qps due to rounding down.")
	spawnEvenly := flag.Bool("spawnEvenly", false, "If true, the test runner will spawn the virtual users evenly across each second duration of the test phase")
	// TODO: add flag for HTTP timeout
	// TODO: add ability to specify multiple phases through CLI
	flag.Parse()

	var numIterations int
	var numVirtualUsers int
	if *virtualUsers == 0 {
		fmt.Println("Number of virtual users must be greater than 0")
		return
	}
	if *qps < 100 {
		fmt.Println("Ignoring virtualUsers flag since QPS is less than 20")
		numIterations = 1
		numVirtualUsers = *qps
	} else {
		numVirtualUsers = *virtualUsers
		numIterations = *qps / numVirtualUsers // TODO: this will round down the QPS to the nearest //10, fix this by adding the remaining QPS to the last virtual user
	}
	fmt.Printf("Executing load test with %d QPS over %d seconds duration: %d virtual users (VUs) spawned per second, %d requests per VU\n", *qps, *duration, numVirtualUsers, numIterations)
	tc := models.TestCaseGetRequest{
		URL:           *url,
		NumIterations: numIterations,
	}
	tp := models.TestPlan{
		Name: fmt.Sprintf("%d QPS Test", *qps),
		Phases: []models.TestPhase{
			{
				Name:            fmt.Sprintf("%d QPS", *qps),
				NumVirtualUsers: numVirtualUsers,
				TestCase:        tc,
				DurationSeconds: time.Duration(*duration),
			},
		},
		SpawnVUsEvenly: *spawnEvenly,
	}
	tr := service.NewTestRunner(tp)
	tr.Start()
}
