# fai-load-tester

# Overview

This loadtester generates load by spawning a specified number of Virtual Users (VUs) each second over a given duration, which each make a specified number of HTTP requests. A desired QPS is achieved by dividing the QPS into VUs and iterations per VU.

The load test is specified using `models.TestPlan`, and this is passed into `service.TestRunner` which executes the test.

Metrics are collected and stored in `service.TestRunner.MetricsByPhase` and `service.TestRunner.OverallMetrics`.



# Usage Instructions

## CLI Documentation
```
ksk@Kaushiks-MacBook-Pro-2 ~/d/loadtester (main)> docker run --network host --rm loadtester --help
Usage of ./cli:
  -duration int
    	Duration of the test in seconds (default 10)
  -qps int
    	Queries per second (default 10)
  -spawnEvenly
    	If true, the test runner will spawn the virtual users evenly across each second duration of the test phase
  -url string
    	URL to test (default "http://google.com")
  -virtualUsers int
    	Number of virtual users. Under 100 QPS, 1 virtual user = 1 query per second. Above 100 QPS, the queries will be distributed evenly across the virtual users. Note that the actual QPS might be slightly under what is specified in -qps due to rounding down. (default 10)
```

## Docker CLI Instructions
```=
docker build -t loadtester .
docker run --network host --rm loadtester -url http://host.docker.internal:8081 -qps 300 ## accessing localhost URL
docker run --network host --rm loadtester -url http://yahoo.com -qps 50 ## accessing external URL
```

# API usage example available in examples/examples.go


# Running Unit Tests
```
~/d/loadtester (main) [1]> go test ./...
?   	loadtester/cmd/cli	[no test files]
?   	loadtester/examples	[no test files]
?   	loadtester/pkg/models	[no test files]
?   	loadtester/pkg/service	[no test files]
ok  	loadtester/pkg/service/tests	17.019s
```
