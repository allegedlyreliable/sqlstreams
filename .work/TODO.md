# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## CLI schema environment variable

- Rename the CLI schema environment variable to `SQLSTREAMS_SCHEMA` and
  update its reference documentation. Flag precedence and the default
  `sqlstreams` schema stay the same. Update existing shell configuration to
  the new name; the previous name is no longer read.
- Preserve historical decision records and the user's THOUGHTS.md.
- Verified: CLI module `go fmt ./...`, `go build ./...`, `go vet ./...`,
  and `go test -race ./internal/cli` pass. Ready for review; close out after commit.
