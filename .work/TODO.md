# TODO

Sliding window of in-flight work only. Future work lives in ROADMAP.md;
shipped work in HISTORY.md; decision rationale in DECISIONS.md ->
.work/decisions/.

## v0.1.4 release follow-through

- Step 4: OTel pins root v0.1.4; standalone tidy/build/vet/race tests pass.
  Await maintainer commit before publishing otel/v0.1.4, then update and
  verify CLI pins before its commit and cmd/sqlstreams/v0.1.4 tag.
- Step 5: verify downloads and installations; record results in HISTORY.
  Release run 34839844935 published all six archives and checksums.
  Chocolatey install/version/uninstall passed; submission returned HTTP 403.
- Step 6: prepare versioned docs, run site verification, ask before deployment.
