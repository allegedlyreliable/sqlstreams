# CLI help review

Reviewed 2026-09-17. The findings below describe the pre-change behavior.
Approved corrections are now implemented in the working tree, including the
metric-filter fix. Build and CLI race tests pass, all 71 help pages render,
and both new filter regressions failed before the fix and pass afterward.
In-flight validation and close-out notes are in .work/TODO.md.

Built the working-tree CLI and rendered all 71 discovered help pages, including
the generated completion commands. Every invocation exited successfully.
Compared command descriptions, examples, arguments, flags, and defaults with
their handlers and the relevant client, controller, and datastore behavior.
Database operations and shell installation instructions were not executed.

## Accuracy findings

1. **P2: RANK 0 does not mean compaction was disabled.**
   `cmd/sqlstreams/internal/cli/stream_key_messages.go:19` makes that claim.
   `ProduceOptions.Compaction` accepts `Enable: true, Rank: 0`; the write keeps
   that zero, while the read also converts SQL NULL to zero
   (`pkg/produce/controller/adapter.go:19`,
   `pkg/compaction/controller/datastore/message.go:68`). A compacted `order-42`
   message using the default rank therefore looks uncompacted by this help's
   definition. Say that zero can represent either case; distinguishing them
   in output would be separate API work.

2. **P2: The key help points to a nonexistent API field.**
   `cmd/sqlstreams/internal/cli/stream_key.go:11` says
   `MessageOptions.MessageKey`; the field is `ProduceOptions.MessageKey`
   (`client/produce_options.go:27`). Correct the owning type.

3. **P2: “Every worker” overstates the system manager's scope.**
   `cmd/sqlstreams/internal/cli/manager_run.go:24` promises to keep every worker
   running without a consumer. Its provisioner list in
   `pkg/systemmanager/systemmanager.go:106` covers maintenance and built-in
   alert consumers, not application message/retry consumers; suspended workers
   also stay suspended. Say it runs deployment maintenance and built-in alert
   processing, and that application consumers run separately. The singleton
   manager claim and replica takeover description match the implementation.

4. **P2: Consumer worker config does not all refresh at ConfigRefreshInterval.**
   `cmd/sqlstreams/internal/cli/consumer_worker.go:17` applies that promise to
   every listed worker. Message and exception consumers refresh their group
   config, but the cursor advancer captures metadata and poll rate when its
   instance is constructed (`pkg/consume/cursoradvancer/instance.go:35`). Limit
   the refresh statement to the workers that implement it; do not imply an
   update deadline for the entire listing.

5. **P2: scheduler messages lists outcomes, not a general message listing.**
   `cmd/sqlstreams/internal/cli/schedule_messages.go:14` omits the consumer-group
   dimension. `pkg/schedule/controller/datastore/messages.go:23` reads the
   newest N messages, then emits one row per matching consumer group. With two
   groups, `--limit 5` can produce ten rows; with no matching groups, it returns
   none even if messages exist. Describe “per-consumer-group outcomes for the
   newest retained messages” and explain that the limit counts messages.

6. **P2: The metric attribute filter does not fulfill “all must match.”**
   `cmd/sqlstreams/internal/cli/metric_read.go:144` promises conjunctive filters,
   but `parseAttributePairs` overwrites duplicate keys and `attributesMatch`
   treats an absent attribute as an empty value. For example,
   `--attribute stream=orders --attribute stream=payments` retains only
   `payments`; `--attribute stream=` can match a series without `stream`.
   This is a behavior/help mismatch, not just prose. Recommend preserving the
   documented contract by rejecting conflicting duplicates and checking key
   presence; review that behavioral fix separately from the text changes.

7. **P2: --quiet promises silence and an absence-only exit code it cannot give.**
   The get commands all say “no output; ... (0 exists, 1 not)”
   (`cmd/sqlstreams/internal/cli/get.go:64`, `consumer_get.go:52`,
   `schedule_get.go:59`, `system_get.go:60`). Operational failures still print
   diagnostics and can exit 1; usage errors exit 2. Verified with an invalid
   URL, without contacting a database. Say “suppress result output; exit 0 if
   found, 1 if absent or the operation fails”; document usage errors separately.

8. **P2: Output-mode restrictions are absent from help.**
   `cmd/sqlstreams/internal/cli/root.go:96` advertises JSON globally, but
   `manager run` rejects it, completion commands emit shell scripts, and help
   remains text. All quiet flags reject `--output json`; all four destroy
   commands require `--yes` with JSON. The relevant flag descriptions never
   disclose those combinations. Document these exceptions where users choose
   the flags, including that manager logs go to stderr.

9. **P3: system get promises config but shows registration timestamps.**
   `cmd/sqlstreams/internal/cli/system_get.go:17` says “Show the singleton system
   config”; `printSystemDetail` prints only CreatedAt and UpdatedAt. Use
   “Show the system's registration” rather than implying it exposes the
   declaration applied by system register.

10. **P3: Binding list and explain descriptions omit their actual output shape.**
    `system_binding_list.go:17` describes each consumer's declared set but also
    returns waiting declarations alongside installed ones
    (`pkg/consume/controller/binding.go:61`). Mention both states.
    `explain.go:17` promises “problem, recovery, fix” across all diagnostic
    kinds, although events, metrics, and alerts have different fields
    (`cmd/sqlstreams/internal/cli/errors.go:188`). Describe the shared purpose
    and say details depend on the diagnostic kind.

## Consistency and completeness

| Location | Issue | Suggested alignment |
| --- | --- | --- |
| `alert_read.go:23`, `metric_read.go:23` | Generated phrases read “history retained evaluations” and “history measurements”; latest is unnecessarily plural for one alert. | Give latest and history distinct descriptions: “Show the latest retained alert” / “List retained alert history, newest first”; “Show the latest measurement per series” / “List retained measurement history per series.” |
| `stream_maintenance.go:24` | Multi-sentence Short text is much longer than sibling summaries. “Next heartbeat” omits the public contract's “successful.” | Keep one brief Short; move stop timing, successful heartbeat, and non-waiting behavior into Long. |
| `schedule_run.go:20`, `schedule_suspend.go:17` | Mixes schedule, row, request, and message; “schedule producer producing a schedule” is awkward. | Use schedule for the resource and message for what it produces. Say suspension stops automatic production; manual run still works. |
| `root.go:92` | Database flag mentions only postgres:// although postgresql:// and keyword/value DSNs work. | Describe a PostgreSQL connection string; retain environment fallback guidance. |
| `alert_read.go:24`, `metric_read.go:23` | Missing retained data exits 1, including empty history; alert help mentions only the missing-owner nonzero case. | State the no-retained-data exit behavior consistently. |
| `stream_config.go:19`, `consumer_binding.go:11`, `system_destroy.go:29` | API references alternate between partial client paths and the internal RegisterSystem name. | Use public client paths when an API reference helps, otherwise describe the CLI action. |
| `migrate_versions.go:14`, `manager_run.go:116`, `explain.go:16` | “THIS”, “EX:”, “N-way”, “log-event”, and “committed advance” depart from the simpler neighboring descriptions. | Use “this binary”, “e.g.”, “multiple replicas”, “log event”, and “cursor advancement.” |

The existing root description is appropriate. Ordinary stream/consumer/schedule
list and get descriptions, migration targets and directions, key history order,
schedule manual-run defaults, schedule destruction, and unsuspend timing agree
with their main execution paths. No command rename or new CLI feature is needed
to fix the help findings.

Implementation applies the factual and wording corrections together. The
metric-filter behavior has its own regression tests. No help snapshot suite
was added; the CLI was rebuilt and its help pages rendered directly.
