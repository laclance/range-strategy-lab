# Next Codex Brief: Inspect SR Boundary Quality

```text
We are in a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md as needed, especially docs/RESEARCH_HELPERS.md and docs/STRATEGY_WORKFLOW.md.

Goal:
Continue building a BTCUSDT 5m offline range-strategy research system from the ground up. Do not reuse strategy/scoring/live-execution logic from the old binance-bot project.

Current state:
- CSV loader, one-position backtest engine, costs, sizing, split metrics, detector diagnostics, detector sweep mode, SR audit mode, and SR boundary-quality audit mode exist.
- Strategy is still lab.EmptyStrategy and currently produces zero trades.
- Detector sweep/audit mode exists:
  - CLI flag: -detector-sweep
  - outputs: detector_sweep.csv and detector_sweep.json
  - balanced baseline: p30_c12_bollinger_on_adx_off
  - balanced baseline full split: 77,231 active bars / 569,451 total bars, 13.5624% duty cycle, 2,996 episodes
- SR audit mode exists:
  - CLI flag: -sr-audit
  - outputs: sr_touch_audit.csv and sr_touch_audit.json
  - default: 5m zone mode, 120-bar lookback, min strength 2, warmup 138 bars
  - latest smoke output: results/sr-audit-smoke/
  - rows: 569,313
  - near-support rows: 126,861
  - near-resistance rows: 128,085
- SR boundary-quality audit mode exists:
  - CLI flag: -sr-boundary-audit
  - outputs: sr_boundary_events.csv/json and sr_boundary_quality.csv/json
  - default horizons: 1, 3, 6, 12 bars
  - default filter: detector_active=true
  - latest smoke output: results/sr-boundary-quality/
  - boundary event rows: 281,080
  - boundary quality rows: 192
- Pinned research helper modules exist in go.mod:
  - github.com/laclance/go-sr v1.0.0
  - github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f
  - nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc

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
- Helper modules may provide feature extraction/audit outputs only; strategy hypotheses, entries, exits, scoring, sizing, and backtest behavior stay inside this lab.

Next recommended milestone:
Inspect SR boundary quality before adding trade entries.

Recommended approach:
- Read results/sr-boundary-quality/sr_boundary_quality.csv.
- Compare support vs resistance outcomes by split, horizon, strength bucket, and distance bucket.
- Identify whether boundary rejection or false-break reclaim has better evidence.
- If more audit code is needed, add compact non-trading outputs under -out-dir; do not add entry/exit trading rules unless explicitly asked.
- Choose at most one first entry template to test later.

Acceptance criteria:
- /usr/local/go/bin/go test ./... passes with GOCACHE=/tmp/range-strategy-lab-go-build if code changes are made.
- Any generated inspection outputs stay under results/.
- Strategy remains empty and trades remain 0 unless the user explicitly asks to add the first trading strategy.
- Update memory/PROGRESS.md with commands, result paths, and short factual outcomes.
```
