package sqlstreams

// Comments are copied from the owning declarations so editor hover and
// Go documentation show them at the public entry point.

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/admin"
	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/common/diagnostic"
	"github.com/allegedlyreliable/sqlstreams/pkg/common/logging"
	"github.com/allegedlyreliable/sqlstreams/pkg/consume"
	"github.com/allegedlyreliable/sqlstreams/pkg/consumer"
	"github.com/allegedlyreliable/sqlstreams/pkg/datastore"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce/batcher"
	"github.com/allegedlyreliable/sqlstreams/pkg/producer"
	"github.com/allegedlyreliable/sqlstreams/pkg/schedule"
	"github.com/allegedlyreliable/sqlstreams/pkg/stream"
	"github.com/allegedlyreliable/sqlstreams/pkg/system"
	"github.com/allegedlyreliable/sqlstreams/pkg/worker"
)

type (
	// Versioned is what every Message type declares: the payload's
	// compatibility version, written on every message row it produces. A
	// value receiver, so the zero value answers.
	//
	// Bump only on a BREAKING change to Message:
	// - a field has a different type
	// - a field has been renamed
	// - a field has been removed
	//
	// A consumer instance reads only rows at its Message type's version.
	Versioned = common.Versioned

	// RawPayload is a message payload kept as the JSON bytes the row stores,
	// for readers with no message type in scope: the CLI, an admin-only
	// script. Its version is 0, which no declared message type may use:
	// Consumer(name).Register refuses it and every produce path refuses to
	// write it.
	RawPayload = common.RawPayload

	// StoredMessage is one message as the log holds it: the payload plus the
	// row's own facts.
	StoredMessage[Message Versioned] = common.StoredMessage[Message]

	// MessageOptions requests processing settings for one message.
	// Unset fields inherit producer defaults, then consumer group defaults.
	// The group's numeric bounds and concurrency override apply when consumed.
	MessageOptions = common.MessageOptions

	// RetryPolicy configures exponential backoff. ClientConfig.Retry applies to
	// supported internal database operations. MessageOptions.Retry applies to
	// message redelivery. Zero fields select defaults, except MaxDelays, where
	// zero means unlimited requested delays.
	RetryPolicy = common.RetryPolicy

	// ConcurrencyPolicy is a message's concurrency policy.
	ConcurrencyPolicy = common.ConcurrencyPolicy

	// Owner is which resource owns a row in a polymorphic table (worker,
	// schedule, migration_log).
	Owner = common.Owner

	// OwnerKind is which resource an Owner names, derived from the ids it holds.
	OwnerKind = common.OwnerKind

	// Logger is exactly *slog.Logger's Context method set. Pass your own
	// *slog.Logger with whatever slog.Handler you want (zap/zerolog/logr all
	// ship one), or anything else that implements these four methods.
	Logger = logging.Logger

	// DiagnosticError describes a named failure with a diagnostic code,
	// recovery classification, problem, fix, and diagnostic queries.
	// With and Wrap attach values and a cause without changing the declaration.
	DiagnosticError = diagnostic.DiagnosticError

	// DiagnosticEvent is a declared operator-actionable log event: the static message
	// a call site logs and the code that rides in its "code" attribute.
	DiagnosticEvent = diagnostic.DiagnosticEvent

	// DiagnosticQuery is one declared diagnostic query: the label names what the query
	// answers, the SQL answers it against the reader's own database. The library
	// never runs it -- the fix says what to change, a query says what to look at.
	DiagnosticQuery = diagnostic.DiagnosticQuery

	// DiagnosticRecovery states whether an unchanged retry of the operation can succeed.
	DiagnosticRecovery = diagnostic.DiagnosticRecovery

	// DiagnosticKind identifies an error, log event, metric, or alert declaration.
	DiagnosticKind = diagnostic.DiagnosticKind

	// Querier is what pool, conn, and tx can all do, minus transaction control
	// (Begin/Commit/Rollback) -- the surface for statements that run inside a
	// boundary the callee doesn't own. *pgxpool.Pool, *pgxpool.Conn, and pgx.Tx
	// all satisfy it.
	Querier = datastore.Querier

	// Tx is the surface handed to a ProducerFunc/TransactionFunc closure: every
	// statement pool, conn, and tx share (Querier), never transaction control --
	// the closure runs inside a transaction it does not own.
	Tx = datastore.Tx

	// TransactionFunc runs inside the transaction InTransaction opened. A nil
	// return commits; an error rolls back and is returned as-is. It is called
	// once -- InTransaction never reruns it.
	TransactionFunc = datastore.TransactionFunc

	// CompactionOptions is ProduceOptions.Compaction: whether the message
	// compacts under its MessageKey, and the rank that decides the key's winner.
	CompactionOptions = produce.CompactionOptions

	// ProducerFunc runs inside the append's transaction and returns the payload to
	// store -- its writes commit or roll back with the message.
	ProducerFunc[Message Versioned] = produce.ProducerFunc[Message]

	// BatcherConfig tunes the shared-transaction batching of concurrent Produce
	// calls on one producer instance.
	BatcherConfig = batcher.BatcherConfig

	// MetricProducerInstance produces custom measurements with routing and
	// message keys derived from their metric name and attributes.
	MetricProducerInstance = producer.MetricProducerInstance

	// ProduceResult is one produce call's outcome.
	ProduceResult[Message Versioned] = producer.ProduceResult[Message]

	// ConsumeOptions is one Consume call's session settings -- how this process
	// runs, free to differ per instance. What the group means lives on
	// ConsumerConfig at Register. Sparse: zero fields take the defaults.
	ConsumeOptions = consumer.ConsumeOptions

	// ConsumerFunc processes one delivered message. It must be safe to repeat
	// after a crash or timeout. Returning nil records success. An ordinary error
	// retries under the message's RetryPolicy. Terminal marks the message dead.
	// Delay requests another attempt later without counting a failure.
	ConsumerFunc[Message Versioned] = consumer.ConsumerFunc[Message]

	// CursorPosition is a place in a stream's message log a group's cursor is set
	// to -- by Register for a group that has no cursor row yet.
	CursorPosition = consume.CursorPosition

	// CursorPositionKind names where a new group's cursor starts.
	CursorPositionKind = consume.CursorPositionKind

	// Consumer is the registered consumer group row. A group is owned by
	// exactly one stream -- names are unique per
	// stream, not globally. Children (cursor, lease, binding) reference Id and
	// carry no stream_id of their own; the stream_id FK cascade is the group's
	// lifecycle -- destroying the stream destroys it.
	Consumer = consume.Consumer

	// Binding is one declarer's newest declaration on a group.
	// BindingInstalled is the group's effective set.
	// BindingJoined is a declarer that found its set already stored.
	// BindingWaiting is a declarer still blocked on changing the effective set.
	Binding = consume.Binding

	// BindingOutcome is where one DeclareBindings attempt ended up.
	BindingOutcome = consume.BindingOutcome

	// MessageMeta is everything about a delivered message besides its payload,
	// read inside consumerFunc via MetaFromContext.
	MessageMeta = consume.MessageMeta

	// JanitorConfig declares retention cleanup settings.
	// A new stream's janitor starts active.
	JanitorConfig = stream.JanitorConfig

	// VacuumConfig declares VACUUM (ANALYZE) settings for the stream's idempotency-key table.
	// A new stream's vacuum starts suspended.
	VacuumConfig = stream.VacuumConfig

	// Stream is the registered stream row; Id addresses this stream's own
	// message_log_<id>. The remaining fields hold StreamConfig's resolved values.
	Stream = stream.Stream

	// DeliveryLogMode selects which delivery outcomes write delivery_log_<id> rows.
	DeliveryLogMode = stream.DeliveryLogMode

	// Schedule is one row of schedule_config joined to its schedule_cursor row.
	// Every schedule is the system's; StreamId is the target stream every produce
	// lands on.
	Schedule = schedule.Schedule

	// ScheduleConsumerGroupSummary is one consumer group's outcomes for one schedule's messages,
	// derived from the target stream's delivery log.
	ScheduleConsumerGroupSummary = schedule.ScheduleConsumerGroupSummary

	// ScheduleMessageStatus is one of a schedule's messages' outcome for one consumer
	// group, newest message first in a SchedulerHandle.Messages listing.
	ScheduleMessageStatus = schedule.ScheduleMessageStatus

	// ScheduleMessageOutcome is where one of a schedule's messages ended up for one
	// consumer group.
	ScheduleMessageOutcome = schedule.ScheduleMessageOutcome

	// ScheduleStoredMessage is a schedule_config row's payload at the row's
	// schema_version, produced as-is: the JSON goes to the message row
	// unchanged and SchemaVersion answers with the stored version, so the
	// producer path needs no Message type at produce time.
	ScheduleStoredMessage = schedule.ScheduleStoredMessage

	// System is the singleton system row, read back by System().Get.
	System = system.System

	// Worker is one row of the worker_config table.
	Worker = worker.Worker

	// InstanceTarget is how many live instances of one worker row may run at once
	// across the deployment
	// positive value   => claim up to that many instances
	// NoInstanceTarget => any number of instances
	// 0                => worker is suspended
	InstanceTarget = worker.InstanceTarget

	// DestroyOptions configures one Destroy call on a stream, consumer, or the
	// system. Every destroy is refused unless ClientConfig.AllowDestroy is set.
	DestroyOptions = admin.DestroyOptions

	// StreamVersionHealth reports whether a payload version can be retired.
	// Safe requires no compaction heads, unread messages, or unresolved exceptions.
	StreamVersionHealth = stream.StreamVersionHealth

	// PartitionCountAlertConfig declares how the partition_count alert is
	// evaluated and the count it alerts at.
	PartitionCountAlertConfig = alert.PartitionCountAlertConfig

	// CompactionReadCostAlertConfig declares how the compaction_read_cost alert is
	// evaluated and the cost it alerts at.
	CompactionReadCostAlertConfig = alert.CompactionReadCostAlertConfig

	// WorkerLivenessAlertConfig declares how the worker_liveness alert is
	// evaluated.
	WorkerLivenessAlertConfig = alert.WorkerLivenessAlertConfig

	// MetricCollectorProgressAlertConfig declares the installation's collector-progress check.
	MetricCollectorProgressAlertConfig = alert.MetricCollectorProgressAlertConfig

	// MetricCollectorWorkerConfig declares the metrics_collector worker row:
	// how often the collector measures the fleet and produces to
	// __system.metrics.
	MetricCollectorWorkerConfig = metric.MetricCollectorWorkerConfig

	// StreamSnapshot is a stream's live metrics, read from its tables at the
	// moment of the call, with every consumer group's snapshot beside it.
	StreamSnapshot = metric.StreamSnapshot

	// WorkerSnapshot reports one worker's operational target, live instances, and recorded failures.
	WorkerSnapshot = metric.WorkerSnapshot

	// WorkerStatus describes suspension, live claims, and recorded failures.
	WorkerStatus = metric.WorkerStatus

	// ConsumerGroupSnapshot is the live, DB-truth picture of one (group, stream),
	// sectioned by the store each number reads -- answers "what's true right now"
	// for state that multiple consumer processes share.
	ConsumerGroupSnapshot = metric.ConsumerGroupSnapshot

	// StreamSchemaVersionSnapshot is one payload version's presence in a stream's log.
	StreamSchemaVersionSnapshot = metric.StreamSchemaVersionSnapshot

	// ConsumerGroupSchemaVersionLag is one consumer group's unread and unresolved rows
	// at one payload version.
	ConsumerGroupSchemaVersionLag = metric.ConsumerGroupSchemaVersionLag

	// ConsumerGroupLag is a group's drain progress -- the retire-relevant distillation
	// of its snapshot.
	ConsumerGroupLag = metric.ConsumerGroupLag

	// CursorSnapshot is the group's read/commit position against the message log.
	CursorSnapshot = metric.CursorSnapshot

	// ExceptionSnapshot counts the group's exception-queue rows by status.
	ExceptionSnapshot = metric.ExceptionSnapshot

	// AbandonedRoutineSnapshot pairs retained abandoned and cleared events on
	// __system.metrics for one stream and consumer group. Counts describe the
	// retained window, not lifetime totals.
	AbandonedRoutineSnapshot = metric.AbandonedRoutineSnapshot

	// Measurement is one value of one metric at one time, on the __system.metrics
	// stream. Names starting with "sqlstreams." are reserved for SQLStreams's own metrics.
	Measurement = metric.Measurement

	// MetricDefinition is one SQLStreams built-in metric's identity and metadata.
	// It exists before any measurement is collected.
	MetricDefinition = metric.MetricDefinition

	// MetricKind is how a series' values read over time: a gauge replaces, a
	// counter accumulates.
	MetricKind = metric.MetricKind

	// MetricScope names the resource or collection described by a built-in metric.
	MetricScope = diagnostic.MetricScope

	// MetricUnit is a metric's UCUM code. A real unit ("ms", "s", "By") carries a
	// dimension a reader may format (47000 ms -> 47s); a braced annotation
	// ("{worker}", via MetricUnitCount) is a dimensionless count whose text is a human
	// label only. "" is no unit.
	MetricUnit = metric.MetricUnit

	// Alert is what one run found for one owner, published to the
	// __system.alerts stream as an ordinary message.
	Alert = alert.Alert

	// AlertDefinition is one SQLStreams built-in alert's identity and metadata. It
	// exists before any alert is published.
	AlertDefinition = alert.AlertDefinition

	// AlertStatus is an alert's lifecycle state -- an active alert and its later
	// resolution are versions of one compacted message key.
	AlertStatus = alert.AlertStatus

	// AlertSeverity is how urgently an operator should act.
	AlertSeverity = alert.AlertSeverity

	// AlertEvaluationSnapshot describes retained evidence under one resolved policy.
	// It is calculated on demand, not a recorded alert or the last scheduled check.
	AlertEvaluationSnapshot = alert.AlertEvaluationSnapshot

	// AlertEvaluationState describes evidence, not a recorded alert's lifecycle.
	AlertEvaluationState = alert.AlertEvaluationState
)

