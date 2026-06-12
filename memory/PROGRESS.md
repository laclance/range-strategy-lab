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

Next implementation:

- Add detector sweep/audit mode before any trade entries.
- Use `memory/NEXT_CODEX_BRIEF.md` as the next-session prompt.
