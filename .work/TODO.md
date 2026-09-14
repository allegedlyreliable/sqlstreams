# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Documentation structure review

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

## PostgreSQL 15–18 compatibility matrix [0812]

One chunk per review; each lands alone in the working tree with its check.

- [x] 1. Seam: `SQLSTREAMS_TEST_POSTGRES_IMAGE` override (default
  `postgres:18`) and a pid suffix on the test schema name in
  `.tests/integration/postgres/postgres.go`; `.env.example` gains the
  variable. Checked 2026-09-14: 103 tests green on the default image and
  on 15.19 (container seen); on a shared server the suffix removes the
  schema collisions, but the consume package still needs `-p 1` because
  the claim's snapshot fence sees other packages' in-flight transactions
  (10 failures parallel, 103 green with `-p 1`).
- [x] 2. Compose image 17 -> 18 in `.tools/database/docker-compose.yaml`.
  Checked 2026-09-14: fresh database on postgres:18, signal e2e four cases
  green, exit 0.
- [x] 3. CI: a `postgres` matrix job over the four major images running
  `just test-integration`, on pushes to main and on `v*` tags only; the
  `verify` job stays as is so pull requests run 18 once. Checked
  2026-09-14: the recipe under the override passes locally on postgres:16;
  the workflow itself is confirmed by the run on the next push to main.
- [x] 4. Evidence at 15a2dde5 (library source clean), 2026-09-14: the
  integration suite with the race detector, 103 tests in 9 packages, green
  on postgres:15.19, 16.15, 17.11, and 18.6; the signal e2e's four cases
  green on 15.19 and 16.15 (one benign Warn on 16 during the forced-exit
  case: the register-time alert pass saw the cancelled ctx). These are the
  tested minors the matrix page cites.
- [ ] 5. Doc page: the matrix on the install or getting-started page, linked
  from the migrations page's release-maintenance section; one sentence each
  for server-version support, PostgreSQL major upgrades, and schema
  upgrades. The page draft is the proposal, reviewed before it lands.
- [ ] 6. Close-out: HISTORY entry citing the evidence, the roadmap item
  removed.
