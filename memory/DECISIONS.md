# Decisions

## Durable Constraints

- This project is offline research only.
- Do not add live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- BTCUSDT only until project scope explicitly changes.
- Use 5m candles first.
- Use confirmed closed-candle decisions only.
- When entries are eventually added, enter on the next bar open.
- Keep one open position max.
- Use stop-first ambiguity.
- Keep every result explainable and reproducible.
- Do not reuse strategy, scoring, or live-execution logic from the old `binance-bot` project.

## Implementation Decisions

- Generated outputs belong under `results/`, which remains ignored by Git.
- Project memory is tracked under `memory/`.
- Future Codex sessions should read `AGENTS.md` and `memory/` before nontrivial work.
- The canonical next-session prompt is `memory/NEXT_CODEX_BRIEF.md`; do not
  keep a duplicate root `CODEX_BRIEF.md`.
- After completing a brief or milestone, Codex should automatically run the
  closeout checks and commit the completed repo changes unless the user
  explicitly says not to commit.
- Detector diagnostics are detector-only and must not create trade signals.
- The initial balanced detector baseline is:
  - percentile: `0.30`
  - min consecutive bars: `12`
  - Bollinger: on
  - ADX: off
- External helper modules may be used for feature extraction and audit outputs
  only; strategy hypotheses, entries, exits, scoring, sizing, and backtest
  behavior stay inside this lab.
- Pinned research helper modules:
  - `github.com/laclance/go-sr v1.0.0`
  - `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
  - `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`
