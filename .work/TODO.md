# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Documentation expansion

Consumer subset review is complete. The four boards, group navigation,
writing rules, and reviewed pages are the accepted baseline.

### Approved Concepts and Guides

This is the complete approved list for this phase. Additional Concepts or
Guides require user approval. Retired sources have been removed. The legacy
source column below records their former paths, not an approved page list.

| Board | Page | Purpose and boundary | Status |
| --- | --- | --- | --- |
| Concepts | Consumer groups | Shared work and independent groups. | Reviewed |
| Concepts | Consumer leases | Shared lease deadlines, queue waiting, expiry, and reclaim. | Reviewed |
| Concepts | Streams and messages | What a stream stores and how retention relates to consumer progress. | Reviewed |
| Concepts | Background managers and workers | Where maintenance runs, what janitors do, and how managers coordinate workers across processes. Replaces Automatic producer batching. | Reviewed |
| Concepts | Retries and failed messages | What happens after processing fails, including retry exhaustion. | Reviewed |
| Concepts | Ordering and concurrency | How concurrent processing and per-key ordering affect which messages can run. | Reviewed |
| Guides | Stop a consumer gracefully | Stop a running consumer through the supported shutdown path. | Reviewed |
| Guides | Size a consumer queue | Configure how much work an instance claims ahead of processing. | Reviewed |
| Guides | Produce messages in a database transaction | Commit an application change and its message together. | Reviewed |
| Guides | Make message processing safe to repeat | One practical pattern for preventing duplicate application effects. | Reviewed |
| Guides | Upgrading SQLStreams | Follow the supported library and database upgrade procedure. | Reviewed |

Six Concepts and five Guides total. “Run background maintenance” is excluded.
Dedicated learning pages for scheduling, routing, compaction, replay, and
payload evolution remain deferred. Their shipped contracts still belong in
Reference. Keep the original Quickstart and Why SQLStreams stickies.

### Work order

- [x] Inventory the remaining shipped public API and diagnostics. Map the
  Reference and Troubleshooting pages, including source files, legacy routes,
  and the supporting pages required by each approved Concept or Guide.
  Keep the accepted resource-and-operation Reference model. Troubleshooting
  covers actionable problems, without automatically making a page per code.
- [x] Complete Reference and Troubleshooting first, one page at a time.
- [x] Create the approved Concepts and Guides after their supporting pages
  exist. Link to Reference and Troubleshooting where useful during drafting.
- [x] After all approved pages exist, review every page's Related topics
  section and fill out all relevant connections. Check terminology across
  pages, navigation, and links. Resolve legacy content and URLs before deployment.

### Reference coverage map

Paths below are under `reference/` and are page targets, not instructions to
restore hidden legacy pages. Existing reviewed Consumer entries stay canonical.
Each semicolon-separated target is a separate page task. Complete the core
stream/producer contracts before the remaining operations and diagnostics.

