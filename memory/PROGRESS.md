# Progress

## 2026-06-13

SR confirmation timing review milestone:

- Added durable review report:
  - `docs/SR_CONFIRMATION_TIMING_REVIEW.md`
- Updated `README.md` docs order to include the review note.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` for the next materially different
  non-trading hypothesis.
- Review verdict:
  - delayed confirmation after SR rejection is not entry-ready
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this audit
  - do not continue broad rejection-confirmation slicing unless the next
    hypothesis changes materially
- Inputs reviewed:
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.csv`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.csv`
- Audit size:
  - candidate rows: `9,692`
  - summary rows: `72`
  - candidate CSV lines including header: `9,693`
  - summary CSV lines including header: `73`
- Compact evidence:
  - broad decision-confirmation cohorts were about `48.9%` to `49.6%` of seed
    rejection candidates
  - side/delay/horizon aggregate favorable-minus-adverse topped near `+0.96bp`
  - recent split support turned flat or negative in several rows
  - FGTA was mostly below `50%`
  - with at least `500` rows in every split, best minimum split diff fell to
    about `+0.16bp`
  - with at least `750` rows in every split, no stable cohorts remained

Latest confirmation-review verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- Verification passed.
- Existing generated artifacts were used; the audit smoke was not rerun because
  `results/sr-confirmation-timing-audit/` was present and current.

SR confirmation timing audit milestone:

- Added CLI flag `-sr-confirmation-timing-audit`.
- Added compact non-trading delayed-confirmation outputs:
  - `sr_confirmation_timing_candidates.csv`
  - `sr_confirmation_timing_candidates.json`
  - `sr_confirmation_timing_summary.csv`
  - `sr_confirmation_timing_summary.json`
- Defaults:
  - confirmation delays: `1`, `2`, `3` bars after the seed rejection candle
  - horizons: `1`, `3`, `6`, `12` bars after the confirmation candle
  - `detector_active=true` seed rows only
- Decision semantics:
  - seed candle must be an existing SR rejection candidate
  - confirmation candle is the decision candle
  - all forward outcome metrics remain `label_*` fields and start after the
    confirmation candle
- Added focused tests for no-lookahead decision features, label-window start,
  support/resistance symmetry, seed filtering, end-of-data skipping, invalid
  config, candidate aggregation, and summary denominators.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest confirmation-audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-confirmation-timing-audit \
  -sr-confirmation-timing-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-confirmation-combined-compat-check \
  -sr-audit \
  -sr-boundary-audit \
  -sr-boundary-inspect \
  -sr-rejection-timing-audit \
  -sr-confirmation-timing-audit
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- New confirmation audit printed:
  - `sr_confirmation_timing_audit candidate_rows=9692 summary_rows=72`
  - `delays=1;2;3`
  - `horizons=1;3;6;12`
  - `detector_active_only=true`
  - `strategy=empty trades=0`
- New confirmation audit CSV lines including header:
  - `sr_confirmation_timing_candidates.csv`: `9,693`
  - `sr_confirmation_timing_summary.csv`: `73`
- Result paths:
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.csv`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.json`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.csv`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.json`
- Combined SR compatibility check preserved existing counts:
  - SR audit rows: `569,313`
  - near-support rows: `126,861`
  - near-resistance rows: `128,085`
  - boundary events: `281,080`
  - boundary quality rows: `192`
  - boundary inspect comparison rows: `192`
  - rejection timing rows: `968` candidates / `24` summary rows
  - confirmation timing rows: `9,692` candidates / `72` summary rows
- Compact aggregate read of `sr_confirmation_timing_summary.csv`:
  - broad side/delay/horizon decision-confirmation cohorts were about
    `48.9%` to `49.6%` of seed rejection candidates
  - decision-candidate favorable-minus-adverse was small-positive across the
    aggregate side/delay/horizon grid, topping out at about `+0.96bp`
  - favorable-greater-than-adverse rates were still mostly below `50%`
- No durable promotion or no-promotion review document was added in this
  milestone; the next step should review split and cohort stability before any
  entries.

SR rejection timing review milestone:

- Added durable review report:
  - `docs/SR_REJECTION_TIMING_REVIEW.md`
- Updated `README.md` docs order to include the review note.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` for the next non-trading
  confirmation-audit milestone.
