---
status: superseded
date: 2026-09-13
phase: pre-v1
---

# Decision records use the board index

Superseded by [0813](0813-documentation-uses-four-boards-and-grouped-articles.md).

## Context

The Regrets page and Decision records board listed the same records.
The separate page added status, decision dates, and newest-first ordering.

## Decision

Use the Decision records board as the single index. Show status and the
record's decision date, with the highest record number first. Keep unread
tracking based on the last Git update. Redirect `/decisions/` to
`/boards/decisions/` and remove the duplicate page and index components.

## Consequences

Readers see the useful metadata without opening a second index. Individual
record URLs remain stable. Other boards still show update dates.
