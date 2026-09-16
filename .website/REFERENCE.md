# Writing reference pages

Use this when creating or revising a Reference page. Give readers an exact
contract they can find, scan, and use while writing code. It should read like
a useful manual, with consistent formatting and code first.
[VOICE.md](VOICE.md) governs general prose and terminology.
[CONVENTIONS.md](CONVENTIONS.md) governs site structure and code.

## Organize around recognizable resources and operations

Follow Stripe's resource-and-operation model. Readers find the resource
they recognize, then the operation they need. Use public names rather than
categories that require knowledge of implementation details.

A resource object entry defines returned fields and links its operations.
An operation entry defines the call and its contract. Give operations direct
links, without imposing one page per handle, type, or symbol. Keep closely
related variants together when comparison helps readers.

Use a descriptive operation title such as “Retrieve a consumer group” and
show the exact API name, `ConsumerHandle[T].Get`, with its signature.
Match the title to the content readers will find, including non-operation
entries. Use the same resource and operation names in navigation, headings,
and links. Format API identifiers such as `Versioned` as code in prose so
readers can connect them to the declarations and examples.

Settings belong beside the operation they configure. `ConsumerConfig` lives
with Register, and `ConsumeOptions` lives with Consume. Link shared definitions
instead of maintaining a second field table on another page.

Reference owns exact behavior, including uncommon cases that change a
caller's implementation. The Guides rule limiting gotchas to common scenarios
does not justify omitting a reference contract. Link to Concepts for causal
explanations, Guides for procedures, and Troubleshooting for investigation
and recovery.

Accuracy alone does not justify including a detail. Keep information that
helps readers understand or use this page's contract. Remove obvious reminders,
incidental implementation details, and peripheral caveats. Keep exceptions
that change what callers must do beside the operation they affect.

## Use a predictable operation layout

1. **Signature.** Start with a code block showing the exact signature and
   identifying its receiver or namespace. The page title supplies the H1.
2. **Introduction.** State what the operation does. Add enough context to
   make its effects clear, without a page overview or motivation section.
3. **Example.** Show a compact, real call with necessary error handling.
4. **Parameters.** Use a Parameter / Type / Required / Description table.
   Put configuration fields underneath with Field / Type / Default /
   Description columns.
5. **Returns.** Use the outcome table described below.
6. **Behavior, when needed.** Include substantial calling semantics such as
   concurrency, repeated calls, persistence, or shutdown. Brief facts belong
   in the introduction or beside the parameter or result they qualify.
7. **CLI, when applicable.** Show a concrete equivalent command, then explain
   its output and exit behavior. Link the full command reference for flags.
8. **Errors, when applicable.** Use the error table described below.

Omit sections that do not apply. Do not invent a CLI equivalent or keep a
Behavior section just to hold one short statement. Object entries use their
type, field table, and operation links instead of this operation outline.

A block listing several methods needs brief explanations of their distinct
purposes. For stream metrics, distinguish a current snapshot, metric definitions,
and selectors for retained measurements. Do not substitute obvious access
instructions such as “Call `.Metrics()` to access these methods.”

## Make the example do useful work

Use the shortest example that demonstrates the call and its result correctly.
Name necessary starting state and substitutions. A reference example can
assume an initialized client, but it must identify that assumption and link
any shared payload definition needed to understand the code.
Link referenced setup instructions directly, such as Quickstart when using
its database, instead of making readers search for them.

Keep names consistent across related entries. The reviewed examples use the
`payments.requested` stream, the `charge-cards` group, and the
`PaymentRequestedV1` payload. Name each resource's kind when introducing it.

Keep generic type arguments tied to the resource they describe. A
`Stream[PaymentRequestedV1]` handle describes payment messages even when the
operation returns an `Alert`. Use `RawPayload` when raw JSON is relevant to
the example, rather than as an unexplained substitute for its payload type.

Do not narrate code that is already clear. “The stream and group names come
from the handle” adds nothing when the example shows both names. Put empty-name
validation in Errors. Preserve prerequisites that affect successful calls,
such as requiring an existing stream or a compatible payload type.

