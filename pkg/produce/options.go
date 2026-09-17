package produce

import (
	"errors"
	"fmt"

	"github.com/allegedlyreliable/sqlstreams/pkg/common"
)

// ProduceOptions sets optional keys, compaction, and processing options for
// one message. The zero value uses producer and consumer defaults.
type ProduceOptions struct {
	// RoutingKey - matched against a consumer group's bindings to decide
	// whether that group receives this message at all.
	// Default: "" (no routing key). A message without a routing key matches
	// no binding, so only groups with no bindings receive it.
	//
	// "" is stored as no routing key, not an empty-string match.
	// Ex: "orders.created", "billing.invoice.paid"
	RoutingKey string

	// MessageKey - identifies related messages for compaction and the
	// exclusive or ordered concurrency policies. A key alone enables neither.
	// Default: "" (no key).
	//
	// Ex: "user:123", "acct-42", a device serial.
	MessageKey string

	// Compaction - opts this message into log compaction under its
	// MessageKey: it becomes one version of the key, and claims only ever
	// return the key's latest version, not every version ever written.
	// Set Enable: true to opt in. Higher schema versions win first, then
	// higher ranks, then higher message ids.
	// Default: nil (not compacted, never superseded).
	//
	// A hot key caps batched throughput: same-key batches commit one after
	// another, and adding producer processes makes a hot key slower, not faster.
	Compaction *CompactionOptions

	// IdempotencyKey - deduplicates message inserts within this stream while
	// the key is retained under StreamConfig.IdempotencyKeyTTL.
	// Default: "" (a fresh key is generated per call, protecting only
	// against retries within that one call).
	//
	// Supply a stable key for protection across calls and process restarts.
	// A string that parses as a UUID is stored as that UUID. Other strings
	// are hashed to a deterministic UUID. Deduplication does not compare
	// payloads or deduplicate application writes in a callback.
	// A caller-supplied key routes the call to a per-call transaction, never a batch.
	//
	// Minting keys yourself on a hot path: prefer time-ordered UUIDv7 strings --
	// random-shaped keys (a hash, a v4) cost extra WAL once the claim table
	// holds millions of unexpired rows.
	// Ex: "evt_9f2c" from a webhook, or a UUIDv7 stored alongside the work.
	IdempotencyKey string

	// Message - requested processing timeout, retry policy, and concurrency.
	// Unset fields inherit producer defaults, then consumer defaults.
	// Consumer group bounds and concurrency overrides still apply.
	// Default: nil.
	Message *common.MessageOptions
}

// Validate checks message options, compaction, and required message keys.
// It does not fill defaults.
func (o ProduceOptions) Validate() error {
	if o.Message != nil && o.Message.Concurrency.HoldsKey() && o.MessageKey == "" {
		return fmt.Errorf("Concurrency %q set without MessageKey -- it runs deliveries one at a time per key, set ProduceOptions.MessageKey", o.Message.Concurrency)
	}
	if o.Message != nil && o.Message.Concurrency == common.ConcurrencyOrdered && o.Compaction != nil && o.Compaction.Enable {
		return errors.New("Concurrency 'ordered' set with Compaction enabled -- compaction supersedes older versions, ordered delivers every one; choose one")
	}
	if o.Compaction != nil && o.Compaction.Enable && o.MessageKey == "" {
		return errors.New("Compaction enabled without MessageKey -- compaction picks a winner per key, set ProduceOptions.MessageKey")
	}
	if err := o.Compaction.Validate(); err != nil {
		return fmt.Errorf("Compaction: %w", err)
	}
	if err := o.Message.Validate(); err != nil {
		return fmt.Errorf("Message: %w", err)
	}
	return nil
}

// CompactionOptions is ProduceOptions.Compaction: whether the message
// compacts under its MessageKey, and the rank that decides the key's winner.
type CompactionOptions struct {
	// Enable - opts the message into compaction. A nil Compaction and Enable
	// false mean the same thing: not compacted.
	Enable bool

	// Rank - selects the winner within one payload schema version. Higher
	// ranks win, with higher message ids breaking ties. A higher schema
	// version takes precedence over rank. Keeping every rank at zero uses
	// message-id order within each schema version.
	//
	// A high rank prevents lower-ranked messages at the same schema version
	// from becoming the head, even when they arrive later.
	// Ex: a source system's row version, a priority tier, epoch micros.
	Rank int64
}

// Validate tolerates a nil receiver -- nil means not compacted.
func (o *CompactionOptions) Validate() error {
	if o == nil {
		return nil
	}
	if !o.Enable && o.Rank != 0 {
		return fmt.Errorf("Rank must be 0 when Enable is false, got %d -- set Enable to true to use a rank", o.Rank)
	}
	return nil
}
