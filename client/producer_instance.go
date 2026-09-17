package sqlstreams

import (
	"context"
	"errors"

	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
	"github.com/allegedlyreliable/sqlstreams/pkg/producer"
)

// ProducerInstance produces messages on the stream resolved at registration.
// Cancellation applies to individual calls. The instance has no shutdown method.
type ProducerInstance[Message Versioned] struct {
	instance *producer.ProducerInstance[Message]
}

func newProducerInstance[Message Versioned](instance *producer.ProducerInstance[Message]) (*ProducerInstance[Message], error) {
	if instance == nil {
		return nil, errors.New("instance must not be nil")
	}
	return &ProducerInstance[Message]{instance: instance}, nil
}

// Produce appends a message and returns after commit.
func (p *ProducerInstance[Message]) Produce(ctx context.Context, message *Message, options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.Produce(ctx, message, (*produce.ProduceOptions)(options))
}

// ProduceBatch appends all items in one transaction.
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

// ProduceFunc runs producerFunc and produces its payload in one transaction.
// The callback may run more than once.
func (p *ProducerInstance[Message]) ProduceFunc(ctx context.Context, producerFunc ProducerFunc[Message], options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.ProduceFunc(ctx, producerFunc, (*produce.ProduceOptions)(options))
}

// ProduceInTx inserts message in the caller's transaction without committing.
func (p *ProducerInstance[Message]) ProduceInTx(ctx context.Context, tx Tx, message *Message, options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.ProduceInTx(ctx, tx, message, (*produce.ProduceOptions)(options))
}

// ProduceFuncInTx runs producerFunc and produces its payload in the caller's
// transaction. It does not commit.
func (p *ProducerInstance[Message]) ProduceFuncInTx(ctx context.Context, tx Tx, producerFunc ProducerFunc[Message], options *ProduceOptions) (*ProduceResult[Message], error) {
	return p.instance.ProduceFuncInTx(ctx, tx, producerFunc, (*produce.ProduceOptions)(options))
}
