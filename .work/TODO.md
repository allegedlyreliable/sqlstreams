# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Runnable example audit

- Set up and execute examples 01–13 against PostgreSQL. Check advertised
  behavior, persisted rows, reruns, and producer, consumer, and database logs.
- Preserve unrelated development data and capture per-example evidence.
- Executed all 13 against PostgreSQL 17.10, then reran the corrected examples
  on a fresh database and retained data. Used a separate Compose project on
  port 55432; the existing quickstart database was not changed.
- Fixed the PostgreSQL 17 volume mount. Reproduced loss of every SQLStreams
  table after the documented `down`/`up`; the corrected mount preserved an
  exact snapshot of messages, cursors, dead rows, keys, and compacted state.
- Fixed shared consumer declarations, missing retry inputs, throughput output,
  slow-handler timing, schedule success logging, and compacted initialization.
  Updated the README and misleading default, head-start, retry-position, and
  lifecycle comments. Real dead-letter and alert warnings remain visible.
- Active shutdown exposed a canceled range-commit context and canceled startup
  alert evaluations logged at WARN. Range recording now uses RecordMargin
  independently of lifecycle cancellation; producer/consumer alert checks stop
  quietly on caller cancellation while unrelated errors still warn.
- Regression evidence: new producer/consumer cancellation tests failed before
  their fixes. The signal drain case failed on the original commit warning.
  Its initial cursor assertion was corrected to a durable success-row assertion:
  the separate manager advances the committed cursor asynchronously. Existing
  tests were not weakened. The final signal suite passes all five cases.
- Checks: examples build/vet and targeted root race tests pass; existing
  PostgreSQL consumer commit/range integration tests pass. Ran the signal
  recipe's actual programs in an isolated Linux container against a fresh
  database, preserving the user's service on localhost:5432.
- Evidence and database dumps: `/tmp/sqlstreams-examples-audit.xuiS2S/`.
  Temporary audit drivers stay outside the repository; no maintained tool added.
- Final run: all 13 programs exited 0. Only 03's intentional dead-letter
  warning remained in application diagnostics; PostgreSQL recorded no warnings
  or errors. The pager also replayed five real worker-liveness activations from
  stopped examples, then their five resolutions after restart. Ready for review.
- Removed all audit containers and volumes after saving evidence; the existing
  quickstart database remains healthy. Changes are unstaged and uncommitted.

| Example | Runtime and stored-data verification |
| --- | --- |
| 01 | Produced the declared video payload and schema version; reruns append. |
| 02 | Consumed uploads and advanced its own cursor; runs alongside 03 without config replacement warnings. |
| 03 | Now produces all three outcomes itself. Corrupt stays dead with its cause; unavailable retries successfully; embargoed delays once then succeeds. |
| 04 | Business rows and both streams commit together. Forced rollback probes left no business or message rows. Rerun append behavior is documented. |
| 05 | One insert followed by a duplicate; reruns insert nothing while the key remains. |
| 06 | 1,000 distinct offsets stored and consumed; ten sample output lines instead of 1,000. |
| 07 | A new group skipped history and consumed a later upload; a restarted group resumed its saved cursor. |
| 08 | Each video's four successful states remained ordered across the injected failure; the other key progressed independently; no unresolved exceptions. |
| 09 | 600ms work used the 1s default; 4.75s work used the requested 10s timeout. Both completed without failures. SIGINT during the longer handler drained successfully and quietly after the library fix. |
| 10 | Real cron ticks produced timestamped reports and persisted success rows; declaration warning removed. Documented the one-minute producer poll. |
| 11 | Live and collected cursor metrics agreed with stored state; collection lag remained visible. |
| 12 | Healthy startup stayed quiet. Controlled threshold changes delivered nine activations and nine resolutions, confirmed in the alert stream. |
| 13 | First run reached count 2, rerun reached 3, then eight concurrent executions reached 11 with every intermediate version retained. |

## Client alias documentation

- Copy the owning declarations' comments to client aliases and re-exports
  so editor hover and Go documentation expose them at the supported entry point.
- Maintain the copies directly under an AGENTS.md rule; no generator or
  sync tool. CONVENTIONS.md retains the owning declaration as the source.
- Add missing source documentation for the four alert evaluation states and
  EventMeasurementsCannotBeExported, then copy it to their re-exports.
- Clarify Delay's time.Duration units in both comments: explicit millisecond,
  second, and minute examples; bare integers mean nanoseconds.
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
