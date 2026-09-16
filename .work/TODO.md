# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Client alias documentation

- Copy the owning declarations' comments to client aliases and re-exports
  so editor hover and Go documentation expose them at the supported entry point.
- Maintain the copies directly under an AGENTS.md rule; no generator or
  sync tool. CONVENTIONS.md retains the owning declaration as the source.
- Add missing source documentation for the four alert evaluation states and
  EventMeasurementsCannotBeExported, then copy it to their re-exports.
- Verified all 197 copied comments against their sources and confirmed
  `gopls` hover documentation for Versioned at a client use site. Root
  `go fmt ./...`, `go build ./...`, and
  `go test -race ./client ./pkg/alert ./pkg/metric` pass. Ready for review;
  close out after commit.

## CLI schema environment variable

- Rename the CLI schema environment variable to `SQLSTREAMS_SCHEMA` and
  update its reference documentation. Flag precedence and the default
  `sqlstreams` schema stay the same. Update existing shell configuration to
  the new name; the previous name is no longer read.
- Preserve historical decision records and the user's THOUGHTS.md.
- Verified: CLI module `go fmt ./...`, `go build ./...`, `go vet ./...`,
  and `go test -race ./internal/cli` pass. Ready for review; close out after commit.
