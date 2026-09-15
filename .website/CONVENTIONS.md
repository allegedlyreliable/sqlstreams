# Website conventions

Rules for the .website/ tree. Violations are bugs, not style nits.

The root CONVENTIONS.md stays authoritative for everything
cross-cutting; these sections bind here by reference and are never
restated: ## Naming & terminology, ## Vocabulary, ## Structure,
## Comments, ## Documentation, and the prose grammar of ## Errors and
## Logging for all UI copy (console errors, empty states, search
results). This file holds only what is specific to the frontend
toolchain.

## Stack

The settled set -- a new tool enters only by replacing a row, never
beside one:

- Astro: static output, content collections, ClientRouter view
  transitions. No Starlight.
- Expressive Code owns Markdown/MDX code fences: highlighting, frames,
  and copying. Configure it before MDX in astro.config.mjs; preserve
  comments when copying and use the board's theme selector.
- Svelte 5 is the ONE island framework; .astro files are page
  scaffolding and layouts only -- anything worth a Storybook story is
  a Svelte component.
- Diagrams use Svelte Flow with ELK layout at build time. Articles supply
  typed resources and connections to the shared Diagram component. ResourceCard
  owns card markup and CSS. Diagrams render statically without client hydration.
- CodeMirror 6 (+ lang-sql) and PGlite power the SQL console.
- Pagefind is site search. Its bundle is imported through a variable
  path -- a literal `import('/pagefind/pagefind.js')` breaks Rollup even
  under `@vite-ignore`.
- Versions are frozen deployments on Pages branch aliases [0601]. No
  build carries the version list: every deployment fetches
  `public/versions.json` from the live origin at read time, and
  `public/_headers` grants CORS on that one path -- without it the
  scheme dies silently.
- Vitest (pure functions), Playwright (flows), Storybook
  (svelte-vite + addon-svelte-csf) are the test surface.

## Dependencies

- The root dependency rule applies unchanged: minimal deps, vendor
  battle-tested code when a whole package is not earned. npm makes
  sprawl cheap -- every new package needs the same justification as a
  Go dep, and a utility that a rune module or twenty lines of CSS can
  carry is written, not installed.

## Browser support

Two lines, moved only by deliberate edit, never by a toolchain
default [0606]:

- Supported -- Baseline Widely Available: Chrome/Edge 121+,
  Firefox 123+, Safari/iOS 17.2+ (87% of tracked global usage,
  2026-08). Everything the site does works here. An enhancement
  needing more (view transitions: Safari 18/Firefox 144;
  requestIdleCallback: Safari 18) ships only with the graceful
  fallback it already carries -- an instant swap, a setTimeout.
- Readable -- the pinned build floor: Chrome/Edge 111+,
  Firefox 114+, Safari/iOS 16.4+, the `vite.build.target` list in
  astro.config.mjs (Vite 8's baseline-widely-available set,
  frozen because the default floats per Vite major). Between the
  two lines every page stays readable and navigable; islands lock
  behind their own failure faces.
- Below the floor the site gets CSS's own error recovery and
  nothing else: no @supports guards, no polyfills, no legacy
  bundle. The structural cliff is Safari 15.4/Chrome 99 (@layer,
  :where(), :focus-visible) -- under 0.2% of usage and falling.

## TypeScript

- `astro/tsconfigs/strictest` + `verbatimModuleSyntax`; `astro check`
  and `svelte-check` run in the build -- the dev server type-checks
  nothing.
- `any` is a lint error (`@typescript-eslint/no-explicit-any`); a
  genuine unknown is `unknown`, narrowed at the boundary.
- Every component declares one named Props type; every field has an
  explicit type and is required -- no `?` fields, no fallback inside
  the component: the caller states every fact, defaults included
  (`width={18}` at the call site; a story file's shared values live
  in its meta `args`). Expected absence of a fact is an explicit
  `| null` union -- the caller passes `null` and the compiler forces
  both branches. Never `?: T | null`, and never an empty-value
  stand-in for null.
- Extending a platform type (`HTMLButtonAttributes & {...}`) is the
  one sanctioned answer to attribute explosion -- never a loose index
  signature.
