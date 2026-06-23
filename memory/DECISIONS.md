# Decisions

## Durable Constraints

- This project is offline research only.
- Do not add live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- BTCUSDT only until project scope explicitly changes.
- Use 5m candles first.
- The active research market is Binance USDT-M futures, not Binance spot.
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
- Candle data source is part of the experiment definition. Record CSV path,
  market type, date coverage, and row counts for any data-dependent verdict.
  Spot-based results are not a promotion basis for the futures trading target
  unless a futures impact review explicitly revalidates them.
- `cmd/rangelab` enforces the active source contract before running audits or
  backtests. The default CSV is
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`; non-default CSV
  paths must pass `-source-product`; spot paths require both
  `-source-product binance-spot` and `-allow-spot-comparison`; every accepted
  run writes `source_manifest.json`. The source guard rejects gaps, duplicates,
  irregular 5m cadence, non-positive OHLC prices, non-finite values, negative
  volume, and invalid high/low containment; zero-volume closed candles are
  allowed and counted in the manifest.
- The canonical next-session prompt is `memory/NEXT_CODEX_BRIEF.md`; do not
  keep a duplicate root `CODEX_BRIEF.md`.
- Keep tracked memory context-budgeted: always-read memory files are a compact
  working index, not a full transcript. Treat size targets for all always-read
  memory files as soft `300-350` line judgment bands, not hard triggers;
  compact or split memory once an always-read file starts feeling bulky or
  repetitive. Handoff briefs should name only task-relevant docs, not require a
  blanket `docs/*.md` read.
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
- The hold-inside directional edge audit has been reviewed. The leading
  `hold_3_inside`/`hold_6_inside` context does not show a split-stable
  directional edge toward the frozen range high or low. No all-bucket row
  passes the review gate of positive worst-split favorable-minus-adverse,
  worst-split favorable-greater-than-adverse above `50%`, and adequate
  candidate counts; no decision-close-position bucket reaches `100` candidates
  in every period split. Do not promote `paper_side=toward_high` or
  `paper_side=toward_low` into entry context; keep `lab.EmptyStrategy`.
- The hold-inside midline transition audit has been reviewed. Broad
  `hold_3_inside`/`hold_6_inside` rows show split-stable midline touch and
  close-across behavior by `12` bars, but the labels are not entry context or
  strategy inputs. Treat the midline as a non-trading observation point for a
  follow-up reindexed midline-event audit; do not promote current midline
  transition labels into entries, exits, scoring, sizing, or strategy logic.
- The futures data impact review revalidated only `hold_3_inside` + first
  `mid_touch` within `12` bars, with event close in the frozen range `mid_50`
  bucket, for a first minimal offline entry prototype on Binance USDT-M futures
  data. This is not strategy promotion or live approval; keep
  `lab.EmptyStrategy` until that prototype exists and is reviewed. The old
  spot-based approval is historical comparison only. `hold_6_inside`,
  `mid_close_across`, side-specific cohorts, and `hold_3_inside_mid_50` remain
  diagnostic after the futures review.
- The minimal futures midline touch prototype has been reviewed. The exact
  close-back template for `hold_3_inside` + first `mid_touch` within `12` bars
  + event close-position bucket `mid_50`, with same-side boundary stop,
  opposite-boundary target, next-bar-open entry, and `6` bar time stop, failed
  on Binance USDT-M futures data. Do not promote, parameter-tune, broaden, or
  live-wire this hold-inside/midline entry family without a materially new
  non-trading premise and fresh review.
- The futures hypothesis pivot inventory has been reviewed. Closed or
  diagnostic families should not be retried, narrowed, or converted into
  entries without a materially new futures hypothesis or data premise. Legacy
  spot-only evidence cannot promote futures work. The reviewed families are now
  an exclusion map plus reusable infrastructure until a new premise is supplied.
- The futures impulse absorption audit has been reviewed. Abnormal OHLCV
  impulse candles on Binance USDT-M futures data are continuation-dominant
  rather than midpoint-reclaim-dominant across every period split and tested
  horizon. Do not convert this impulse absorption surface into an entry
  prototype, retune, paper/testnet/live path, or strategy replacement without a
  materially new futures hypothesis or data premise.
- The futures scope pivot remains range-strategy-only unless the user
  explicitly changes the project objective. Higher-timeframe, multi-symbol, or
  sibling-repo context may be considered only as range-source, range-premise,
  process, or exclusion evidence. Non-range trend/volatility paths, broad
  multi-pair mining, and cross-exchange execution are not authorized by a
  range-only pivot review.
- The futures scope pivot review paused further BTCUSDT 5m range mining. The
  next authorized range-only move is a docs/source-spec task for BTCUSDT
  higher-timeframe futures bars derived from the accepted 5m source, with
  `15m`, `1h`, and `4h` as candidate intervals, before any audit or prototype.
  BTC/ETH expansion stays deferred until that higher-timeframe BTCUSDT
  source/premise review is complete or the user explicitly changes scope.
- Higher-timeframe BTCUSDT range work must derive `15m`, `1h`, and `4h`
  candidate bars from the accepted Binance USDT-M futures BTCUSDT 5m parent
  source by closed UTC resampling until an explicit source or scope change is
  reviewed. No higher-timeframe audit may start until generated coverage and
  row counts are documented and a materially different range premise with a
  closed-candle event and falsification rule is explicit.
- External helper modules may be used for feature extraction and audit outputs
  only; strategy hypotheses, entries, exits, scoring, sizing, and backtest
  behavior stay inside this lab.
- Pinned research helper modules:
  - `github.com/laclance/go-sr v1.0.0`
  - `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
  - `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`
