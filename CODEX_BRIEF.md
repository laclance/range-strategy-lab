# Codex Brief: Detector Sweep Audit

```text
We are in a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md as needed.

Goal:
Continue building a BTCUSDT 5m offline range-strategy research system from the ground up. Do not reuse strategy/scoring/live-execution logic from the old binance-bot project.

Current state:
- CSV loader, one-position backtest engine, costs, sizing, split metrics, and tests exist.
- Strategy is still lab.EmptyStrategy and must remain no-trade for this milestone.
- Detector-only diagnostics exist:
  - ATR14 normalized
  - Donchian20 width
  - Bollinger20 width
  - optional ADX14
  - CompressionRangeDetector
  - detector_duty_cycle.csv/json
  - range_episodes.csv/json
- Latest smoke run on ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv:
  - full_2021_2026 duty cycle: 13.5624%
  - full active bars: 77,231 / 569,451
  - full episodes: 2,996
  - strategy=empty trades=0

Non-negotiables:
- Offline research only.
- No live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- BTCUSDT only.
- 5m candles first.
- Confirmed closed-candle decisions only.
- Enter on next bar open when entries are eventually added.
- One open position max.
- Stop-first ambiguity.
- Keep every result explainable and reproducible.

Next implementation milestone:
Add detector sweep/audit mode before adding trade entries.

Implement:
1. CLI flag -detector-sweep.
2. Sweep a compact grid:
   - percentile: 0.20, 0.30, 0.40
   - min consecutive bars: 6, 12, 24
   - Bollinger on/off
   - ADX off by default, plus one ADX-on comparison using the balanced profile.
3. Write under the selected -out-dir:
   - detector_sweep.csv
   - detector_sweep.json
4. Include per-split metrics already used by detector diagnostics:
   - active bars
   - total bars
   - duty cycle
   - episodes
   - average episode length
   - median episode length
   - longest episode length
5. Mark or clearly identify the balanced baseline:
   - percentile 0.30
   - min consecutive bars 12
   - Bollinger on
   - ADX off

Acceptance criteria:
- /usr/local/go/bin/go test ./... passes with GOCACHE=/tmp/range-strategy-lab-go-build if needed.
- Smoke command writes detector_sweep.csv/json on the BTCUSDT 5m CSV.
- Existing -detector output still works.
- Strategy remains empty and trades remain 0.
- Summarize which detector profiles look usable: prefer nonzero episodes in every split, full duty cycle roughly 5%-25%, and no wildly unstable split behavior.
```
