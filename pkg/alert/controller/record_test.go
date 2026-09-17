package controller

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/alert"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/common/diagnostic"
)

// Behavior: activation logs use the built-in alert's severity, while resolutions use INFO.
func TestAlertTransitionsLogAtTheirSeverity(t *testing.T) {
	for _, test := range []struct {
		name     string
		status   alert.AlertStatus
		severity string
		level    string
		detail   string
	}{
		{name: "worker_liveness", status: alert.AlertStatusActive, severity: "info", level: "INFO"},
		{name: "partition_count", status: alert.AlertStatusActive, severity: "warn", level: "WARN", detail: "partition evidence"},
		{name: "compaction_read_cost", status: alert.AlertStatusActive, severity: "warn", level: "WARN"},
		{name: "metrics_collector_progress", status: alert.AlertStatusActive, severity: "warn", level: "WARN"},
		{name: "worker_liveness", status: alert.AlertStatusResolved, severity: "info", level: "INFO"},
		{name: "partition_count", status: alert.AlertStatusResolved, severity: "warn", level: "INFO"},
	} {
		t.Run(test.name+"/"+string(test.status), func(t *testing.T) {
			// setup
			var output bytes.Buffer
			var controller AlertController
			controller.Logger = slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelInfo}))
			owner, err := common.NewStreamOwner(1, 42, "signup.welcome-email")
			if err != nil {
				t.Fatal(err)
			}
			if test.name == "metrics_collector_progress" {
				owner, err = common.NewSystemOwner(1)
				if err != nil {
					t.Fatal(err)
				}
			}
			definition, ok := diagnostic.GetAlert(test.name)
			if !ok {
				t.Fatalf("GetAlert(%s) = absent, want a built-in alert", test.name)
			}
			published, err := alert.NewAlert(test.name, owner, test.status, alert.AlertSeverity(definition.Severity), "condition observed", time.Now(), &alert.AlertOptions{Detail: test.detail})
			if err != nil {
				t.Fatal(err)
			}

			// test
			controller.logAlerts(t.Context(), published)

			// verify
			if string(published.Severity) != test.severity {
				t.Errorf("logAlerts(%s) severity = %s, want %s", test.name, published.Severity, test.severity)
			}
			var record map[string]any
			if err := json.Unmarshal(output.Bytes(), &record); err != nil {
				t.Fatalf("logAlerts(%s) output = %q, want one JSON record: %v", test.name, output.String(), err)
			}
			if record["level"] != test.level || record["alert"] != test.name {
				t.Errorf("logAlerts(%s) level, alert = %v, %v, want %s, %s", test.name, record["level"], record["alert"], test.level, test.name)
			}
			detail, present := record["detail"]
			if present != (test.detail != "") || (present && detail != test.detail) {
				t.Errorf("logAlerts(%s) detail = %v, present %t, want %q, present %t", test.name, detail, present, test.detail, test.detail != "")
			}
		})
	}
}

// Behavior: insufficient evidence logs its cause without accessing persistence.
func TestRecordClassifiesEvidenceWithoutWritingAlerts(t *testing.T) {
	for _, test := range []struct {
		name    string
		system  bool
		invalid bool
		reason  string
		level   string
	}{
		{name: "absent", reason: "no retained measurement", level: "DEBUG"},
		{name: "stale", reason: "newest measurement exceeds MaximumGap", level: "DEBUG"},
		{name: "uncovered", system: true, reason: "no current manager lease coverage", level: "DEBUG"},
		{name: "invalid completion", system: true, invalid: true, reason: "completion timestamp must be a positive integer within int64 range, got -1", level: "WARN"},
		{name: "invalid", invalid: true, reason: "partition count must be a non-negative integer within int64 range", level: "WARN"},
		{name: "display text does not classify", invalid: true, reason: "no retained measurement", level: "WARN"},
	} {
		t.Run(test.name, func(t *testing.T) {
			// setup
			var output bytes.Buffer
			var controller AlertController
			controller.Logger = slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
			owner, err := common.NewStreamOwner(1, 42, "signup.welcome-email")
			if err != nil {
				t.Fatal(err)
			}
			name := alert.AlertPartitionCount.Name
			if test.system {
				owner, err = common.NewSystemOwner(1)
				if err != nil {
					t.Fatal(err)
				}
				name = alert.AlertMetricsCollectorProgress.Name
			}
			result, err := alert.NewAlertEvaluationSnapshot(alert.AlertEvaluationStateInsufficientEvidence, nil, &alert.AlertEvaluationSnapshotConfig{
				EvidenceInvalid: test.invalid, Reason: test.reason,
			})
			if err != nil {
				t.Fatal(err)
			}

			// test
			outcome, err := controller.Record(t.Context(), name, owner, result)

			// verify
			if err != nil || outcome != alert.RecordOutcomeNothing {
				t.Fatalf("Record(%s) = %s, %v, want nothing, nil", test.name, outcome, err)
			}
			var record map[string]any
			if err := json.Unmarshal(output.Bytes(), &record); err != nil {
				t.Fatalf("Record(%s) output = %q, want one JSON record: %v", test.name, output.String(), err)
			}
			for key, want := range map[string]any{
				"level": test.level, "alert": name, "detail": test.reason,
				"owner": owner.Name, "owner_kind": string(owner.Kind()), "system_id": float64(owner.SystemId), "stream_id": float64(owner.StreamId), "group_id": float64(owner.ConsumerGroupId),
			} {
				if record[key] != want {
					t.Errorf("Record(%s) %s = %v, want %v", test.name, key, record[key], want)
				}
			}
			if test.invalid && record["code"] != alert.EventAlertEvidenceInvalid.GetCode() {
				t.Errorf("Record(%s) code = %v, want %s", test.name, record["code"], alert.EventAlertEvidenceInvalid.GetCode())
			}
		})
	}
}
