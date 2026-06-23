# Decisions

## Durable Constraints

- This project is offline research only.
- Do not add live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- Symbol scope is BTCUSDT by default, but the explicit 2026-06-26 range-universe
  discovery scope also allows local ETHUSDT and SOLUSDT Binance USDT-M futures
  sources for offline range-first discovery only.
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
- The project is now range-first broad rather than narrowly range-only. Do not
  treat prior no-promotion reviews as permanent bans on `5m`, buy/sell-touch,
  single-candle reaction, boundary rejection, failed-break re-entry, or
  breakout-continuation ideas. The banned action is rerunning the exact failed
  template under a new name; materially reframed BTCUSDT futures range ideas
  may compete in a broad discovery audit and should move quickly to a baseline
  backtest brief if the discovery gate passes.
- The futures range candidate discovery audit has been reviewed. Only clean
  breakout continuation passed the balanced discovery gate. The next authorized
  baseline backtest is limited to the top non-duplicative `4h` up-breakout
  `h12` and `1h` all-side clean-breakout `h12` candidates. Boundary touch
  rejection, single-candle wick rejection, failed breakout re-entry, and mature
  balance persistence did not pass this discovery gate and should not be
  backtested from this milestone.
- The futures clean breakout baseline has been reviewed. The independent
  `4h` up-breakout `h12` and `1h` all-side clean-breakout `h12` candidates
  both failed after costs on Binance USDT-M futures data, with negative full
  net P&L and negative `2023_2024_oos` plus `2025_2026_recent` splits. Do not
  optimize, live-wire, paper/testnet, automatically expand to `15m`, or combine
  these clean breakout candidates into a portfolio-style stream from this
  result. A portfolio stream or `15m` comparison requires a new user-approved
  premise.
- The futures range universe discovery spec explicitly broadens this lab from
  BTCUSDT-only to a local BTC/ETH/SOL Binance USDT-M futures range-discovery
  universe. This authorizes source validation and non-trading discovery only:
  it does not authorize optimization, live wiring, paper/testnet, data
  downloads, broad symbol mining, sibling repo mutation, or importing sibling
  strategy results as evidence. Any passing surface must still earn a fixed-rule
  baseline backtest before optimization.
- The futures range universe discovery audit has been reviewed. The next
  authorized fixed-rule baseline backtest is limited to the top
  non-duplicative structured-compression surfaces:
  `4h structured_compression_expansion all h6` and
  `1h structured_compression_expansion all h12`, evaluated across local
  BTCUSDT, ETHUSDT, and SOLUSDT Binance USDT-M futures sources.
  `breakout_retest_acceptance` is secondary evidence only for now, and
  `boundary_touch_rejection`, `single_candle_wick_rejection`,
  `failed_breakout_reentry`, and `mature_balance_persistence` are not approved
  for baseline backtest from this audit. This is not optimization, live wiring,
  paper/testnet, broader symbol mining, data download, or strategy promotion.
- The futures range universe structured-compression baseline has been reviewed.
  The `4h structured_compression_expansion all h6` aggregate passed after costs
  and is authorized for bounded offline optimization/robustness only, with BTC
  weakness, stress-split fragility, and ETH/SOL dependence treated as explicit
  constraints. The `1h structured_compression_expansion all h12` surface failed
  after costs and is not approved for optimization or promotion from this
  result. No live, paper/testnet, exchange API, deployment, data download,
  broad symbol mining, grid, martingale, averaging down, or two-exchange path
  is approved.
- The futures range universe structured-compression optimization has been
  reviewed. The selected `4h` configuration is
  `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`: ETHUSDT and SOLUSDT are
  the authority symbols, while BTCUSDT is diagnostic-only and remains negative.
  This authorizes a first offline candidate strategy spec for the ETH/SOL
  universe stream only. It does not authorize BTC strategy promotion,
  additional grid search, live wiring, paper/testnet, exchange API, deployment,
  data download, broad symbol mining, martingale, averaging down, or
  two-exchange work.
- The futures range universe structured-compression strategy spec freezes the
  selected ETH/SOL authority stream for one fixed offline replay/backtest:
  closed UTC `4h`, detector `p30_c12_bollinger_on_adx_off`, first closed
  breakout within `24` bars, confirmation window `2`, max hold `12`, target
  `1.0` range width, stop buffer `0.0`, next-bar-open entry, ETHUSDT and
  SOLUSDT authority only, BTCUSDT diagnostic-only. The replay must stop for
  review on source mismatch, material result mismatch, BTC promotion, or any
  grid-style retuning.
- The futures range universe structured-compression strategy replay has been
  reviewed and passed. It authorizes a bounded offline walk-forward robustness
  pass only, using the already declared `4h` structured-compression grid for
  forward-selection checks and the frozen ETH/SOL authority replay as the
  candidate stream. It does not authorize BTCUSDT promotion, new grid
  dimensions, new symbols, live/paper/testnet, exchange API, deployment, data
  download, martingale, averaging down, or two-exchange work.
- The futures range universe structured-compression walk-forward robustness
  pass has been reviewed and is fragile. It does not authorize a candidate
  strategy package for `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`:
  only one of three folds selected the exact frozen config and passed, one
  selected same-shape ETH/SOL authority config tested worse than frozen, and
  one had no selectable training config under the multi-split `100` trade gate.
  Do not retune around this result, relax gates, promote BTCUSDT, add new grid
  dimensions, or move to live/paper/testnet/deploy from this stream without a
  materially new user-approved premise and fresh review.
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
- The futures scope pivot review paused automatic narrow BTCUSDT 5m mining and
  produced the higher-timeframe source spec. That source contract remains valid
  for `15m`, `1h`, and `4h` resamples, but the active next move is now the
  range-first broad candidate discovery audit. BTC/ETH expansion stays
  deferred until the user explicitly changes scope.
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
