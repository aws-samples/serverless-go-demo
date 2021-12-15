package types

import "context"

type Event struct {
	Source     string
	Detail     string
	DetailType string
	Resources  []string
}

type FailedEvent struct {
	Event
	FailureCode    string
	FailureMessage string
}

type Bus interface {
	Put(context.Context, []Event) ([]FailedEvent, error)
}
