---
status: accepted
date: 2026-09-13
phase: pre-v1
---

# Runtime alerts classify invalid evidence

## Context

The runtime policy in [0796] needs to distinguish unavailable evidence from
invalid measurements without parsing the display-only Reason [0704]. The
user reviewed and approved the alerts reference proposal, including the
snapshot field, shared logging, and retained read/decode-error behavior.

## Decision

- Add EvidenceInvalid to AlertEvaluationSnapshot and its config. True requires
  insufficient_evidence; false never establishes health. Keep existing states.
- The existing semantic checks classify decoded measurements and name their
  rejected requirement. EvaluateHistory preserves the newest sample's flag
  and reason. Future timestamps are invalid; absent/stale evidence and absent
  manager coverage are insufficient without establishing invalidity.
- Freshness takes precedence over a stale sample's invalid value. Older
  invalid samples break pending spans without overriding current evidence.
  Read/decode errors keep their error return and existing retry behavior.
- AlertController.Record logs insufficient evaluations at DEBUG, or the
  declared invalid-evidence event at WARN when EvidenceInvalid is true.
  Include alert, owner identity, and detail. No alert is written or resolved.
- Remove collector progress's separate warning. This supersedes [0703]'s
  blanket insufficient-evidence WARN only; its coverage and activation
  semantics remain. Keep failed-stream counts and registration DEBUG [0796].
- Snapshot stays read-only. Use existing suppression; no new datastore read,
  collection path, pending state, timer, or suppression key.

## Consequences

Routine evidence absence is quiet at the default level. Invalid evidence
remains visible even while collection completes. Suppression can collapse
different owners into one warning; snapshots and failed-stream metrics remain
the detailed inspection surface. No internal check runs when the entire
installation is stopped. Existing snapshot users gain one boolean field.

Grafana distinguishes no-data results from evaluation errors:
https://grafana.com/docs/grafana/latest/alerting/fundamentals/alert-rule-evaluation/nodata-and-error-states/.
Prometheus recommends actionable symptoms and monitoring the monitoring path:
https://prometheus.io/docs/practices/alerting/.
