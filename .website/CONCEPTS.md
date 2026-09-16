# Writing concept pages

Use this when creating or revising a Concepts page. Start with the reader's
question, then choose the shape that explains it. This is an editing guide,
not a required set of headings. [VOICE.md](VOICE.md) governs the general
prose style. [CONVENTIONS.md](CONVENTIONS.md) governs site structure, code,
and the diagram implementation.

## Define the job before drafting

Write down what the reader should be able to predict after reading. For
consumer groups: which instances share work, which groups progress
independently, and what happens when an instance stops.

Name the assumed knowledge. Introduce new terms at the point where the
example needs them. A definition that requires three more definitions is
starting too deep in the implementation.

Keep one central question per article. Include a failure or limitation when
it changes the answer. API variants, field catalogues, setup instructions,
and recovery procedures belong on their own pages. Link at the relevant
sentence instead of summarizing those pages at the end.
Describe stored records only to the depth the explanation needs. Readers
may need to know that a failure is recorded without needing an inventory
of the exception record's fields.

Apply the same scope discipline to each section. A fact can be accurate and
still interrupt the explanation. Once a section has explained manager scopes,
it can stop without adding a separate qualification about which application
code a system manager cannot run. Keep that qualification only when the
reader needs it to understand the behavior being explained.

Review scope across related pages too. Give each explanation one primary
home and link directly to the relevant section from other pages. Consumer
groups owns independent group progress. Streams and messages owns message
retention. Background managers and workers explains the janitor's role and
links to retention instead of repeating its mechanics.

Keep the local setup an example needs without repeating the underlying
concept. A retention example needs the group names, their recorded progress,
message ages, and retention period. It does not need another explanation of
how consumer groups work. A link supplies deeper background, not missing
facts needed to follow the current example.

## Build the explanation around a case

- Open with the surrounding situation and the need that makes the concept
  useful, then name the concept. Let this establish the page's intent naturally.
  Avoid an isolated definition or “this page explains” preamble.
- When application code triggers the behavior, establish that interaction
  before the domain example. For retries, start with the consumer handler's
  return value and SQLStreams' response, then introduce the failed payment.
  A generic statement that payment providers can be unavailable skips the
  connection the reader needs first.
- When the introduction asks readers to picture unfamiliar relationships,
  ownership, or simultaneous states, add a small visual alongside it. Show
  the participants and their relationships before introducing timing math or
  failure cases. Do not make readers wait for a later diagram to understand
  the opening. Use prose alone when the situation is already easy to picture.
- Establish normal behavior with a named application, concrete messages,
  and explicit assumptions. Explain what each participant knows or owns.
  Use the smallest case that demonstrates the mechanism. One message can
  explain retries without introducing other messages, concurrency, instance
  changes, or a separate approval-delay scenario.
- Change one condition at a time: another instance, another group, a restart,
  or a crash. Carry the same names and ids through each variation.
- State the mechanism, then the consequence. Distinguish in-memory results,
  durable state, and external effects when that boundary matters.
- Put the limitation beside the behavior it qualifies. End when the question
  is answered. Do not add a recap that repeats the opening.
- When closing with links to related topics, put them under a “Related topics”
  heading. Use a short list with one link and a reason to follow it per item.
  Keep links needed to understand the explanation beside the relevant prose.

For example, the Consumer groups page moves through shared work, independent
progress, persistence across restarts, and repeated processing after a crash. A page
about ordering might instead compare two execution traces. Choose the
sequence from the explanation rather than copying an outline mechanically.

## Keep terminology and names grounded in the example

Apply the terminology and naming requirements in VOICE.md as the explanation
unfolds. Establish what each term means in the worked case before using it
as shorthand. An earlier incidental mention is not enough.

For example, “the range’s lease” asks the reader to reconstruct which range
is meant. “A’s lease on message ids 1–2” connects the mechanism to the current
scenario. Introduce “range” explicitly when understanding ranges is part of
the page’s job, rather than adding a term only to shorten a later sentence.

Keep distinct steps distinct. In the Consumer groups example, processing a
message and recording its completion are separate actions. Calling both
“delivery” hides the boundary that explains why processing can repeat.

