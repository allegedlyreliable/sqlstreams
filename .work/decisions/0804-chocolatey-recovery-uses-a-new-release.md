---
status: accepted
date: 2026-09-14
phase: pre-v1
---

# Chocolatey recovery uses a new release

## Context

The manual v0.1.4 retry rebuilt archives with different checksums and stopped
before submission. Recovering a package against existing archives adds more
packaging work than the user wants for this release.

## Decision

Supersede [0803](0803-chocolatey-retries-skip-release-publication.md).
Remove the manual retry trigger and restore the original tag-driven release
workflow. Publish v0.1.5 through the normal pipeline now that Chocolatey's
first package is approved. Keep published tags and archives unchanged.

## Consequences

The archive and Chocolatey installer are generated together again. There is
no manual retry mode or second packaging path to maintain. This costs one
patch release; v0.1.4 will not be submitted to Chocolatey.
