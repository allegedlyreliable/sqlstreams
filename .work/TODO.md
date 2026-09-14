# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## Chocolatey v0.1.4 submission [0803]

- v0.1.2 approval confirmed at 2026-09-14T11:57:05Z. Retry v0.1.4 now that
  the first-version moderation restriction no longer applies.
- Manual retry added to release.yml: select a tag, skip GoReleaser publication,
  reuse Chocolatey install/version/uninstall checks, then submit.
- Await maintainer commit and push, dispatch with tag=v0.1.4, monitor the
  result, then verify package status and record evidence in HISTORY.
