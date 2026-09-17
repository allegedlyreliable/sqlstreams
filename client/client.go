// Package sqlstreams provides message streams, consumer groups, and schedules
// backed by Postgres. A Client uses an application-owned pgx pool.
package sqlstreams

import (
	"context"
	"errors"

	"github.com/allegedlyreliable/sqlstreams/pkg/admin"
	"github.com/allegedlyreliable/sqlstreams/pkg/consumer"
	"github.com/allegedlyreliable/sqlstreams/pkg/datastore"
	"github.com/allegedlyreliable/sqlstreams/pkg/producer"
	"github.com/allegedlyreliable/sqlstreams/pkg/scheduler"
	"github.com/allegedlyreliable/sqlstreams/pkg/systemmanager"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client provides access to one SQLStreams installation in a PostgreSQL schema.
type Client struct {
	disableManager bool

	ds        *datastore.PostgresDatastore
	admin     *admin.MessageAdmin
	consumer  *consumer.Consumer
	producer  *producer.Producer
	scheduler *scheduler.Scheduler
	manager   *systemmanager.SystemManager
}

// NewClient builds every registration object over pool and pings it once, so a wrong
// address or credential fails here instead of at the first query. The pool
// stays the caller's. SQLStreams never closes it. cfg may be nil or sparse.
// Settings are captured at construction, including a copy of Retry.
// Later edits to cfg do not reconfigure the client. The supplied logger is shared.
func NewClient(ctx context.Context, pool *pgxpool.Pool, cfg *ClientConfig) (*Client, error) {
	if pool == nil {
		return nil, errors.New("pool must not be nil")
	}
	if cfg == nil {
		cfg = &ClientConfig{}
	}
	cfg.WithDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var retry *RetryPolicy
	if cfg.Retry != nil {
		policy := *cfg.Retry
		retry = &policy
	}

	ds, err := datastore.NewPostgresDatastore(ctx, pool, &datastore.PostgresDatastoreConfig{
		Schema: cfg.Schema,
		Logger: cfg.Logger,
		Retry:  retry,
	})
	if err != nil {
		return nil, err
	}

	messageAdmin, err := admin.NewMessageAdmin(ds, &admin.MessageAdminConfig{AllowDestroy: cfg.AllowDestroy})
	if err != nil {
		return nil, err
	}

	messageConsumer, err := consumer.NewConsumer(ds)
	if err != nil {
		return nil, err
	}
	messageProducer, err := producer.NewProducer(ds)
	if err != nil {
		return nil, err
	}
	messageScheduler, err := scheduler.NewScheduler(ds)
	if err != nil {
		return nil, err
	}

	systemManager, err := systemmanager.NewSystemManager(ds, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		disableManager: cfg.DisableManager,
		ds:             ds,
		admin:          messageAdmin,
		consumer:       messageConsumer,
		producer:       messageProducer,
		scheduler:      messageScheduler,
		manager:        systemManager,
	}, nil
}

// InTransaction opens one transaction, runs transactionFunc against it, and
// commits. Use ProduceInTx to produce to multiple streams atomically.
//
// It does not retry. A commit error can leave the outcome uncertain.
// Caller-supplied IdempotencyKeys protect message inserts across calls while
// those keys are retained. They do not deduplicate other callback work.
// Retrying requires every part of the callback to be safe to repeat.
func (c *Client) InTransaction(ctx context.Context, transactionFunc TransactionFunc) error {
	return datastore.InTransaction(ctx, c.ds, transactionFunc)
}
