package controller

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
)

// Behavior: semantic rejection survives history evaluation while decode failures remain errors.
func TestHistoryClassifiesInvalidMeasurementsAndPreservesDecodeErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		value      float64
		kind       metric.MetricKind
		metadata   string
		wantState  alert.AlertEvaluationState
		wantReason string
		wantError  bool
	}{
		{name: "valid", metadata: `{"compaction_status":"uncompacted"}`, wantState: alert.AlertEvaluationStateHealthy},
		{name: "wrong kind", kind: metric.MetricKindCounter, wantReason: "partition measurement must be the stream partitions gauge with the declared unit"},
		{name: "invalid count", value: -1, wantReason: "partition count must be a non-negative integer within int64 range"},
		{name: "missing metadata", wantReason: "partition measurement must include compaction metadata"},
		{name: "unknown status", metadata: `{"compaction_status":"unknown"}`, wantReason: "compaction status must be compacted or uncompacted"},
		{name: "decode error", metadata: `{"compaction_status":42}`, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			// setup
			current := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
			owner, err := common.NewStreamOwner(1, 42, "signup.welcome-email")
			if err != nil {
				t.Fatal(err)
			}
			measurement, err := metric.NewBuiltInMeasurement(metric.MetricStreamPartitions, test.value, map[string]string{"stream": owner.Name}, current)
			if err != nil {
				t.Fatal(err)
			}
			measurement.Metadata = json.RawMessage(test.metadata)
			if test.kind != "" {
				measurement.Kind = test.kind
			}
			history := &metric.MeasurementHistory{EvaluatedAt: current, Messages: []*common.StoredMessage[metric.Measurement]{
				{Id: 7104, CreatedAt: current, Message: measurement},
			}}
			policy, err := alert.NewJobPayload(0, 0, 0, true)
			if err != nil {
				t.Fatal(err)
			}
			var controller CompactionReadCostController

			// test
			result, err := controller.evaluateHistory(owner, policy, history)

			// verify
			if test.wantError {
				var decodeError *json.UnmarshalTypeError
				if !errors.As(err, &decodeError) || result != nil {
					t.Fatalf("evaluateHistory(%s) = %v, %v, want nil and a decode error", test.name, result, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("evaluateHistory(%s) = %v, want nil", test.name, err)
			}
			wantState := test.wantState
			if test.wantReason != "" {
				wantState = alert.AlertEvaluationStateInsufficientEvidence
			}
			if result.State != wantState || result.EvidenceInvalid != (test.wantReason != "") || result.Reason != test.wantReason {
				t.Errorf("evaluateHistory(%s) = %s, %v, %q, want %s, %v, %q", test.name, result.State, result.EvidenceInvalid, result.Reason, wantState, test.wantReason != "", test.wantReason)
			}
		})
	}
}
