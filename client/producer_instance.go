package sqlstreams

import (
	"context"
	"errors"

	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
	"github.com/allegedlyreliable/sqlstreams/pkg/producer"
)

// ProducerInstance is a registered producer: it appends messages to the
// stream its ProducerHandle.Register resolved.
type ProducerInstance[Message Versioned] struct {
	instance *producer.ProducerInstance[Message]
}

func newProducerInstance[Message Versioned](instance *producer.ProducerInstance[Message]) (*ProducerInstance[Message], error) {
	if instance == nil {
		return nil, errors.New("instance must not be nil")
	}
	return &ProducerInstance[Message]{instance: instance}, nil
}

// Produce appends message, returning once it is durably committed.
func (p *ProducerInstance[Message]) Produce(ctx context.Context, message *Message, options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.Produce(ctx, message, (*produce.ProduceOptions)(options))
}

// ProduceBatch appends every item in one transaction -- none land unless all do.
func (p *ProducerInstance[Message]) ProduceBatch(ctx context.Context, items ...*ProduceItem[Message]) ([]*ProduceResult[Message], error) {
	converted := make([]*producer.ProduceItem[Message], len(items))
	for i, item := range items {
		if item != nil {
			converted[i] = &producer.ProduceItem[Message]{
				Message: item.Message,
				Options: produce.ProduceOptions(item.Options),
			}
		}
	}
	return p.instance.ProduceBatch(ctx, converted...)
}

// ProduceFunc appends the message producerFunc returns from inside the
// message's own transaction.
func (p *ProducerInstance[Message]) ProduceFunc(ctx context.Context, producerFunc ProducerFunc[Message], options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.ProduceFunc(ctx, producerFunc, (*produce.ProduceOptions)(options))
}

// ProduceInTx appends message inside a transaction the caller owns.
func (p *ProducerInstance[Message]) ProduceInTx(ctx context.Context, tx Tx, message *Message, options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.ProduceInTx(ctx, tx, message, (*produce.ProduceOptions)(options))
}

// ProduceFuncInTx appends the message producerFunc returns, inside a
// transaction the caller owns.
func (p *ProducerInstance[Message]) ProduceFuncInTx(ctx context.Context, tx Tx, producerFunc ProducerFunc[Message], options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.ProduceFuncInTx(ctx, tx, producerFunc, (*produce.ProduceOptions)(options))
}