- Review verdict:
  - boundary-rejection timing audit is not entry-ready
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this audit alone
- Compact evidence:
  - broad support decision candidates were flat-to-negative by
    favorable-minus-adverse across most horizons
  - resistance was better but still tiny in aggregate, topping out at
    `+0.45bp` decision-candidate favorable-minus-adverse at `12` bars
  - top OOS/recent cohorts were split-specific or side-specific rather than one
    simple support/resistance-symmetric template
  - common h12 in-zone strength-2 shape was unstable across pierced state and
    splits
- Result paths reviewed:
  - `results/sr-rejection-timing-audit/sr_rejection_timing_summary.csv`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_candidates.csv`
- Next recommended work:
  - stay non-trading
  - test delayed confirmation after an SR rejection candle, re-indexed so the
    confirmation candle is the decision candle and future `label_*` outcomes
    start after that confirmation candle

Latest review verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-rejection-timing-audit \
  -sr-rejection-timing-audit
```

Timing audit smoke result:

- `sr_rejection_timing_audit candidate_rows=968 summary_rows=24`
- `strategy=empty trades=0`

Boundary-rejection timing audit milestone:

- Added CLI flag `-sr-rejection-timing-audit`.
- Added compact non-trading outputs:
  - `sr_rejection_timing_candidates.csv`
  - `sr_rejection_timing_candidates.json`
  - `sr_rejection_timing_summary.csv`
  - `sr_rejection_timing_summary.json`
- Candidate cohorts group decision-candle features separately from forward
  labels:
  - side, horizon, close location, touched/pierced/closed-back state
  - wick-beyond, strength, and distance buckets
  - balanced detector context
  - all forward outcome metrics use `label_` prefixes
- Added tests for no-lookahead decision features, support/resistance symmetry,
  touch/pierce/closed-back behavior, bucket behavior, detector-active
  filtering, missing-future skipping, and summary denominators.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest timing-audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection-check \
  -sr-boundary-inspect

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-rejection-timing-audit \
  -sr-rejection-timing-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-combined-compat-check \
  -sr-audit \
  -sr-boundary-audit \
  -sr-boundary-inspect \
  -sr-rejection-timing-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-both-check \
  -detector \
  -detector-sweep
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- Existing boundary inspect check printed:
  - `sr_boundary_inspect events=281080 comparison_rows=192`
  - `strategy=empty trades=0`
- New timing audit printed:
  - `sr_rejection_timing_audit candidate_rows=968 summary_rows=24`
  - `strategy=empty trades=0`
- New timing audit CSV lines including header:
  - `sr_rejection_timing_candidates.csv`: `969`
  - `sr_rejection_timing_summary.csv`: `25`
