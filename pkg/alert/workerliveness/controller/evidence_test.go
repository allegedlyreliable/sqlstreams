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
		{name: "valid", metadata: `{"workers":[]}`, wantState: alert.AlertEvaluationStateHealthy},
		{name: "wrong kind", kind: metric.MetricKindCounter, wantReason: "worker measurement must be the stream unclaimed workers gauge with the declared unit"},
		{name: "invalid count", value: -1, wantReason: "unclaimed worker count must be a finite non-negative integer"},
		{name: "missing metadata", wantReason: "worker measurement must include worker metadata"},
		{name: "missing workers", metadata: `{}`, wantReason: "worker details must account for the unclaimed worker count"},
		{name: "count mismatch", value: 1, metadata: `{"workers":[]}`, wantReason: "worker details must account for the unclaimed worker count"},
		{name: "missing identity", value: 1, metadata: `{"workers":[null]}`, wantReason: "each worker detail must include a name and owner"},
		{name: "wrong stream", value: 1, metadata: `{"workers":[{"worker":"janitor","owner":{"system_id":1,"stream_id":43,"owner":"other"},"target_instances":1}]}`, wantReason: "each worker detail must belong to the evaluated stream and have a nonzero target"},
		{name: "suspended worker", value: 1, metadata: `{"workers":[{"worker":"janitor","owner":{"system_id":1,"stream_id":42,"owner":"signup.welcome-email"},"target_instances":0}]}`, wantReason: "each worker detail must belong to the evaluated stream and have a nonzero target"},
		{name: "decode error", metadata: `{"workers":42}`, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			// setup
			current := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
			owner, err := common.NewStreamOwner(1, 42, "signup.welcome-email")
			if err != nil {
				t.Fatal(err)
			}
			measurement, err := metric.NewBuiltInMeasurement(metric.MetricStreamUnclaimedWorkers, test.value, map[string]string{"stream": owner.Name}, current)
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
			var controller WorkerLivenessController

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