An example that prints a payment's id does not charge a card. Identify
simulated work and keep error handling real. Do not repeat “check the error
first” after an example and return table that already make this clear.

## Describe parameters through their effect

Show a configuration struct's complete field list in one table, in source
declaration order. Follow it with separate tables for nested structs, in
the order their fields introduce them. Link each nested type to its table
and document a shared type once. Do not split the parent struct into themed
field tables. Standard-library types do not need their own field catalogue.

Check completeness and ordering against the actual struct, including nested
types. A shared struct can have different defaults in different fields.
State which use a nested table describes: `MessageOptions` supplies defaults
in `Message`, lower bounds in `MessageMin`, and upper bounds in `MessageMax`.
Do not present one use's defaults as universal or repeat its entire table
for each field.

Keep field names and types unbroken and readable. Let a wide table scroll
within its container instead of wrapping identifiers into narrow fragments.
Check explanatory columns too: long identifiers must not squeeze Description
or Value into narrow strips. Give prose columns sufficient minimum width and
let the table scroll. List independent prerequisites on separate lines or
as bullets.

State types, required values, defaults, zero and nil meanings, constraints,
scope, precedence, and effective timing where applicable. A derived default
names its dependencies and when it resolves. Say whether the call copies,
mutates, or retains a supplied configuration.

Separate first-registration behavior from later updates. Say what an update
preserves, when running instances observe it, and whether changing the Go
struct after the call has any effect. Keep that timing in one place and link
to it from the relevant settings.

When defaults and bounds interact, use a small input / effective-value table
with named settings and concrete values. Verify the resolution order in the
implementation, including any helpers that fill defaults or apply overrides.
A type's field list alone does not establish that precedence.

Describe context through cancellation and deadlines. “Bounds the read's I/O”
is ambiguous and can sound like a performance setting. For a database read,
write “Cancelling this context or reaching its deadline interrupts the
database read.” For a running consumer, explain that cancellation starts
graceful shutdown. Match the description to the operation's actual behavior.

Short does not mean context-free. “Creates no registration and changes no
group settings or progress” leaves readers to infer the subject and effects.
Explain that Get does not create a missing group, and reading an existing
group does not update its settings or change which messages it processes next.

## Make return values visible

Use a **Condition / Return value** table. Give success, expected absence or
empty results, and failure separate rows. Spell literal tuples, `nil`, zero,
and booleans in code formatting. Identify returned objects with a link to
their field definitions.

For Get:

| Condition | Return value |
| --- | --- |
| Group is registered | `(group, nil)`, where `group` is a `Consumer` object. |
| Stream or group is not registered | `(nil, nil)`. |
| Read fails | `(nil, err)`. |

Keep this contract in Returns. Do not also repeat `(nil, nil)` in the
introduction and a paragraph beneath the table. An empty collection and a
missing parent resource may have different results, so describe each.

Use ThreadAside for a supporting qualification that needs a note. Give it a
specific title and place it beside the result it qualifies. The reviewed
Get page uses “Running instances” to explain that a group can remain
registered with no running instances, and its object does not report activity.
Do not substitute a bold “Note:” paragraph for the shared component.

Keep the actual return contract in the table. Notes add useful context,
not another copy of the table or a home for unnecessary warnings.

Explain distinctions through a concrete cause and consequence. For example,
a saved active alert can remain after new measurements stop arriving.
Explain what the reader should inspect next instead of relying on an abstract
heading such as “Recorded state is not current health.”

## Standardize errors and CLI behavior

Use an **Error / When** table. Put exact declared error names in code
formatting. For errors without a stable public name, use a descriptive
category such as validation error or database error. State the triggering
condition in the second column. Avoid vague entries such as “the operation's
error” when the kind of failure can be named.

Separate materially different conditions. Do not classify expected absence
as an error when the API returns `(nil, nil)`. Do not promise a particular
error identity or wrapping behavior without checking the implementation.
If the function returns no error, omit Errors.

Include handling requirements only when they add to the example and tables.
State partial-write or irreversible effects beside the failure or operation
they qualify. Leave diagnostic procedures to Troubleshooting.

