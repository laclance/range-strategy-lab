# Entry Readiness Review

Date: 2026-06-13

## Findings

No open critical findings remain after this review.

Fixed findings:

- Medium: `CompressionRangeDetector{}` did not fully match
  `DefaultCompressionRangeDetectorConfig()` because boolean defaults cannot be
  inferred field by field. A zero-value config now returns the declared default,
  preserving `UseBollinger=true` for default detector use
  (`internal/lab/detector.go:105`).
- Medium: `SRBoundaryAuditConfig{}` did not fully match
  `DefaultSRBoundaryAuditConfig()` and could include inactive detector rows when
  called directly with a zero-value config. A wholly zero config now returns the
  declared default, preserving `DetectorActiveOnly=true`
  (`internal/lab/sr_boundary_audit.go:131`).

## Review Result

The project passes the entry-readiness code review gate for the next
non-trading timing audit. It is not ready for first trade entries yet because
the current boundary-rejection evidence is still an ex-post label built from
future bars, not a closed-candle decision rule.

Reviewed areas:

- CLI behavior remains offline-only and continues to use `lab.EmptyStrategy`.
- CSV loading supports Binance and normalized shapes and now has focused tests
  for open/read/parse failures.
- The one-position engine preserves next-bar-open entry, stop-first ambiguity,
  costs, sizing cap, time stop, force close, and short-side behavior.
- Detector and detector sweep logic use prior-window thresholds and deterministic
  output ordering.
- SR audit remains prefix-only: each row calls `sr.Compute` on
  `srCandles[:i+1]`, so no future candles are available to the SR adapter
  (`internal/lab/sr_audit.go:91`).
- SR boundary audit forward fields are labels only. They are calculated from
  candles after the audit row index and must not be treated as entry inputs
  (`internal/lab/sr_boundary_audit.go:161`).

## Coverage

`internal/lab` statement coverage increased from `84.5%` to `99.8%`.

The only remaining uncovered internal statements are defensive error returns in
`RunSRAudit`:

- detector classification error propagation from a known-valid default detector
  config (`internal/lab/sr_audit.go:80`).
- third-party `go-sr` compute error propagation
  (`internal/lab/sr_audit.go:92`).

Those branches are intentionally left uncovered because forcing them would
require test-only seams around deterministic internal defaults or the external
`go-sr` call. The directly controllable SR validation, warmup, metadata,
prefix-only behavior, fallback behavior, and boundary aggregation paths are
covered.

Known tooling limitation:

- `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./cmd/rangelab`
  fails with the snap-confine error.
- `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./...`
  fails for the same `cmd/rangelab` coverage reason, while internal coverage is
  reported and `/usr/local/go/bin/go test ./cmd/rangelab` passes.

## Verification

Commands run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./cmd/rangelab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./internal/lab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -coverprofile=/tmp/range-strategy-lab-internal.cover ./internal/lab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go tool cover -func=/tmp/range-strategy-lab-internal.cover
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv -out-dir results/sr-boundary-inspection-check -sr-boundary-inspect
```

Observed smoke result:

- `sr_boundary_inspect events=281080 comparison_rows=192`
- `loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z`
- `strategy=empty trades=0`

## Next Step

This review gate was followed by the non-trading boundary-rejection timing
audit and `docs/SR_REJECTION_TIMING_REVIEW.md`.

Current verdict after that timing review: boundary rejection is not
entry-ready, and the strategy should remain `lab.EmptyStrategy`.
