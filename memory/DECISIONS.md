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
- The current `p30_c12_bollinger_on_adx_off` detector is not approved as
  context for future entry hypotheses until detector/context refinement is
  reviewed; its durability weakness repeats across splits and it quick
  invalidates too often after episode end.
- The current `DefaultDetectorSweepProfiles` detector durability sweep has
  been reviewed, and no profile is approved as future entry context. The ADX
  comparison profile `p30_c12_bollinger_on_adx_on` is diagnostic only, not a
  promoted detector.
- The detector context refinement audit has been reviewed. The delayed
  `hold_3_inside` and `hold_6_inside` context rules are the leading context
  refinement: they materially and split-stably reduce quick invalidation and
  trend leakage with adequate candidates, and the hold condition is
  closed-candle knowable at the decision candle. They are still not approved as
  entry context, because the gain is a heavy survivorship/conditioning effect,
  residual `12` bar trend leakage stays material, and the `label_*` fields are
  regime-durability outcomes, not P&L. No profile or context rule is promoted;
  keep `lab.EmptyStrategy`.
- External helper modules may be used for feature extraction and audit outputs
  only; strategy hypotheses, entries, exits, scoring, sizing, and backtest
  behavior stay inside this lab.
- Pinned research helper modules:
  - `github.com/laclance/go-sr v1.0.0`
  - `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
  - `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`
