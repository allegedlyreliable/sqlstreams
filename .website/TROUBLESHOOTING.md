# Writing troubleshooting pages

Use this when creating or revising a Troubleshooting page. Help readers
identify their problem, inspect the evidence, and follow a short procedure
to fix it. Keep the page simple, direct, and procedure driven.
[VOICE.md](VOICE.md) governs general prose and terminology.
[CONVENTIONS.md](CONVENTIONS.md) governs site structure and code.

The reviewed [Queued message cannot start (SQL0105)](src/content/docs/troubleshooting/queued-message-cannot-start.mdx)
is the local model. Earlier drafts failed because they made readers assemble
a procedure from explanations, scattered outcomes, and unrelated suggestions.

## Start with the symptom and get to the action

Name the recognizable symptom and diagnostic code in the title. Show the
exact warning or error text near the top. Follow with one or two sentences
explaining what happened and its immediate effect.

For a log warning or error, show the full line with its level, exact message,
code when emitted, and fields. Use named placeholders for variable values:

```text
level=WARN msg="produce exceeded the duration threshold" code=SQL0038 stream={stream} duration={duration} threshold={threshold}
```

Use this format everywhere a diagnostic log is presented. A message-only
excerpt leaves later references to attributes unexplained. Merely saying
“the warning's duration field” does not fix that disconnect. Show the field
in the log before interpreting it, and use placeholders so no fictional
scenario or explanation of sample values is needed. Verify the emitted fields
and level against source. Returned errors and recorded failure text keep their
actual format rather than acquiring invented log attributes.

Give only the context required to understand the next action. For SQL0105,
readers need to know that a message waited in the instance's local queue
and that its consumer handler is skipped for this attempt. They do not need
a lesson on the entire lease lifecycle before running the query.

Cut special-case instructions that add no decision or action. If the procedure
already tells readers to use the warning's worker name, it does not also need
to name the collector worker as a separate example.

Do not repeat the title with “Use this page when your logs show SQL0105.”
Do not open with background, a catalogue of causes, or a worked fictional
scenario. Link the relevant Concepts page when deeper understanding helps.

## Follow evidence → result → action

Use this reading order, adapting the headings to the problem:

1. **Inspect the affected resource.** Supply the actual diagnostic SQL or
   command. Put a **Result / Action** table immediately below it.
2. **Fix the continuing problem.** Give numbered steps in execution order.
   Put the condition for taking a branch at the start of that step.
3. **Related topics.** Link useful mechanisms and exact settings without
   repeating their documentation.

Each step identifies what to do, where to do it, and what finding determines
the next action. Keep code inside the step that uses it. State when a change
takes effect if the procedure depends on it, including a required restart.
Readers should be able to scan the headings, result column, and step openings
to find their next action without reading every paragraph.

Distinguish an affected resource's current state from a recurring problem.
A message already in exception processing needs a different action from
newly claimed messages repeatedly running out of lease time.

Keep all relevant query outcomes in the table, including no rows and null
values. Do not leave two outcomes in a paragraph beneath a four-row table.
Use precise conditions so readers can identify the applicable row.

Do not append a routine “Verify recovery” section. Readers do not need a
reminder to check whether their problem stopped. Include a specific check
within a step only when its result determines what to do next or prevents
a misleading conclusion. Author validation does not become reader homework.

## Use diagnostics that actually exist

Check diagnostic messages, fields, queries, and configuration behavior against
shipped code. Name the exact log code and fields readers should inspect.
If logging must be enabled first, show how before asking readers to use it.

Do not assume application traces, dashboards, metrics, or timing records
exist. “Check consumer handler durations in application traces” was an
unsupported instruction. The approved procedure enables `SlowDispatchThreshold`
and reads the actual `SQL0039` fields instead.

Provide SQL rather than telling readers to “locate the affected message.”
Look for an existing query and placeholder convention in the repository
before inventing a new presentation. Follow diagnostic placeholders: `{schema}`, `{stream_id}`,
`{message_id}`, and `'{group}'`. Preserve SQL quoting where required.
Do not fill diagnostic queries with fictional resource names or ids and
then add a paragraph explaining how to replace them. Explain substitutions
only when their source or meaning is unclear.

