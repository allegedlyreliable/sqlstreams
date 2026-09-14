# Release

## 1. Verify

Start with a fresh development DB ([recipe](AGENTS.md#verification)).
Pin compatibility to the prior tag ([setup](.tools/compat/go.mod)).

```sh
just verify
just signal-e2e
just compat-lab # Use refused when the registry requires it.
```

## 2. Record

Update the [migration table](.website/src/content/docs/guides/migrations.mdx)
and [HISTORY](.work/HISTORY.md) with test results and release notes.

## 3. Publish

Commit and push the release changes, wait for CI
to pass.

From that commit, publish next version tag `vX.Y.Z`:

```sh
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin refs/tags/vX.Y.Z
```

Let GoReleaser create the GitHub release. Creating it manually first makes
it immutable and blocks artifact uploads. Never move a published tag.

[Automation](.github/workflows/release.yml) publishes archives and, for
stable releases with credentials, Homebrew and a Chocolatey submission.

## 4. Go modules — when updated

`root` → `otel/vX.Y.Z` → `cmd/sqlstreams/vX.Y.Z`

At each step: wait for dependencies to resolve, update pins, tidy/build/test
with `GOWORK=off`, then have the maintainer commit before tagging.

## 5. Check

Verify downloads, installation, and CLI version; record outcomes in HISTORY.
Chocolatey public-feed installation waits for approval.

## 6. Docs — when versioning

Update `site.ts` and `public/versions.json` under `.website/`.

```sh
just site-verify
just site-deploy
just site-freeze <slug>
```

Never reuse a frozen alias.

## Post-release checklist

- [ ] Confirm published root → OTel → CLI pins/tags.
- [ ] Bump `examples/`, `.tests/`, `.bench/`, `.tools/` module pins; tidy and build with `GOWORK=off`.
- [ ] Bump and run the [Go quickstart](https://github.com/allegedlyreliable/sqlstreams-quickstart-go/blob/main/go.mod).
- [ ] Update install versions and release links across READMEs and Markdown/MDX, including the quickstart repo.
- [ ] Add migration-table evidence; preserve historical versions.
- [ ] When versioning docs: update `site.ts` / `versions.json`, deploy, freeze, verify.
- [ ] Verify Homebrew/Chocolatey installs; update HISTORY, ROADMAP, and TODO.

## Gotchas

- Push tags; let GoReleaser create releases. Published releases are immutable.
- A failed workflow may have published assets. Check before retrying.
- Chocolatey can return 403 until the first package is approved.
- Rebuilds can change checksums. Use a new release, not the removed retry path.
- Verify pins with `GOWORK=off`; keep compatibility pinned to the prior client.
- Free port 5432 for database checks; use a fresh site build for browser checks.