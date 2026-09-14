---
status: accepted
date: 2026-09-14
phase: pre-v1
---

# The producer's register-time pass drops worker_liveness

Partly supersedes [0627]: the alert itself and its schedule stand; only
the producer's register-time evaluator goes.

## Context

[0627] gave `Producer().Register` a fourth evaluator so a produce-only
deployment would learn at registration that nothing runs its stream's
rows. The consumer's pass never carried it: `Register` runs before
`Consume` claims anything, so it would warn about its own rows.

Running the quickstart showed the cost. A user runs `produce`, then
`consume`, stops it, and runs `produce` again: the second produce logs
SQL0063 for `worker_liveness` naming email-sender's three rows. The
producer is warning about a different process being stopped -- the
rows are the consumer group's, the producer holds none of them and
will claim none of them.

## Decision

The producer's evaluator set is `partition_count` and
`compaction_read_cost`, the same two the consumer's pass runs. The
scheduled check ([0627]) is the only reporter of unclaimed worker rows.
SQL0063 stays declared for the two remaining conditions.

## Consequences

- A produce-only deployment no longer hears at `Register` that nothing
  runs its rows. The site tells it to run `sqlstreams manager run`
  beside the producer; the `worker_liveness` alert on `__system.alerts`
  is where the fact lands, published by the manager that is missing --
  the accepted gap [0627] already carried for deployments without one.
- The SQL0063 page, the producer and manager reference pages, and the
  producer's imports change; no test observed the producer's liveness
  line, so none is edited.
