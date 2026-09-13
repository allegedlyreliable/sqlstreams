package evaluation

import (
	"testing"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
)

// Invariant: only the current evidence determines invalidity; older invalid samples break pending.
func TestHistoryPreservesCurrentEvidenceClassification(t *testing.T) {
	for _, test := range []struct {
		name        string
		states      []alert.AlertEvaluationState
		ages        []time.Duration
		invalidAt   int
		wantState   alert.AlertEvaluationState
		wantInvalid bool
		wantReason  string
	}{
		{name: "empty", invalidAt: -1, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantReason: "no retained measurement"},
		{name: "invalid newest", states: []alert.AlertEvaluationState{alert.AlertEvaluationStateInsufficientEvidence}, ages: []time.Duration{0}, invalidAt: 0, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantInvalid: true, wantReason: "invalid measurement value"},
		{name: "stale invalid", states: []alert.AlertEvaluationState{alert.AlertEvaluationStateInsufficientEvidence}, ages: []time.Duration{3 * time.Minute}, invalidAt: 0, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantReason: "newest measurement exceeds MaximumGap"},
		{name: "future", states: []alert.AlertEvaluationState{alert.AlertEvaluationStateHealthy}, ages: []time.Duration{-time.Minute}, invalidAt: -1, wantState: alert.AlertEvaluationStateInsufficientEvidence, wantInvalid: true, wantReason: "measurement storage time is after evaluation time"},
		{name: "fresh recovery", states: []alert.AlertEvaluationState{alert.AlertEvaluationStateHealthy, alert.AlertEvaluationStateInsufficientEvidence}, ages: []time.Duration{0, time.Minute}, invalidAt: 1, wantState: alert.AlertEvaluationStateHealthy},
		{name: "invalid breaks span", states: []alert.AlertEvaluationState{alert.AlertEvaluationStateActive, alert.AlertEvaluationStateInsufficientEvidence, alert.AlertEvaluationStateActive}, ages: []time.Duration{0, time.Minute, 2 * time.Minute}, invalidAt: 1, wantState: alert.AlertEvaluationStatePending},
	} {
		t.Run(test.name, func(t *testing.T) {
			// setup
			current := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
			samples := storedSamples(t, current, test.states, test.ages)
			if test.invalidAt >= 0 {
				samples[test.invalidAt].Message.EvidenceInvalid = true
				samples[test.invalidAt].Message.Reason = "invalid measurement value"
			}
			policy, err := alert.NewJobPayload(0, 0, 0, false)
			if err != nil {
				t.Fatal(err)
			}
			original := make([]alert.AlertEvaluationSnapshot, len(samples))
			for i, sample := range samples {
				original[i] = *sample.Message
			}

			// test
			result, err := EvaluateHistory(samples, current, policy)

			// verify
			if err != nil {
				t.Fatalf("EvaluateHistory(%s) = %v, want nil", test.name, err)
			}
			if result.State != test.wantState || result.EvidenceInvalid != test.wantInvalid || result.Reason != test.wantReason {
				t.Errorf("EvaluateHistory(%s) = %s, %v, %q, want %s, %v, %q", test.name, result.State, result.EvidenceInvalid, result.Reason, test.wantState, test.wantInvalid, test.wantReason)
			}
			for i, sample := range samples {
				if *sample.Message != original[i] {
					t.Errorf("EvaluateHistory(%s) changed sample %d, want read-only evaluation", test.name, i)
				}
			}
		})
	}
}

// Invariant: an invalid flag cannot accompany a healthy assessment.
func TestInvalidEvidenceCannotDescribeHealthyState(t *testing.T) {
	// setup
	result, err := alert.NewAlertEvaluationSnapshot(alert.AlertEvaluationStateInsufficientEvidence, nil, &alert.AlertEvaluationSnapshotConfig{EvidenceInvalid: true})
	if err != nil {
		t.Fatal(err)
	}
	result.State = alert.AlertEvaluationStateHealthy

	// test
	err = result.Validate()

	// verify
	if err == nil {
		t.Error("Validate(healthy with invalid evidence) = nil, want an inconsistent-state error")
	}
}