Respect what the evidence can establish. A current database row does not
reconstruct historical queue waiting. An absent lease does not prove that
processing succeeded. A warning about one attempt does not establish what
happened on every attempt. State a limitation where it changes the next
action, without accumulating defensive caveats.

Accuracy alone does not justify a caveat or reminder. Keep it when it changes
how readers interpret evidence or choose their next action. Remove details
that do neither, even when they are technically correct.

When interpreting duration, identify what is timed, what waits and where,
and when measurement starts and ends. Distinguish queue time, database work,
and application callbacks when that distinction affects the investigation.
Do not make readers infer those meanings from a list of timing terms.

## Make corrections conditional and ordered

Choose a path through investigation and correction. “Try a smaller queue”
does not tell readers why or when to do it. First establish whether consumer
handler calls are unexpectedly slow. Then explain when reducing claimed
work is appropriate, and what to change if the warning continues.

Use conditional step openings such as “If processing times are acceptable…”
and “If SQL0105 continues on new claims…”. Numbering a bag of independent
suggestions does not turn it into a procedure.

Identify the configuration object and exact fields to edit. Keep a material
cost beside its change, such as more database claims or slower reclaim after
a crash. Do not add alternative tuning strategies merely because they exist.

## Remove the patterns that failed review

- **Essays before action.** Cut context that does not help identify the
  problem, interpret a result, or execute the next step.
- **Vague instructions.** Replace “investigate waiting” or “inspect logs”
  with the relevant resource, diagnostic, fields, and action. Establish what
  is waiting, where, and for what before referring to “the waiting.”
- **Random remedies.** Do not offer queue size, concurrency, timeouts, and
  service performance as equal choices without evidence or an order.
- **Repeated facts.** State each result, constraint, and cost where it is
  needed. Do not repeat it in the intro, procedure, note, and closing check.
- **Competing visual hierarchy.** Use headings for phases, a table for query
  outcomes, and numbered steps for execution. Do not pile headings, bold
  paragraphs, lists, tables, and notes onto the same decision.
- **Notes that preserve clutter.** ThreadAside can hold a brief, common
  qualification. It cannot rescue an unnecessary paragraph. Required actions
  stay in the procedure.
- **Exhaustive edge cases.** Include branches relevant to the diagnostic
  results. Move unrelated failure modes and detailed contracts elsewhere.
- **Invented certainty.** Do not fabricate output, imply unsupported
  observability, or treat missing evidence as proof of success.

Shortening means deleting unnecessary topics and repetition. Do not compress
the same essay into dense sentences or remove the conditions that make a
step actionable.

## Revise the whole procedure after feedback

Read the whole page before fixing a criticized section. Identify the question
that section must answer, what readers already learned above it, and where
its result leads. Do not add context locally when it already exists elsewhere.
Move or remove content to repair the sequence.

A formatting change cannot fix an undefined purpose. Earlier revisions
turned the same unordered advice into paragraphs, bullets, and tables without
answering “Why should I change this, and when?” Settle that decision before
choosing its presentation.

Preserve the user's deletions across revisions. Do not restore removed
substitution instructions, background, or recovery checks under new wording.
The accepted page's omissions are part of the pattern, not gaps to fill from
a generic troubleshooting template.

## Review before calling it ready

- Can readers recognize the symptom and reach the first action immediately?
- Does every diagnostic instruction use a real query, command, or log field?
- Are all relevant query outcomes and their actions together in the table?
- Can readers follow the numbered steps without choosing among unexplained
  remedies or reconstructing missing setup?
- Does each branch say when to take it, with effects and costs beside the edit?
- Are terms and resource names consistent, including “consumer handler” in prose?
- Can any paragraph, note, or routine verification step be deleted without
  making diagnosis or correction harder?

Apply the revision checklist in [VOICE.md](VOICE.md). Validate the documented
queries and behavior as the author, keeping checks proportional to the change.
