# Decisions

## Durable Constraints

- This project is offline research only.
- Do not add live orders, exchange API keys, private account endpoints, deploy
  scripts, martingale, averaging down, or two-exchange execution.
- The active research market is Binance USDT-M futures, not Binance spot.
- BTCUSDT `5m` remains the default source and default CLI identity unless a
  reviewed scope brief explicitly expands it.
- Local ETHUSDT and SOLUSDT Binance USDT-M futures files may be used only where
  an approved offline range-universe or context brief explicitly allows them.
- Use confirmed closed-candle decisions only.
- When entries are approved by a future spec, enter on the next bar open.
- Keep one open position max unless a later explicit engine spec changes that.
- Use stop-first ambiguity.
- Keep every result explainable and reproducible.
- Do not reuse strategy, scoring, order-management, live-execution, credential,
  deploy, or portfolio coordinator logic from the old `binance-bot` project.
- Generated outputs belong under `results/`, which remains ignored by Git.
- Project memory is tracked under `memory/` and should stay compact.
- `memory/NEXT_CODEX_BRIEF.md` is the canonical next-session prompt; do not keep
  a duplicate root `CODEX_BRIEF.md`.

## Source And Verification Decisions

- Candle data source is part of the experiment definition. Record CSV path,
  market type, date coverage, row count, gap count, duplicate count, zero-volume
  count, comparison-only status, and validation status for any data-dependent
  verdict.
- The default accepted source is
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`, Binance USDT-M
  futures `BTCUSDT` `5m`, `573,984` loaded candles from
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`, `gap_count=0`,
  `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`,
  `validation_status=accepted`.
- Spot-based results are historical context only unless a futures impact review
  explicitly revalidates them.
- Source validation must reject gaps, duplicates, irregular `5m` cadence,
  non-positive OHLC prices, non-finite values, negative volume, invalid high/low
  containment, wrong symbol identity, wrong interval identity, and comparison-only
  sources where promotion is required. Zero-volume closed candles may be counted
  and allowed.
- After completing a brief or milestone, Codex should run closeout checks and
  commit completed repo changes unless the user explicitly says not to commit.

## Current Research Decision

- The next implementation-ready research direction is the futures range-state
  construction loop in `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md`.
- That spec authorizes only a non-trading BTCUSDT state audit from the accepted
  `5m` futures source and closed UTC `15m`, `1h`, and `4h` resamples.
- The audit must combine range geometry, volatility state, trend state, impulse
  state, and OHLCV liquidity/participation proxies before any new entry is
  specified.
- The audit must produce zero trades and keep common outputs zero-trade
  compatible.
- Passing the audit may authorize only a later documentation-only router or
  strategy-premise spec. It does not authorize entries, exits, P&L backtests,
  optimizer grids, fixed replay, walk-forward, packaging, source expansion,
  symbol expansion, live-adjacent work, or closed-family retuning.
- The intended stop state for the current docs milestone is
  `range_state_construction_loop_spec_ready_for_audit_implementation`.

## Parked Future Direction Decisions

The following specs are parked and not implementation-ready from current state:

- `docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md`: may start only if the range-state
  audit passes and explicitly authorizes router work.
- `docs/FUTURES_VOLATILITY_AWARE_EXIT_MODEL_SPEC.md`: may start only after a
  materially new entry template first shows gross edge before costs.
- `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md`: may start only with user
  scope approval or a later range-state review recommendation; it is a
  non-trading context audit first, not ETH/SOL strategy promotion.
- `docs/FUTURES_SPREAD_RANGE_STRATEGY_SPEC.md`: may start only with explicit
  engine/source approval; spread trading requires a separate multi-leg engine
  spec before any P&L strategy work.
- `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md`: may start only
  with explicit user source approval; it is market-data context only and does
  not permit API keys, private endpoints, live/paper/testnet, or exchange order
  paths.

## Exclusion Decisions

Reviewed failed or fragile families must not be retuned, renamed, gate-relaxed,
or promoted from their reviewed forms:

- structured compression, including the fragile ETH/SOL authority stream;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline;
- impulse absorption after abnormal OHLCV candles;
- higher-timeframe nested range rotation;
- `range_occupancy_rotation_v1`;
- range quality, UTC session, and failure-mode triage cohorts by themselves;
- legacy spot-only SR rejection, confirmation, false-break, and compression
  promotion evidence.

These branches remain useful as exclusion evidence and infrastructure only.
Reusable infrastructure includes source guards, closed UTC resampling,
closed-candle semantics, event labels as labels only, split metrics, artifact
patterns, and detector/range episode helpers as feature extraction only.

## Historical Reviewed Decisions

- The futures range-context triage audit passed source/resampling but failed to
  find a gated strategy premise. Do not create a strategy, baseline, optimizer,
  replay, walk-forward, package, retune, source expansion, symbol expansion, or
  live-adjacent path from that audit.
- The futures range-first occupancy rotation V1 optimizer evaluated the declared
  `1,152` row grid, the fixed baseline lost after costs, and `0` grid rows
  passed. Do not create a fixed replay or retune the grammar.
- The structured-compression universe stream selected ETH/SOL authority with BTC
  diagnostic-only, but walk-forward robustness was fragile. Do not package,
  retune, promote BTC, or broaden the grid from that result.
- The breakout-retest/acceptance baseline failed after costs on BTCUSDT,
  ETHUSDT, and SOLUSDT. Do not optimize or robustness-review that branch.
- The higher-timeframe nested range rotation audit produced only `3` valid
  events across the full BTCUSDT sample. Do not build a baseline from it.
- The impulse absorption audit found continuation-dominant behavior after
  abnormal OHLCV impulse candles. Do not convert that surface into an entry.

## Helper Module Decisions

External helper modules may be used for feature extraction and audit outputs
only. Strategy hypotheses, entries, exits, scoring, sizing, and backtest behavior
stay inside this lab.

Pinned research helper modules:

- `github.com/laclance/go-sr v1.0.0`
- `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
- `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`