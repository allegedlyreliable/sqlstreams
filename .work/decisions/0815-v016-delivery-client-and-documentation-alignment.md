---
status: accepted
date: 2026-09-17
phase: pre-v1
---

# v0.1.6 delivery, client, and documentation alignment

## Context

Ordered deferrals advanced attempt counters before a consumer handler ran.
Runnable examples exposed cancellation during shutdown recording. Public
config hovers exposed owning-package types, and CLI help and documentation
had drifted from the implementation. These changes shipped in root v0.1.6.

Partly supersedes [0813](0813-documentation-uses-four-boards-and-grouped-articles.md)
on board count and [0764](0764-avatar-image-sizes.md) on avatar files and sizing.

## Decision

- Store the current or next zero-based delivery attempt. Failures, requested
  delays, and expired leases advance it. Fresh claims and deferrals do not.
  Requested delays remain outside MaxRetries and its backoff position.
  First failures on either path use ExceptionInitialBackoff.
- Record completed consumer work under RecordMargin independently of lifecycle
  cancellation. Canceled startup alert checks stop quietly. The signal suite
  verifies active drain as well as idle, producer, crash, and forced shutdown.
- Mirror composite client inputs when gopls exposes owning-package field types.
  Keep public field names, methods, defaults, validation, and comments aligned.
  Convert at API boundaries. Copy owning comments onto aliases and re-exports
  directly, without a generator. Keep comments brief and useful to callers.
- CLI help describes actual scope, results, limits, exits, and output modes.
  Reject conflicting values for a metric attribute. An empty-value filter
  still requires the attribute to exist. The schema variable is SQLSTREAMS_SCHEMA.
- Keep all 13 examples runnable and verify their stored outcomes and reruns.
  Correct the example PostgreSQL volume mount so restarts retain data.
- Add Overview as the first board, with Quickstart, Why SQLStreams, Benchmark,
  and Roadmap. The first two remain homepage stickies. Keep Concepts, Guides,
  Reference, and Troubleshooting and their existing article-group navigation.
  Phone navigation shows Overview, Concepts, and Guides.
- Replace the cat avatar with i-just-woke-up-like-this.png and 66px/132px WebP
  derivatives. The user chose 1.5 times display width for the sharper appearance.
  Review root Markdown, public comments, and site content against actual API
  behavior. Use consistent terminology and short prose without semicolons.

## Consequences

v0.1.6 requires a fresh database. Historical exception counters are ambiguous,
so there is no automatic conversion or supported mixed-consumer rollout.
Both migration scopes stay at v1 and cannot enforce this restriction.
The prior-client round-trip does not establish existing-database compatibility.

Mirrored public structs no longer share owning-package type identity. Callers
using those packages must use sqlstreams config types at client boundaries.
The mirrors require manual upkeep. The docs now contain 96 articles across
five boards. HISTORY records the release verification and publication evidence.