- Result paths:
  - `results/sr-rejection-timing-audit/sr_rejection_timing_candidates.csv`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_candidates.json`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_summary.csv`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_summary.json`
- Combined SR compatibility check preserved existing counts:
  - SR audit rows: `569,313`
  - near-support rows: `126,861`
  - near-resistance rows: `128,085`
  - boundary events: `281,080`
  - boundary quality rows: `192`
  - boundary inspect comparison rows: `192`
  - timing audit rows: `968` candidates / `24` summary rows
- Detector compatibility check preserved existing counts:
  - detector active bars: `77,231`
  - detector total bars: `569,451`
  - detector episodes: `2,996`
  - detector sweep profiles: `19`
  - detector sweep rows: `76`
- Every smoke/compatibility run printed `strategy=empty trades=0`.

Follow-up:

- Superseded by the SR rejection timing review milestone above.
- The review found this boundary-rejection timing shape is not entry-ready.

Entry-readiness review gate:

- Added durable review report:
  - `docs/ENTRY_READINESS_REVIEW.md`
- Added focused tests across CSV loading, detector/default behavior, detector
  sweep, engine, indicators, SR audit helpers, SR boundary audit helpers, and
  empty strategy behavior.
- Fixed two defaulting bugs found during review:
  - zero-value `RangeDetectorConfig` now returns
    `DefaultCompressionRangeDetectorConfig()`, including `UseBollinger=true`
  - zero-value `SRBoundaryAuditConfig` now returns
    `DefaultSRBoundaryAuditConfig()`, including `DetectorActiveOnly=true`
- Review verdict:
  - the codebase passes the review gate for the next non-trading timing audit
  - first trade entries should still wait until the timing audit proves a
    closed-candle rejection signal can be identified without future bars
- Updated `memory/NEXT_CODEX_BRIEF.md` for the boundary-rejection timing audit.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest review verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./cmd/rangelab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./internal/lab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -coverprofile=/tmp/range-strategy-lab-internal.cover ./internal/lab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go tool cover -func=/tmp/range-strategy-lab-internal.cover
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection-check \
  -sr-boundary-inspect
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- `go test ./cmd/rangelab` passed.
- `internal/lab` coverage: `99.8%`.
- Remaining uncovered internal statements are justified defensive SR audit error
  propagation branches in `RunSRAudit`.
- `go test -cover ./cmd/rangelab` and `go test -cover ./...` still fail with
  the known snap-confine issue for `cmd/rangelab` coverage.
- SR boundary inspection check printed:
  - `sr_boundary_inspect events=281080 comparison_rows=192`
  - `loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z`
  - `strategy=empty trades=0`

Repository setup milestone:

- Set up tracked project memory in `memory/`.
- Added root `AGENTS.md` so future Codex sessions read and maintain memory.
- Replaced the starter Codex brief with the detector sweep/audit brief.

Detector-only milestone completed before repository setup:

- Added indicator helpers:
  - normalized ATR14
  - Donchian20 width
  - Bollinger20 width
  - optional ADX14
- Added `CompressionRangeDetector` with detector-only diagnostics.
- Added detector outputs:
  - `detector_duty_cycle.csv`
  - `detector_duty_cycle.json`
  - `range_episodes.csv`
  - `range_episodes.json`
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest detector smoke run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-smoke \
  -detector
```

Observed `full_2021_2026` detector metrics:

- Active bars: `77,231`
- Total bars: `569,451`
- Duty cycle: `13.5624%`
- Episodes: `2,996`

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
```

Result: passed.

Detector sweep/audit milestone:

- Added CLI flag `-detector-sweep`.
- Added a compact detector sweep:
  - percentile: `0.20`, `0.30`, `0.40`
  - min consecutive bars: `6`, `12`, `24`
  - Bollinger on/off
  - ADX off for the grid, plus one balanced ADX-on comparison
- Added outputs:
  - `detector_sweep.csv`
  - `detector_sweep.json`
- Marked the balanced baseline:
  - `p30_c12_bollinger_on_adx_off`
  - percentile: `0.30`
  - min consecutive bars: `12`
  - Bollinger: on
  - ADX: off
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest detector sweep run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep \
  -detector-sweep
```

Observed detector sweep facts:

- Output rows: `76` (`19` profiles x `4` splits)
- Balanced baseline full split:
  - Active bars: `77,231`
  - Total bars: `569,451`
  - Duty cycle: `13.5624%`
  - Episodes: `2,996`
- All profiles had nonzero episodes in every period split.
- Profiles that roughly fit the first-pass usability screen (`5%`-`25%` full duty, nonzero episodes in every split, no obviously unstable split duty):
  - all ADX-off profiles except `p20_c24_bollinger_on_adx_off`, `p40_c06_bollinger_on_adx_off`, and `p40_c06_bollinger_off_adx_off`
  - `p40_c12_bollinger_off_adx_off` is near the upper edge
- The balanced ADX-on comparison was too restrictive for the first-pass screen:
  - `p30_c12_bollinger_on_adx_on`
  - full duty cycle: `4.36%`

Latest detector compatibility check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector
```

Result:

- Existing detector outputs still write.
- `strategy=empty trades=0`.

Latest combined detector/sweep check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-both \
  -detector \
  -detector-sweep
```

Result:

- Existing detector outputs and detector sweep outputs write in one run.
- `strategy=empty trades=0`.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
```

Result: passed.

Research helper dependency milestone:

- Added pinned pure-Go helper modules:
  - `github.com/laclance/go-sr v1.0.0`
  - `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
  - `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`
- Added `docs/RESEARCH_HELPERS.md`.
- Updated docs to keep helper modules behind adapters and audit outputs.
- No strategy entries, exits, scoring, live code, or generated result artifacts
  were added.