A CLI section begins with a usable command using the same example names:

```sh
sqlstreams consumer get payments.requested charge-cards
```

Describe what it prints and when it exits with each relevant status. The Go
API's absence result does not establish the CLI's exit behavior. Link to
the command's full reference for output formats, flags, and connection options.

## Adapt the manual to the reference kind

Keep formatting consistent within each kind without forcing Go headings onto
every reference page.

| Kind | Contract to document |
| --- | --- |
| Resource object | Type, fields, meanings, representation, and related operations |
| CLI command | Syntax, arguments, flags, defaults, output, and exit status |
| Metric | Name, calculation, units, attributes, scope, and freshness |
| Log | Emission conditions, level, fields, and suppression |
| Alert | Trigger, thresholds, severity, evaluation, and resolution |

Use code blocks for exact syntax and representative output, and tables for
comparable fields or outcomes. Keep important settings and operations directly
linkable. A diagnostic code retains one canonical page.

For duration fields, explain what is timed, what waits and where, and when
measurement starts and ends. Separate distinct components such as the producer's
batch queue, database operations, and an application callback. Name the methods
whose timing differs, including whether a caller-owned transaction's later
commit falls outside the measurement.

A CLI installation section includes the supported package managers as well
as Go, with commands verified against the release configuration. Keep general
CLI documentation focused on installation, connection, and useful command
contracts rather than cataloguing incidental output details.

## Show complete log templates

When presenting a log event, show its full line with level, exact message,
code when emitted, and event fields. Use named placeholders for variable
values, not invented names, durations, or counts:

```text
level=WARN msg="produce exceeded the duration threshold" code=SQL0038 stream={stream} duration={duration} threshold={threshold}
```

Apply this throughout log references and diagnostic examples. A message-only
excerpt or an attributes table does not substitute for the line. Readers must
see the fields that the following prose interprets. Do not add a fictional
scenario just to explain sample values. Check the level, message, and fields
against the emitting code, including conditional fields. Logger-added timestamps
and application context can vary. Keep returned errors and recorded failure text
in their own format rather than inventing a log level or attributes for them.

Keep field meanings and emission conditions beside the template. Avoid repeating
the message and field list in a table when the template already shows both.

## Review before calling it ready

- Can readers find the operation using recognizable names and a direct link?
- Are the signature and a useful example near the top?
- Are examples and contracts checked against shipped code, including absence,
  cancellation, repeated calls, and side effects where applicable?
- Are parameters explicit about defaults, scope, and timing without narrating
  obvious setup or using ambiguous shorthand?
- Does the parent table list every field in source order, with linked nested
  tables and defaults scoped to the operation? Are names and types readable
  at desktop and phone widths without horizontal page overflow?
- Do Returns and Errors use their standard tables, with absence visible once?
- Are useful qualifications in ThreadAside, and duplicate reminders removed?
- Does CLI code use the same names and accurately describe output and exits?
- Are terms and resource names consistent across the entry and its links?
- Can a section or paragraph be removed without losing a useful contract?

Then apply the revision checklist in [VOICE.md](VOICE.md). Check rendered
tables, code, and notes when their layout changes. Keep verification targeted
to the change and the review-depth marker honest.

## Basis for this guide

[Stripe's Customers reference](https://docs.stripe.com/api/customers) supplies
the resource-and-operation model, with consistent examples and parameter and
return definitions. [PostgreSQL's INSERT reference](https://www.postgresql.org/docs/current/sql-insert.html)
provides the manual precedent for exact syntax and complete calling contracts.

The review of [Retrieve a consumer group](src/content/docs/reference/get-consumer.mdx)
established the local format: code first, contextual descriptions, outcome
and error tables, ThreadAside qualifications, a concrete CLI example, and
no repeated advice that the example already demonstrates.

The review of [Register a consumer group](src/content/docs/reference/consumer-registration.mdx)
established complete struct tables in source order, nested-type tables,
scoped defaults, separate prerequisite lines, and readable identifier columns.