Carry the same resource identities through prose, code, diagrams, captions,
and alternative text. The `payments.requested` stream stays a stream,
`charge-cards` stays a consumer group, and A and B stay consumer instances.
Reuse those identities on related pages that continue the same example.

## Explain before compressing

State the inputs or conditions that determine behavior directly. For retention,
explain that cleanup considers message age and group progress, rather than
opening with “age alone is not enough.” Correct an assumption explicitly only
when that correction helps the reader understand the mechanism.

Give a worked example its own paragraph when it follows a general rule.
“For example…” makes the transition visible. In the janitor explanation, the
reader first learns what controls cleanup, then sees why a lagging group's
messages can remain after 30 days, then learns that the janitor attempts
removal once the conditions are met. Keep that sequence connected without
packing the rule, example, and resulting action into one dense paragraph.

A short statement can still ask the reader to supply several missing steps.
When a consequence is surprising, continue the current example until the
reader can see why it follows. A diagram does not replace that explanation.
“A message can become too late to start before the lease expires” needs the
steps connecting waiting, elapsed lease time, and the check before starting
the next consumer handler. Show what A does and why before relying on that shorthand.

Explain the reason for an application decision before naming the API that
expresses it. For a permanent rejection: the consumer handler checks the
order, finds it cancelled, determines that retrying the payment would not
help, and returns `Terminal` to stop retries. Naming `Terminal` and its
result alone leaves out the reasoning. Spend words on that causal chain
while cutting unrelated detail elsewhere.

When explaining a duration or budget built from several parts, give each
part its own line with its purpose and a concrete value. Then show the total,
when it starts, and which messages share it. Avoid a sentence that lists
several unfamiliar allowances followed by another sentence repeating them
with numbers. Distinguish configured allowances from actual time spent.
Introduce configuration names after their purpose is clear.

## Choose the form that does the work

| What the reader needs to understand | Use |
| --- | --- |
| Who owns what, or which components share state | A small labeled diagram with explicit boundaries |
| What is happening at the same moment, including work waiting to start | A snapshot showing the participants, work, and shared constraints |
| What changes over time | Numbered causal steps or a timeline |
| How two cases differ | A compact comparison using the same attributes |
| Why a behavior exists or what it costs | Connected prose beside the example |
| An API expression essential to the mechanism | The smallest real code expression that proves the point |

Use headings that name behavior: “Different groups keep independent
progress” tells the reader more than “Group behavior.” Keep paragraphs on
one idea, but retain the connective words that explain cause and effect.
A concept page should be understandable without running its code.

## Give each visual one job

Write the takeaway caption first. Draw only what supports it. Put the visual
next to the explanation, with its caption underneath. Reuse names, shapes,
and direction across related figures. Introduce detail progressively.
Remove a diagram when it mainly explains a neighboring topic. The original
retry diagram explained concurrency, while a short numbered sequence better
explained one message's failure and retry. A concept page does not need a diagram
when prose, a sequence, or a table already makes the mechanism clear.

Show shared storage once. Draw group boundaries explicitly. Label arrows
with their meaning when direction alone is ambiguous. Mark snapshots and
possible assignments as examples so a drawing does not invent an ordering,
load-balancing rule, or delivery guarantee.

Give each resource a separate header containing its kind and name. Put its
data or current state in the body below that header. Use the same structure
for streams, groups, instances, and other resources. Nested resources keep
their own headers within the containing resource's body.

Put supporting commentary in the caption rather than among a resource's
data or state. In the stream diagram, message ids belong in the body.
“Stored in Postgres” belongs in the caption because it explains the storage
context rather than representing another item in the stream.

Keep rendered text sizes, header heights, and node padding consistent across
related diagrams. Choose each diagram's display width to preserve that scale.
The overall canvas can be wider or taller as the content changes. Do not fit
every diagram to the same width if that makes its text and nodes smaller.
Sibling cards should have equal body heights, with their content vertically
centered and their headers unchanged. Use the shared diagram components,
not blank lines or per-diagram layout patches. At phone widths, preserve
readable text and contain horizontal scrolling within the figure.

