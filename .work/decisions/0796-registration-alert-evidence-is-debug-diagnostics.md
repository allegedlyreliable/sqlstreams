---
status: accepted
date: 2026-09-13
phase: pre-v1
---

# Registration alert evidence is debug diagnostics

## Context

Producer and consumer registration evaluate retained measurements before
collection necessarily runs. Missing evidence produced WARN messages across
user and system streams, discarded the evaluator's reason, and manufactured
an error. Successful quickstart startup looked degraded.

## Decision

- Registration logs insufficient evidence at DEBUG with the existing Reason
  in detail. Evaluation errors and pending/active findings retain WARN.
- Registration remains log-only and successful after an unsuccessful check.
  Evaluation, recording, thresholds, and scheduled checks remain unchanged.
- The accepted runtime direction is DEBUG for expected absence, the existing
  collector-progress alert for sustained collection failure, and WARN for
  unusable evidence. Complete a classification design before changing runtime
  levels; the first implementation changes registration only.

## Consequences

Default startup output stays quiet for missing evidence. DEBUG retains the
actual reason, including absent, stale, or unusable retained measurements.
No evaluator interface or public snapshot field changes are needed.

The runtime review found that stream checks count insufficient evidence as
failed evaluations without warning. Collector progress warns for both invalid
completion evidence and absent manager coverage. Its existing behavior [0703]
remains until the follow-up implements the accepted distinction.

Reason is display-only [0704], so runtime classification must not parse it.
The bounded stream-history read can represent expired samples as an empty
history, and the shared evaluator replaces invalid samples' detail with a
generic reason. Successful collection can still leave unusable evidence;
collector progress alone therefore cannot replace every runtime diagnostic.
