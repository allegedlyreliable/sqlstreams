package producer

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/common/logging"
	iDatastore "github.com/allegedlyreliable/sqlstreams/pkg/datastore"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce/batcher"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce/controller"
	"github.com/allegedlyreliable/sqlstreams/pkg/stream"
)

// ProducerInstance produces messages on the stream resolved at registration.
// Cancellation applies to individual calls. The instance has no shutdown method.
type ProducerInstance[Message common.Versioned] struct {
	Stream *stream.Stream  // the stream row Register resolved
	Config *ProducerConfig // resolved settings retained from registration
	Logger logging.Logger  // bound to the stream; the batcher logs through it

	controller *controller.ProduceController
	batcher    *batcher.Batcher[Message]
}

// cfg is already resolved (WithDefaults + Validate) by Register; logger is
// its per-instance pipeline over the datastore's logger.
func NewProducerInstance[Message common.Versioned](resolvedStream *stream.Stream, produceController *controller.ProduceController, cfg *ProducerConfig, logger logging.Logger) (*ProducerInstance[Message], error) {
	if resolvedStream == nil {
		return nil, errors.New("stream must not be nil")
	}
	if produceController == nil {
		return nil, errors.New("controller must not be nil")
	}
	if cfg == nil {
		return nil, errors.New("config must not be nil")
	}
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}

	streamBatcher, err := batcher.NewBatcher[Message](produceController, resolvedStream.Id, resolvedStream.PartitionSize, &cfg.Batch, logger)
	if err != nil {
		return nil, err
	}

	return &ProducerInstance[Message]{
		Stream:     resolvedStream,
		Config:     cfg,
		Logger:     logger,
		controller: produceController,
		batcher:    streamBatcher,
	}, nil
}

// Produce appends a message and returns after commit.
func (p *ProducerInstance[Message]) Produce(ctx context.Context, message *Message, options *produce.ProduceOptions) (*ProduceResult[Message], error) {
	defer p.warnSlowProduce(ctx, time.Now())
	ctx = logging.WithLogBuffer(ctx)

	resolved := produce.ProduceOptions{}
	if options != nil {
		resolved = *options
	}

	resolved.Message = resolved.Message.Fill(p.Config.Message)
	if err := resolved.Validate(); err != nil {
		return nil, err
	}

	// caller keys can collide -- a collision inside a shared txn stalls the
	// whole batch, so keyed calls take a per-call transaction
	if resolved.IdempotencyKey != "" {
		passthrough := func(context.Context, iDatastore.Tx) (*Message, error) { return message, nil }
		appended, err := p.controller.AppendMessage(ctx, p.Stream.Id, p.Stream.PartitionSize, passthrough, resolved)
		if err != nil {
			return nil, err
		}
		return NewProduceResult(appended.Message, appended.Id, appended.Duplicate)
	}

	appended, err := p.batcher.Produce(ctx, message, resolved)
	if err != nil {
		return nil, err
	}
	return NewProduceResult(appended.Message, appended.Id, appended.Duplicate)
}

