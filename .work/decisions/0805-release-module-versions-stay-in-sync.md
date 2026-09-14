---
status: accepted
date: 2026-09-14
phase: pre-v1
---

# Release module versions stay in sync

## Context

Skipping unchanged OTel and CLI modules left their tags and CLI installation
instructions at v0.1.4 after root v0.1.5 shipped.

## Decision

Every release publishes matching root, OTel, and CLI versions, even without
code changes. Publish in dependency order: root, OTel, then CLI. Update each
module's SQLStreams dependency pins before its tag and validate outside
go.work. Update all current installation references and verify the published
CLI go install command and version output.

## Consequences

Each release needs staged maintainer commits and module tags. Historical
records, frozen docs, and the compatibility lab's prior-version pin stay intact.
