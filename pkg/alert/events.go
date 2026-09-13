package alert

import (
	"github.com/allegedlyreliable/sqlstreams/pkg/common/diagnostic"
)

// EventAlertConditionHolds means a Register-time pass measured the stream
// against the built-in alert conditions and one of them held. The pass is
// log-only -- Register never fails on it.
var EventAlertConditionHolds = diagnostic.NewDiagnosticEvent("SQL0063",
	"alert condition holds",
	"nothing was published; the scheduled check is what publishes and resolves an alert")

// EventAlertEvidenceInvalid means a scheduled check cannot assess its owner
// because retained evidence fails semantic validation.
var EventAlertEvidenceInvalid = diagnostic.NewDiagnosticEvent("SQL0109",
	"alert evidence does not meet measurement requirements",
	"recorded alert unchanged")
