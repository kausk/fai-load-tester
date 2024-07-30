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
	vSP := "N/A"
	if v.NumVUsCreated > 0 {
		vSP = fmt.Sprintf("%f%%", 100*float64(v.NumVUsSucceeded)/float64(v.NumVUsCreated))
	}
	return fmt.Sprintf("NumVUsCreated: %d, NumVUsSucceeded: %d, NumVUsFailed: %d, VUSuccessPercentage: %s, AvgVUDuration: %dms, P50VUDuration: %dms, P95VUDuration: %dms, P99VUDuration: %dms, MaxVUDuration: %dms\n", v.NumVUsCreated, v.NumVUsSucceeded, v.NumVUsFailed, vSP, v.AvgVUDuration, v.P50VUDuration, v.P95VUDuration, v.P99VUDuration, v.MaxVUDuration)
}

// OverallMetrics are averaged across all phases
type OverallMetrics struct {
	NumVUsCreated   int
	NumVUsSucceeded int
	NumVUsFailed    int
	AvgVUDuration   int
	MaxVUDuration   int
}

func (t OverallMetrics) String() string {
	vSP := "N/A"
	if t.NumVUsCreated > 0 {
		vSP = fmt.Sprintf("%f%%", 100*float64(t.NumVUsSucceeded)/float64(t.NumVUsCreated))
	}
	return fmt.Sprintf("Overall Metrics: NumVUsCreated: %d, NumVUsSucceeded: %d, NumVUsFailed: %d, VUSuccessPercentage: %s, AvgVUDuration: %dms, MaxVUDuration: %dms\n", t.NumVUsCreated, t.NumVUsSucceeded, t.NumVUsFailed, vSP, t.AvgVUDuration, t.MaxVUDuration)
}
