---
status: accepted
date: 2026-09-13
phase: pre-v1
---

# Threads declare a slop level

## Context

Every thread on the doc site is LLM-drafted against VOICE.md, and a
blanket "AI generated" warning tells a reader nothing about any one page.
What a reader needs is how far the author reviewed that page before
publishing it. The user's design review settled the scale, the placement,
and the two exemptions.

## Decision

- A `slop` frontmatter field on the docs collection, `z.enum` over
  `none | low | medium | high`, declared once in the slop-notice
  component's types module and imported by the content schema.
- The levels are facts about the author's review, not vibes: low is read
  line by line against the code and edited; medium is structure checked
  and code samples run, prose skimmed; high is unreviewed. Each level's
  body copy is fixed in one table beside the component.
- A new thread starts at `high`, the honest floor; the author lowers it
  when they review the page.
- `none` is a hand-checked thread and renders no notice; today only the
  quickstart. The diagnostics reference and the SQL-code threads carry no
  field at all -- diagnostics because the user exempted it, the code
  threads because they are declaration-derived and render through their
  own layout. Decision records render as-is and never carry one.
- ThreadLayout renders the notice as the first child of the thread body
  from the field alone; no MDX import, so the MDX component whitelist does
  not grow. The box is `data-pagefind-ignore` so "LLM" matches no search.
- The look is a thread-aside variant: same box, chip, and header strip; a
  `data-level` attribute tints the chip through six slop-named semantic
  tokens per board (the solved, amber, and refused primitives) and fills a
  three-block pixel meter.

## Consequences

Forty-five threads carry `slop: high` until reviewed, and the site says
so on each. Lowering a level is a one-line frontmatter edit with no
copy to write. A fourth rendered level or per-page wording would need a
new record. The Storybook done-checklist is the three rendered levels.
