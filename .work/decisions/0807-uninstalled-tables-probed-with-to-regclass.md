---
status: accepted
date: 2026-09-14
phase: pre-v1
---

# Uninstalled tables are probed with to_regclass, never caught as 42P01

## Context

The catalog reads that tolerate an uninstalled database -- system.get,
stream.get, stream.getById, stream.list, migrate.SystemOwner -- ran their
SELECT and mapped SQLSTATE 42P01 (undefined_table) to "not registered"
([0347], [0624]). The client never saw the error, but Postgres logs every
statement error server-side, so each fresh install's first Register left
`ERROR: relation "sqlstreams.system_config" does not exist` in the
database log. The quickstart surfaced it on the first produce.

## Decision

A read that tolerates an uninstalled database asks the catalog first:
`SELECT to_regclass($1) IS NOT NULL` with the schema-qualified name built
in Go, and returns its absence answer without running the read when the
table is missing. The 42P01 branches are gone; after a positive probe a
missing table is an ordinary error.

The probe is a private on each datastore that needs it (system, stream)
and a helper in the migrate datastore, whose SystemOwner is a free func.
Duplication over a shared seam, per CONVENTIONS ## Structure.

## Consequences

- One extra round trip on those reads. Every caller is a Register,
  provision, admin, or collector path; no per-message path resolves a
  stream by name or id.
- The migration_log reads keep mapping only pgx.ErrNoRows: every id they
  take comes from a row in a table created in the same transaction as
  migration_log, so the table cannot be missing when they run.
- The existing integration tests (system Delete then Get, stream List
  and Get on an unregistered database, migrate SystemOwner on a database
  with no tables) pin the behavior unchanged. A fresh Postgres 17 run of
  examples/01-produce-only twice shows zero ERROR lines in the server log.
