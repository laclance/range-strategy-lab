# Backtest-First Session Opening-Range Expansion Lane Selection

Date: 2026-07-01

## Verdict

Stop state:

```text
session_opening_range_expansion_lane_selected_for_candidate_packet
```

Selected next research lane:

```text
session-based opening-range expansion
```

Next bounded approval gate:

```text
docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_CANDIDATE_PACKET.md
```

This is a docs-only lane selection after the failed
`btc_15m_trend_pullback_continuation_v1` backtest. It does not authorize Go
code, CLI flags, generated results, backtests, optimizers, source expansion,
derivatives-veto interaction, exchange APIs, credentials, deployment,
paper/testnet/live flow, martingale, averaging down, two-exchange logic, or
promotion.

## Current Stop State

The approved trend-pullback continuation baseline failed at:

```text
btc_15m_trend_pullback_continuation_backtest_failed_no_usable_strategy
```

It produced `3,816` executed trades on the accepted BTCUSDT Binance USDT-M
futures source, but failed gross edge, net edge, profit factor, and drawdown
gates:

| Metric | Full `2021_2026` Value |
| --- | ---: |
| Gross P&L | `-123.427844` |
| Net P&L | `-956.566589` |
| Profit factor | `0.837958` |
| Max drawdown | `0.958438` |

The fixed trend-pullback candidate is closed in reviewed form. Do not rescue it
with alternate EMA lengths, slope lookbacks, pullback windows, EMA-band
definitions, continuation triggers, stop buffers, target R values, time stops,
side selection, session filters, volume filters, volatility filters,
derivatives-veto interaction, source expansion, replay, walk-forward, or
optimizer grids.

The current range-reversion / midpoint / edge-fade / previous-day range /
bounded range-optimization path is also closed in reviewed form.

## Candidate Direction Review

| Direction | Decision | Reason |
| --- | --- | --- |
| Volatility expansion / generic breakout continuation | Defer | It remains interesting, but it sits close to closed clean-breakout and post-compression expansion evidence. A future version would need especially strict separation from those failed families. |
| Session-based opening-range expansion | Select | It can be expressed as one deterministic time-box expansion baseline using only accepted BTCUSDT futures OHLCV. It is not an EMA pullback retry and does not fade range edges or target midpoints. |
| Liquidity sweep plus reclaim | Defer | It is conceptually different, but OHLCV-only sweep proxies create too many ambiguous knobs for the next simplest fixed baseline. It also risks drifting into closed false-break or SR timing work. |
| Cross-asset or regime filter before entries | Reject for now | It is not a standalone entry lane and would risk becoming a post-hoc filter for failed candidates. It also conflicts with the current no source-expansion boundary. |

## Why This Lane Is Materially Different

Session-based opening-range expansion uses a predeclared time box as the event
reference and tests continuation away from that box after a closed-candle break.
The future candidate packet must select one fixed UTC session anchor, one opening
range length, one decision timeframe, one continuation trigger, and one fixed
exit model before any code or backtest exists.

This differs from the closed range-reversion families because it must not:

- fade outer range zones;
- target a midpoint, VWAP, value area, or prior-day mean;
- use previous-day range reversion;
- use range-edge exhaustion as the entry reason; or
- search session anchors after seeing results.

This differs from the failed trend-pullback baseline because it must not:

- define trend with EMA stacks or EMA slopes;
- require a pullback into an EMA band;
- use the failed continuation-close trigger; or
- reuse the failed trend-pullback stop, target, time stop, or side-selection
  variants as rescue knobs.

The selected lane is still BTCUSDT-only, Binance USDT-M futures-only,
closed-candle-only, offline, and local-source-only.

## Next Candidate Packet Requirements

If the operator approves the next gate, the candidate packet may add exactly one
docs-only artifact:

```text
docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_CANDIDATE_PACKET.md
```

That packet must define exactly one fixed baseline, likely named:

```text
btc_15m_session_opening_range_expansion_v1
```

The packet must lock:

- accepted BTCUSDT Binance USDT-M futures source;
- one exact decision timeframe;
- one fixed UTC session anchor;
- one fixed opening-range length;
- one closed-candle expansion or acceptance trigger;
- next-bar-open entry;
- one structural or ATR-based invalidation stop;
- one fixed target and time stop;
- fixed sizing, fee, and slippage assumptions;
- side-separated reporting if both sides are tested;
- split, trade-count, gross, net, profit-factor, drawdown, leakage, and
  no-contamination gates; and
- explicit no-rescue boundaries.

The candidate packet must not include Go code, CLI flags, generated results,
backtests, optimizers, source expansion, derivatives-veto interaction,
paper/testnet/live flow, exchange API work, credentials, deployment, martingale,
averaging down, two-exchange logic, or promotion.
