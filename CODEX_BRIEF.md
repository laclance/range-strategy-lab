# Next Codex Brief: First Entry Template Prep

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
- Strategy is still lab.EmptyStrategy and currently produces zero trades.
- Detector-only diagnostics exist:
  - ATR14 normalized
  - Donchian20 width
  - Bollinger20 width
  - optional ADX14
  - CompressionRangeDetector
  - detector_duty_cycle.csv/json
  - range_episodes.csv/json
- Detector sweep/audit mode exists:
  - CLI flag: -detector-sweep
  - outputs: detector_sweep.csv and detector_sweep.json
  - grid: percentile 0.20/0.30/0.40, min consecutive bars 6/12/24, Bollinger on/off
  - ADX is off for the grid, plus one balanced ADX-on comparison
- Balanced baseline:
  - profile_id: p30_c12_bollinger_on_adx_off
  - percentile: 0.30
  - min consecutive bars: 12
  - Bollinger: on
  - ADX: off

Latest detector sweep run:

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep \
  -detector-sweep

Observed detector sweep facts:
- Output rows: 76 (19 profiles x 4 splits)
- Balanced baseline full split: 77,231 active bars / 569,451 total bars, 13.5624% duty cycle, 2,996 episodes
- All profiles had nonzero episodes in every period split.
- First-pass usable profiles are the ADX-off profiles with full duty roughly 5%-25%, except p20_c24_bollinger_on_adx_off is too sparse and p40_c06_bollinger_on_adx_off / p40_c06_bollinger_off_adx_off are too broad.
- p40_c12_bollinger_off_adx_off is near the upper edge.
- The balanced ADX-on comparison is too restrictive for the first-pass screen: p30_c12_bollinger_on_adx_on full duty cycle 4.36%.
- strategy=empty trades=0.

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

Next recommended milestone:
Choose one detector profile and one simple first entry template before implementing trades.

Recommended default:
- Use the balanced baseline p30_c12_bollinger_on_adx_off unless there is a clear reason to test a broader/narrower detector first.
- Start with exactly one explainable entry template.
- Do not combine multiple entry styles, grids, averaging down, martingale, or multi-exchange ideas.

Acceptance criteria for the next coding milestone:
- Strategy changes are isolated and still use confirmed closed candles with next-bar-open entries.
- Only one entry template is added.
- Results include the existing split metrics and preserve zero-lookahead behavior.
- /usr/local/go/bin/go test ./... passes with GOCACHE=/tmp/range-strategy-lab-go-build if needed.
- Generated outputs stay under results/.
```
