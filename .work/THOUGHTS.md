# Public API

Need to change this Warning to include 'or destroy consumer group with example command'
or just rethink if we should be sending this alert of not. It could get rather noisy and annoying while testing things, its tough
```
time=2026-09-16T19:40:01.459-04:00 level=WARN msg="alert active" schema=sqlstreams worker=alert.worker_liveness group=alert.worker_liveness alert=worker_liveness alert_message="stream \"signup.welcome-email\" has no live instance on 6 of its worker rows" detail="Nothing is running: email-sender-beginning (exception_consumer, manager, message_consumer), email-sender-head (exception_consumer, manager, message_consumer). A worker row with no live instance does no work: expired partitions are not dropped, exceptions are not retried, and the group's cursor stops advancing." hint="Run \"sqlstreams manager run\" in a process that stays up, or start a consumer on the stream -- either one claims these rows." owner=signup.welcome-email severity=warn
```

# Docs

Get back roadmap and benchmark doc and create Overview board (quickstart, why sqlstreams, benchmark, roadmap)

Consider changing crying cat profile pic, its a bit distracting and out of place

Really need a single Overview page that goes into the main concepts to understand:
- stream
- producer
  - auto batching
  - idempotency
- consumer
  - batching
  - queue
  - claim range
  - lease
  - timeouts

Need to tweak and visually improve the subsection headers in reference board
- The '> Consumer Related Pages' drop down should do the same thing

Improve the docsite benchmark chart/graph it is not very **clear**

Roadmap feels lifeless and not exciting. Not sure what to do about that but something should be done

keyword highlighting and or linking (stream, produce, consume, dead)

thought bubbles (hover over text and I interject my random thoughts)

Announcments need to be curated

Should add a nerd 'But what about!' note or section in doc for when obvious point or counter argument pops up to address it

roadmap later item for review code for interesting design decisions to write articles on (big or small things)
- diagnostic code system
- suppression logger
- doc site
- client builder pattern choice
- writing your own commit messages (and small commit size)
- Decision tracking and index
- Code quality is even MORE important now
  - bad code and patterns snowball, just as good code and patterns do
  - stray unused code of patterns dilutes context
  - skimmable code (ie readable) is more important than ever
- AI Native projects
  - faster but lose context
  - you can learn but its worse and must be disciplined? (does learning even matter)
  - Is it right or fair?
  - How could this be maintained long term
  - amount of corrections for frontier models (50%)
  - When to heavily review the code vs skim or not at all (what is the right balance)
- The importance of keeping notes (previously the slower pace would allow for easier tracking, things change so quick now it is easier to forget)
- rules vs conventions and the tradeoffs of each
  - rules are enforced via code, scripts etc. Have maintainenance and overly aggressive rules can be annoying and brittle
  - conventions easy and work well with workflows but easily accumlate drift overtime (if large enough project)
- Hitting a flow state with two sessions. The new heads down coding joy. Thought I lost the joyouse moment but it still does feel good.
- If you are not cursing out your llms I am concerned about your coding capabilities
- the evolution and stages of testing
  - should you start out with unit tests and increase token costs or wait till code is closer to finalization
  - what are valuable tests in the agent era
  - is testing validation logic valuable, setting up integration test that you don't understand?
  - unit tests are dead (few small mistakes with llm, id vs name)
  - regression tests from debug sessions
- A new world and the case of low dependencies
- the rot within (VPs and middle managers), fuck em those shady ass mf
- When to use kafka vs sqlstreams
  - kafka is far more efficient
- The difference between Opus 4.6 and Fable 5.1 is both HUGE and miniscule
  - its ability to one shot things and solve for inputs and outputs
  - but its effect on day to day maintaining and extending code is minimal
- Showing a feature coding session with setup and prompts and how much back and forth there is
  - would need snapshots of git diffs after each prompt or manual change
- why topic per table instead of single table with LIST and RANGE subpartitions (0757 decision)
- why cursor claim ranges, instead of singular or bit map with holes

# Review

Make sure all our payload structs in examples, docs and quickstart suffix with V1 for good patterns

## Manual

Probably should have one more table name and column review (this will be hard to change later)

need to make sure we do some manual testing for cli, metrics and alerts

manual review of public user facing comments :(. I don't want to but its got to be done

review of most important website docs

# Other

Metrics don't quite accurately reflect user expectations with compacted messages ie
```
brandonlouiscate@Brandons-MacBook-Air sqlstreams-quickstart % go run ./cmd/metric
live: head 13, committed 9, backlog 4, dead 0
collected: backlog 3, 21s ago
```
For above case backlog = 4 However there is really only one compact message to consume its just that cursor is 4 messages behind head NOT 4 compacted messages behind head.