| Coverage | Page targets and boundaries | Source files | Legacy source |
| --- | --- | --- | --- |
| Stream identity and registration | `stream-object`: returned fields; `stream-registration`: Register and complete StreamConfig; `get-stream`; `list-streams`; `rename-stream`; `destroy-stream` | `client/stream.go`, `pkg/admin/stream.go`, `pkg/stream/stream.go`, `pkg/stream/stream_config.go`, `pkg/stream/controller/` | `stream.mdx` |
| Producing | `producer-registration`: instance setup and complete ProducerConfig/BatcherConfig; `produce`: Produce, ProduceOptions, ProduceResult, idempotency; `produce-batch`: ProduceBatch and NewProduceItem; `produce-func`: callback-produced payloads | `client/producer*.go`, `pkg/producer/`, `pkg/produce/`, `pkg/produce/batcher/` | `producer.mdx` |
| Transactions | `transactions`: InTransaction, TransactionFunc, Tx and Querier; `produce-in-transaction`: ProduceInTx and ProduceFuncInTx, shared transaction ownership | `client/client.go`, `client/producer_instance.go`, `pkg/datastore/`, `pkg/producer/` | `client.mdx`, `producer.mdx`, `guides/transactional-produce.mdx` |
| Payloads and stored messages | `message-payload`: Versioned and RawPayload; `stored-message`: returned message fields. Preserve `message-metadata` and the reviewed MessageOptions/RetryPolicy/ConcurrencyPolicy tables in `consumer-registration` | `pkg/common/`, `client/alias.go`, `pkg/consume/meta.go` | `message-options.mdx`, `concepts/api-shape.mdx` |
| Consumer outcomes and bindings | Preserve `consume#message-results` as the outcome contract, including Terminal and Delay. `consumer-bindings`: Binding, Get and System.Bindings; `consumer-workers`: Workers and Worker fields. Add CLI coverage to the existing `consumer-cli` page | `client/binding.go`, `client/consumer.go`, `pkg/consume/`, `pkg/worker/worker.go` | `consumer.mdx`, `concepts/handler-outcomes.mdx`, `concepts/routing.mdx` |
| Message keys | `key-messages`: Messages and limits; `compaction-head`: single-key and stream-wide reads, LockCompactionHead and transaction requirement | `client/key.go`, `client/stream.go`, `pkg/compaction/`, `pkg/admin/` | `key.mdx` |
| Connection and lifecycle | `pool`: NewPostgresPool and PostgresConnectionConfig; `client`: NewClient and ClientConfig; `lifecycle-context`: LifecycleContext and signal behavior | `client/pool.go`, `client/postgres_connection_config.go`, `client/client.go`, `client/client_config.go`, `pkg/common/` | `pool.mdx`, `client.mdx` |
| System and migrations | `system-object`; `system-registration`; `get-system`; `destroy-system`; `schema-migrations`: system/stream Migrate, MigrationVersion and MigrateStreams; `stream-version-health`: Health and StreamVersionHealth; `supported-postgresql-versions` | `client/system.go`, `client/stream.go`, `pkg/admin/`, `pkg/system/`, `pkg/migrate/`, migration registries and compatibility declarations | `system.mdx`, `stream.mdx`, `supported-postgresql-versions.mdx`, `guides/migrations.mdx` |
| Maintenance | `manager`: Manager.Run; `maintenance`: Suspend, Unsuspend and Status for Janitor/Vacuum. Configuration stays with registration | `client/manager.go`, `client/maintenance.go`, `pkg/systemmanager/`, `pkg/worker/`, `pkg/stream/` | `manager.mdx`, `maintenance.mdx` |
| Schedules | `schedule-object`; `schedule-registration`: Register, SchedulerConfig and instance.Schedule; `get-schedule`; `list-schedules`; `suspend-schedule`: Suspend/Unsuspend; `run-schedule`: Run and options; `schedule-messages`: Status/Messages and returned objects; `destroy-schedule` | `client/scheduler*.go`, `pkg/scheduler/`, `pkg/schedule/` | `scheduler.mdx` |
| Metrics | `read-metrics`: Definitions, Latest, History, Measurement and MetricDefinition; `consumer-metrics`, `stream-metrics`, `system-metrics`: snapshots and series contracts, linking existing cursor backlog; `publish-metrics`: NewMeasurement and metric producer; `metrics-exporter`: otel constructors, config and serving lifecycle | `client/*metric*.go`, `pkg/metric/`, `otel/` | `metrics.mdx`, code declarations |
| Alerts | `read-alerts`: Definitions, Latest, History, Snapshot and returned objects; `partition-count-alert`; `compaction-read-cost-alert`; `metric-collector-progress-alert`. Preserve `worker-liveness-alert` | `client/*alert*.go`, `pkg/alert/` | `alerts.mdx`, code declarations |
| Logs and diagnostics | `diagnostics`: code lookup, recovery classifications, explain command and logging config; `producer-logs`, `consumer-logs`, `worker-logs`, `system-logs`: exact events/fields, linking dedicated troubleshooting and existing consumer-stopped entry | `pkg/*/events.go`, `pkg/*/errors.go`, `pkg/common/diagnostic/`, `pkg/common/logging/` | `diagnostics.mdx`, hidden `errors/SQL*.md` |
| CLI | `cli`: installation, connection flags and global behavior; `stream-cli`; `system-cli`; `schedule-cli`; `manager-cli`; `metric-cli`; `alert-cli`. Preserve `consumer-cli`. Migration commands belong in schema-migrations and explain in diagnostics | `cmd/sqlstreams/internal/cli/` | operation pages and examples |

