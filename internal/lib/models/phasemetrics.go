package models

// VirtualUserAggregatedMetrics can be aggregated across a phase or an entire test runß
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