Dependency add command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go get \
  github.com/laclance/go-sr@v1.0.0 \
  github.com/markcheno/go-talib@v0.0.0-20250114000313-ec55a20c902f \
  nproject.io/gitlab/libraries/talib-cdl-go@v0.0.0-20211217160304-2ed8176448cc
```

Verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

SR audit milestone:

- Added CLI flag `-sr-audit`.
- Added `go-sr` zone-mode adapter behind lab-owned audit output:
  - timeframe: `5m`
  - lookback bars: `120`
  - min strength: `2`
  - warmup bars: `138`
- Added outputs:
  - `sr_touch_audit.csv`
  - `sr_touch_audit.json`
- Included balanced detector context on each SR row:
  - `detector_profile_id=p30_c12_bollinger_on_adx_off`
  - `detector_raw_active`
  - `detector_active`
- Refreshed the next Codex brief for the next SR boundary-inspection step.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest SR audit smoke run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-audit-smoke \
  -sr-audit
```

Observed SR audit facts:

- Output rows: `569,313`
- CSV lines including header: `569,314`
- Near-support rows: `126,861`
- Near-resistance rows: `128,085`
- Result paths:
  - `results/sr-audit-smoke/sr_touch_audit.csv`
  - `results/sr-audit-smoke/sr_touch_audit.json`
  - `results/sr-audit-smoke/summary.csv`
  - `results/sr-audit-smoke/summary.json`
  - `results/sr-audit-smoke/trades.json`
- `strategy=empty trades=0`.

Latest detector compatibility check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector
```

Result:

- Existing detector outputs still write.
- Full split active bars: `77,231`
- Full split total bars: `569,451`
- Full split duty cycle: `13.56%`
- Episodes: `2,996`
- `strategy=empty trades=0`.

Latest detector sweep compatibility check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep-check \
  -detector-sweep
```

Result:

- Existing detector sweep outputs still write.
- Profiles: `19`
- Rows: `76`
- Balanced baseline active bars: `77,231`
- Balanced baseline total bars: `569,451`
- Balanced baseline duty cycle: `13.56%`
- Balanced baseline episodes: `2,996`
- `strategy=empty trades=0`.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

Implementation verification rerun:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-audit-smoke \
  -sr-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep-check \
  -detector-sweep
```

Result:

- Tests passed and `git diff --check` was clean.
- SR audit wrote:
  - `results/sr-audit-smoke/sr_touch_audit.csv`
  - `results/sr-audit-smoke/sr_touch_audit.json`
- SR audit rows: `569,313`.
- SR audit CSV lines including header: `569,314`.
- Near-support rows: `126,861`.
- Near-resistance rows: `128,085`.
- Detector compatibility still writes `results/detector-check/detector_duty_cycle.csv/json`.
- Detector full split: `77,231` active bars / `569,451` total bars, `13.56%` duty cycle, `2,996` episodes.
- Detector sweep compatibility still writes `results/detector-sweep-check/detector_sweep.csv/json`.
- Detector sweep full baseline: `77,231` active bars / `569,451` total bars, `13.56%` duty cycle, `2,996` episodes.
- Every smoke/compatibility run printed `strategy=empty trades=0`.

SR boundary-quality audit milestone:

- Added CLI flag `-sr-boundary-audit`.
- Added non-trading SR boundary quality outputs:
  - `sr_boundary_events.csv`
  - `sr_boundary_events.json`
  - `sr_boundary_quality.csv`
  - `sr_boundary_quality.json`
- Defaults:
  - horizons: `1`, `3`, `6`, `12` bars
  - `detector_active=true` rows only
  - one event per near boundary side
  - skip event/horizon pairs without enough future candles
- Metrics include favorable/adverse forward move, wick break, close break,
  reclaim-after-break, rejection, strength bucket, and distance bucket.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest SR boundary-quality smoke run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-quality \
  -sr-boundary-audit
```

Observed SR boundary-quality facts:

- Boundary event rows: `281,080`.
- Boundary event CSV lines including header: `281,081`.
- Boundary quality rows: `192`.
- Boundary quality CSV lines including header: `193`.
- Result paths:
  - `results/sr-boundary-quality/sr_boundary_events.csv`
  - `results/sr-boundary-quality/sr_boundary_events.json`
  - `results/sr-boundary-quality/sr_boundary_quality.csv`
  - `results/sr-boundary-quality/sr_boundary_quality.json`
