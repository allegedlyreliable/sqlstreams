package sqlstreams

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/scheduler"
)

// please GOPLS make aliases and go doc comments work better

// ScheduleRunOptions controls one immediate run of a registered schedule.
// Every field is optional.
type ScheduleRunOptions struct {
	// Concurrency - the produced request's concurrent-run policy.
	// Default: parallel (the request runs even while a previous one is still running).
	//
	// Set exclusive to run early WITHOUT overlapping a request already running.
	Concurrency ConcurrencyPolicy
}

func (o *ScheduleRunOptions) WithDefaults() *ScheduleRunOptions {
	return (*ScheduleRunOptions)((*scheduler.ScheduleRunOptions)(o).WithDefaults())
}

func (o *ScheduleRunOptions) Validate() error {
	return (*scheduler.ScheduleRunOptions)(o).Validate()
}