const (
	// AlertEvaluationStateHealthy means the evidence establishes a healthy condition.
	AlertEvaluationStateHealthy = alert.AlertEvaluationStateHealthy
	// AlertEvaluationStatePending means an unhealthy condition has not met the required pending duration.
	AlertEvaluationStatePending = alert.AlertEvaluationStatePending
	// AlertEvaluationStateActive means an unhealthy condition meets the policy's activation requirements.
	AlertEvaluationStateActive = alert.AlertEvaluationStateActive
	// AlertEvaluationStateInsufficientEvidence means the evidence cannot establish the condition's health.
	AlertEvaluationStateInsufficientEvidence = alert.AlertEvaluationStateInsufficientEvidence

	RecoveryTransient = diagnostic.RecoveryTransient // attempt unchanged -> retry can succeed
	RecoveryPermanent = diagnostic.RecoveryPermanent // attempt unchanged -> retry cannot succeed

	DiagnosticKindError  = diagnostic.DiagnosticKindError  // a declared error value (Err*)
	DiagnosticKindEvent  = diagnostic.DiagnosticKindEvent  // a declared log event (Event*)
	DiagnosticKindMetric = diagnostic.DiagnosticKindMetric // a built-in metric
	DiagnosticKindAlert  = diagnostic.DiagnosticKindAlert  // a built-in alert

	ConcurrencyParallel  = common.ConcurrencyParallel  // same-key deliveries may overlap
	ConcurrencyExclusive = common.ConcurrencyExclusive // one delivery per key at a time: a same-key message finding the key busy is deferred and runs when the key frees, oldest first (with compaction: the key's current head)
	ConcurrencyOrdered   = common.ConcurrencyOrdered   // exclusive, and a keyed message runs only after every earlier same-key message is resolved for the group -- a failed predecessor's retry goes first, dead does not hold the key

	DeliveryLogModeOff      = stream.DeliveryLogModeOff      // no rows at all
	DeliveryLogModeFailures = stream.DeliveryLogModeFailures // every outcome except success
	DeliveryLogModeAll      = stream.DeliveryLogModeAll      // every outcome, including a 'success' row per success

	// OwnerAny lifts an owner-kind guard -- any kind is admitted, like the
	// manager worker every owner declares.
	OwnerAny           = common.OwnerAny
	OwnerSystem        = common.OwnerSystem        // SystemId only
	OwnerStream        = common.OwnerStream        // SystemId and StreamId
	OwnerConsumerGroup = common.OwnerConsumerGroup // all three ids

	// CursorPositionBeginning is the zero value: the oldest retained message.
	CursorPositionBeginning = consume.CursorPositionBeginning

	// CursorPositionHead is MAX(id) of the message log when the cursor row is
	// written; a produce with a lower id that commits later is never read.
	CursorPositionHead = consume.CursorPositionHead
	BindingInstalled   = consume.BindingInstalled // the declared set is now the group's effective set
	BindingJoined      = consume.BindingJoined    // the declared set was already stored
	BindingWaiting     = consume.BindingWaiting   // a live instance still declares a different stored set

	ScheduleMessagePending    = schedule.ScheduleMessagePending    // produced, not yet run
	ScheduleMessageDeferred   = schedule.ScheduleMessageDeferred   // waiting for a previous message to finish running
	ScheduleMessageSucceeded  = schedule.ScheduleMessageSucceeded  // ran to a 'success' delivery log row
	ScheduleMessageFailed     = schedule.ScheduleMessageFailed     // raised without ever succeeding
	ScheduleMessageSuperseded = schedule.ScheduleMessageSuperseded // dropped unrun -- a newer message replaced it

	// NoInstanceTarget lifts the claim gate -- any number of instances can run.
	NoInstanceTarget = worker.NoInstanceTarget

	WorkerSuspended = metric.WorkerSuspended // target_instances = 0
	WorkerClaimed   = metric.WorkerClaimed   // live instances without a recorded failure streak
	WorkerFailing   = metric.WorkerFailing   // live instances report consecutive failures
	WorkerUnclaimed = metric.WorkerUnclaimed // no live instance row and not suspended

	MetricKindCounter          = metric.MetricKindCounter              // a running total, each measurement carries the new total
	MetricKindGauge            = metric.MetricKindGauge                // a point-in-time level, each measurement replaces the last
	MetricScopeSystem          = diagnostic.MetricScopeSystem          // one series per installation
	MetricScopeStream          = diagnostic.MetricScopeStream          // one series per stream
	MetricScopeConsumerGroup   = diagnostic.MetricScopeConsumerGroup   // one series per consumer group
	MetricScopeConsumerSession = diagnostic.MetricScopeConsumerSession // one series per Consume call
	MetricScopeExporter        = diagnostic.MetricScopeExporter        // one series per export collection, not stored
	MetricUnitMilliseconds     = metric.MetricUnitMilliseconds         // milliseconds in UCUM
	AlertStatusActive          = alert.AlertStatusActive               // the condition holds
	AlertStatusResolved        = alert.AlertStatusResolved             // a later run found the condition gone
	AlertSeverityInfo          = alert.AlertSeverityInfo               // informational -- no immediate operator action is required
	AlertSeverityWarn          = alert.AlertSeverityWarn               // degraded, not down -- an operator should learn of it eventually

	// MetricStreamName is __system.metrics
	MetricStreamName = metric.MetricStreamName

	// ScheduleStreamName is __system.schedules -- the target stream of the system-owned
	// schedules (the built-in alert checks); user schedules target their own.
	ScheduleStreamName = schedule.ScheduleStreamName

	// AlertStreamName is __system.alerts
	AlertStreamName = alert.AlertStreamName
)

