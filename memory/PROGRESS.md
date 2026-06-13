# Progress

## 2026-06-13

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
- Refreshed `CODEX_BRIEF.md` and `memory/NEXT_CODEX_BRIEF.md` for the next
  SR boundary-inspection step.
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

Next implementation:

- Inspect `results/sr-boundary-quality/sr_boundary_quality.csv`.
- Compare support vs resistance outcomes by split, horizon, strength bucket,
  and distance bucket.
- Choose at most one first entry template if the audit supports it, likely
  false-break reclaim or boundary rejection.
- Keep the next milestone research-first unless explicitly adding trades.
