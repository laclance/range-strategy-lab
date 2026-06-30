# Decisions

## Durable Constraints

- This project is offline Binance USDT-M futures research only.
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
- For workbench/optimization outputs, verification must preserve prior run
  artifacts. Use immutable run directories under
  `results/range-optimization-workbench-v1/runs/<run_id>/`; reruns must create a
  new run ID instead of deleting the canonical workbench results parent.

## Current Research Decision

- The three backtest-first fixed baselines from
  `docs/BACKTEST_FIRST_CANDIDATE_PACKET.md` failed as no usable strategies:
  `btc_5m_rolling_value_area_reversion_v1`,
  `btc_15m_previous_day_range_reversion_v1`, and
  `btc_15m_range_edge_exhaustion_fade_v1`.
- These failures are closed in reviewed form. Do not rescue them with alternate
  thresholds, alternate windows, side selection, added filters, derivatives
  context, replay, walk-forward, or optimizer grids as if they were the same
  fixed-baseline test.
- The new docs-only workbench spec in
  `docs/RANGE_OPTIMIZATION_WORKBENCH_SPEC.md` allows combining and optimizing
  related range-family components only as controlled offline discovery.
- Stop state:
  `range_optimization_workbench_spec_ready_for_implementation_approval`.
- The workbench may search bounded combinations of range context, volatility,
  trend/impulse, activity, and session features, but every trial must be emitted
  and preserved. No failed/ugly trials may be hidden or deleted.
- Optimizer output alone cannot authorize paper/shadow/live, promotion,
  production changes, exchange API work, credentials, deployment, martingale,
  averaging down, or two-exchange logic.
- At most one workbench candidate may be selected, and selection may only stop at
  `range_optimization_workbench_candidate_selected_needs_fixed_validation`; a
  later locked fixed-validation lane is required before stronger claims.
- If no robust candidate passes the declared filters, stop at
  `range_optimization_workbench_failed_no_candidate` or
  `range_optimization_workbench_rejected_overfit_risk`.

## Derivatives And Context Decisions

- Durable Binance public Data Vision USDT-M futures `markPriceKlines`,
  `indexPriceKlines`, and optional `premiumIndexKlines` `5m` source files for
  `BTCUSDT`/`ETHUSDT`/`SOLUSDT` exist under
  `../binance-bot/data/derivatives/` (`729` checksum-verified raw zips, `9`
  normalized CSVs, `5` manifests). These files are outside this repo and are not
  tracked by Git.
- Derivatives source usage must apply the conservative one-native-interval lag
  (`source_close_time + 5m <= decision_candle_close_time`); no forward fill,
  interpolation, or nearest-future joins.
- The canonical veto `btc_15m_basis_discount_no_trade_veto_v1` is preserved as a
  future veto candidate only. It may only annotate an already-approved candidate
  stream as skipped or retained; it may not create entries, change entry logic,
  choose side, act as an exit, rank trades, score P&L, replay, walk forward,
  optimize, promote, expand sources/symbols, or reopen closed families.

## Historical Decisions

- Detailed historical decisions remain in focused `docs/` reviews and git
  history. This always-read memory intentionally carries only the compact current
  decision state and durable constraints.