- `strategy=empty trades=0`.

Combined SR audit/boundary check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-combined-check \
  -sr-audit \
  -sr-boundary-audit
```

Result:

- `sr_touch_audit.csv` lines including header: `569,314`.
- `sr_boundary_events.csv` lines including header: `281,081`.
- `sr_boundary_quality.csv` lines including header: `193`.
- `strategy=empty trades=0`.

Latest compatibility checks after SR boundary-quality work:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-audit-smoke \
  -sr-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep-check \
  -detector-sweep
```

Result:

- SR audit still writes `sr_touch_audit.csv/json`.
- SR audit rows: `569,313`.
- Near-support rows: `126,861`.
- Near-resistance rows: `128,085`.
- Detector full split remains `77,231` active bars / `569,451` total bars,
  `13.56%` duty cycle, `2,996` episodes.
- Detector sweep full baseline remains `77,231` active bars / `569,451` total
  bars, `13.56%` duty cycle, `2,996` episodes.
- Every compatibility run printed `strategy=empty trades=0`.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

SR boundary-inspection milestone:

- Added CLI flag `-sr-boundary-inspect`.
- Added compact non-trading candidate comparison outputs:
  - `sr_boundary_candidate_comparison.csv`
  - `sr_boundary_candidate_comparison.json`
- Grouping:
  - split
  - side
  - horizon bars
  - strength bucket
  - distance bucket
- Metrics include counts and rates for close breaks, rejections, reclaimed
  breaks, and favorable-vs-adverse cohorts for all, rejected, and reclaimed
  events.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest SR boundary-inspection run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection \
  -sr-boundary-inspect
```

Observed SR boundary-inspection facts:

- Boundary event rows inspected in memory: `281,080`.
- Candidate comparison rows: `192`.
- Candidate comparison CSV lines including header: `193`.
- Result paths:
  - `results/sr-boundary-inspection/sr_boundary_candidate_comparison.csv`
  - `results/sr-boundary-inspection/sr_boundary_candidate_comparison.json`
- Compact inspect mode did not write `sr_boundary_events.*` or
  `sr_boundary_quality.*`; only the standard empty-strategy summary/trades
  files were also written.
- `strategy=empty trades=0`.

Factual comparison outcome:

- Boundary rejection has better current evidence than false-break reclaim.
- By side/horizon across all splits:
  - rejection rates ranged from `10.52%` to `36.77%`
  - rejected-cohort favorable-minus-adverse ranged from `14.75bp` to `27.73bp`
  - rejected-cohort favorable-greater-than-adverse was `96.34%` to `98.16%`
- Reclaim-after-break was mostly a longer-horizon conditional outcome:
  - reclaim event rate reached about `20%` at `12` bars
  - reclaim given close break reached about `43%` at `12` bars
  - reclaimed-cohort favorable-minus-adverse was negative at `3`, `6`, and
    `12` bars for both support and resistance in this event definition
- The `10_20bp` distance bucket remains sparse:
  - resistance events across all splits/horizons: `208`
  - support events across all splits/horizons: `96`
- Do not treat rejection as a ready entry rule yet; it is still an ex-post
  audit label. False-break reclaim needs a later post-reclaim timing audit
  before becoming a first entry template.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

Next implementation:

- Start with an entry-readiness review before adding trades.
- Use `memory/NEXT_CODEX_BRIEF.md` as the next prompt.
- Review code correctness, lookahead safety, coverage gaps, docs, and memory
  before the first entry implementation.
- After the review is clean, consider one more compact timing audit that asks
  whether rejection can be identified on the decision candle without using
  future bars.

Entry-readiness handoff:

- Committed SR boundary inspection mode:
  - commit: `1e26695 Add SR boundary inspection mode`
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with a paste-ready prompt for:
  - whole-project review before entries
  - test coverage gap investigation
  - docs/memory readiness
  - next non-trading boundary-rejection timing audit
- Removed duplicate root `CODEX_BRIEF.md`; `memory/NEXT_CODEX_BRIEF.md` is the
  canonical next-session prompt.
- Coverage check:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./internal/lab`
  - result: `84.5%` statement coverage
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./...`
  - result: failed in `cmd/rangelab` with the known snap-confine issue, while
    `/usr/local/go/bin/go test ./cmd/rangelab` passed
