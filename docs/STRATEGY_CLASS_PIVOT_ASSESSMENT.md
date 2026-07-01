# Strategy Class Pivot Assessment

Date: 2026-07-01

## Verdict

Stop state:

```text
strategy_class_pivot_assessment_recommends_trend_pullback
```

Recommended next research lane:

```text
trend-pullback continuation
```

Exact next allowed artifact after user approval:

```text
docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_CANDIDATE_PACKET.md
```

This is a docs-only assessment. It does not authorize Go code, CLI flags,
backtests, optimizers, generated results, source mutation, exchange APIs,
credentials, deployment files, paper/testnet/live flow, martingale, averaging
down, two-exchange logic, strategy promotion, or a renamed rescue of failed
range-reversion work.

## Current Project State

The active research source remains the accepted local BTCUSDT Binance USDT-M
futures `5m` CSV:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
```

Source facts used by the current research contract:

| Field | Value |
| --- | --- |
| Product | Binance USDT-M futures |
| Symbol | BTCUSDT |
| Base interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume rows | `66` |
| Comparison-only | `false` |
| Validation status | `accepted` |

The project has now tested three fixed backtest-first range candidates:

| Candidate | Stop State | Gross P&L | Net P&L | Full Trades | Primary Failure |
| --- | --- | ---: | ---: | ---: | --- |
| `btc_5m_rolling_value_area_reversion_v1` | `btc_5m_rolling_value_area_reversion_backtest_failed_no_usable_strategy` | `-489.495423` | `-999.999999` | `20,144` | Gross edge, net edge, drawdown |
| `btc_15m_previous_day_range_reversion_v1` | `btc_15m_previous_day_range_reversion_backtest_failed_no_usable_strategy` | `-506.146104` | `-935.263472` | `1,956` | Gross edge, net edge, drawdown |
| `btc_15m_range_edge_exhaustion_fade_v1` | `btc_15m_range_edge_exhaustion_fade_backtest_failed_no_usable_strategy` | `-154.405286` | `-261.595256` | `156` | Gross edge, net edge, drawdown |

The bounded range optimization workbench then ran:

| Field | Value |
| --- | --- |
| Review doc | `docs/RANGE_OPTIMIZATION_WORKBENCH_IMPLEMENTATION_REVIEW.md` |
| Run id | `20260630T200041Z-78f9a9e` |
| Immutable run path | `results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/` |
| Total trials | `112` |
| Passing candidates | `0` |
| Rejected candidates | `112` |
| Selected candidate | none |
| Stop state | `range_optimization_workbench_failed_no_candidate` |

## Why The Current Range Path Is Closed For Now

The fixed range-reversion candidates did not merely fail after costs. All three
failed gross edge, net edge, and drawdown gates. The value-area and previous-day
range baselines produced large enough trade counts to avoid a "not enough data"
excuse, while the range-edge exhaustion fade still passed the minimum
first-baseline trade-count gate and failed across every primary split.

The workbench then tested a controlled combination space around the failed
range-family components. It emitted all `112` trials, rejected all `112`, and
selected no candidate for locked fixed validation.

That closes the current range-reversion / midpoint / edge-fade / previous-day
range / bounded range-optimization path for now. Continuing by changing VWAP
windows, edge zones, previous-day definitions, fade thresholds, exit timings,
session choices, side selection, filters, derivatives veto interaction, or
optimizer grids would be a rescue of the failed family rather than a materially
different next lane.

## Scoring Summary

Scores are `0` to `5`; higher is better. The implementation-complexity column is
scored as simplicity, so `5` means very small implementation and `0` means large
new architecture.

| Rank | Strategy Class | Orthogonality | Existing Source Feasibility | Clear Event Definition | Cost/Slippage Realism | Validation Clarity | Implementation Simplicity | Total |
| ---: | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 1 | Trend-pullback continuation | 5 | 5 | 4 | 4 | 4 | 4 | 26 |
| 2 | Volatility expansion / breakout continuation | 4 | 5 | 4 | 4 | 4 | 4 | 25 |
| 3 | Session-based opening-range expansion | 4 | 5 | 4 | 4 | 4 | 4 | 25 |
| 4 | Liquidity sweep plus reclaim | 4 | 3 | 3 | 3 | 3 | 3 | 19 |
| 5 | Cross-asset or regime filter before entries | 3 | 2 | 3 | 1 | 2 | 3 | 14 |

No extra strategy class is added. The five requested classes are enough to
choose a clean next lane.

## Class Assessments

### Trend-Pullback Continuation

| Criterion | Score | Explanation |
| --- | ---: | --- |
| Orthogonality to failed work | 5 | The premise follows an established directional trend after a controlled pullback. It does not depend on range midpoint reversion, prior-day range reversion, or edge-fade exhaustion. |
| Existing source feasibility | 5 | A fixed BTCUSDT OHLCV baseline can use only the accepted `5m` futures source, with optional exact `15m` resampling already consistent with the lab's infrastructure. |
| Clear event definition | 4 | A simple first candidate can be fixed as EMA or moving-average slope, pullback into a predefined band, continuation close, and next-bar-open entry. The exact packet must lock one timeframe and one rule set. |
| Cost/slippage realism | 4 | If an edge exists, a continuation move should offer more room than tight midpoint fades. It is still not immune to churn, so the first packet should prefer selective events and gross-edge-first gates. |
| Validation clarity | 4 | Split metrics, side reporting, trade-count gates, gross/net gates, drawdown gates, and no-rescue boundaries are straightforward. |
| Implementation complexity | 4 | It needs only bounded indicator calculations and the existing backtest engine. No new source, optimizer, cross-asset join, or derivatives context is needed. |

Assessment: this is the cleanest next research lane. It is materially different
from the failed range path, can be assessed with the current source, has enough
movement potential to justify costs if the signal exists, and fits the
backtest-first workflow as one fixed baseline after a short docs-only candidate
packet.

### Volatility Expansion / Breakout Continuation

| Criterion | Score | Explanation |
| --- | ---: | --- |
| Orthogonality to failed work | 4 | It is directional continuation rather than range reversion, midpoint targeting, or range-edge fading. The score is not `5` because prior clean-breakout and post-compression expansion work is already closed and must not be renamed or rescued. |
| Existing source feasibility | 5 | Range expansion, ATR expansion, recent high/low breaks, candle body expansion, and volume confirmation can all be computed from the accepted BTCUSDT OHLCV source. |
| Clear event definition | 4 | A fixed event such as close above a prior `N`-bar high with ATR or body expansion is easy to specify, but the packet must avoid broad threshold search. |
| Cost/slippage realism | 4 | Breakout continuation can produce enough movement to survive the existing `0.0004` fee and `0.000116` slippage per side if selective. False breaks can still make it fragile. |
| Validation clarity | 4 | Clean pass/fail gates are available: gross edge, OOS/recent gross edge, net edge, drawdown, trade count, side splits, and no adjacent-cell rescue. |
| Implementation complexity | 4 | A single fixed baseline is small with current infrastructure, provided it avoids optimizer expansion and closed post-compression templates. |

Assessment: promising, but second to trend-pullback because the repo already has
failed adjacent breakout/expansion families. A future packet would need to state
precisely why it is not clean-breakout, post-compression, or breakout-retest
rescue work.

### Session-Based Opening-Range Expansion

| Criterion | Score | Explanation |
| --- | ---: | --- |
| Orthogonality to failed work | 4 | It trades expansion away from a fixed opening range rather than fading range edges. It remains range-boundary language, so it is less orthogonal than trend-pullback. |
| Existing source feasibility | 5 | A UTC opening range can be built from the accepted BTCUSDT futures OHLCV source without new data. |
| Clear event definition | 4 | A simple fixed packet is possible, but BTCUSDT trades continuously and the chosen "open" is partly conventional. The first packet would need one predeclared UTC anchor and range length. |
| Cost/slippage realism | 4 | Opening-range expansion can have enough movement if it catches real session impulse. It can also overtrade noisy breaks if the opening window is arbitrary. |
| Validation clarity | 4 | Split/trade/drawdown gates are clean. The major guardrail is forbidding calendar/session mining after seeing results. |
| Implementation complexity | 4 | It is small using current candles and resampling, with no new source or engine change. |

Assessment: a viable backup lane. It loses the tie-break to volatility expansion
and trend-pullback because session anchoring in a 24/7 market invites calendar
mining unless the candidate packet is very strict.

### Liquidity Sweep Plus Reclaim

| Criterion | Score | Explanation |
| --- | ---: | --- |
| Orthogonality to failed work | 4 | A sweep and reclaim event can be different from range midpoint reversion, especially if it uses a failed break and continuation acceptance instead of fading to the midpoint. It still touches edge/reclaim language and must avoid old false-break or SR timing rescue work. |
| Existing source feasibility | 3 | The current OHLCV source can approximate sweeps with highs/lows and closes, but it cannot observe order-book liquidity, liquidation flow, or actual stop clusters. |
| Clear event definition | 3 | A bounded event is possible, but swing selection, sweep depth, reclaim window, and invalidation rules create parameter ambiguity. |
| Cost/slippage realism | 3 | A true sweep-reclaim can move enough, but tight invalidation and frequent failed reclaims can become cost-sensitive. |
| Validation clarity | 3 | It can be falsified, but guardrails must be strong because small changes to sweep depth and reclaim timing can become post-hoc rescue knobs. |
| Implementation complexity | 3 | More bookkeeping is needed for swing references, sweep windows, reclaim windows, and duplicate-event suppression. |

Assessment: conceptually interesting but not the cleanest next step. It should
wait unless the operator prefers a microstructure-inspired lane despite the OHLCV
proxy limitation.

### Cross-Asset Or Regime Filter Before Entries

| Criterion | Score | Explanation |
| --- | ---: | --- |
| Orthogonality to failed work | 3 | A filter/context lane is different from range entries, but it is not a standalone entry premise. It risks becoming a way to rescue failed entries after the fact. |
| Existing source feasibility | 2 | BTC-only regime filters can use the accepted source, but cross-asset context cannot be assessed using only the accepted BTCUSDT CSV. |
| Clear event definition | 3 | Regime states can be defined, but a filter-before-entry class still needs an independent entry stream to evaluate. |
| Cost/slippage realism | 1 | A filter does not produce movement or trades by itself. Cost realism cannot be judged until it annotates an approved candidate stream. |
| Validation clarity | 2 | It can be tested cleanly only after a fixed entry baseline exists. Before that, it invites post-hoc conditioning. |
| Implementation complexity | 3 | A BTC-only regime labeler is manageable, but cross-asset joins or source expansion add review burden. |

Assessment: not recommended as the next lane. It should remain deferred until
there is a separate entry premise worth filtering, consistent with the existing
derivatives no-trade filter boundary.

## Ranking

1. `trend-pullback continuation`
2. `volatility expansion / breakout continuation`
3. `session-based opening-range expansion`
4. `liquidity sweep plus reclaim`
5. `cross-asset or regime filter before entries`

The second and third classes tie on total score. Volatility expansion ranks
above opening-range expansion because it can be specified without choosing an
arbitrary session open in a continuous market. Opening-range expansion remains a
reasonable backup if the operator wants a session-structured baseline.

## Recommended Next Research Lane

Create a docs-only backtest-first candidate packet for trend-pullback
continuation.

The packet should select exactly one simple BTCUSDT futures baseline. A suitable
shape for user approval would be:

```text
btc_15m_trend_pullback_continuation_v1
```

The candidate packet should lock:

- accepted BTCUSDT Binance USDT-M futures source;
- one decision timeframe, likely exact closed UTC `15m` bars from the current
  `5m` source;
- one trend definition;
- one pullback definition;
- one continuation trigger;
- fixed stop, target, time stop, sizing, fee, and slippage assumptions;
- side-separated reporting;
- split/trade/gross/net/drawdown gates; and
- explicit no-rescue boundaries.

Do not implement the baseline in the assessment PR. The packet is the next
bounded docs-only gate, and implementation would require a later explicit user
approval.

## Next Allowed Artifact After User Approval

After the operator approves this recommendation, the next task may add exactly
one docs-only artifact:

```text
docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_CANDIDATE_PACKET.md
```

That artifact may define the first fixed trend-pullback continuation baseline.
It may not add Go code, CLI flags, generated results, a backtest run, an
optimizer, source expansion, derivatives-veto interaction, paper/testnet/live
flow, exchange API work, credentials, deployment files, martingale, averaging
down, two-exchange logic, or promotion.
