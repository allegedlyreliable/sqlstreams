---
status: superseded
date: 2026-09-13
phase: pre-v1
---

# Troubleshooting uses the board index

Superseded by [0813](0813-documentation-uses-four-boards-and-grouped-articles.md).

## Context

The Error codes page duplicated the Troubleshooting board but listed only
57 of its 109 code pages. Its useful additions were lookup guidance,
grouping, and recovery, level, or metric-kind metadata.

## Decision

Use the Troubleshooting board as the single index, grouped into errors,
log events, metrics, and alerts. Show recovery, level, kind, or severity
beside each code. Read declaration metadata from the existing export and
log levels from code-page frontmatter. Keep each section in code order.
Move the short code-lookup guidance onto the board and redirect `/errors/`
to `/boards/troubleshooting/`. Remove the manually maintained index.

## Consequences

New code pages join the appropriate section automatically. Individual code
URLs and update-based read tracking remain unchanged. The board includes
alerts and metrics as well as errors and log events.
