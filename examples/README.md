# Examples

Every example you should run from repo root.
Read the file header comment first, it helps explain things.

## Install

- Go 1.27 or newer.
- Docker with the Compose plugin (`docker compose version`).

## Set up the database

From the repo root:

```sh
docker compose -f examples/docker-compose.yaml up -d --wait
```

Stop it, keeping the data:

```sh
docker compose -f examples/docker-compose.yaml down
```

Start over with an empty database:

```sh
docker compose -f examples/docker-compose.yaml down -v
```

## Run a scenario

```sh
go run ./examples/01-produce-only
```

Scenarios that run a consumer, a scheduler, or a manager keep running until
you press Ctrl-C; the rest print and exit. Ctrl-C drains current work before
exiting; a second Ctrl-C forces exit. A "Run first" column means the
scenario reads a stream or group another scenario creates.

| # | Scenario | Command | Run first | Runs until |
| --- | --- | --- | --- | --- |
| 01 | produce-only service | `go run ./examples/01-produce-only` | | exits |
| 02 | consume-only service | `go run ./examples/02-consume-only` | 01 | Ctrl-C |
| 03 | consume with retry and dead-lettering | `go run ./examples/03-consume-retry-dead` | | Ctrl-C |
| 04 | produce inside the caller's own transaction | `go run ./examples/04-produce-in-tx` | | exits |
| 05 | idempotent produce with a caller-supplied key | `go run ./examples/05-idempotent-produce` | | exits |
| 06 | tuning a stream for throughput | `go run ./examples/06-throughput` | | Ctrl-C |
| 07 | a consumer that starts at the head of the stream | `go run ./examples/07-consume-from-head` | 01 | Ctrl-C |
| 08 | keyed ordering | `go run ./examples/08-keyed-ordering` | | Ctrl-C |
| 09 | a longer timeout for a slow message | `go run ./examples/09-slow-handler` | | Ctrl-C |
| 10 | a schedule that produces on a cron expression | `go run ./examples/10-schedule-produce` | | Ctrl-C |
| 11 | reading what the system measures about itself | `go run ./examples/11-metrics-read` | 01, then 02 | exits |
| 12 | consuming `__system.alerts` as a pager feed | `go run ./examples/12-alert-consumer` | | Ctrl-C |
| 13 | a compacted stream used as a key/value store | `go run ./examples/13-compacted-kv` | | exits |

## What to expect

- **03** produces corrupt, temporarily unavailable, and embargoed uploads.
  The corrupt upload stays dead in the exception queue and emits one expected
  warning. The other two succeed after a retry or when the embargo expires
  about five seconds after startup.
  Its group is separate from 02, so running both does not replace their configs.
- **04** appends messages on every run, even when the business rows already
  exist. **05** demonstrates message deduplication: the first run appends once;
  reruns report duplicates until the stored idempotency key expires.
- **06** produces and consumes 1,000 thumbnails, printing ten sample results.
- **07** waits for new uploads on its first run. Start it, then run 01 again
  in another terminal. Later runs resume the group's saved cursor.
- **09** finishes both simulated videos in about six seconds, then waits for
  more work. The longer video uses its requested ten-second timeout.
- **10** schedules one report per UTC minute. The producer polls once a
  minute, so the first output can take up to two minutes. Each report prints
  its scheduled time, and the stream records successes as well as failures.
- **11** prints a live cursor snapshot and the latest collected backlog.
  Collected measurements can lag the live snapshot by a collector poll.
- **12** prints the current alert count, then waits. A healthy database can
  remain quiet; the feed prints alert activations and resolutions when they occur.
  Stopped example consumers can raise worker-liveness alerts; restarting them
  resolves those alerts. A newly created pager group also reads retained alerts.
- **13** initializes a missing key and increments it under a row lock. Reruns
  keep incrementing the current value rather than resetting it.

Consumer cursors, dead messages, idempotency keys, and compacted values survive
restarts. To repeat every scenario from empty state, use the database reset
command above. When upgrading an existing playground from the old volume
mount, export any data you need before using `down -v`, then start the database
and restore the export. PostgreSQL 17 previously wrote to an anonymous volume,
which Compose does not reattach after `down`.
