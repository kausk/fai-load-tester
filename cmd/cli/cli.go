package cli

import (
	"flag"
	"fmt"
	"time"

	"loadtester/internal/lib/models"
	"loadtester/internal/lib/service"
)

func main() {
	qps := flag.Int("qps", 10, "Queries per second")
	duration := flag.Int("duration", 10, "Duration of the test in seconds")
	url := flag.String("url", "http://google.com", "URL to test")
	virtualUsers := flag.Int("virtualUsers", 10, "Number of virtual users. Under 100 QPS, 1 virtual user = 1 query per second. Above 100 QPS, the queries will be distributed evenly across the virtual users.")
	flag.Parse()

	var numIterations int
	var numVirtualUsers int
	if *qps < 100 {
		numIterations = 1
		numVirtualUsers = *qps
	} else {
		numVirtualUsers = *virtualUsers
		numIterations = *qps / numVirtualUsers
	}
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
				DurationSeconds: time.Duration(*duration) * time.Second,
			},
		},
	}
	tr := service.NewTestRunner(tp)
	tr.Start()
}
