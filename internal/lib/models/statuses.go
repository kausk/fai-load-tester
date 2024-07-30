package models

type Status int

const (
	WaitingForStart Status = iota
	Running
	Finished
	CompletedWithFailure
	CompletedWithSuccess
)
