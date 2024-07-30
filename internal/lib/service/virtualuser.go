package service

import (
	"fmt"
	"time"

	"loadtester/internal/lib/models"
)

type VirtualUser struct {
	status   models.Status
	duration time.Duration
	client   HTTPClientInterface
}

func NewVirtualUser(client HTTPClientInterface) *VirtualUser {
	return &VirtualUser{
		status: models.WaitingForStart,
		client: client,
	}
}

// ExecuteTestCase executes a test case and saves the result to the struct
func (v *VirtualUser) ExecuteTestCase(tc models.TestCaseGetRequest) {
	start := time.Now()
	for i := 0; i < tc.NumIterations; i++ {
		resp, err := v.client.Get(tc.URL)
		if err != nil {
			fmt.Printf("Encountered error when requesting %s: %s\n", tc.URL, err)
			v.status = models.CompletedWithFailure
			v.duration = time.Now().Sub(start) / time.Millisecond
			return
		}
		success := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !success {
			fmt.Printf("Received non-2xx status code %d when requesting %s\n", resp.StatusCode, tc.URL)
			v.status = models.CompletedWithFailure
			v.duration = time.Now().Sub(start) / time.Millisecond
			return
		}
	}
	v.duration = time.Now().Sub(start) / time.Millisecond
	v.status = models.CompletedWithSuccess
}

// TODO: interrupt test case execution if it takes it
