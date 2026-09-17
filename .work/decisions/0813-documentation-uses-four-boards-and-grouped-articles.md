---
status: accepted
date: 2026-09-16
phase: pre-v1
---

# Documentation uses four boards and grouped articles

## Context

The previous site mixed explanations, procedures, API contracts, and code
catalogues. One reference page per handle made readers learn the API's
structure before finding an operation. Review of a consumer subset established
the navigation and writing patterns before expanding the rest of the site.

Supersedes [0596](0596-decision-records-publish-as-a-board.md),
[0679](0679-the-doc-site-splits-by-page-kind-a-reference-board-one-thread-per-handle.md),
[0712](0712-docs-start-with-a-runnable-example-and-separate-delivery-outcomes.md),
[0799](0799-decision-records-use-the-board-index.md), and
[0800](0800-troubleshooting-uses-the-board-index.md).
Retains the Reference board index [0801] and local board routes [0802].

## Decision

- Keep four boards: Concepts explains mechanisms, Guides completes tasks,
  Reference defines contracts, and Troubleshooting diagnoses and fixes symptoms.
  The original Quickstart and Why SQLStreams are homepage stickies above the
  sandbox and boards. Decision records remain repository-only.
- Markdown `group` frontmatter owns grouping in Concepts, Guides, and Reference.
  Breadcrumbs include the group. A collapsible tree below the article title
  connects its group across boards without reducing content width. No subject
  page or separate full-map route duplicates that navigation.
- Reference follows recognizable resources and operations, informed by Stripe.
  Use code-first entries, complete configuration tables in declaration order,
  separate nested-struct tables, and consistent returns, errors, and CLI sections.
  Metrics, logs, and alerts define observations. Troubleshooting owns procedures.
- Keep general prose rules in VOICE.md and page-specific rules in CONCEPTS.md,
  GUIDES.md, REFERENCE.md, and TROUBLESHOOTING.md under .website/. Introduce terms
  and resource identities before relying on them. Preserve example names and
  exact diagnostic literals. Show diagnostic log fields with placeholders.
- Concepts use concrete causal examples and focused visuals. Guides follow
  setup, execution, then useful qualifications. Troubleshooting follows evidence,
  result, then action. Use ThreadAside for supporting notes and Related topics
  for useful onward links. Each explanation has one canonical home.
- Construct diagrams with Svelte Flow and ELK, using consistent resource headers,
  body alignment, node and text sizes. Choose the visual form for the mechanism,
  including queue timelines when a flowchart cannot explain timing clearly.
- Approve the minimal Concepts and Guides list first. Complete Reference and
  Troubleshooting one page at a time before writing the learning pages. Finish
  with accuracy, consistency, Related topics, and URL reviews. Preserve honest
  review markers. Additional learning topics require approval.

## Consequences

The site has 94 articles: two stickies, six Concepts, five Guides, 71 Reference,
and ten Troubleshooting. Retired sources and per-code rendering are removed.
Astro redirects retain 109 diagnostic URLs and 27 retired article URLs alongside
the two existing board redirects. Only current articles enter navigation and
search. Implementation is committed in b06a5d2c, and the user confirmed deployment.

The four-board choice is superseded by [0815](0815-v016-delivery-client-and-documentation-alignment.md),
which adds Overview and retains the article-group structure and writing rules.
