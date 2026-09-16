package sqlstreams

// Every declared log event, under its own name, for callers that filter
// on a code.

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/consume"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
	"github.com/allegedlyreliable/sqlstreams/pkg/produce"
	"github.com/allegedlyreliable/sqlstreams/pkg/schedule"
	"github.com/allegedlyreliable/sqlstreams/pkg/stream"
	"github.com/allegedlyreliable/sqlstreams/pkg/system"
	"github.com/allegedlyreliable/sqlstreams/pkg/worker"
)

var (
	// EventAlertConditionHolds means a Register-time pass measured the stream
	// against the built-in alert conditions and one of them held. The pass is
	// log-only -- Register never fails on it.
	EventAlertConditionHolds = alert.EventAlertConditionHolds

	// EventAlertEvidenceInvalid means a scheduled check cannot assess its owner
	// because retained evidence fails semantic validation.
	EventAlertEvidenceInvalid = alert.EventAlertEvidenceInvalid

	// EventConsumerStopped is the session summary a consumer instance logs on
	// every exit.
	EventConsumerStopped = consume.EventConsumerStopped

	// EventExceptionDeadLettered marks one exception written as terminal.
	//
	// Diagnose queries: sqlstreams explain SQL0030
	EventExceptionDeadLettered = consume.EventExceptionDeadLettered

	// EventGroupConfigNotRefreshed means an instance could not read its worker
	// row's stored config back, so it keeps running on the copy it already has.
	//
	// Diagnose queries: sqlstreams explain SQL0060
	EventGroupConfigNotRefreshed = consume.EventGroupConfigNotRefreshed

	// EventKillBackstopFired means the crash-loop backstop marked a group's
	// exceptions dead after repeated consumer crashes on the same rows.
	//
	// Diagnose queries: sqlstreams explain SQL0031
	EventKillBackstopFired = consume.EventKillBackstopFired

	// EventLeaseReclaimed means a range lease's worker stopped renewing and the
	// range went back to the claimable pool.
	//
	// Diagnose queries: sqlstreams explain SQL0026
	EventLeaseReclaimed = consume.EventLeaseReclaimed

	// EventMessageDeadLettered marks one delivery written as terminal.
	//
	// Diagnose queries: sqlstreams explain SQL0029
	EventMessageDeadLettered = consume.EventMessageDeadLettered

	// EventMessagesDeadLettered marks a commit that wrote terminal outcomes for
	// a batch of messages.
	//
	// Diagnose queries: sqlstreams explain SQL0028
	EventMessagesDeadLettered = consume.EventMessagesDeadLettered

	// EventMessagesNotClaimed means a claim did not land for a reason an
	// unchanged retry can fix; a permanent cause ends Consume instead.
	EventMessagesNotClaimed = consume.EventMessagesNotClaimed

	// EventRangeQuarantined means a range hit MaxRangeReclaims and is treated as
	// poison instead of being handed out again.
	//
	// Diagnose queries: sqlstreams explain SQL0027
	EventRangeQuarantined = consume.EventRangeQuarantined

	// EventSlowDispatch means one delivery's dispatch ran past the group's
	// SlowDispatchThreshold, whatever the delivery's outcome.
	EventSlowDispatch = consume.EventSlowDispatch

	// EventQueuedRangeStale means a queued message cannot start within its range's lease budget.
	EventQueuedRangeStale = consume.EventQueuedRangeStale

	// EventStoredOptionsClamped means a stored message's options fell outside
	// this consumer's MessageMin/MessageMax bounds.
	EventStoredOptionsClamped = consume.EventStoredOptionsClamped

	// EventGoRoutineEventsDropped means abandoned/cleared routine events were
	// discarded -- the queue filled between flush ticks, or their batch could
	// not land.
	EventGoRoutineEventsDropped = metric.EventGoRoutineEventsDropped

	// EventMeasurementsCannotBeExported means the current collection omits a
	// rejected metric family while healthy families still export.
	EventMeasurementsCannotBeExported = metric.EventMeasurementsCannotBeExported

	// EventPartitionCreatedOnInsert means an insert found no partition for its
	// id and created one itself: create-ahead did not run, or a burst
	// outran its triggers.
	EventPartitionCreatedOnInsert = produce.EventPartitionCreatedOnInsert

	// EventPartitionNotCreatedAhead means the create-ahead pass gave up on the
	// next partition; the write path still covers it.
	EventPartitionNotCreatedAhead = produce.EventPartitionNotCreatedAhead

	// EventSlowProduce means one produce call ran past the producer's
	// SlowProduceThreshold, whatever the call's outcome.
	EventSlowProduce = produce.EventSlowProduce

	// EventMessageAlreadyProduced means a schedule producer tick found its
	// message already in the stream: an earlier tick's commit confirmation was
	// lost after the produce landed, so this tick produces nothing.
	EventMessageAlreadyProduced = schedule.EventMessageAlreadyProduced

	// EventScheduleConfigReplaced means a declaration overwrote a schedule row's
	// differing config -- two declarers disagree about the schedule.
	//
	// Diagnose queries: sqlstreams explain SQL0062
	EventScheduleConfigReplaced = schedule.EventScheduleConfigReplaced

	// EventTargetKeepsNoSuccessRows means the schedule's target stream keeps
	// failure rows only, so ScheduleStatus can never count a success.
	EventTargetKeepsNoSuccessRows = schedule.EventTargetKeepsNoSuccessRows

	// EventSystemManagerStopped is a Run life ending on its own -- a spawned
	// worker declared itself unrunnable, or the manager row could not be claimed.
	// The caller blocked in Run has no error value coming, so this line is where
	// an operator learns of it.
	EventSystemManagerStopped = system.EventSystemManagerStopped

	// EventStreamConfigReplaced means a declaration overwrote a stream row's
	// differing mutable config -- two declarers disagree about the stream.
	//
	// Diagnose queries: sqlstreams explain SQL0061
	EventStreamConfigReplaced = stream.EventStreamConfigReplaced

	// EventInstanceLost means this instance's worker_instance row was claimed by
	// a replacement while it was still running.
	//
	// Diagnose queries: sqlstreams explain SQL0034
	EventInstanceLost = worker.EventInstanceLost

	// EventManagerRowSuspended means the manager's own row has target_instances
	// 0, so its workers stop being reconciled.
	//
	// Diagnose queries: sqlstreams explain SQL0035
	EventManagerRowSuspended = worker.EventManagerRowSuspended

	// EventSlowTick means one tick ran longer than the row's own poll_rate --
	// the worker is behind its own schedule.
	EventSlowTick = worker.EventSlowTick

	// EventTickBackoffCurveExhausted means a tick loop's failure streak passed
	// its TickRetry cap -- the failure is no longer self-healing.
	//
	// Diagnose queries: sqlstreams explain SQL0036
	EventTickBackoffCurveExhausted = worker.EventTickBackoffCurveExhausted

	// EventWorkerConfigReplaced means a declaration overwrote a worker row's
	// differing stored config -- two declarers disagree about the worker.
	//
	// Diagnose queries: sqlstreams explain SQL0059
	EventWorkerConfigReplaced = worker.EventWorkerConfigReplaced
)