### Troubleshooting coverage map

Keep the approved SQL0105 page. Create these procedures from the actual
emission paths and diagnostic queries. Simple validation/absence contracts
stay in Reference. Do not invent a redrive API or database timing history.

| Page target under `troubleshooting/` | Problem and supporting sources |
| --- | --- |
| `messages-not-claimed` | Recurring SQL0108. `pkg/consume/messageconsumer/`, `pkg/consume/events.go`, consumer registration/binding/worker reads. |
| `messages-dead-lettered` | SQL0028–SQL0031. Inspect recorded failures and distinguish terminal results from retry exhaustion. `pkg/consume/events.go`, exception processing and message options. |
| `consumer-handler-exceeds-timeout` | SQL0039 and consumer handler timeout/cancellation errors. Inspect duration and cancellation behavior using actual logs and configured timeouts. `pkg/consume/base/`, ConsumeOptions. |
| `lease-repeatedly-reclaimed` | SQL0003/SQL0026/SQL0027. Diagnose repeated expiry or range quarantine. Lease and exception queries plus instance lifecycle. |
| `commit-confirmation-lost` | SQL0019. Inspect consumer outcome/lease state without treating absent rows as proof of success. `pkg/common/errors.go`, recording path. |
| `produce-is-slow` | SQL0038 and partition-boundary SQL0033/SQL0057/SQL0056. Timing evidence, batch settings and partition state. `pkg/produce/`, `pkg/producer/`. |
| `schema-version-mismatch` | SQL0022/SQL0023. Determine whether the database or binary must advance. `pkg/migrate/errors.go`, registries, schema migration API/CLI. |
| `migration-lock-timeout` | SQL0053. Identify the blocking session and select a justified action. Migration execution and declared diagnostic. |
| `worker-not-running` | Missing upkeep or collector progress. Worker status, suspension, ownership and actual failures. `pkg/worker/`, `pkg/systemmanager/`, maintenance/alert contracts. |

### Dependencies and legacy handling

| Approved learning page | Supporting coverage |
| --- | --- |
| Streams and messages | Stream object/registration, stored message, stream metrics and retention settings |
| Background managers and workers | Manager execution, consumer lifecycle, worker claims, stream janitor, maintenance controls, and worker-not-running procedure |
| Retries and failed messages | Existing Consume outcomes/options, consumer metrics/logs, dead-letter and timeout procedures |
| Ordering and concurrency | Existing concurrency options/metadata, key reads, binding reads, reclaim and timeout procedures |
| Produce messages in a database transaction | Produce, transactions, both produce-in-transaction variants and their actual failure contracts |
| Make message processing safe to repeat | Metadata, Consume outcomes, transaction contracts, reclaim and commit-confirmation procedures |
| Upgrade SQLStreams | Supported versions, system/stream migrations, version health, mismatch and lock procedures |

Reuse existing legacy paths when their final purpose is unchanged. For split
pages, keep the legacy source hidden during drafting and record its eventual
redirect target when the replacements are complete. Route new pages through
the existing board mapping and group frontmatter. Do not link unfinished pages.
Keep each diagnostic's contract in one canonical location and link procedures
rather than copying their steps into the code catalogue. Final Related topics
and URL review includes original sticky links and the hidden code-page routes.

### Progress

- Reference and Troubleshooting complete for this phase: 71 Reference pages
  and 10 Troubleshooting pages in navigation. The coverage maps above record
  each page's scope and source ownership. Existing reviewed Consumer pages
  received a consistency and correctness pass.
- Validation: 54 extracted Go examples compile, including the separate otel
  module. 51 struct tables match source fields and declaration order.
  CLI examples were checked against a locally built current binary.
- All 12 diagnostic SQL blocks were validated against a uniquely named local
  fixture schema. Administrative statements were planned without executing
  them. Read queries also passed with exception states and exclusive lease
  boundaries. The fixture was removed.
- Targeted Remark, Vale, navigation ESLint, and site build passed. The site
  indexes 87 pages. Local links and anchors were checked on the built pages.
