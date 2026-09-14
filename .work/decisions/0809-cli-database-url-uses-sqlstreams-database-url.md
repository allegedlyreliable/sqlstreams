---
status: accepted
date: 2026-09-14
phase: pre-v1
---

# CLI database URL uses SQLSTREAMS_DATABASE_URL

## Context

The user requested SQLSTREAMS_DATABASE_URL for every current use of the
CLI database URL environment variable.

## Decision

Rename SQLSTREAMS_ADMIN_DATABASE_URL to SQLSTREAMS_DATABASE_URL in CLI
connection resolution, help, documentation, and environment examples.
This supersedes the database URL variable spelling in [0727]. The
--database-url flag still takes precedence over the environment variable.

## Consequences

Existing deployments must rename the variable; the old spelling is no
longer read. When no URL flag or new variable is provided, the CLI's
usage error names SQLSTREAMS_DATABASE_URL. SQLSTREAMS_ADMIN_SCHEMA and
the required database privileges are unchanged.
