package sqlstreams

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
)

// please GOPLS make aliases and go doc comments work better

// ProduceOptions holds per-message knobs that are optional and rarely set --
// the zero value means "neither is set," so a caller who doesn't need them
// never has to name them.
type ProduceOptions struct {
	// RoutingKey - matched against a consumer group's bindings to decide
	// whether that group receives this message at all.
	// Default: "" (no routing key; a keyless message matches no binding, so
	// only groups with no bindings receive it).
	//
	// "" is stored as no routing key, not an empty-string match.
	// Ex: "orders.created", "billing.invoice.paid"
	RoutingKey string

	// MessageKey - the entity this message is about. On its own it is stored
	// and nothing more; Compaction and ConcurrencyExclusive both read it.
	// Default: "" (no key).
	//
	// Ex: "user:123", "acct-42", a device serial.
	MessageKey string

	// Compaction - opts this message into log compaction under its
	// MessageKey: it becomes one version of the key, and claims only ever
	// return the key's latest version, not every version ever written.
	// Set Enable: true to opt in; Rank: 0 uses message-id order.
	// Default: nil (not compacted; delivered independently, never superseded).
	//
	// A hot key caps batched throughput: same-key batches commit one after
	// another, and adding producer processes makes a hot key slower, not faster.
	Compaction *CompactionOptions

	// IdempotencyKey - protects a retried Produce (after a blip) from double-publishing.
	// Default: "" (a fresh key is generated per call, protecting only
	// against retries within that one call).
	//
	// Supply your own for protection across your OWN retries too -- e.g. your
	// process crashes and restarts before learning whether a publish landed,
	// and you call Produce again with the same key. Any stable string works:
	// one that parses as a UUID is stored verbatim, anything else is hashed
	// to a deterministic UUID first, so the same string always dedups.
	// A caller-supplied key routes the call to a per-call transaction, never a batch.
	//
	// Minting keys yourself on a hot path: prefer time-ordered UUIDv7 strings --
	// random-shaped keys (a hash, a v4) cost extra WAL once the claim table
	// holds millions of unexpired rows.
	// Ex: "evt_9f2c" from a webhook; a UUIDv7 persisted alongside the work.
	IdempotencyKey string

	// Message - per-message MessageOptions: what this message REQUESTS from
	// whoever consumes it (work timeout, redelivery policy, concurrency).
	// Default: nil (defaults to Producer Defaults > Consumer Defaults).
	Message *MessageOptions
}

// Validate rejects nonsensical option combinations.
// Must be called after Fill().
func (o ProduceOptions) Validate() error {
	return produce.ProduceOptions(o).Validate()
}
