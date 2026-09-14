---
status: accepted
date: 2026-09-13
phase: pre-v1
---

# Reference uses the board index

## Context

The API reference index repeated construction and method lists from the
individual reference pages. The Shape of the API already explains the
structure, and the Reference board provides navigation.

## Decision

Remove the duplicate API reference page and its board entry. Redirect
`/reference/` to `/boards/reference/`. Keep construction and method details
on the individual reference pages and the explanation on the concept page.

## Consequences

There is one fewer method list to maintain. Existing index links still
reach the Reference board, and individual reference URLs remain unchanged.
