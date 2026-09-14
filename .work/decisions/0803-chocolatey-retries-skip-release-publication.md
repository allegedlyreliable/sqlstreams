---
status: superseded
date: 2026-09-14
phase: pre-v1
---

# Chocolatey retries skip release publication

Superseded by [0804](0804-chocolatey-recovery-uses-a-new-release.md).

## Context

The v0.1.4 archives and Homebrew cask published successfully, but Chocolatey
submission returned HTTP 403 while the first package awaited approval.
Version 0.1.2 was approved minutes later. The runner retained no package
artifact, and the publishing key is available only in GitHub Actions.
Repeating normal publication would attempt to update an immutable release.

## Decision

Add a manual tag input to the existing release workflow. Check out that
tag and run GoReleaser with --skip=publish on manual dispatch. Reuse the
existing Chocolatey install, version, uninstall, and submission step.
Tag pushes retain normal publication.

The regenerated package must install from the existing release archive
and pass its checksum verification before submission. A mismatch fails
the retry; it does not rewrite checksums or replace published archives.

## Consequences

No second packaging implementation or new publishing credential is needed.
Retries rebuild all configured targets and depend on reproducible packaging.
The workflow must be committed to main before it can be dispatched.

References: https://goreleaser.com/getting-started/quick-start/ and
https://docs.github.com/en/actions/how-tos/manage-workflow-runs/manually-run-a-workflow
