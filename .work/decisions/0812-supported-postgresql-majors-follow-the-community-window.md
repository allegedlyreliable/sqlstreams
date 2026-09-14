---
status: accepted
date: 2026-09-14
phase: pre-v1
---

# Supported PostgreSQL majors follow the community support window

## Context

The integration suite passed on 15.19, 16.15, 17.11, and 18.6 against a
frozen snapshot, but nothing repeatable exists: the integration seam pins
`postgres:18`, the compose stack that signal-e2e and compat-lab use runs
`postgres:17`, and CI runs only 18. No SQL in the library gates a version
floor (SKIP LOCKED, RANGE partitioning, gen_random_uuid; UUIDv7 is minted
in Go), so the floor is a policy choice, and a published matrix without a
drop rule rots.

## Decision

The supported set is the PostgreSQL majors the community still supports,
verified per release. 14 leaves the window in November 2026, so the first
matrix is 15 to 18. A tested minor is evidence, never a floor: any minor of
a listed major is supported.

Verification is the integration suite, because its subject is the lock
manager and MVCC. The seam reads `SQLSTREAMS_TEST_POSTGRES_IMAGE` (default
`postgres:18`); CI runs a four-image matrix on pushes to main and on tags,
and pull requests stay on 18. The signal e2e runs by hand on 15 and 16
once at this checkpoint for the evidence line and is not part of the
matrix. compat-lab tests client version against tables and stays out.
The compose stack moves to 18 so every local surface matches CI.

The shared-server override's schema collision is a seam bug (per-process
`test_<n>` counter), fixed by a pid suffix. `-p 1` stays for that path:
the claim's snapshot fence declines while any other package's transaction
is in flight, which is the mechanism working, not a seam bug.

The matrix is published where the reader asks "can I run this on my
Postgres?": the install or getting-started page, linked from the
migrations page's release-maintenance section, never on the migrations
page itself. One sentence each distinguishes server-version support,
PostgreSQL major upgrades (operator-owned; plain tables and partitions,
no extensions), and SQLStreams schema upgrades (the migration registry).

## Consequences

The tag-triggered matrix job is the first concrete gate for repeatable
release assurance, which depends on it rather than inventing one. The
schema-upgrade table [0794] stays a separate adjacent table. A major is
dropped from the matrix when the community drops it, in the next release.
