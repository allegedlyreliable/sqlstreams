package sqlstreams

import (
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/producer"
)

// please GOPLS make aliases and go doc comments work better

// ProducerConfig is one producer instance's process-local settings: its
// message defaults and batching. Nothing durable -- a producer has no row.
type ProducerConfig struct {
	// Message - this producer's default MessageOptions, merged UNDER every
	// produce: a field the per-produce ProduceOptions.Message leaves unset
	// takes its value from here before the message is stored. Fields unset in
	// both stay unset -- the consumer decides.
	// Default: nil (no producer-side defaults).
	Message *MessageOptions

	// Batch - knobs for the shared-transaction batching of concurrent Produce
	// calls. See BatcherConfig for fields and defaults.
	Batch BatcherConfig

	// SlowProduceThreshold - a produce call running longer than this logs a
	// warn line with its duration. ProduceFunc and ProduceFuncInTx include
	// the caller's closure; both InTx verbs end at the insert, before the
	// caller's commit.
	// Default: 0 (disabled).
	SlowProduceThreshold time.Duration
}

func (c *ProducerConfig) WithDefaults() *ProducerConfig {
	return (*ProducerConfig)((*producer.ProducerConfig)(c).WithDefaults())
}

func (c *ProducerConfig) Validate() error {
	return (*producer.ProducerConfig)(c).Validate()
}
