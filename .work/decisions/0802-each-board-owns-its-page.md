---
status: accepted
date: 2026-09-13
phase: pre-v1
---

# Each board owns its page

## Context

A single board route and configuration accumulated special cases for
Troubleshooting sections, introductory content, and decision metadata.
Changing a board required following generic dispatch instead of its page.

## Decision

Give each board its own Astro page and local board definition. Each page
composes the existing layout, breadcrumb, sections, rows, and navigation.
Troubleshooting owns its intro and code grouping; Decision records owns
its status and decision-date column labels.

Remove the dynamic board route and shared section dispatcher. Keep a small
registry of the local definitions for the home page, jump navigation,
thread membership, previous/next links, and unread tracking. That registry
contains no intro content or presentation decisions.

## Consequences

Board URLs, thread membership, ordering, redirects, and read tracking stay
stable. Some page scaffolding repeats, letting each board change without
adding another branch or field to the shared rendering system.
