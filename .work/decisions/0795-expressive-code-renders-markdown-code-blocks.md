---
status: accepted
date: 2026-09-13
phase: pre-v1
---

# Expressive Code renders Markdown code blocks

## Context

Readers need to copy documentation examples, starting with the Quickstart.
Astro's default Shiki output provides highlighting without copy controls.
A browser island that moves rendered blocks and mounts buttons adds layout
changes and component cleanup to each page.

## Decision

Use astro-expressive-code for Markdown and MDX fences, before the MDX
integration. It replaces the default fence renderer and owns copy controls.
Keep the existing CopyButton for page links and diagnostic queries.

Move the existing code palettes into Expressive Code and select them using
data-board-style. Keep the board's fonts, spacing, square borders, and plain
frames. Disable filename-comment extraction and terminal-comment removal
so examples retain their comments. No per-page imports are required.

Astro documents the integration at
https://docs.astro.build/en/guides/syntax-highlighting/;
installation and configuration are documented at
https://expressive-code.com/installation/ and
https://expressive-code.com/reference/configuration/.

## Consequences

Code-block markup is produced during the build. Expressive Code supplies
the clipboard script and its navigation handling; no custom island wraps
the blocks after the page loads. The website gains one direct dependency
and uses its generated styles through supported configuration options.