- Accuracy corrections include the `sqlstreams scheduler` command prefix,
  success/deferred counter timing, producer-owned partition creation, and
  SQL0031's crash-loop backstop meaning.
- Editorial sweep reviewed all 81 Reference and Troubleshooting pages. Trimmed
  repetition and unrelated commentary in 25 pages and preserved moved JSON
  field definitions in two canonical tables. The other 54 pages were unchanged.
  Examples were unchanged. Targeted prose checks, build, and links passed.
  This was an agent review. The user has not reviewed the new pages.
- Consistency sweep reviewed all 81 pages. Aligned operation and resource names,
  producer terminology, return formatting, duration defaults, and JSON columns.
  Corrected exception-row terminology and the reclaim-limit boundary. Checked
  55 struct tables against source, including 25 JSON mappings. Prose checks,
  site build, table structure, and 3,445 local links and anchors passed.
  Runnable examples and review markers remain unchanged.
- Accuracy sweep reviewed all 81 pages against implementation. Corrected seven
  claims in five Reference pages: message retry counts, reclaim-limit boundary,
  committed cursor and backlog interpretation, alert check counters, and paired
  manager failure handling. Troubleshooting pages and code examples unchanged.
  Targeted prose checks, site build, source struct tables, and 3,445 local links
  passed. The optional live retry probe failed database authentication before
  creating its fixture. Retry semantics were verified from source instead.
- Applied the spot-check lessons across all 81 Reference and Troubleshooting
  pages. Updated 67 pages with identifier formatting, clearer method purposes
  and timing, and fewer peripheral caveats. Preserved executable examples,
  diagnostic queries, and review markers. Prose checks, build, 55 source struct
  tables, and 3,451 local links passed.
- Reference and Troubleshooting are complete for now following the user's
  spot checks. Further user reviews will happen later.
- Background managers and workers replaces Automatic producer batching and
  is drafted under Concepts → System. Source verification, prose checks,
  navigation lint, site build, and its 47 local links/anchors passed. The
  process diagram was checked at desktop and phone widths in both themes.
  User review is complete. The page's reviewed edits informed CONCEPTS.md
  and VOICE.md.
- Streams and messages is drafted under Concepts → Stream. Checked the
  payload and retention claims against source. Prose checks, navigation lint,
  build, JSON validation, and 42 local links/anchors passed. Diagram labels
  and page overflow checked at desktop and phone widths in both themes.
  User review is complete. Group-progress explanation stays on Consumer
  groups, retention on Streams and messages, and the janitor section links
  there. CONCEPTS.md records the review lessons.
- Retries and failed messages is drafted under Concepts → Consumer. The
  initial draft was rejected for context gaps and unnecessary material.
  Revised around one message's failure, retry, and terminal outcome, with
  separate context, scope, terminology, and accuracy reviews. Removed the
  concurrency diagram and requested-delay scenario. Prose checks, build,
  and 51 local links/anchors passed. The preview indexes 90 pages.
  User review is complete. CONCEPTS.md and VOICE.md record its lessons.
- Ordering and concurrency replaces the hidden legacy page at its existing
  route. Explains three policies with order-status updates sharing a message
  key, then follows a failure through ordered processing. Source verification,
  prose checks, navigation lint, build, and 53 local links/anchors passed.
  The preview indexes 91 pages. User review is complete, finishing all six
  approved Concepts.
- Produce messages in a database transaction continues the Quickstart producer
  with a user insert and welcome-email request in ProduceFunc. Ran the edited
  producer against a disposable database, confirmed successful commit, repeat
  deduplication, and callback-error rollback, then removed the database.
  Prose checks, navigation lint, build, and 31 local links/anchors passed.
  Checked desktop and phone rendering in both themes. The preview indexes
  92 pages. User review is complete. GUIDES.md and VOICE.md record the lessons
  about scenario setup, explicit outcomes, and focused notes. Restored the
  required context import instruction and clarified the external-call note.
- Make message processing safe to repeat uses a signup reward id to guard a
  points update in the same transaction. Ran both Go programs in a disposable
  database. Two messages left one reward record and 10 points. A rejected
  points update rolled back the reward record, and the automatic retry after
  recovery credited it once. The database was removed. Prose checks,
  navigation lint, build, and 53 local links/anchors passed. Checked desktop
  and phone rendering in both themes. The preview indexes 93 pages.
  User review is complete. The introduction distinguishes producer deduplication
  from repeated processing. The outcome breakdown and supporting notes were refined.
