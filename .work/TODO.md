# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Runtime alert evidence diagnostics [0796]

- Review the Proposed section in the alerts reference before implementation.
  Keep the existing evaluation states; add EvidenceInvalid to the snapshot and
  its config. Reason remains display-only. False does not establish health.
- Mark semantic rejection at the existing measurement checks and preserve its
  reason through EvaluateHistory. Missing/stale evidence and absent manager
  coverage stay insufficient with EvidenceInvalid false. Future timestamps
  and decoded but unusable measurements set it true. Existing read/decode
  errors retain their error return; no compatibility change to that contract.
- Classify scheduled-check logs in AlertController.Record's existing
  insufficient-evidence branch. Declare the invalid-evidence WARN event with
  code, alert, owner identity, and detail, plus its diagnostic page. Remove
  the separate collector-progress warning. Preserve failed-check counts,
  recorded alert state, and the existing collector-progress activation policy.
- Keep registration at DEBUG for all insufficient results and Snapshot
  read-only. Use the existing logger suppression, with its documented shared
  window; no new persistence, SQL, timers, or collection path.
- Verify semantic rejection and reason propagation in all four evaluators,
  history freshness precedence and pending-span breaks, and Record's log
  levels/attributes without database access. Run the existing tests unchanged,
  root build, and targeted race tests for touched packages and client aliases.
- After review settles the shape, record the decision in the same session.
  Update [0703]'s blanket warning decision and record that [0704]'s error
  contract and [0796]'s registration behavior remain in force.

Research: Grafana distinguishes successful empty reads from evaluation errors
(https://grafana.com/docs/grafana/latest/alerting/fundamentals/alert-rule-evaluation/nodata-and-error-states/).
Prometheus recommends actionable symptom alerts and monitoring the monitoring
path (https://prometheus.io/docs/practices/alerting/). The existing collector
progress alert owns missing collection; snapshots and WARN identify invalid
retained measurements without creating a second alert stream or state machine.
