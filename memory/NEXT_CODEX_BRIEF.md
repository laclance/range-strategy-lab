# Next Codex Brief: Support/Resistance Audit Mode

```text
We are in a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md as needed, especially docs/RESEARCH_HELPERS.md.

Goal:
Continue building a BTCUSDT 5m offline range-strategy research system from the ground up. Do not reuse strategy/scoring/live-execution logic from the old binance-bot project.

Current state:
- CSV loader, one-position backtest engine, costs, sizing, split metrics, detector diagnostics, and detector sweep mode exist.
- Strategy is still lab.EmptyStrategy and currently produces zero trades.
- Detector sweep/audit mode exists:
  - CLI flag: -detector-sweep
  - outputs: detector_sweep.csv and detector_sweep.json
  - balanced baseline: p30_c12_bollinger_on_adx_off
  - balanced baseline full split: 77,231 active bars / 569,451 total bars, 13.5624% duty cycle, 2,996 episodes
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
Add a go-sr support/resistance audit mode before implementing trades.

Recommended default:
- Add a CLI flag such as -sr-audit.
- Use github.com/laclance/go-sr v1.0.0 through a small adapter from lab.Candle to sr.Candle.
- Start with 5m zone mode and an explainable fixed lookback.
- Write generated outputs under -out-dir, for example sr_zones.csv/json or sr_touch_audit.csv/json.
- Include enough fields to inspect boundary quality: nearest support/resistance, distance to close, near flags, zone strength/score, and split/time context.
- Do not add entry/exit trading rules in this audit milestone.

Acceptance criteria:
- /usr/local/go/bin/go test ./... passes with GOCACHE=/tmp/range-strategy-lab-go-build if needed.
- Smoke run on ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv writes SR audit CSV/JSON under results/.
- Existing detector and detector-sweep outputs still work.
- Strategy remains empty and trades remain 0.
- Update memory/PROGRESS.md with commands, result paths, and short factual outcomes.
```
