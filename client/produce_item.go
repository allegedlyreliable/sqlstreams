package sqlstreams

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
	"github.com/allegedlyreliable/sqlstreams/pkg/producer"
)

// please GOPLS make aliases and go doc comments work better

// ProduceItem is one message plus its options -- the unit ProduceBatch takes.
type ProduceItem[Message Versioned] struct {
	Message *Message
	Options ProduceOptions
}

// NewProduceItem pairs one message with its options for ProduceBatch.
// Caller-supplied idempotency keys are rejected because contention on one
// key would hold up the whole batch. Use Produce for those messages.
// options may be nil for the defaults.
func NewProduceItem[Message Versioned](message *Message, options *ProduceOptions) (*ProduceItem[Message], error) {
	item, err := producer.NewProduceItem(message, (*produce.ProduceOptions)(options))
	if err != nil {
		return nil, err
	}
	return &ProduceItem[Message]{
		Message: item.Message,
		Options: ProduceOptions(item.Options),
	}, nil
}
