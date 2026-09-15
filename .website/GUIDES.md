# Writing guides

Use this when creating or revising a Guides page. Help the reader complete
one task with a clear path from setup to a working result.
[VOICE.md](VOICE.md) governs general prose and terminology.
[CONVENTIONS.md](CONVENTIONS.md) governs site structure and code.

## Keep the task in focus

Be direct and straightforward. Keep information that helps the reader take
an action, make a necessary choice, or recognize the result. Cut material
that does not add value to the task at hand.

Repeat information when it helps readers grasp a difficult concept. Avoid
repeating straightforward checks or advice across setup, running, and tuning.
Put each check beside the action or decision it supports.

Name the task in the title. Open with enough surrounding context to explain
when the task matters and what the reader will accomplish. Avoid a long
motivation, an isolated definition, or a catalogue of everything involved.

Assume basic competence in the tools the guide requires. State concrete
prerequisites without turning the guide into an introductory tutorial.
Explain a step's purpose when it affects how readers carry it out.

Keep the endpoint narrow. A guide to stopping a consumer does not need to
become a stop-and-restart recovery exercise. Include follow-on tasks only
when they are necessary to accomplish the task named in the title.

## Follow setup → how to run → more information

This is the reading order, not a mandatory set of literal headings.

1. **Setup.** Establish the required starting state, identify the file or
   configuration to change, and provide the code needed for the task.
2. **How to run.** Put the commands and expected behavior immediately after
   setup. Let readers try the result before asking them to read further.
3. **More information.** Follow with useful adaptation guidance, common
   qualifications, and links to deeper explanations or exact contracts.

Do not separate setup from execution with several sections of background,
tuning advice, or deployment considerations. Information required to run
the example belongs before the command. Optional detail can follow it.

Use action headings and number steps when their order matters. Keep each
step focused on one useful action. A heading does not need an introductory
sentence that merely repeats it.

## Make examples usable

- Say where code belongs and whether it replaces or extends existing code.
  Include the surrounding setup needed to understand dependencies and order.
  A complete file is useful when fragments would make readers assemble the
  solution themselves. Do not require a new project for every small task.
- When continuing another example, name that starting point and preserve its
  resource names. Link prerequisites without making readers repeatedly switch
  pages to reconstruct the procedure.
- Use the ordinary commands readers would use. Do not add a build, wrapper,
  or setup step merely because it made the author's testing easier. Verify
  claims that a particular command is required before prescribing it.
- Do not add environment-isolation steps, such as stopping other instances,
  solely to make the demonstration easier to control. Require them only when
  the reader's task depends on them.
- State necessary substitutions. Avoid narrating details already obvious
  from the code unless misunderstanding them would prevent completion.
- Keep error handling real and identify simulated work. A handler that waits
  and prints output must not be presented as sending an email or charging a
  payment provider.

## Show what actually happens

Run the documented commands and capture their output. Verify the same entry
point and interaction readers will use. Testing a compiled binary does not
by itself verify instructions that use `go run` and terminal Ctrl-C.

Distinguish library logs from messages printed by example code. Preserve the
actual log format and wording. Label an excerpt as an excerpt and explain
variable fields such as timestamps or message ids. Do not present selected
application prints as the full output or construct an idealized transcript.

Give readers a short, observable success check beside the command. State
what the observation establishes without claiming more than it proves.
Keep the check proportional to the task. The author's validation may require
database inspection, exit-status checks, or repeated runs without making all
of those activities part of the reader's procedure.

## Use ThreadAside notes for common qualifications

ThreadAside is the preferred pattern for brief supporting information or a
qualification that most readers are likely to need. Give the note a specific
title and place it beside the relevant action. Use a note when the information
is useful but does not warrant its own numbered step or section.

For example, “Press Ctrl-C once” belongs beside the stop instruction because
a second signal forces an immediate exit. A short deployment note can point
to shutdown timing details without becoming a separate tuning section.

Reserve what-ifs and gotchas for very common situations that most users
would want to know about while doing this task. Technical possibility alone
does not justify inclusion. Do not add speculative failure tables, exhaustive
warnings, or branches for every supported configuration. A note is not an
excuse to retain unnecessary content. Put detailed diagnosis in Troubleshooting,
mechanisms in Concepts, and settings or contracts in Reference.

Required actions stay in the main steps. A note must not hide a prerequisite
or contain the only instruction needed to make the example work. If several
notes accumulate, reconsider the scope and order of the guide.

## Choose formatting that supports action

Use code blocks for edits and commands, followed by the relevant result.
Use a short list for parallel actions or choices. Add a visual only when it
clarifies what to change, the required order, or the expected state.

Finish useful onward links under **Related topics**, with one link and its
purpose per item. Keep links needed during a step beside that step. Do not
add a recap that repeats the procedure.

## Review before calling it ready

- Can readers identify the task and required starting state immediately?
- Does setup lead directly into how to run, before optional information?
- Can readers place and run the code without reconstructing missing pieces?
- Do commands and captured output match the example as written?
- Is the success check useful and limited to what the task needs?
- Would most readers benefit from each what-if or gotcha while doing this
  task? Are brief supporting details in ThreadAside notes, with required
  actions in the main steps?
- Can any paragraph or extra step be removed without making the task harder?
- Does the guide stop at its intended result rather than adding a second task?
- Are names and terminology consistent across code, prose, and output?

Then apply the revision checklist in [VOICE.md](VOICE.md). Review the rendered
page for readable code, clear step order, and notes that support the flow.

## Basis for this guide

[Diátaxis: How-to guides](https://diataxis.fr/how-to-guides/) supplies the focus
on real tasks, logical sequence, and practical scope.
[Google: Procedures](https://developers.google.com/style/procedures) supplies
the action, command, and result pattern.
[DigitalOcean reader feedback](https://news.ycombinator.com/item?id=18658709)
and [Stripe reader feedback](https://news.ycombinator.com/item?id=31342008)
highlight usable instructions and meaningful examples. These are qualitative
signals, not proof that every feature of those sites should be copied.

The review of [Stop a consumer gracefully](src/content/docs/guides/stop-a-consumer.mdx)
established the local pattern: setup, run, then supporting information, with
brief ThreadAside notes and checks that stay focused on the task.