- Every console query gets a named row type listing the columns it
  returns -- the TS sibling of the datastore `db:`-tagged row
  structs. PGlite results are typed at the query call, nowhere else.

## Components

- One folder per component under `src/components/<name>/`,
  kebab-case files: `thread-row/thread-row.svelte` +
  `thread-row.css` + `thread-row.stories.svelte`, plus `types.ts`
  and a `<name>-state.svelte.ts` state module when the component is
  stateful (the `-state` suffix keeps the module from colliding with
  `<name>.svelte` in import resolution). No index.ts barrels --
  callers import the file directly.
- A component used by one route lives beside it in that route's
  `_components/` directory, not in `src/components/` (the placement
  law).
- A page's data logic lives in a route-local `_<page>/` folder,
  one concern per file: `model.ts` holds the typed row shapes the
  page renders (the sibling of a datastore's model.go), the hand-
  curated content that feeds them gets its own named file
  (each board's `board.ts`), and the functions that build rows from the
  collection live in `rows.ts`. The `.astro` frontmatter is
  fetch-call-render only.
- A component's styles live in its sibling `<name>.css`, reached ONLY
  through the one-line `<style src="./<name>.css"></style>` tag --
  svelte-preprocess inlines the file before the compiler runs, so
  scoping is identical to an inline block. A `.css` file is never
  imported (the ESLint no-restricted-imports guard bans it -- an
  imported stylesheet is global css in disguise); the global
  stylesheet's sanctioned import sites are the layout and the
  Storybook preview. An oversized style file splits the component,
  never escapes scoping. Astro cannot scope an external stylesheet
  (a css import in `.astro` is global), so a page carries NO
  `<style>` block at all -- page styling belongs to components, and
  the shared content frame to the layout. The layout alone keeps
  one small scoped block.
- Stateful vs presentational is a file split: logic is a runes class
  in `<name>-state.svelte.ts`, the .svelte file is a thin renderer.
- `$derived` over `$effect` -- `$effect` is for real side effects
  (DOM, storage), never for deriving state (lint-enforced,
  `svelte/prefer-writable-derived`).
- Callback props, never event dispatchers; snippets, never slots;
  `onclick=`, never `on:click`.
- Shared state is a runes module (`.svelte.ts`) imported by whoever
  needs it; no store library. localStorage reads and writes live in
  the owning state module, wrapped so absence or denial falls back to
  defaults.

## Islands & loading

- Static is the default. Every island names its trigger at the use
  site -- `client:visible` or `client:idle`; `client:load` is banned.
  An island whose content is viewport-gated uses `client:media` with
  the breakpoint from ## CSS, paired with CSS that hides its
  server-rendered markup below the same width -- the JS never loads
  where the content never shows [0602].
- Island props are serializable values against the component's Props
  type -- functions never cross the boundary.
- Heavy chunks (PGlite wasm, CodeMirror) load only through dynamic
  import inside their island, never in the initial payload. The
  console's shell is server-rendered -- example SQL as build-time
  highlighted text over real seeded result rows -- and the editor
  swaps in sized identically.
- A Playwright flow test sums the script bytes the homepage loads
  below the sandbox gate -- the island never hydrates there, so the
  count is stable -- against a declared ceiling, and asserts no PGlite
  chunk is requested; past the ceiling is a failing build. The
  build-time import-graph walk stays rejected [0594].

## Errors

The prose grammar for every reader-facing message binds from the root
## Errors by reference (see the preamble); this section owns which
surface a failure uses. Four tiers, disruptiveness matched to scope
[0597]:

- Inline at the source is the default: the failure renders beside the
  control that caused it, an `errorMessage: string | null` prop into a
  `role="alert"` block. A DOM handler or async call reports through its
  own try/catch at the call site -- no other tier sees those throws.
- A section whose rendering can fail wraps its markup in
  island-boundary's `<svelte:boundary>`: render and effect throws
  inside the children swap in the fallback face with a working retry.
  The boundary never catches handler or async throws -- tier 1 stays
  mandatory under it.
- The page-level notice (src/state/site-notice.svelte.ts, the ONE
  site-wide surface) is fed only by the global nets registered once in
  BoardLayout's bundled script: window 'error', 'unhandledrejection',
  and 'vite:preloadError'. A banner for faults the reader can wave
  away; the modal ONLY for the reload-required stale-chunk case. The
  full-page face was cut -- a prerendered shell always leaves readable
  prose [0598]. A component never show()s its own notice: a failure a
  component can name belongs to the tiers above.
- Error toasts do not exist here: an auto-dismissing error is missable
  and inaccessible, so the persistent banner is the only roaming
  surface.

What a message shows:

- Reader-typed SQL fails with the real Postgres message, verbatim --
  the console is a terminal and the error is its output. Every other
  failure speaks the problem -- fix grammar in the site's own words.
- A caught value reaches the page only through helpers/caught-message;
  String() on a thrown object renders "[object Object]".
- Every failure state a component supports is a story -- ## Storybook's
  done-checklist rule applied to error faces.

## CSS

Vanilla CSS only -- native nesting, custom properties; no
preprocessor, no utility framework, no third-party token pack.

- One global stylesheet opens with
  `@layer reset, tokens, base, compositions, utilities;` and is the
  only handwritten global CSS. Expressive Code generates its own code-block
  styles, configured through its style overrides. Scoped component styles
  stay unlayered, so they always win over the layers.
- Tokens are two tiers, one way: primitives (the raw design sheet)
  feed semantic names; components consume semantic tokens ONLY.
  Token names spell out per the root naming rules.
- Raw values -- hex colors, font stacks, z-index integers, magic
  spacing -- exist only in the tokens layer (stylelint +
  declaration-strict-value enforces it). z-index and spacing are
  closed token scales, not ad-hoc integers.
- Two breakpoints, and only these two [0602] -- a media query cannot
  read a custom property, so the values are convention: `640px` is
  the layout collapse (multi-column rows stack, phone-only removals),
  `760px` is the sandbox gate (`max-width: 760px` hides, the island
  hydrates on `min-width: 761px`). Touch-target and input sizing
  adjust under `(pointer: coarse)`, never under a width query.
- base styles the classless HTML that rendered MDX produces.
  compositions are the hand-written layout primitives (.flow,
  .cluster, ...) and carry no color or type. A utility is generated
  from a token and exists only once three call sites need it; until
  then the declaration stays in the block.
- A component's scoped style block stays under ~80 lines; states are
  `data-` attributes (`[data-folder='new']`), never modifier
  classes.
- `:global()` is a sanctioned-crossing list: the base layer's MDX
  styling, view-transition names, and post-frame's reset of the
  classless children a post body receives (`.post-body :global(p)`).
  Anywhere else it is a bug.
- A board style is a tokens-layer redefinition -- its own primitive
  sheet, then a re-point of the semantic names -- selected by
  `data-board-style` on the root element; a component never branches on
  it. The ONE crossing is the code fence: Expressive Code highlights both
  palettes at build time and selects them through `data-board-style`.
  All motion sits behind `prefers-reduced-motion`.
- A page-level shove is a scroll, never a transform: a transform on the
  page becomes the containing block for every `position: fixed` child
  and flings a fixed bar to the document bottom. `animationend`
  bubbles -- a listener on `body` waiting for the page's own animation
  checks `event.target`.

## Content

- Frontmatter is enforced by the collection's zod schema and nothing
  else; enum-shaped fields are `z.enum`; every field carries
  `.describe()` prose.
- The body starts at H2 (the H1 is the frontmatter title); heading
  levels never skip (remark-lint enforces both).
- Markdown first: a component appears in MDX only from the
  whitelisted set (aside, tabs, the console, the compat matrix, diagrams); anything prose can carry stays prose. A
  component earns a place on the list by rendering something prose
  cannot state -- live query results, a computed grid -- never by
  presenting prose more nicely.
- Decision records are internal repository documentation. The docsite does
  not load, render, link to, count, or index them.
- Concepts articles follow CONCEPTS.md for scope, introductions, worked
  cases, visual choices, and review. Guides follow GUIDES.md for task scope,
  setup and execution order, examples, notes, and review. VOICE.md governs
  prose style across all page types.
- Each page does ONE job -- tutorial, how-to, reference, or
  explanation; a guide that starts explaining links to the concept
  page instead of drifting. Four boards organize those purposes: Concepts (explanation), Guides
  (how-to), Reference
  (contract lookup), and Troubleshooting (diagnosis and recovery). Article dropdown trees connect related pages across boards without
  repeating their content. Articles have one owning board except the two homepage
  stickies, which link directly to the Board Index; contextual links can reach it
  from other boards. A diagnostic code keeps one canonical page, reachable
  from both reference and troubleshooting.
- Concepts, Guides, and Reference group their articles beneath group
  headings. Consumer is the current example. Reference's six subsections sit
  inside Consumer; Concepts and Guides list their articles directly under it.
  Article breadcrumbs include the group after the board, linking to that
  group's anchor on the board. Standalone stickies and Troubleshooting keep
  their existing trails. No standalone group route is created.
- The Consumer map connects related articles across the four boards. The
  article frontmatter owns group membership (`group: consumer`). Board
  groups, dropdown trees, and breadcrumbs derive it from the content collection;
  no separate group membership list is maintained. Groups use lowercase
  hyphen-separated identifiers; articles on grouped boards require one. There is no standalone group page or full-map link. A native disclosure row directly beneath the page title
  starts collapsed at every width. Its map highlights the current page and
  opens that board group. The dropdown is a vertical tree of board names
  and indented article links at every width. Expanding it pushes the article down and never
  reduces the reading width. The row is excluded from the article search body. No global Groups tab or fifth
  board is introduced. Contextual prose links still carry their own meaning.
- Page size has mechanical triggers, not taste [0679]: a guide or
  concept thread splits past six H2s or roughly 1,500 words; a section
  under three sentences folds into its neighbor; a reference thread may
  run long but never mixes kinds. A thread that has become several kinds
  of page (a tutorial, an explanation, a reference, a changelog) splits
  by kind, and changelog prose ("X has been removed") is deleted, never
  kept -- the decision records hold history. A moved or split thread
  leaves an astro.config `redirects` entry at its old URL.
- Reference has six sections: Go API, Configuration, CLI, Metrics, Logs,
  and Alerts. Article dropdown trees provide another route to those pages.
  The route's section mapping owns membership and order; it does not move
  a linked guide or code page out of its owning board.
- A reference entry answers one independently useful lookup. Neither a
  handle nor a symbol forces a page boundary. Keep closely related variants
  together when readers need to compare them. Open with the answer and the
  exact name. State applicable inputs, outputs, absence, errors, side effects,
  cancellation, and concurrency beside the operation they qualify.
- Configuration owns settings' types, defaults, zero/nil meanings, validity,
  precedence, scope, and effective timing. Derived defaults name their
  dependencies and when they resolve. Go API owns calling contracts and
  returned values; CLI owns syntax, flags, output, and exits. Metrics owns
  measurement definitions, calculations, units, attributes, and freshness.
  Logs owns emission conditions, levels, fields, and suppression. Alerts owns
  triggering conditions, thresholds, severity, evaluation, and resolution.
  Troubleshooting owns investigation and recovery. Link to shared definitions instead of maintaining a second table. Task-specific
  command inputs can stay with their operation.
- Use descriptive titles and exact symbols in headings or labeled entries.
  Important settings and operations have directly linkable anchors. Examples
  clarify contracts; procedures link to guides. No required Verbs/Config/
  Gotchas skeleton; omit inapplicable sections and place limitations locally.
  A changed contract updates its definition in the same change.
- During the documentation structure review, only the new Consumer subset
  appears on the site alongside the user's original Quickstart and Why
  SQLStreams articles. Other existing article sources remain
  on disk but have no generated article routes or search entries. Board
  membership and the homepage sticky list control visibility, navigation, and counts through
  siteThreads; article dropdown trees consume that same visible collection.
  Do not redirect hidden articles into the preview or list the archive.
  This review exception replaces the moved-page redirect rule until migration.
- The homepage opens with one Start Here section: Quickstart and Why SQLStreams
  stickies, then the sandbox. They are standalone pages, not a Getting Started
  board. These stickies use the original quickstart.mdx
  and why-sqlstreams.mdx without rewriting their prose. Navigation metadata
  may be added to frontmatter. Their existing
  links into hidden articles remain as written during the preview. The board listing follows. Keep board-purpose
  explanations on their boards and group discovery in the article dropdown trees;
  do not add a separate introductory panel above the stickies.
- Every board introduction is one sentence stating its purpose. Omit
  cross-board placement rules and repeated navigation from the intro.
  The homepage stickies own product orientation and a first success. Concepts
  owns mechanisms and causal examples; Guides owns task steps, practical
  choices, and result checks; Reference owns exact contracts; Troubleshooting
  owns symptoms, evidence, cause, recovery, and verification. A code's one
  canonical page can belong to Reference (measurement definition) or
  Troubleshooting (diagnosis); cross-links do not create a second copy.
- The same example group crosses every board. Keep its names, ids, and
  assumptions consistent. Article dropdown trees link to canonical pages grouped by board;
  they do not become another reference manual. Never link into hidden articles to complete an example.
- Code samples show real error handling -- `if err != nil { return err }`
  or `_` for an unused value -- never a `must()` helper: it hides the
  path readers copy and is not a real API. Pages that still carry one
  (schema-versions, new-group-start, replay) are swept when touched.
- An MDX aside inside a list item closes at the item's indent.
- Every thread through the thread route declares `slop` in its
  frontmatter, the author's review depth of the LLM-drafted body:
  `high` (unreviewed, the value a new thread starts at), `medium`
  (structure checked, samples run), `low` (read line by line), `none`
  (hand-checked, renders no notice). The level's copy is fixed in the
  slop-notice component; a page never writes its own. The diagnostics
  reference and the code threads carry no field.
- Contextual internal links sit where the reader needs them. When another
  thread owns a prerequisite, the detailed mechanism, or a relevant
  contrast, link its first useful mention with anchor text that names what
  the reader will find (`See [Message Lifecycle](/concepts/lifecycle/) for
  the delivery states`, never `learn more`). Review this as prose: there is
  no link quota, generated related-thread box, or mechanically added link.
- Vale runs in CI with the Google developer-docs style plus the
  SQLStreams style; the SQLStreams substitution rule mirrors the root
  ## Vocabulary table, and a new vocabulary row updates the Vale rule
  in the same change. Code blocks are exempt (IgnoredScopes = code).

## Storybook

- `.stories.svelte` (addon-svelte-csf) is the ONE story format,
  co-located with its component. Route-local components (`_components/`)
  are outside the story glob and carry no stories.
- An Astro project has no vite.config.ts, so `.storybook/main.ts`'s
  `viteFinal` prepends `svelte()` -- without it every addon parses raw
  .svelte as JS.
- Every supported state is a story -- the story list is the
  component's done-checklist; story count tracks states, never usage
  sites.
- Everything editable is an arg; titles follow `Board/<Name>`. The
  board style is the one toolbar global -- it is the document's state,
  not any component's prop, and every component has as many looks as
  the board has styles.

## Data from the library

- The compat widget reads a build-time JSON export of the migration
  registry -- the gate logic is never reimplemented in TS.
- An embedded SQL literal keeps its `-- sqlstreams: <package>.<method>`
  first line and is drift-checked against the Go source in CI; a
  diverged literal is a failing build.
- Every number on the site is real or absent -- no invented member
  counts, view counts, or activity. Performance claims follow the
  root ## Documentation rule: benchmark records only.

## Verification

One command runs the whole floor (the site's sibling of
`just verify`): Prettier (svelte + astro plugins), ESLint (svelte +
astro + typescript-eslint recommended), stylelint, `astro check` +
`svelte-check`, Vale + remark-lint, Vitest, Playwright. A rule in
this file that a tool can check runs as a check -- review carries
only what machines cannot.
