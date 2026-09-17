# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Release v0.1.6

- Source 2bd417e1 passes `just verify` and CI's PostgreSQL 15–18 matrix.
  Five signal cases and the v0.1.5 compatibility round-trip pass with race
  detection on an isolated, fresh PostgreSQL 18.6 database.
- Compatibility now pins v0.1.5. Release notes and the migration table state
  the fresh-database requirement for changed exception attempt counters.
- Root v0.1.6 is published at 2732e695. Pre-tag CI 35276192181, tag CI
  35276629543, and release run 35276629662 pass. All six archive checksums,
  the local archive binary, and the Homebrew upgrade are verified. Windows
  package checks pass and Chocolatey is submitted, pending moderation.
  The published notes include the fresh-database restriction.
- OTel and all four development modules pin root v0.1.6. Standalone
  tidy/build and OTel vet/race tests pass. The external Go Quickstart's
  updated pin builds and its fresh-database produce/consume walkthrough passes.
- OTel v0.1.6 is published at 723a4c9e after CI 35278931468 passed.
  The CLI's root/OTel pins now use v0.1.6 and pass standalone tidy, build,
  vet, and race tests. Await the maintainer's commit and push before tagging
  cmd/sqlstreams/v0.1.6, then verify go install and binary module versions.
- The Go Quickstart still needs its go.mod and go.sum committed and pushed.
- Site checks pass, including 83 unit tests and 45 browser tests against a
  fresh preview on a separate port. Version references are committed.
- Obtain approval before site deployment. Publish and freeze v0-1-6 after
  the CLI module is available. Verify both deployments and preserve old aliases.
- Committed implementation tasks are closed out in [0815]. Finish this release
  entry with CLI publication, installation, and documentation evidence.