// ProduceBatch appends all items in one transaction.
func (p *ProducerInstance[Message]) ProduceBatch(ctx context.Context, items ...*ProduceItem[Message]) ([]*ProduceResult[Message], error) {
	if len(items) == 0 {
		return nil, errors.New("items must not be empty")
	}
	defer p.warnSlowProduce(ctx, time.Now())
	ctx = logging.WithLogBuffer(ctx)

	appends := make([]*controller.Append[Message], 0, len(items))
	for i, item := range items {
		appendItem, err := p.toAppend(item)
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		appends = append(appends, appendItem)
	}

	appendedRows, failedIdx, err := p.controller.AppendMessageBatch(ctx, p.Stream.Id, p.Stream.PartitionSize, p.Config.Batch.AttemptTimeout, appends)
	if err != nil {
		if failedIdx >= 0 {
			return nil, fmt.Errorf("item %d: %w", failedIdx, err)
		}
		return nil, err
	}

	results := make([]*ProduceResult[Message], 0, len(appendedRows))
	for _, appendedRow := range appendedRows {
		result, err := NewProduceResult(appendedRow.Message, appendedRow.Id, appendedRow.Duplicate)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// ProduceFunc runs producerFunc and produces its payload in one transaction.
// The callback may run more than once.
func (p *ProducerInstance[Message]) ProduceFunc(ctx context.Context, producerFunc produce.ProducerFunc[Message], options *produce.ProduceOptions) (*ProduceResult[Message], error) {
	defer p.warnSlowProduce(ctx, time.Now())

	resolved := produce.ProduceOptions{}
	if options != nil {
		resolved = *options
	}

	resolved.Message = resolved.Message.Fill(p.Config.Message)
	if err := resolved.Validate(); err != nil {
		return nil, err
	}

	appended, err := p.controller.AppendMessage(ctx, p.Stream.Id, p.Stream.PartitionSize, producerFunc, resolved)
	if err != nil {
		return nil, err
	}
	return NewProduceResult(appended.Message, appended.Id, appended.Duplicate)
}

// ProduceInTx inserts message in the caller's transaction without committing.
func (p *ProducerInstance[Message]) ProduceInTx(ctx context.Context, tx iDatastore.Tx, message *Message, options *produce.ProduceOptions) (*ProduceResult[Message], error) {
	passthrough := func(context.Context, iDatastore.Tx) (*Message, error) { return message, nil }
	return p.ProduceFuncInTx(ctx, tx, passthrough, options)
}

// ProduceFuncInTx runs producerFunc and produces its payload in the caller's
// transaction. It does not commit.
func (p *ProducerInstance[Message]) ProduceFuncInTx(ctx context.Context, tx iDatastore.Tx, producerFunc produce.ProducerFunc[Message], options *produce.ProduceOptions) (*ProduceResult[Message], error) {
	defer p.warnSlowProduce(ctx, time.Now())

	resolved := produce.ProduceOptions{}
	if options != nil {
		resolved = *options
	}

	resolved.Message = resolved.Message.Fill(p.Config.Message)
	if err := resolved.Validate(); err != nil {
		return nil, err
	}

	appended, err := p.controller.AppendMessageInTx(ctx, tx, p.Stream.Id, p.Stream.PartitionSize, producerFunc, resolved)
	if err != nil {
		return nil, err
	}
	return NewProduceResult(appended.Message, appended.Id, appended.Duplicate)
}

// warnSlowProduce logs one line when a produce entry point ran past the
// configured threshold -- slowness is its own fact, logged whatever the
// call's outcome.
func (p *ProducerInstance[Message]) warnSlowProduce(ctx context.Context, start time.Time) {
	duration := time.Since(start)
	if p.Config.SlowProduceThreshold <= 0 || duration <= p.Config.SlowProduceThreshold {
		return
	}
	p.Logger.WarnContext(ctx, produce.EventSlowProduce.Message(), "code", produce.EventSlowProduce.GetCode(), "stream", p.Stream.Name, "duration", duration, "threshold", p.Config.SlowProduceThreshold)
}

// toAppend shapes one batch item for the controller: fills message options
// and generates the key the datastore's ambiguous-commit rerun dedups on.
func (p *ProducerInstance[Message]) toAppend(item *ProduceItem[Message]) (*controller.Append[Message], error) {
	if item == nil {
		return nil, errors.New("item must not be nil")
	}
	if item.Options.IdempotencyKey != "" {
		return nil, errors.New("IdempotencyKey is not supported in a batch -- produce keyed messages individually")
	}

	options := item.Options
	options.Message = options.Message.Fill(p.Config.Message)
	options.IdempotencyKey = uuid.NewV7().String()
	return controller.NewAppend(item.Message, options)
}
