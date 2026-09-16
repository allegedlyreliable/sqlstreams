package sqlstreams

import (
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/scheduler"
)

// please GOPLS make aliases and go doc comments work better

// SchedulerConfig is a schedule's declared delivery semantics, stored on
// its row: how each message it produces runs.
type SchedulerConfig struct {
	// Timeout - how long one message's delivery may run.
	// Default: 30s.
	Timeout time.Duration

	// Concurrency - whether a message runs while a previous one is still
	// running (parallel) or waits for it (exclusive).
	// Default: parallel.
	Concurrency ConcurrencyPolicy

	// Metadata - marshaled to opaque JSON stored on the row and shown by
	// `sqlstreams scheduler get`; it is not part of the produced message.
	// Default: nil, stored as {}.
	Metadata any
}

func (c *SchedulerConfig) WithDefaults() *SchedulerConfig {
	return (*SchedulerConfig)((*scheduler.SchedulerConfig)(c).WithDefaults())
}

func (c *SchedulerConfig) Validate() error {
	return (*scheduler.SchedulerConfig)(c).Validate()
}
