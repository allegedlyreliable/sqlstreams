# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Overview board

- Restore benchmarks.mdx and roadmap.mdx from b06a5d2c^ unchanged; both had
  no redirect and README's /benchmarks/ link was a live 404.
- Add a fifth board, Overview, first in the board list, listing quickstart,
  why-sqlstreams, benchmarks, and roadmap. Quickstart and Why SQLStreams stay
  homepage stickies and are also the board's first two articles.
- Site CONVENTIONS updated: five boards, stickies own an Overview home,
  Benchmark/Roadmap page rules. Close-out record supersedes 0813's "no fifth
  board" line; HISTORY entry at close-out.
- Phone header: Board Index, Reference, and Troubleshooting links hide under
  the 639px collapse (`inPhoneNav` on each board); nav and version bar side
  padding match the page content there.
- Avatar swap: `public/i-just-woke-up-like-this.png` (500x500 RGBA) replaces
  cat.png as the original; derivatives `-66.webp` (1,434 bytes) and
  `-132.webp` (3,388 bytes) via Sharp resize + webp quality 75, effort 6,
  alphaQuality 55. The user chose 1.5x display width over the 2x rule in
  0764: a hint of pixelation on 2x screens suits the site's look, 1x was too
  blocky, and 2x files at any quality only went soft. The removebg edge is
  16% partial-alpha pixels, so alphaQuality matters as much as quality.
  cat.png and its two derivatives deleted. Close-out record supersedes
  0764's file names and 2x rule.

## Delivery attempt numbering

- Implement approved option A: attempts stores the current/next zero-based
  delivery attempt. A fresh exception claim does not increment it; failures,
  requested delays, and expired leases advance it. Deferrals do not.
- Keep the existing delivery budget: MaxRetries retries after the initial
  attempt, with handler-requested delays excluded. First failures on either
  path use ExceptionInitialBackoff; later failures use the retry curve.
- Pre-v1 baseline change: use a fresh database, not mixed old/new consumers
  over existing exception rows. No automatic conversion of ambiguous historical
  counters or edits to the user's quickstart database.
- Existing tests whose expected counters encode claim-time increments will
  change to the approved outcome-time semantics; lease and ordering assertions
  stay intact. Record the new regression failing before the implementation.
- The ordered-successor regression failed at attempt 1 before the change in
  all three delivery-log modes; it now passes at 0, including repeated deferral.
  Existing exception outcome assertions now expect the unchanged claim number
  and increments on recorded delays; RecordFailure calls pass initial backoff.
- Verification: consumer and stream PostgreSQL integration suites pass with
  race detection, including outcome/backoff/expiry history checks. Targeted
  root race tests and vet pass. A separate PostgreSQL instance reproduced the
  quickstart handler calls user-1:0, user-1:1, user-2:0, user-2:1 and stored
  failure:0/success:1 for both messages, with no remaining exception rows.
  Temporary driver: `/tmp/sqlstreams-attempt-check.z635ajza/`.
- Expand attempt integration coverage from expected delivery behavior:
  first outcomes record zero; failure/delay advance the next attempt;
  deferral, success, terminal, and supersession do not advance it. Delays
  preserve the failure budget and backoff position. Each expired lease
  advances once; live polling, renewal, stale results, and duplicate outcomes
  cannot advance it. Counters belong to each message and consumer group.
- Cover cursor retries, partial commits, range replays/surrender/quarantine,
  ordered successors, zero/final retry limits, and stored delivery history
  across off/failures/all logging modes in attempt_test.go.
- Expanded attempt_test.go from 3 to 15 scenarios (60 leaf cases across
  logging modes). The complete consume PostgreSQL integration suite passes
  with `go test -race -count=1 ./integration/consume`. No production changes
  or changes to existing test expectations were needed for this expansion.
  The new cursor-success test initially supplied an outcome the production
  caller omits outside all mode; corrected its input, retaining the expected
  empty success history in failures/off modes.

## Public config hover

- Accepted option B: explicit public ProducerConfig,
  ConsumerConfig, StreamConfig, SystemConfig, SchedulerConfig,
  ScheduleRunOptions, and ProduceOptions structs with public field names.
  ClientConfig uses its public Logger and RetryPolicy aliases too.
- AGENTS.md specifies when to mirror a struct for gopls; CONVENTIONS.md
  records the exception to the one-declaration rule. Each mirrored struct's
  file carries the requested GOPLS comment below imports.
  Keep leaf types and read-models aliased. MessageOptions itself remains an
  alias, so this does not remove every internal package name from all hovers.
- Boundary pointer conversions preserve nil inputs and sharing; existing
  owners still implement defaults, validation, DeepCopy, and ToStream.
  Conversion compilation catches incompatible field changes, but comments
  and newly added methods still need manual synchronization. No generator.
- ProduceItem also becomes a public struct so its Options accepts the public
  ProduceOptions; ProduceBatch allocates converted items at its boundary.
- Public types no longer have their underlying packages' type identities.
  Changed the stream integration setup's one public Register call to use
  sqlstreams.StreamConfig; no assertion or behavior was changed.
- Verified ten gopls hovers from a separate consumer module: every changed
  composite's checked field uses its sqlstreams type, including
  ProducerConfig.Message as *sqlstreams.MessageOptions. Temporary LSP probe
  and evidence live outside the repository.
- Checks pass: root build/vet; client race tests (default mutation, independent
  deep copies, option capture, validation, and nil batch options); targeted
  alias/import/closure checks; all stream PostgreSQL integration tests.
  Examples, CLI, exporter, integration/e2e callers, and benchmarks compile;
  example vet passes. No existing test assertion changed.
- Eight public declaration/forwarding files plus boundary edits. Leaves and
  read-models retain their existing type identity. Accepted; changes remain
  uncommitted.

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
