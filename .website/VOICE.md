# Voice

How AI-drafted prose on the doc site sounds like its author. Binds
the prose written for .website/ — guides, concept pages, homepage and
board copy. It does not touch error/log message grammar (root
## Errors, ## Logging), decision records (append-only, rendered
as-is), or code comments (root ## Comments). The root ## Vocabulary
registry outranks everything here.

## How to use this file

- Write against the samples: match their structure, stance, diction,
  and rhythm. The samples are verbatim and keep the author's real
  typos — correct spelling, punctuation, and apostrophes in what you
  write; reproduce everything else.
- Samples outrank rules. Where a rule and a sample disagree, the
  samples win — the rules are the checkable summary, not the source.
- After drafting, run the ## Revision checklist as its own separate
  pass. Rule-following decays over a long draft; the checklist is the
  enforcement, not the generation prompt.

## The register

Terse, lowercase-plain, zero ceremony. The answer arrives in the
first sentence, the mechanism after it, and the piece stops when the
content runs out — no wrap-up sentence, no restated thesis. Judgments
come as a flat verdict with the cost attached. Uncertainty is said
plainly the moment it exists. Humor is deadpan, self-directed, and
bolted on as a tail clause — never a performed joke. Enthusiasm is
one understated clause, never an exclamation point.

## Samples

Explaining a mechanism:

> "A long running job in phase1 would hold open a transaction the
> entire lifecycle. With high concurrency a huge number of
> connections would remain open which is not scalable."

> "Worker crashes mid-range -> Lease is 'lost' -> Lease expires ->
> worker reclaims on new claim cycle."

> "backlog - the classic consumer lag metrics. Means you are trailing
> behind head which is normally not good."

> "just to confirm we would still only ever claim the most recent
> 'deferred' message for the key once its finishes correct? ie:
> key busy -> defer added 1 -> defer added 2 -> key finish ->
> claim defer 2, mark 1 as superseded ?"

> "We must take a read only snapshot of xmax before our cursorSql.
> If we did some kind of insert first that insert itself would submit
> an inflight transaction shifting xmax to a number that could never
> be completed by the time cursorSql was run such that the fresh pair
> would never win"

> "The way I'm thinking of this is a user has k8s producer 'app' with
> a few replicas that has autoscaling enabled. Scale ups and downs
> (important one) happen frequently and often throughout the day. The
> potential for ambiguous commit outcomes just seems very bad. But
> I'm not sure if I'm overthinking it because I don't know what exact
> scenario it would be bad in"

Judging and deciding:

> "a seperate janitor process is a no-go. One of the important parts
> of this projects is no additional servers or runners or cronjobs
> you have to manage just producers and consumers."

> "imo disk space is cheap, what about throughput concerns or cpu /
> mem overhead for the db itself?"

> "I think option A makes the most sense I just worry about our
> users. To be frank I'm not sure many people understand transactions
> to well. They have heard the term but rarely think about them or
> their implications. And A fully leans into a transaction and its
> implication"

> "nah I rolled back these changes - I don't like having to
> coalesce... not worth it too me"

Humor:

> "profile page and the descent into madness"

> "I spent waay too long on this, but I love it"

> "mobile friendly design or whatever"

> "wtf? nobody picks this option. Are you a psycopath? You're trying
> to tell me that for every website you go to you manage the cookie
> preferences?"

## Same passage, two ways

The pairs below are the same content written generically and written
in this voice. When drafting, name the differences between your draft
and the voice version before revising — the comparison is the point.
Replace these pairs over time with real before/after edits from
shipped pages; a hand-corrected pair beats a constructed one.

Generic:

> Leases are the heart of SQLStreams's delivery guarantees. When a
> consumer instance crashes, its lease doesn't simply disappear — it
> expires. This isn't a failure mode; it's the design working as
> intended. Expired leases are reclaimed on the next claim cycle,
> ensuring no message is ever lost and processing continues
> seamlessly.

This voice:

> A lease is how claimed work survives a crash. If an instance crashes
> mid-range -> lease eventually expires -> the next claim cycle reclaims
> it. And so no message is lost. The cost is a redelivery window as long
> as the lease, so anything already processed in that range can technically
> run twice.

What changed: the aphorism opener became a definition; "isn't a
failure mode; it's the design" (a contrastive binary) was cut;
"seamlessly" was cut; the passage ends on the cost, not on
reassurance. The mechanism became an arrow chain — but one carried
inside a grammatical sentence ("If an instance crashes..."), not a
bare fragment. Connective tissue stays: "And so" marks the
consequence, "eventually" names the reclaim lag, "technically" flags
true-but-rare. Compression cuts ceremony, never syntax.

Generic:

> It's important to note that calling sendEmailConfirmation() inside
> your transaction can lead to unexpected behavior. If the
> transaction rolls back, the email has already been sent — creating
> an inconsistency between your database state and the real world. To
> avoid this issue, always defer side effects until after the commit
> is confirmed.

This voice:

> Don't call sendEmailConfirmation() before the transaction is known
> to commit. If a later step fails and rolls back the database will
> undo its half but that email can't be 'unsent'. Side effects must
> wait for the commit.

What changed: "It's important to note" and "can lead to unexpected
behavior" were replaced by the actual instruction and the actual
failure. The asymmetry is stated as one plain sentence with "but" —
never a mirrored epigram ("the database undid its half, the real
world can't" is the banned cadence even when it lands). The cause is
spelled before the mechanism ("fails and rolls back"), the abstract
"real world" becomes the concrete "that email", a coined word gets
single quotes ("'unsent'"), and the closing rule carries its modal
("must wait") — then the passage stops dead.

## Rules

Measurable only — anything a checklist pass can verify.

- The answer is the first sentence. Definition or verdict first,
  mechanism second, implications last.
- No closing sentence. When the content ends, the piece ends.
- Causal chains may be written as arrow steps in prose:
  `crash -> lease expires -> reclaim`.
- Comparisons may use the paired-label shape: "for X - ... for
  Y - ...", including naming a symmetry instead of re-deriving it
  ("it is mostly the opposite").
- A judgment carries its cost or reason in the same sentence:
  "not worth it - <cost>", "a no-go - <reason>".
- Real uncertainty is stated in first person where the doc's frame
  allows it ("I'm not sure this holds when...") and never faked as
  confidence.
- Zero exclamation points. Zero emoji. At most one CAPS word per
  page for stress.
- No semicolons in newly written or revised doc prose. Split the sentence
  or connect its clauses with words. Existing content does not need a
  retroactive sweep. This rule governs prose, not code syntax.
- Use paragraph breaks a little more often to emphasize a key point or
  consequence. A short paragraph can give that point room. Keep connected
  reasoning together rather than putting every sentence on its own line.
- Identify what a name refers to on first use: “the `payments.requested`
  stream,” “the `charge-cards` group,” or “consumer instance A.” Formatting
  and naming conventions do not establish that context. Repeat the noun
  when a bare name could refer to more than one thing, including in captions
  and alternative text. Keep application, group, and instance identities clear.
- Explain what a technical term denotes before relying on it as shorthand.
  An incidental mention does not establish a concept, even when the term is
  an ordinary word such as “range.” Tie it to the concrete example when it
  first matters. After a change of section or scenario, restore enough context
  that readers do not have to reconstruct the referent. Prefer the concrete
  description when introducing another term would not help the explanation.
- Sentences vary in length; a short verdict sentence may sit beside a
  long mechanism sentence. Do not sand every sentence to the same
  medium length.
- Spelling, apostrophes, and punctuation are corrected — the typos in
  the samples are fingerprints of the source, not the style.

## Revision checklist

Run over the finished draft as its own pass; fix what it catches.

Strip (generic AI tells):

- em-dash aphorisms and mirrored-clause cadence ("the check is the
  action")
- "not X but Y" contrastive binaries
- wrap-up closers ("In summary", a restated thesis, a reassuring
  final clause)
- significance-signaling ("crucially", "importantly", "it's worth
  noting")
- ornamental diction: robust, seamless, elegant, leverage, delve,
  comprehensive, powerful
- rule-of-three lists where two or four items are the honest count
- hedges that assert nothing ("can lead to unexpected behavior" —
  name the actual failure)

Check (this voice's marks):

- first sentence answers; last sentence is content, not ceremony
- at least one concrete worked case with real names where the page
  explains a mechanism (a named group, real ids — the root AGENTS
  rule)
- costs stated next to the verdicts they belong to
- any real open question admitted plainly
- humor, if present, is a deadpan tail clause and would survive being
  deleted

## Sample sources

Future samples come only from the author's hand-typed writing:
session messages, git commit subjects, and original explain-it-back answers
supplied by the author (the retired archive is no longer in the repository).
Never from decision records, ROADMAP, site
prose, or the sharpened explain-it-back answers — those carry the
AI-drafted register, and sampling them teaches the AI to imitate
itself. Hand-corrected before/after pairs from shipped pages outrank
everything above.

## Concept article playbook

Use this when creating or revising a Concepts page. Start with the reader's
question, then choose the shape that explains it. This is an editing guide,
not a required set of headings.

### Define the job before drafting

Write down what the reader should be able to predict after reading. For
consumer groups: which instances share work, which groups progress
independently, and what happens when an instance stops.

Name the assumed knowledge. Introduce new terms at the point where the
example needs them. A definition that requires three more definitions is
starting too deep in the implementation.

Keep one central question per article. Include a failure or limitation when
it changes the answer. API variants, field catalogues, setup instructions,
and recovery procedures belong on their own pages; link at the relevant
sentence instead of summarizing those pages at the end.

### Build the explanation around a case

- Open with the surrounding situation and the need that makes the concept
  useful, then name the concept. Let this establish the page's intent naturally;
  avoid an isolated definition or “this page explains” preamble.
- Establish normal behavior with a named application, concrete messages,
  and explicit assumptions. Explain what each participant knows or owns.
- Change one condition at a time: another instance, another group, a restart,
  or a crash. Carry the same names and ids through each variation.
- State the mechanism, then the consequence. Distinguish in-memory results,
  durable state, and external effects when that boundary matters.
- Put the limitation beside the behavior it qualifies. End when the question
  is answered; do not add a recap that repeats the opening.

For example, the Consumer page moves through shared work, independent
progress, persistence across restarts, and redelivery after a crash. A page
about ordering might instead compare two execution traces. Choose the
sequence from the explanation rather than copying an outline mechanically.

### Choose the form that does the work

| What the reader needs to understand | Use |
| --- | --- |
| Who owns what, or which components share state | A small labeled diagram with explicit boundaries |
| What changes over time | Numbered causal steps or a timeline |
| How two cases differ | A compact comparison using the same attributes |
| Why a behavior exists or what it costs | Connected prose beside the example |
| An API expression essential to the mechanism | The smallest real code expression that proves the point |

Use headings that name behavior: “Different groups keep independent
progress” tells the reader more than “Group behavior.” Keep paragraphs on
one idea, but retain the connective words that explain cause and effect.
A concept page should be understandable without running its code.

### Give each visual one job

Write the takeaway caption first. Draw only what supports it. Put the visual
next to the explanation, with its caption underneath. Reuse names, shapes,
and direction across related figures; introduce detail progressively.

Show shared storage once. Draw group boundaries explicitly. Label arrows
with their meaning when direction alone is ambiguous. Mark snapshots and
possible assignments as examples so a drawing does not invent an ordering,
load-balancing rule, or delivery guarantee.

Give each resource a separate header containing its kind and name. Put its
data or current state in the body below that header. Use the same structure
for streams, groups, instances, and other resources. Nested resources keep
their own headers within the containing resource's body.

Keep rendered text sizes, header heights, and node padding consistent across
related diagrams. Choose each diagram's display width to preserve that scale.
The overall canvas can be wider or taller as the content changes. Do not fit
every diagram to the same width if that makes its text and nodes smaller.

Use diagrams to expose relationships that prose makes hard to hold in mind.
Avoid decorative icons, a diagram of every implementation layer, and motion
that does not explain a state change. Start with static figures; add
interaction only when changing an input helps the reader reason about the
result. Keep labels readable on phones, include meaningful alternative text,
and make the figure understandable without color. Check both site themes.

### Review before calling it ready

- Can a reader answer the central question from the opening and headings?
- Does the worked case actually demonstrate every claim attached to it?
  If it introduces four message ids, do those ids have a purpose?
- Do the prose and diagram agree about scope, state, timing, and ownership?
- Are durable completion and external effects distinguished where needed?
- Does every link name useful detail at the point the reader needs it?
- Could a paragraph be moved unchanged into Reference or a Guide? If so,
  decide whether it is needed to explain this page's central question.
- Have claims been checked against shipped code? Keep the slop level honest.
- Does the page work at phone width and in both themes, with readable labels
  and no horizontal page overflow?

Then run the existing voice revision checklist. Keep only the visuals and
prose that help the reader explain or predict the behavior themselves.

### Basis for this playbook

[Diátaxis: Explanation](https://diataxis.fr/explanation/) supplies the scope:
understanding, connections, reasons, and bounded explanation.
[Google: Illustrating](https://developers.google.com/tech-writing/two/illustrations)
supplies the visual guidance: takeaway captions, limited information, and
progressive detail. [RabbitMQ: Work Queues](https://www.rabbitmq.com/tutorials/tutorial-two-python)
provides an example of introducing normal work before a failure that motivates
a mechanism. Borrow that explanatory sequence, not its tutorial scaffolding
or RabbitMQ-specific delivery behavior.
