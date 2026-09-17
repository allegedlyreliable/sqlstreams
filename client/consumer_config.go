package sqlstreams

import (
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/consumer"
)

// please GOPLS make aliases and go doc comments work better

// ConsumerConfig is the group's declaration: what the group means, identical
// for every instance of the group. Session settings -- how one process runs
// -- live on ConsumeOptions at Consume.
type ConsumerConfig struct {
	// Message - default MessageOptions: fills any option the produced message left unset.
	// Default: Timeout 30s; Retry MaxRetries 3 with the default curve.
	Message *MessageOptions

	// MessageMin - per-option floors: raises any resolved option below these.
	// Concurrency is not orderable and must stay unset.
	// Default: nil (no floors).
	MessageMin *MessageOptions

	// MessageMax - per-option ceilings: lowers any resolved option above these.
	// Concurrency is not orderable and must stay unset.
	// Default: Message's values -- messages cannot request above the group's
	// defaults unless raised here.
	MessageMax *MessageOptions

	// ConcurrencyOverride - this group runs every message under this policy,
	// beating whatever the message requested.
	// Default: "" (honor each message's own policy).
	ConcurrencyOverride ConcurrencyPolicy

	// Start - where a group's cursor is placed when Register creates it;
	// a group that already has a cursor row keeps its position.
	// Default: Beginning() -- the oldest retained message.
	Start CursorPosition

	// Bindings - the group's whole pattern set, declared on every Register.
	// Default: nil (the whole stream).
	Bindings []string

	ExceptionInitialBackoff time.Duration // can_run_after delay after the first failed delivery, including an initially deferred message -- Message.Retry takes over on later retries. Default: 5s.
	MaxRangeReclaims        int           // past this many reclaims a range is POISON -- quarantined into the exception window instead of handed out again. Default: 3.
}

func (c *ConsumerConfig) WithDefaults() *ConsumerConfig {
	return (*ConsumerConfig)((*consumer.ConsumerConfig)(c).WithDefaults())
}

func (c *ConsumerConfig) Validate() error {
	return (*consumer.ConsumerConfig)(c).Validate()
}

// DeepCopy returns a copy sharing nothing mutable with the receiver --
// WithDefaults fills the message options in place, so a captured
// declaration must not alias them.
func (c *ConsumerConfig) DeepCopy() *ConsumerConfig {
	return (*ConsumerConfig)((*consumer.ConsumerConfig)(c).DeepCopy())
}
