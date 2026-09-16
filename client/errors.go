package sqlstreams

// Every declared error, under its own name: the same value, so errors.Is
// holds whether a caller spells the root or this package.

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/compaction"
	"github.com/allegedlyreliable/sqlstreams/pkg/consume"
	"github.com/allegedlyreliable/sqlstreams/pkg/migrate"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
	"github.com/allegedlyreliable/sqlstreams/pkg/schedule"
	"github.com/allegedlyreliable/sqlstreams/pkg/stream"
	"github.com/allegedlyreliable/sqlstreams/pkg/system"
	"github.com/allegedlyreliable/sqlstreams/pkg/worker"
)

var (
	// ErrAlreadyConsuming means Consume ran twice at once on one instance -- an
	// instance runs one Consume at a time.
	ErrAlreadyConsuming = common.ErrAlreadyConsuming

	// ErrCommitConfirmationLost means the connection died at Commit with
	// outcomes already queued: whether they landed is unconfirmable, so a retry
	// could record duplicates -- the lease's expiry sorts the truth out.
	//
	// Diagnose queries: sqlstreams explain SQL0019
	ErrCommitConfirmationLost = common.ErrCommitConfirmationLost

	// ErrLeaseLost means the row was reclaimed by another consumer between the
	// claim and the write; the delivery machinery handles the redelivery.
	ErrLeaseLost = common.ErrLeaseLost

	// ErrLifecycleContextNotCancellable means Consume's ctx can never be
	// cancelled (e.g. context.Background()), so shutdown could never be
	// requested.
	ErrLifecycleContextNotCancellable = common.ErrLifecycleContextNotCancellable

	// ErrPayloadNotEncodable means encoding/json rejected the payload -- a NaN
	// float, a channel or func field, a MarshalJSON that returned an error. The
	// payload is encoded in Go before it reaches pgx because pgx's own encode
	// failure prints the value it could not encode.
	ErrPayloadNotEncodable = common.ErrPayloadNotEncodable

	// ErrCompactionHeadNotFound means no message produced under the key is its
	// current compaction head. A lockable row may exist without a head.
	//
	// Diagnose queries: sqlstreams explain SQL0066
	ErrCompactionHeadNotFound = compaction.ErrCompactionHeadNotFound

	// ErrDeliveryDelayed is what Delay returns: the handler asked for a later
	// run, so the delivery waits out the delay and no failure is counted.
	ErrDeliveryDelayed = consume.ErrDeliveryDelayed

	// ErrDeliveryTerminal is what Terminal returns: the handler declared that no
	// retry could succeed, so the delivery dead-letters on this attempt.
	ErrDeliveryTerminal = consume.ErrDeliveryTerminal

	// ErrConsumerGroupDeliveriesPending means Destroy was called while the group still
	// holds delivery rows, without a force override. Deleting them discards:
	//   - ready/inflight/deferred rows -> failures promised a retry
	//   - dead rows                    -> the dead-letter record
	//
	// Diagnose queries: sqlstreams explain SQL0016
	ErrConsumerGroupDeliveriesPending = consume.ErrConsumerGroupDeliveriesPending

	// ErrConsumerGroupLive means Destroy was called while a worker instance still runs
	// on the group, without a force override.
	//
	// Diagnose queries: sqlstreams explain SQL0015
	ErrConsumerGroupLive = consume.ErrConsumerGroupLive

	// ErrConsumerNotFound means the named group has no row on that stream.
	//
	// Diagnose queries: sqlstreams explain SQL0014
	ErrConsumerNotFound = consume.ErrConsumerNotFound

	// ErrNotRegistered means the queried owner has no baseline record -- the system
	// or stream was never registered, or migration_log is missing.
	ErrNotRegistered = migrate.ErrNotRegistered

	// ErrSchemaNewerThanBuild means the database was migrated past this build by
	// a step whose MinCompatibleVersion is above it.
	//
	// Diagnose queries: sqlstreams explain SQL0023
	ErrSchemaNewerThanBuild = migrate.ErrSchemaNewerThanBuild

	// ErrSchemaOlderThanBuild means the stored schema version is below what this
	// build requires.
	//
	// Diagnose queries: sqlstreams explain SQL0022
	ErrSchemaOlderThanBuild = migrate.ErrSchemaOlderThanBuild

	// ErrStepLockTimeout reclassifies a lock_timeout expiry (55P03) on the txn
	// step path: lock contention is what the step retry exists to ride out, while
	// IsTransientPgError alone would stop the run.
	ErrStepLockTimeout = migrate.ErrStepLockTimeout

	// ErrPartitionCreationBehind is the heal loop's exhaustion: every rerun of
	// the insert drew ids past the partition the previous heal created.
	ErrPartitionCreationBehind = produce.ErrPartitionCreationBehind

	// ErrPartitionLockTimeout reclassifies a lock_timeout expiry (55P03) on the
	// create-ahead path: lock contention is exactly what the run's backoff
	// schedule exists to ride out, while the heal path keeps its fail-fast read.
	ErrPartitionLockTimeout = produce.ErrPartitionLockTimeout

	// ErrScheduleDeclarationInterrupted means the schedule row was deleted between the
	// declaration's insert attempt and its update; an unchanged retry re-creates
	// the row, so DatastoreRetry heals the race.
	ErrScheduleDeclarationInterrupted = schedule.ErrScheduleDeclarationInterrupted

	// ErrScheduleNotFound means the named schedule has no row.
	ErrScheduleNotFound = schedule.ErrScheduleNotFound

	// ErrSchemaNotCreatable means RegisterSystem could not create the namespace
	// sqlstreams's tables live in -- the connecting role has no CREATE privilege.
	ErrSchemaNotCreatable = system.ErrSchemaNotCreatable

	// ErrSystemLive means DestroySystem was refused because a worker instance is
	// still live -- a manager or consumer is running somewhere.
	ErrSystemLive = system.ErrSystemLive

	// ErrStreamsRegistered means DestroySystem was refused because non-system
	// streams are still registered.
	ErrStreamsRegistered = system.ErrStreamsRegistered

	// ErrDestroyDisabled means a Destroy* call ran without AllowDestroy set
	// on the admin's config.
	ErrDestroyDisabled = stream.ErrDestroyDisabled

	// ErrReservedStreamName means Register/Rename touched a name under
	// SystemStreamPrefix -- reserved for the admin's own system streams.
	ErrReservedStreamName = stream.ErrReservedStreamName

	// ErrStreamConfigMismatch means Register was called with a PartitionSize the stream wasn't created with.
	// Every other config field can be changed by registering again.
	ErrStreamConfigMismatch = stream.ErrStreamConfigMismatch

	// ErrStreamDeclarationInterrupted means the stream row was destroyed between
	// the declaration's config write and its re-read; an unchanged retry
	// registers the stream fresh, so DatastoreRetry heals the race.
	ErrStreamDeclarationInterrupted = stream.ErrStreamDeclarationInterrupted

	// ErrStreamNameTaken means Rename's target name already belongs to another stream.
	ErrStreamNameTaken = stream.ErrStreamNameTaken

	// ErrStreamNotEmpty means Destroy was called on a stream that still holds
	// messages, without an explicit force override.
	//
	// Diagnose queries: sqlstreams explain SQL0006
	ErrStreamNotEmpty = stream.ErrStreamNotEmpty

	// ErrStreamNotFound means the named stream has no row.
	//
	// Diagnose queries: sqlstreams explain SQL0005
	ErrStreamNotFound = stream.ErrStreamNotFound

	// ErrStreamPartitionsRemain means Destroy kept finding new partitions after
	// its drop-pass limit -- a producer is likely still writing.
	//
	// Diagnose queries: sqlstreams explain SQL0020
	ErrStreamPartitionsRemain = stream.ErrStreamPartitionsRemain

	// ErrInstanceLost means the instance row expired or was removed mid-work:
	// stop -- a replacement may already be running.
	ErrInstanceLost = worker.ErrInstanceLost

	// ErrWorkerDeclarationInterrupted means the worker row was deleted between the
	// declaration's insert attempt and its update; an unchanged retry re-creates
	// the row, so DatastoreRetry heals the race.
	ErrWorkerDeclarationInterrupted = worker.ErrWorkerDeclarationInterrupted

	// ErrWorkerNotFound means the requested worker has not been declared.
	ErrWorkerNotFound = worker.ErrWorkerNotFound

	// ErrWorkerSuspended stops a running execution after its operational target becomes zero.
	ErrWorkerSuspended = worker.ErrWorkerSuspended
)
