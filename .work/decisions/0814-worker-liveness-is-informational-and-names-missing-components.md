---
status: accepted
date: 2026-09-17
phase: pre-v1
---

# Worker liveness is informational and names missing components

Partly supersedes [0627]'s worker-liveness warning policy and [0485]'s
activation log levels. Alert scope and status-change logging remain intact.

## Context

Stopping a test consumer leaves its group registered. Worker liveness cannot
infer whether the stop was intentional, so WARN made normal testing noisy.
The report listed internal rows and consequences that did not apply to every
finding. Its manager-run hint could not restart application consumers.

## Decision

- Worker liveness has INFO severity. Activation logs use the alert's severity,
  resolutions stay INFO, and the other built-in alerts remain WARN. The public
  vocabulary includes AlertSeverityInfo and the declaration check accepts it.
- Keep the existing condition, retained evidence, pending policy, publication,
  repeat interval, and resolution behavior. This changes reporting, not health.
- Use one sentence per resource listing missing components: message consumer,
  retry consumer, consumer manager, or named maintenance workers. Do not infer
  whole-group failure or suppress rows to produce a special summary.
- Pair application-consumer findings with starting a consumer or deleting a
  group no longer needed. Maintenance findings suggest sqlstreams manager run.
  Mixed findings include both actions; unrecognized workers retain their names
  and an inspection hint without assuming a manager can claim them.
- Keep worker-row evidence in Data. Detail is empty for worker liveness, and
  activation logs omit empty detail. A short coordinator and summary methods
  separate grouping, component wording, hints, and evidence formatting.
- Keep delayed-processing and overdue-retry alerts in ROADMAP's Later section.
  Their design must account for backoff, intentional delays, bindings,
  compaction, retained evidence, measurement cost, and actionable thresholds.

## Consequences

- The default WARN logger hides liveness activations, but alert consumers and
  API/CLI readers still receive them. INFO is not proof that a stop is harmless.
- New worker-liveness routing keys use the info severity segment. Consumers
  filtering on warn must change their bindings to keep receiving these alerts.
  Stored reports retain their old contents until a new version is published.
- Labels and recovery hints stay accurate for partial and mixed failures at
  the cost of listing each missing component rather than inferring group state.
- The alert condition and evidence storage need no datastore or migration change.
