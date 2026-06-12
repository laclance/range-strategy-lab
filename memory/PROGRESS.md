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

Next implementation:

- Add a `go-sr` support/resistance audit mode before adding trade entries.
- Use the balanced detector profile as the default context unless a sweep result
  gives a better reason.
- After SR audit outputs are inspectable, choose one detector profile and one
  simple first entry template.
- Use `memory/NEXT_CODEX_BRIEF.md` as the next-session prompt.
