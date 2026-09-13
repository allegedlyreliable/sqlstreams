package controller

import (
	"testing"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
	"github.com/allegedlyreliable/sqlstreams/pkg/worker"
)

// Behavior: invalid completion warns while absence uses manager coverage and the existing pending policy.
func TestCollectorHistoryDistinguishesInvalidEvidenceFromMissingCoverage(t *testing.T) {
	for _, test := range []struct {
		name          string
		completionAge time.Duration
		storageAge    time.Duration
		missing       bool
		invalidValue  bool
		wrongKind     bool
		coverage      time.Duration
		wantState     alert.AlertEvaluationState
		wantInvalid   bool
	}{
		{name: "fresh completion", wantState: alert.AlertEvaluationStateHealthy},
		{name: "absent without coverage", missing: true, wantState: alert.AlertEvaluationStateInsufficientEvidence},
		{name: "overdue without coverage", completionAge: 3 * time.Minute, wantState: alert.AlertEvaluationStateInsufficientEvidence},
		{name: "absent while starting", missing: true, coverage: time.Minute, wantState: alert.AlertEvaluationStatePending},
		{name: "absent while running", missing: true, coverage: 2 * time.Minute, wantState: alert.AlertEvaluationStateActive},
		{name: "overdue while running", completionAge: 5 * time.Minute, coverage: 2 * time.Minute, wantState: alert.AlertEvaluationStateActive},
		{name: "invalid value", invalidValue: true, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantInvalid: true},
		{name: "wrong kind", wrongKind: true, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantInvalid: true},
		{name: "future completion", completionAge: -time.Minute, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantInvalid: true},
		{name: "future storage", storageAge: -time.Minute, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantInvalid: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			// setup
			current := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
			owner, err := common.NewSystemOwner(1)
			if err != nil {
				t.Fatal(err)
			}
			policy, err := alert.NewJobPayload(0, 0, 0, false)
			if err != nil {
				t.Fatal(err)
			}
			var completion *common.StoredMessage[metric.Measurement]
			if !test.missing {
				measurement, err := metric.NewBuiltInMeasurement(metric.MetricCollectorCompletedTimestamp, float64(current.Add(-test.completionAge).Unix()), nil, current)
				if err != nil {
					t.Fatal(err)
				}
				if test.invalidValue {
					measurement.Value = -1
				}
				if test.wrongKind {
					measurement.Kind = metric.MetricKindCounter
				}
				completion = &common.StoredMessage[metric.Measurement]{Id: 7104, CreatedAt: current.Add(-test.storageAge), Message: measurement}
			}
			history := &worker.WorkerInstanceHistory{EvaluatedAt: current}
			if test.coverage > 0 {
				history.Instances = []worker.WorkerInstanceSnapshot{{CreatedAt: current.Add(-test.coverage), ExpiresAt: current.Add(time.Minute)}}
			}
			var controller CollectorProgressController

			// test
			result, err := controller.evaluateHistory(owner, completion, history, 2*time.Minute, policy)

			// verify
			if err != nil {
				t.Fatalf("evaluateHistory(%s) = %v, want nil", test.name, err)
			}
			if result.State != test.wantState || result.EvidenceInvalid != test.wantInvalid {
				t.Errorf("evaluateHistory(%s) = %s, %v, want %s, %v", test.name, result.State, result.EvidenceInvalid, test.wantState, test.wantInvalid)
			}
			if (result.Reason != "") != (test.wantState == alert.AlertEvaluationStateInsufficientEvidence) {
				t.Errorf("evaluateHistory(%s).Reason = %q, want a reason only for insufficient evidence", test.name, result.Reason)
			}
		})
	}
}
