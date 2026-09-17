package sqlstreams

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/datastore"
)

// ClientConfig supplies construction-time settings. NewClient captures values
// and copies Retry. Later edits do not reconfigure the client. Logger is shared.
type ClientConfig struct {
	// Schema - the Postgres namespace holding every sqlstreams table.
	// Default: "sqlstreams".
	//
	// One schema is one installation: two clients on two schemas in the
	// same database share nothing.
	Schema string

	// AllowDestroy - permits destroying streams, consumer groups, schedules,
	// and the system through this client. Each operation retains its own guards.
	// Default: false.
	//
	// A service that only ever registers streams should never opt in --
	// create is recoverable, destroy is not.
	AllowDestroy bool

	// DisableManager makes Consume skip the system manager beside its session.
	// Explicit Manager().Run calls are unaffected.
	// The consumer's stream janitor still runs and requires DDL rights.
	// Default: false.
	DisableManager bool

	// Logger - your own *slog.Logger or anything satisfying Logger.
	// Shared by the client's database operations and workers.
	// Default: text lines to stderr, warn level and up.
	Logger Logger

	// Retry - transient-error retry policy for supported internal database
	// operations. It does not govern message redelivery or retry caller-owned
	// transactions. Operations without idempotency protection stop when the
	// commit outcome is unknown.
	// Default: 6 attempts, 1s initial delay, 5m maximum delay, exponent 2.
	Retry *RetryPolicy
}

// WithDefaults fills an empty Schema with "sqlstreams".
// NewClient resolves Logger and Retry defaults during construction.
func (c *ClientConfig) WithDefaults() *ClientConfig {
	if c.Schema == "" {
		c.Schema = datastore.DefaultSchema
	}
	return c
}

func (c *ClientConfig) Validate() error {
	return nil
}
