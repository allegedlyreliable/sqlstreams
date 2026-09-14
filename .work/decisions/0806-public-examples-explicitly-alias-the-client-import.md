---
status: accepted
date: 2026-09-14
phase: pre-v1
---

# Public examples explicitly alias the client import

## Context

The client import path ends in `client`, but its Go package name is
`sqlstreams`. Public samples omitted the alias, while the user's editor
adds it when importing the package.

## Decision

Use `sqlstreams "github.com/allegedlyreliable/sqlstreams/client"` in all
public-facing examples and documentation. Record the rule in CONVENTIONS.md.

## Consequences

The import spells the name callers use in the sample. This adds one word
per import without changing API behavior. Historical records remain as written.
