package models

import "fmt"

// VirtualUserAggregatedMetrics are collected in 1 phase
type VirtualUserAggregatedMetrics struct {
	NumVUsCreated   int
	NumVUsSucceeded int
	NumVUsFailed    int
	AvgVUDuration   int
	P50VUDuration   int
	P95VUDuration   int
	P99VUDuration   int
	MaxVUDuration   int
}

func (v VirtualUserAggregatedMetrics) String() string {
	return fmt.Sprintf("NumVUsCreated: %d, NumVUsSucceeded: %d, NumVUsFailed: %d, VUSuccessPercentage: %f%%, AvgVUDuration: %dms, P50VUDuration: %dms, P95VUDuration: %dms, P99VUDuration: %dms, MaxVUDuration: %dms\n", v.NumVUsCreated, v.NumVUsSucceeded, v.NumVUsFailed, float64(100*v.NumVUsSucceeded/v.NumVUsCreated), v.AvgVUDuration, v.P50VUDuration, v.P95VUDuration, v.P99VUDuration, v.MaxVUDuration)
}

// TestRunAggregatedMetrics are averaged across all phases
type TestRunAggregatedMetrics struct {
	NumVUsCreated   int
	NumVUsSucceeded int
	NumVUsFailed    int
	AvgVUDuration   int
	MaxVUDuration   int
}

func (t TestRunAggregatedMetrics) String() string {
	return fmt.Sprintf("NumVUsCreated: %dms, NumVUsSucceeded: %dms, NumVUsFailed: %dms, VUSuccessPercentage: %f%%, AvgVUDuration: %dms, MaxVUDuration: %dms\n", t.NumVUsCreated, t.NumVUsSucceeded, t.NumVUsFailed, float64(100*t.NumVUsFailed/t.NumVUsCreated), t.AvgVUDuration, t.MaxVUDuration)
}
