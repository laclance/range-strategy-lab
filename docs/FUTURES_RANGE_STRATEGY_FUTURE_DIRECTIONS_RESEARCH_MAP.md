# Futures Range Strategy Future Directions Research Map

Date: 2026-06-27

## Verdict

Stop state:
`future_range_strategy_directions_map_ready_for_range_state_spec`.

This is a documentation-only research map. It records the next range-strategy
research direction set requested by the user after the BTCUSDT range-context
triage audit failed to produce a gated strategy premise.

It does not add a strategy, optimizer, replay, walk-forward run, CLI flag,
generated result directory, artifact writer, data download, paper/testnet/live
path, exchange API, credentials, deploy file, martingale, averaging down, broad
symbol mining, or two-exchange logic.

## Central Pivot

The next arc should move away from:

```text
range event predicts direction
```

and toward:

```text
range context + volatility state + trend state + impulse/liquidity state
  -> closed-candle state classification
  -> route to no-trade, rotation, continuation, or later strategy research
```

The project should treat a range as a context object, not as a standalone entry
edge. A range may be tradable, toxic, continuation-biased, or useful only as a
filter. The next implementation-ready step is therefore a non-trading state
audit, not another entry grammar.

## Current Exclusion Evidence

The following reviewed families remain exclusion evidence in their reviewed
forms:

- structured compression, including the fragile ETH/SOL authority stream;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline;
- impulse absorption after abnormal candles;
- higher-timeframe nested range rotation;
- `range_occupancy_rotation_v1`;
- range quality, UTC session, and failure-mode triage cohorts by themselves.

Future work may reuse infrastructure from these branches, such as source guards,
closed UTC resampling, range episode helpers, event labels, split metrics, and
artifact patterns. It must not retune, rename, or relax gates around the failed
entry premises.

## Direction Ladder

| Priority | Direction | Spec | Implementation Status |
| ---: | --- | --- | --- |
| 1 | Range state construction loop | `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md` | Next implementation-ready documentation spec. Non-trading audit only. |
| 2 | Range context router | `docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md` | Parked until the range-state audit identifies usable states. |
| 3 | Volatility-aware exits | `docs/FUTURES_VOLATILITY_AWARE_EXIT_MODEL_SPEC.md` | Parked until a materially new entry has gross edge before costs. |
| 4 | BTC regime plus ETH/SOL range context | `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md` | Parked scope expansion. Non-trading audit first if approved. |
| 5 | Spread-range / pair-range strategy | `docs/FUTURES_SPREAD_RANGE_STRATEGY_SPEC.md` | Parked engine/source expansion. Non-trading spread audit first. |
| 6 | Derivatives context source expansion | `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md` | Parked source-expansion spec. Market-data only if later approved. |

Only the first direction is the next Codex implementation target. The other
specs are deliberately parked so future sessions do not improvise scope or
accidentally reopen failed families.

## Direction 1: Range State Construction Loop

Purpose: build a closed-candle, non-trading audit that combines range geometry,
volatility, trend, impulse, and OHLCV liquidity proxies before attempting any
entry.

Default scope:

- BTCUSDT only;
- Binance USDT-M futures `5m` parent source;
- closed UTC `15m`, `1h`, and `4h` resamples;
- no ETH/SOL authority;
- no funding, open interest, order book, or external market-data source;
- no P&L, entries, exits, optimizer, replay, or walk-forward.

The audit should answer whether any pre-entry state family separates usable
range context from toxic range context across the existing period splits.

## Direction 2: Range Context Router

Purpose: if Direction 1 finds stable usable states, convert those states into a
router:

```text
no_trade | tradable_rotation | trend_continuation | diagnostic_only
```

The router is not an entry strategy. It is a gate that decides which kind of
later strategy spec may be written. It must not choose trade direction from
future labels, and it must not transform failed occupancy rotation,
structured-compression, breakout-retest, or impulse-absorption rules into new
names.

## Direction 3: Volatility-Aware Exit Model

Purpose: once a materially new entry template shows gross edge before costs,
compare exit behavior under volatility-aware rules.

This is intentionally not first. An exit model cannot rescue a missing entry
edge. The repo workflow still requires gross evidence before cost-sensitive
variants.

## Direction 4: BTC Regime Plus ETH/SOL Range Context

Purpose: test whether BTCUSDT is more useful as market-regime context than as an
authority trade symbol, while ETHUSDT and SOLUSDT provide authority rows only if
explicitly approved.

This direction is motivated by prior BTC weakness and ETH/SOL strength in the
structured-compression branch, but it must not reopen structured compression.
The first step would be a non-trading context audit, not a strategy.

## Direction 5: Spread-Range / Pair-Range Strategy

Purpose: define ranges on relative price or residual spreads instead of single
instrument price boxes.

Examples:

- `ETHUSDT / BTCUSDT` relative price;
- `SOLUSDT / BTCUSDT` relative price;
- `SOLUSDT / ETHUSDT` relative price;
- beta-adjusted ETH or SOL residual versus BTC.

This requires multi-series synchronization and, before trading, multi-leg P&L,
fee, slippage, notional, and margin accounting. It is therefore parked until a
separate engine/source expansion is approved.

## Direction 6: Derivatives Context Source Expansion

Purpose: evaluate whether range states become more useful when combined with
funding, basis, open interest, long/short, taker-flow, or order-book context.

This is market-data source expansion only. It does not permit exchange order
APIs, keys, live trading, deployment, or data downloads without an explicit later
brief.

## Ordering Rule

The next automatic implementation target is only:

```text
FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md
```

A later direction may start only after one of these happens:

1. the range-state audit passes and explicitly authorizes a router or strategy
   spec;
2. the user explicitly chooses one parked direction and accepts its source or
   engine implications;
3. the user asks for a review-only inventory update without implementation.

## Stop States

Research-map stop states:

- `future_range_strategy_directions_map_ready_for_range_state_spec`;
- `future_range_strategy_directions_map_needs_user_scope_choice`;
- `future_range_strategy_directions_map_rejected_closed_family_reslice`.

Downstream stop states are defined in the individual specs.