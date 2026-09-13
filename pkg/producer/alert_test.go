package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/stream"
)

type registrationEvaluator func(ctx context.Context, owner *common.Owner, policy *alert.JobPayload) (*alert.AlertEvaluationSnapshot, error)

func (e registrationEvaluator) Evaluate(ctx context.Context, owner *common.Owner, policy *alert.JobPayload) (*alert.AlertEvaluationSnapshot, error) {
	return e(ctx, owner, policy)
}

// Regression: registration keeps absent evidence diagnostic-only without hiding evaluation errors.
func TestRegistrationLogsInsufficientEvidenceAtDebugAndFailuresAtWarn(t *testing.T) {
	for _, test := range []struct {
		name      string
		err       error
		wantLevel string
	}{
		{name: "insufficient evidence", wantLevel: "DEBUG"},
		{name: "evaluation failed", err: errors.New("measurement read failed"), wantLevel: "WARN"},
	} {
		t.Run(test.name, func(t *testing.T) {
			// setup
			var output bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
			current := &stream.Stream{SystemId: 1, Id: 42, Name: "signup.welcome-email"}
			evaluator := registrationEvaluator(func(ctx context.Context, owner *common.Owner, policy *alert.JobPayload) (*alert.AlertEvaluationSnapshot, error) {
				if test.err != nil {
					return nil, test.err
				}
				return alert.NewAlertEvaluationSnapshot(alert.AlertEvaluationStateInsufficientEvidence, nil, &alert.AlertEvaluationSnapshotConfig{
					Reason: "no retained measurement",
				})
			})
			var p Producer

			// test
			p.logAlerts(t.Context(), current, logger, []alert.Evaluator{evaluator})

			// verify
			var record map[string]any
			if err := json.Unmarshal(output.Bytes(), &record); err != nil {
				t.Fatalf("logAlerts(%s) output = %q, want one JSON record: %v", test.name, output.String(), err)
			}
			if record["level"] != test.wantLevel {
				t.Errorf("logAlerts(%s) level = %v, want %s", test.name, record["level"], test.wantLevel)
			}
			if record["stream"] != current.Name {
				t.Errorf("logAlerts(%s) stream = %v, want %s", test.name, record["stream"], current.Name)
			}
			if test.err == nil {
				if record["detail"] != "no retained measurement" {
					t.Errorf("logAlerts(%s) detail = %v, want no retained measurement", test.name, record["detail"])
				}
				if _, exists := record["error"]; exists {
					t.Errorf("logAlerts(%s) error = %v, want absent", test.name, record["error"])
				}
			} else if _, exists := record["error"]; !exists {
				t.Errorf("logAlerts(%s) error = absent, want the evaluation error", test.name)
			}
		})
	}
}