Use diagrams to expose relationships that prose makes hard to hold in mind.
Avoid decorative icons, a diagram of every implementation layer, and motion
that does not explain a state change. Start with static figures. Add
interaction only when changing an input helps the reader reason about the
result. Keep labels readable on phones, include meaningful alternative text,
and make the figure understandable without color. Check both site themes.

### Make the visual explain the mechanism

Treat diagram design as part of writing the explanation. Budget time for
several iterations and reader review. A technically correct figure can still
leave readers unsure what they are seeing. Review the rendered figure in the
article, where its scale and surrounding explanation matter.

- Choose the visual form from the reader's difficulty. A snapshot can show
  processing and waiting together more directly than a flowchart of steps.
- Give each participant a visible role and relationship. A stream name
  tucked into an instance header does not explain how they relate. Show the
  stream separately and identify the claim connecting it to the instance's work.
- Build on patterns readers have already learned. If an earlier diagram
  puts instances inside a consumer-group frame, preserve that nesting when
  expanding an instance to show its queue. Enlarge the frame and canvas to
  make room while keeping the established text and node scale.
- Make spatial order agree with the action shown. For a queue moving toward
  processing on the right, show waiting messages as 4, 3, 2, then processing
  message 1. Label this as the example's state, without implying that every
  consumer processes messages sequentially.
- Draw scope precisely. A lease bracket spans exactly the messages it covers,
  rather than the full card width. Boundaries, alignment, and spacing carry
  meaning just as labels do.
- Anchor time to an event. State that the lease clock starts at the claim
  and continues while messages wait. A clock or deadline without its starting
  event leaves the central mechanism unexplained.

The opening figure in [Consumer leases](src/content/docs/concepts/consumer-leases.mdx)
is the worked example for these rules. A reader should be able to identify
the stream, group, instance, processing message, waiting messages, lease
scope, and clock's starting event from the figure itself.

## Review before calling it ready

- Make a dedicated terminology and naming pass. Check that each term keeps
  one meaning and each named resource keeps one identity across every form
  of the explanation. Resolve inconsistencies before calling the draft ready.
- Can a reader answer the central question from the opening and headings?
- Does the worked case actually demonstrate every claim attached to it?
  If it introduces four message ids, do those ids have a purpose?
- Does any short claim require the reader to reconstruct an unexplained
  causal step? Expand the example where that understanding is needed.
- Do the prose and diagram agree about scope, state, timing, and ownership?
- Are durable completion and external effects distinguished where needed?
- Does every link name useful detail at the point the reader needs it?
- Could a paragraph be moved unchanged into Reference or a Guide? If so,
  decide whether it is needed to explain this page's central question.
- Does each section stop once its question is answered? Are examples easy
  to distinguish from the rules they illustrate, without unrelated caveats?
- Does another page already explain this mechanism better? Keep only the
  context needed here and link to that explanation.
- After removing or moving a section, check what depended on it. Phrases
  such as “both groups” or “the audit group's cursor” need an introduction
  that still exists. Restore the necessary setup without restoring the
  duplicated explanation. Check smaller cuts too: after removing a record's
  field inventory, “after the recorded time” may need to become “after the
  backoff period” to refer to the explanation that remains.
- Have claims been checked against shipped code? Keep the slop level honest.
- Does the page work at phone width and in both themes, with readable labels
  and no horizontal page overflow?

Then run the revision checklist in [VOICE.md](VOICE.md). Keep only the visuals and
prose that help the reader explain or predict the behavior themselves.

## Basis for this playbook

[Diátaxis: Explanation](https://diataxis.fr/explanation/) supplies the scope:
understanding, connections, reasons, and bounded explanation.
[Google: Illustrating](https://developers.google.com/tech-writing/two/illustrations)
supplies the visual guidance: takeaway captions, limited information, and
progressive detail. [RabbitMQ: Work Queues](https://www.rabbitmq.com/tutorials/tutorial-two-python)
provides an example of introducing normal work before a failure that motivates
a mechanism. Borrow that explanatory sequence, not its tutorial scaffolding
or RabbitMQ-specific delivery behavior.
