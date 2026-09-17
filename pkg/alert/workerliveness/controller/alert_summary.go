package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
)

type workerLivenessSummary struct {
	messages      []string
	startConsumer bool
	runManager    bool
	inspectWorker bool
}

func newWorkerLivenessSummary(unclaimed [][]*metric.UnclaimedWorkerMetadata) (*workerLivenessSummary, error) {
	if len(unclaimed) == 0 {
		return nil, errors.New("unclaimed must not be empty")
	}

	summary := &workerLivenessSummary{}
	for _, workers := range unclaimed {
		summary.describeOwner(workers)
	}
	return summary, nil
}

func (s *workerLivenessSummary) describeOwner(workers []*metric.UnclaimedWorkerMetadata) {
	owner := workers[0].Owner
	names := make([]string, 0, len(workers))
	for _, worker := range workers {
		names = append(names, s.describeWorker(owner, worker.Name))
	}

	resource := "Stream"
	if owner.Kind() == common.OwnerConsumerGroup {
		resource = "Consumer group"
	}
	s.messages = append(s.messages, fmt.Sprintf("%s %q has no running %s.", resource, owner.Name, joinWorkerNames(names)))
}

func (s *workerLivenessSummary) describeWorker(owner *common.Owner, name string) string {
	group := owner.Kind() == common.OwnerConsumerGroup
	switch {
	case group && name == "message_consumer":
		s.startConsumer = true
		return "message consumer"
	case group && name == "exception_consumer":
		s.startConsumer = true
		return "retry consumer"
	case group && name == "manager":
		s.startConsumer = true
		return "consumer manager"
	case group && name == "delivery_consumer":
		s.startConsumer = true
		return "delivery consumer"
	case owner.Kind() == common.OwnerStream && name == "stream_janitor":
		s.runManager = true
		return "retention cleanup worker"
	case owner.Kind() == common.OwnerStream && name == "stream_vacuum":
		s.runManager = true
		return "vacuum worker"
	case group && name == "cursor_advancer":
		s.runManager = true
		return "cursor advancement worker"
	case group && (name == "alert.partition_count" || name == "alert.compaction_read_cost" || name == "alert.worker_liveness" || name == "alert.metrics_collector_progress"):
		s.runManager = true
		return fmt.Sprintf("%q alert evaluator", strings.TrimPrefix(name, "alert."))
	default:
		s.inspectWorker = true
		return fmt.Sprintf("worker %q", name)
	}
}

func (s *workerLivenessSummary) message() string {
	return strings.Join(s.messages, " ")
}

func (s *workerLivenessSummary) hint() string {
	hints := make([]string, 0, 3)
	if s.startConsumer {
		hints = append(hints, "Start a consumer for each affected group, or delete any group you no longer need.")
	}
	if s.runManager {
		hints = append(hints, "Run \"sqlstreams manager run\" to resume maintenance.")
	}
	if s.inspectWorker {
		hints = append(hints, "Check the process responsible for the named workers.")
	}
	return strings.Join(hints, " ")
}

// ***************
// *** HELPERS ***
// ***************

func joinWorkerNames(names []string) string {
	if len(names) == 1 {
		return names[0]
	}
	separator := " or "
	if len(names) > 2 {
		separator = ", or "
	}
	return strings.Join(names[:len(names)-1], ", ") + separator + names[len(names)-1]
}