- Upgrade SQLStreams covers preparation, version checks, required migrations,
  and deployment. Preserved the compatibility table verbatim. Tested the
  published v0.1.4 → v0.1.5 modules and matching CLI against a disposable
  database. The new consumer processed old and new producer messages, and
  schema status and no-op migration commands passed. This pair has no upgrade
  DDL. Explicit database URL flags work with the published and current CLI.
  Removed the database. Prose checks, navigation lint, build, 52 local
  links/anchors, and desktop/phone checks in both themes passed. The preview
  indexes 94 pages. User review is complete, finishing all five Guides.
- Final cross-page pass complete across all 94 visible articles. Updated
  Related topics on 41 pages, connected the original stickies to canonical
  pages, and applied the established consumer handler term to older learning
  pages. Corrected the CLI example to use Quickstart's database credentials.
  All 92 board articles now have Related topics sections.
- Added 27 legacy article redirects and 109 diagnostic redirects through
  Astro's existing configuration. Benchmark, Roadmap, and Table Design have no
  approved replacements or routes.
- Validation: targeted remark, Vale, and config lint passed. Site build passed.
  All 4,644 generated local links and 138 redirects, including their target
  anchors, passed. Pagefind indexes exactly the 94 approved articles and the
  sitemap excludes redirects. One browser smoke check verified article and
  diagnostic redirects, Quickstart → Consumer groups → repeat-processing
  navigation, SQL0105 search, and mobile width, with no page errors.
  Changes remain uncommitted. No deployment performed.

- Retired-source cleanup removed 138 unreferenced documents (29 articles and
  109 diagnostic pages), plus the obsolete diagnostic renderer and its helpers.
  Redirects remain canonical. Removed frontmatter and count handling specific
  to per-code pages. Replaced per-page diagnostic tests with redirect coverage
  and updated the obsolete browser 404 assertions. Targeted lint and formatting,
  Astro type checking, seven unit tests, and build passed. All 94 articles,
  138 redirects, and 4,644 generated local links remain intact. The updated
  browser test for redirects and sitemap exclusions passed.

- Final consistency review covered all 94 remaining articles. Standardized
  inline API and field formatting, removed instance/worker terminology drift,
  and aligned stream alert examples with the shared payment payload. Moved the
  stream-owned worker-liveness alert into the Stream navigation group. Existing
  guide and troubleshooting procedures, resource names, diagnostic literals,
  and review markers are preserved. Prose checks passed across all 94 pages.
  Build and validation of all 4,626 generated local links and 138 redirects
  passed. Search still indexes exactly 94 articles.

### Per-page workflow


Current task: **Documentation expansion complete, awaiting the user's commit**.
All approved Concepts and Guides have completed user review. Reference and
Troubleshooting, cross-page consistency, Related topics, and URL migration are
complete for this phase. Further user spot checks may happen later. Consolidate
this task's decision record and HISTORY entry only after the user commits it.

Use one focused task per page, preferably a fresh session. Record the current
page's purpose, boundaries, source files, dependencies, status, and unresolved
review feedback here when starting it. Keep this a compact working checklist.

1. Read the page's writing document, VOICE.md, conventions, and accepted
   exemplar. Inspect the relevant shipped implementation and dependencies.
2. Draft only that page, including useful visuals and links where needed.
3. Make separate accuracy and usefulness passes. Verify terms, resource names,
   examples, contracts, and whether every section serves the page's purpose.
4. Validate the documented code or SQL. Check rendering when layout changes.
   Keep checks targeted and avoid repeated browser runs for small prose edits.
5. Complete the page's dedicated review and validation before starting another.
   Reference and Troubleshooting do not require user approval or review pauses.
   Preserve user deletions and review the whole procedure after changes.
   Concepts and Guides retain their approved scope and individual user review.

Concepts, Guides, Reference, and Troubleshooting follow their respective
writing documents in .website/. No bulk drafting or per-page decision records.
At close-out, consolidate the work's decision and history entry after the user
commits. Agents do not commit or deploy without the required authorization.
