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
