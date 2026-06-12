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
- Detector diagnostics are detector-only and must not create trade signals.
- The initial balanced detector baseline is:
  - percentile: `0.30`
  - min consecutive bars: `12`
  - Bollinger: on
  - ADX: off
