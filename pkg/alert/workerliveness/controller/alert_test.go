package controller

import (
	"strings"
	"testing"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
)

// Behavior: liveness reports name the missing components and its matching recovery action.
func TestWorkerLivenessNamesMissingComponentsAndActions(t *testing.T) {
	for _, test := range []struct {
		name          string
		workers       [][]string // stream, first consumer group, second consumer group
		message       string
		startConsumer bool
		runManager    bool
		inspectWorker bool
	}{
		{
			name: "stopped groups",
			workers: [][]string{nil,
				{"exception_consumer", "manager", "message_consumer"},
				{"exception_consumer", "manager", "message_consumer"}},
			message: `Consumer group "email-sender-beginning" has no running retry consumer, consumer manager, or message consumer. Consumer group "email-sender-head" has no running retry consumer, consumer manager, or message consumer.`, startConsumer: true,
		},
		{
			name: "one stopped group", workers: [][]string{nil, {"message_consumer", "exception_consumer"}},
			message: `Consumer group "email-sender-beginning" has no running message consumer or retry consumer.`, startConsumer: true,
		},
		{
			name: "retries only", workers: [][]string{nil, {"exception_consumer"}},
			message: `Consumer group "email-sender-beginning" has no running retry consumer.`, startConsumer: true,
		},
		{
			name: "new messages only", workers: [][]string{nil, {"message_consumer"}},
			message: `Consumer group "email-sender-beginning" has no running message consumer.`, startConsumer: true,
		},
		{
			name: "consumer manager only", workers: [][]string{nil, {"manager"}},
			message: `Consumer group "email-sender-beginning" has no running consumer manager.`, startConsumer: true,
		},
		{
			name: "delivery consumer", workers: [][]string{nil, {"delivery_consumer"}},
			message: `Consumer group "email-sender-beginning" has no running delivery consumer.`, startConsumer: true,
		},
		{
			name: "retention cleanup", workers: [][]string{{"stream_janitor"}},
			message: `Stream "signup.welcome-email" has no running retention cleanup worker.`, runManager: true,
		},
		{
			name: "vacuum", workers: [][]string{{"stream_vacuum"}},
			message: `Stream "signup.welcome-email" has no running vacuum worker.`, runManager: true,
		},
		{
			name: "cursor advancement", workers: [][]string{nil, {"cursor_advancer"}},
			message: `Consumer group "email-sender-beginning" has no running cursor advancement worker.`, runManager: true,
		},
		{
			name: "alert evaluation", workers: [][]string{nil, {"alert.partition_count"}},
			message: `Consumer group "email-sender-beginning" has no running "partition_count" alert evaluator.`, runManager: true,
		},
		{
			name: "consumers and maintenance", workers: [][]string{{"stream_janitor"}, {"exception_consumer", "manager", "message_consumer"}},
			message: `Stream "signup.welcome-email" has no running retention cleanup worker. Consumer group "email-sender-beginning" has no running retry consumer, consumer manager, or message consumer.`, startConsumer: true, runManager: true,
		},
		{
			name: "multiple maintenance workers", workers: [][]string{{"stream_janitor", "stream_vacuum"}},
			message: `Stream "signup.welcome-email" has no running retention cleanup worker or vacuum worker.`, runManager: true,
		},
		{
			name: "consumer and maintenance in one group", workers: [][]string{nil, {"message_consumer", "exception_consumer", "manager", "cursor_advancer"}},
			message: `Consumer group "email-sender-beginning" has no running message consumer, retry consumer, consumer manager, or cursor advancement worker.`, startConsumer: true, runManager: true,
		},
		{
			name: "unrecognized worker", workers: [][]string{nil, {"custom_worker"}},
			message: `Consumer group "email-sender-beginning" has no running worker "custom_worker".`, inspectWorker: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			// setup
			owner, err := common.NewStreamOwner(1, 42, "signup.welcome-email")
			if err != nil {
				t.Fatal(err)
			}
			beginning, err := common.NewConsumerGroupOwner(1, 42, 101, "email-sender-beginning")
			if err != nil {
				t.Fatal(err)
			}
			head, err := common.NewConsumerGroupOwner(1, 42, 102, "email-sender-head")
			if err != nil {
				t.Fatal(err)
			}
			owners := []*common.Owner{owner, beginning, head}
			var unclaimed []*metric.UnclaimedWorkerMetadata
			for i, names := range test.workers {
				for _, name := range names {
					metadata, err := metric.NewUnclaimedWorkerMetadata(name, owners[i], -1)
					if err != nil {
						t.Fatal(err)
					}
					unclaimed = append(unclaimed, metadata)
				}
			}

			// test
			finding, err := newWorkerLivenessAlert(owner, unclaimed, time.Date(2026, 9, 17, 12, 0, 0, 0, time.UTC))

			// verify
			if err != nil {
				t.Fatal(err)
			}
			if finding.Message != test.message {
				t.Errorf("newWorkerLivenessAlert(%s) message = %q, want %q", test.name, finding.Message, test.message)
			}
			for hint, want := range map[string]bool{
				"Start a consumer for each affected group, or delete any group you no longer need.": test.startConsumer,
				`Run "sqlstreams manager run" to resume maintenance.`:                               test.runManager,
				"Check the process responsible for the named workers.":                              test.inspectWorker,
			} {
				if got := strings.Contains(finding.Hint, hint); got != want {
					t.Errorf("newWorkerLivenessAlert(%s) contains hint %q = %t, want %t", test.name, hint, got, want)
				}
			}
			if finding.Detail != "" || finding.Severity != alert.AlertSeverityInfo {
				t.Errorf("newWorkerLivenessAlert(%s) detail, severity = %q, %s, want empty, info", test.name, finding.Detail, finding.Severity)
			}
			rows := finding.Data["workers"].([]map[string]any)
			if finding.Data["unclaimed_count"] != len(unclaimed) || len(rows) != len(unclaimed) {
				t.Fatalf("newWorkerLivenessAlert(%s) count, workers = %v, %d, want %d, %d", test.name, finding.Data["unclaimed_count"], len(rows), len(unclaimed), len(unclaimed))
			}
			for i, row := range rows {
				if row["worker"] != unclaimed[i].Name || row["owner"] != unclaimed[i].Owner.Name || row["owner_kind"] != string(unclaimed[i].Owner.Kind()) || row["target_instances"] != unclaimed[i].TargetInstances {
					t.Errorf("newWorkerLivenessAlert(%s) worker %d = %v, want retained evidence for %+v", test.name, i, row, unclaimed[i])
				}
			}
		})
	}
}
