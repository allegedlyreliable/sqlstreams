package controller

import (
	"testing"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/metric"
)

// Regression: invalid partition evidence retains its classification and reason through history.
func TestPartitionHistoryPreservesInvalidEvidence(t *testing.T) {
	// setup
	current := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
	owner, err := common.NewStreamOwner(1, 42, "signup.welcome-email")
	if err != nil {
		t.Fatal(err)
	}
	policy, err := alert.NewJobPayload(0, 0, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	history := &metric.MeasurementHistory{EvaluatedAt: current, Messages: []*common.StoredMessage[metric.Measurement]{
		{Id: 7104, CreatedAt: current, Message: partitionsGauge(-1)},
	}}
	var controller PartitionCountController

	// test
	result, err := controller.evaluateHistory(owner, policy, 200, history)

	// verify
	if err != nil {
		t.Fatalf("evaluateHistory(negative partition count) = %v, want nil", err)
	}
	if result.State != alert.AlertEvaluationStateInsufficientEvidence || !result.EvidenceInvalid {
		t.Errorf("evaluateHistory(negative partition count) = %s, %v, want insufficient_evidence, true", result.State, result.EvidenceInvalid)
	}
	if result.Reason != "partition count must be a non-negative integer within int64 range" {
		t.Errorf("evaluateHistory(negative partition count).Reason = %q, want the partition count requirement", result.Reason)
	}
}
