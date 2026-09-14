# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Documentation structure review

- [ ] Review the rewritten Consumer groups and delivery concept page and its
  three diagrams. Reusable writing guidance is in VOICE.md, Concept article
  playbook. Settle this example before applying it to other Concepts pages.
  Svelte Flow with ELK is the accepted diagram renderer. All three examples
  use the shared Diagram and ResourceCard components. Mermaid, Graphviz, and
  the original SVG comparison assets have been removed.

- [ ] Review the local preview: four purpose boards; Consumer groups in
  Concepts, Guides, and Reference; group breadcrumbs; dropdown tree.
  Membership now comes from article `group` frontmatter.
- Original Quickstart and Why SQLStreams remain homepage stickies above the
  sandbox. Legacy sources stay on disk but outside navigation and search.
- Latest build/lint checks pass; 16 articles are indexed. Browser checks covered
  navigation and responsive layouts; original sticky prose is unchanged (Quickstart adds group frontmatter).
- [ ] Settle remaining layout feedback, then map legacy articles and migration
  order. Resolve old URLs before deployment. Decision records are excluded
  from the docsite entirely; internal records remain in `.work/`.
- [ ] At task close-out when committed, write one consolidated decision and
  history entry. No records for intermediate tweaks. Agents do not commit.
