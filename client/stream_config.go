package sqlstreams

import (
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/stream"
)

// please GOPLS make aliases and go doc comments work better

// StreamConfig is Register's spec -- separate from Stream so Register can grow
// (retention, etc.) without a signature change.
type StreamConfig struct {
	// PartitionSize - message-id interval per partition. Gaps in allocated
	// ids mean a partition can contain fewer rows.
	// Default: 1_000_000.
	//
	// Lower values give finer-grained retention drops at the cost of more
	// partitions to maintain. Tune down for low-volume streams, up for
	// high-throughput ones.
	// Ex: 10_000 for a low-volume audit stream, 5_000_000 for high-throughput ingest.
	PartitionSize int64

	// RetentionTTL - how long a message survives before the janitor may drop
	// or sweep it.
	// Default: 0 (keep every message indefinitely).
	//
	// Set this once a stream has a real expiry requirement.
	// Ex: 30 * 24 * time.Hour for a 30-day event stream.
	RetentionTTL time.Duration

	// AllowDropPastCommitted - if true, retention can drop data a lagging
	// consumer group hasn't committed yet (Kafka's default behavior).
	// Default: false.
	//
	// Set true only if a badly-lagging consumer should lose data rather than
	// block cleanup.
	// Ex: true for a metrics stream where staleness beats unbounded disk growth.
	AllowDropPastCommitted bool

	// IdempotencyKeyTTL - how long a produce-retry claim survives in
	// idempotency_key before the janitor sweeps it.
	// Default: 24h.
	//
	// Zero means the default, not "forever" -- WithDefaults resolves it
	// before the stream is ever registered. TTL only needs to cover your retry horizon,
	// not a retention window: 24h covers a webhook provider's retry day.
	// Every produce writes a claim row, minted key or not, so the TTL is
	// the claim table's size -- lower it for a stream whose producers never
	// retry across a restart.
	// Ex: 10 * time.Minute.
	IdempotencyKeyTTL time.Duration

	// EmptyCompactionHeadTTL - how long a compaction-head row with no current
	// head may stay idle before the stream janitor sweeps it.
	// Default: 1h.
	//
	// Zero means the default, not "forever" -- WithDefaults resolves it
	// before the stream is ever registered. Locking an empty head refreshes its activity;
	// the TTL never applies to a row that points at a head.
	EmptyCompactionHeadTTL time.Duration

	// DeliveryLogMode - which delivery outcomes write to delivery_log_<id>, the
	// per-attempt audit trail.
	// Default: DeliveryLogModeFailures (every outcome except success).
	//
	// DeliveryLogModeOff for a stream whose failure volume would make the extra
	// per-attempt write not worth paying for.
	// DeliveryLogModeAll when successes must be auditable per message
	// each success txn then also writes its 'success' row.
	DeliveryLogMode DeliveryLogMode

	// Janitor - retention cleanup settings.
	// Default: its own defaults.
	Janitor *JanitorConfig

	// Vacuum - key-table vacuum settings. New streams start with vacuum suspended.
	// Default: its own defaults.
	Vacuum *VacuumConfig
}

func (c *StreamConfig) WithDefaults() *StreamConfig {
	return (*StreamConfig)((*stream.StreamConfig)(c).WithDefaults())
}

func (c *StreamConfig) Validate() error {
	return (*stream.StreamConfig)(c).Validate()
}

// ToStream builds the stream row a registration writes from the resolved config.
func (c *StreamConfig) ToStream(id int64, systemId int64, name string) *Stream {
	return (*stream.StreamConfig)(c).ToStream(id, systemId, name)
}