var (
	// LifecycleContext returns the application-lifetime context to pass to the
	// blocking verbs -- Consume, Manager().Run, SchedulerInstance.Schedule:
	// cancelled on the first SIGINT/SIGTERM, which starts
	// graceful wind-down (new work refused, queued work drains). A SECOND exit
	// signal during the drain force-exits immediately (status 128+signum).
	//
	// log may be nil -- an info-level logger reports graceful shutdown starting
	// and completing, and warns when a second signal forces an exit.
	//
	//	ctx, stop := LifecycleContext(nil)
	//	defer stop()
	LifecycleContext = common.LifecycleContext

	// MetaFromContext retrieves MessageMeta from context within consumerFunc.
	MetaFromContext = consume.MetaFromContext

	// Terminal dead-letters this delivery now instead of retrying: cause stays
	// reachable through errors.Is/As and renders after the code in last_error.
	// Any diagnostic Permanent error classifies the same way; Terminal is the
	// spelling for a cause that carries no classification of its own.
	Terminal = consume.Terminal

	// Delay runs this delivery again after delay without counting a failure: the
	// row's can_run_after moves out by delay and its delays count goes up by one.
	// delay is a time.Duration: use 500*time.Millisecond, 5*time.Second, or
	// time.Minute. A bare integer is interpreted as nanoseconds.
	// Zero or less runs it on the next poll.
	Delay = consume.Delay

	// Beginning places a new group's cursor at the oldest retained message, so
	// it reads history before live traffic. The default.
	Beginning = consume.Beginning

	// Head places a new group's cursor at the log's MAX(id) when its row is
	// written, so it reads only messages produced after that.
	Head = consume.Head

	// NewMeasurement builds a custom measurement for the system metrics stream.
	// name, a valid kind and unit, and a non-zero at are required; attributes
	// may be nil. Names under the "sqlstreams." prefix are refused at produce time.
	NewMeasurement = metric.NewMeasurement
)
