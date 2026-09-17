package producer

import (
	"errors"

	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
)

// ProduceItem is one message plus its options -- the unit ProduceBatch takes.
type ProduceItem[Message common.Versioned] struct {
	Message *Message
	Options produce.ProduceOptions
}

// NewProduceItem pairs one message with its options for ProduceBatch.
// Caller-supplied idempotency keys are rejected because contention on one
// key would hold up the whole batch. Use Produce for those messages.
// options may be nil for the defaults.
func NewProduceItem[Message common.Versioned](message *Message, options *produce.ProduceOptions) (*ProduceItem[Message], error) {
	if message == nil {
		return nil, errors.New("message must not be nil")
	}

	resolved := produce.ProduceOptions{}
	if options != nil {
		resolved = *options
	}
	if resolved.IdempotencyKey != "" {
		return nil, errors.New("IdempotencyKey is not supported in a batch -- produce keyed messages individually")
	}

	return &ProduceItem[Message]{
		Message: message,
		Options: resolved,
	}, nil
}
