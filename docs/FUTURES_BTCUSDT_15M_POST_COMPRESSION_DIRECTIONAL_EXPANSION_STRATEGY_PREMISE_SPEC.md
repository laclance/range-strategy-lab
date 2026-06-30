# Futures BTCUSDT 15m Post-Compression Directional Expansion Strategy Premise Spec

Date: 2026-06-30

## Verdict

Stop state:
`post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval`.

The user explicitly approved this docs-only strategy-premise spec for the
passed zero-trade audit:

```text
btc_15m_post_compression_directional_expansion_v1
```

The zero-trade evidence is sufficient to request a later docs-only offline
backtest spec, but only for one tightly bounded representative candidate from
the passing pocket. It is not a P&L result and it does not authorize a backtest
implementation, trade simulation, optimizer, replay, walk-forward,
derivatives-veto interaction, paper/testnet/live path, or promotion.

Selected later backtest-spec candidate:

```text
btc_15m_post_compression_l192_q20_m020_none_long_h48_v1
```

This candidate is the conservative representative cell from the passing
zero-trade pocket: long side only, `48` closed `15m` bar evidence horizon,
`192`-bar compression lookback, bottom `20%` compression threshold, `0.2`
prior-bar `ATR(14)` upside breakout threshold, and no volume confirmation.

## Evidence Status

The approved zero-trade audit passed at:

```text
btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review
```

That pass means the predeclared label-separation gates cleared. It does not
mean the entry earns money after costs, survives stop/target rules, has good
fill quality, or is durable under replay or walk-forward.

The only eligible evidence pocket was:

| Field | Eligible evidence |
| --- | --- |
| Side | `long` only |
| Horizon | `48` closed `15m` bars only |
| Compression lookback | `192` prior closed `15m` bars |
| Compression threshold | bottom `20%` of prior range width |
| Breakout threshold | `0.1`, `0.2`, `0.3` prior-bar `ATR(14)` |
| Volume mode | `none`, `above_prior_96_median`, `above_prior_96_p60` |

Everything outside that pocket is rejected for this strategy-premise line:
short-side evidence, `16`/`32`-bar labels, `48`/`96` lookbacks, `30%`/`40%`
compression thresholds, derivatives-veto facts, old detector episodes, router
rows, occupancy rotation, midline/hold-inside states, BTC/ETH/SOL context, and
source expansion.

## Why One Representative Cell

The later backtest request should use a single conservative representative cell,
not the full `81`-cell grid and not a P&L-ranked cluster.

The selected cell is:

| Dimension | Value |
| --- | --- |
| Compression lookback | `192` prior closed `15m` bars |
| Compression threshold | bottom `20%` |
| Breakout threshold | `0.2` prior-bar `ATR(14)` |
| Volume confirmation | `none` |
| Side | `long` |
| Evidence horizon | `48` closed `15m` bars |

Rationale:

- `0.2 ATR` is the center breakout threshold of the three passing breakout
  values, avoiding the loosest and strictest edge cases.
- `none` avoids adding a volume filter that did not materially change the
  zero-trade evidence and would add another tuning surface.
- `192`/`q20` is not optional: it is the only lookback and compression threshold
  surface that passed the full gate.
- The adjacent passing cells remain supporting robustness evidence only. They
  may be cited in the later spec or review, but they may not become a P&L
  optimizer or rescue a failed representative-cell backtest.

## Candidate Construction For Later Spec

The later backtest spec may request only this candidate stream.

For each closed UTC `15m` decision candle `d`, using exact `15m` candles
resampled from complete `5m` children:

1. Compute the prior range over closed `15m` bars `[d-192, d-1]`.
2. Let `range_high(d)` be the maximum high in that prior window.
3. Let `range_low(d)` be the minimum low in that prior window.
4. Let `range_width_pct(d)` be
   `(range_high(d) - range_low(d)) / close[d-1]`.
5. Build the compression reference from the previous `1,920` valid
   `range_width_pct` observations ending at `d-1`.
6. Require `range_width_pct(d)` at or below the `20%` percentile of that prior
   reference set.
7. Compute `ATR(14)` on the resampled `15m` candles and use `ATR(14)[d-1]`.
8. Require `close[d] >= range_high(d) + 0.2 * ATR(14)[d-1]`.
9. Do not apply a volume confirmation filter.
10. Label the candidate side as `long`.
11. The later backtest entry timing, if separately approved, is the next `15m`
    bar open after the decision candle.

Rows must be skipped, not filled, when any required prior range, percentile
reference, ATR value, current close, next open, or future path needed by the
later approved test is unavailable.

## Source Facts

The later spec must carry forward the accepted source contract:

| Field | Value |
| --- | ---: |
| Source path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol | BTCUSDT |
| Interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume rows | `66` |
| Comparison-only | `false` |
| Validation status | `accepted` |

The accepted `15m` resample facts from the zero-trade audit were:

| Field | Value |
| --- | ---: |
| Row count | `191,328` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:45:00Z` |
| Last close | `2026-06-16T23:59:59Z` |
| Expected child bars | `3` |
| Missing child opens | `0` |
| Validation status | `accepted` |

Any later backtest spec or implementation must fail closed if the source,
market type, symbol, interval, coverage, or resample identity drifts without a
new approved source-impact review.

## No-Lookahead Rules

The later spec and any later implementation must preserve these rules:

- prior range uses `[d-192, d-1]`;
- compression percentile references exclude `d` and require exactly `1,920`
  prior valid observations;
- ATR breakout threshold uses `ATR(14)[d-1]`;
- the decision close is known only at `d` close;
- the earliest entry timing is `open[d+1]`;
- forward returns, high/low paths, stops, targets, and trade outcomes may never
  shape candidate construction;
- missing context produces skip/missingness rows, never silent defaults.

## Fixed Risk And Exit Question

This strategy-premise spec does not choose trade exits or run P&L. It defines
the only risk/exits question the next docs-only backtest spec may answer:

Can this one long-only representative entry stream survive a single fixed,
non-optimized offline backtest model with:

- one open position max;
- next-`15m`-open entries;
- stop-first ambiguity;
- repo-standard fees, slippage, `1%` risk sizing, and `1x` notional cap unless
  the later backtest spec explicitly restates the current repo defaults;
- a fixed maximum holding window tied to the zero-trade evidence horizon of
  `48` closed `15m` bars;
- one predeclared stop/target or stop/time-exit model selected in the later
  docs-only backtest spec before any implementation.

The later backtest spec must not run a stop grid, target grid, holding-period
grid, ATR-multiple optimizer, volume-filter optimizer, side optimizer, or
parameter-cell optimizer. If one fixed exit model cannot be selected without
user preference, the later docs-only backtest spec must stop at a user-choice
gate.

## Fail Gates For Later Backtest Spec

The later docs-only backtest spec must reject or stop before implementation if
any of these conditions hold:

| Gate | Stop condition |
| --- | --- |
| Source identity | source or `15m` resample contract cannot be preserved exactly |
| Candidate identity | the requested stream uses anything other than `L=192`, `q20`, `M=0.2`, volume `none`, long side |
| Evidence misuse | zero-trade labels are treated as P&L, executable fill proof, or promotion evidence |
| Exit ambiguity | one fixed risk/exit model cannot be declared before implementation |
| Optimizer creep | the full `81`-cell grid, neighboring cells, stops, targets, holds, sides, volume modes, or veto rows are used for P&L selection |
| Closed-family reslice | the premise depends on old detector episodes, structured compression, clean breakout, router rotation, occupancy rotation, midline/hold-inside, HTF nested rotation, or BTC/ETH/SOL context |
| Veto contamination | `btc_15m_basis_discount_no_trade_veto_v1` shapes entries, exits, side, ranking, P&L, or pass/fail scoring |

If a later implementation is eventually approved, its review must be allowed to
fail the line even if the zero-trade evidence passed. A P&L backtest pass would
need positive net evidence after costs across the full sample and primary
splits under the one fixed model, without optimizer rescue.

## Output Expectations For Later Backtest Spec

The next milestone, if approved, should still be docs-only. It should create a
future implementation brief, not code.

Expected docs-only backtest-spec outputs:

- a single backtest candidate id and event definition;
- exact source and resample contract;
- fixed entry timing;
- one fixed risk/exit model or a user-choice stop state;
- P&L metrics and split gates for a later implementation;
- artifact names for a future ignored `results/` directory;
- stop states for pass, fail, source failure, closed-family reslice, and
  optimizer contamination;
- explicit exclusions for derivatives veto interaction and all non-passing
  post-compression cells.

No generated CSV/JSON outputs are expected from this strategy-premise spec.

## Boundaries

This spec authorizes no Go code, CLI flags, generated result directory,
source download, network request, data write, audit run, backtest run, replay,
walk-forward, optimizer selection, source/symbol expansion, portfolio
construction, paper/testnet/live path, exchange API, credentials, deploy file,
martingale, averaging down, two-exchange logic, derivatives-veto interaction,
or strategy promotion.

The canonical derivatives veto remains parked as future skip/retain evidence
only:

```text
btc_15m_basis_discount_no_trade_veto_v1
```

It may not shape candidate rows, create entries, choose side, rank, score P&L,
replay, walk forward, optimize, promote a strategy, or reopen closed families.
Any future veto interaction audit requires separate approval after an
independent backtest candidate exists.

## Next Gate

The next allowed step is explicit user approval for a docs-only offline
backtest spec for:

```text
btc_15m_post_compression_l192_q20_m020_none_long_h48_v1
```

Until that approval is given, do not add backtest code, CLI flags, generated
outputs, audit/backtest runs, source downloads, P&L artifacts, replay,
walk-forward, or derivatives-veto interaction.

## Stop States

Used:

```text
post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval
```

Reserved for the later docs-only backtest spec:

```text
post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval
post_compression_directional_expansion_backtest_spec_needs_user_exit_choice
post_compression_directional_expansion_backtest_spec_rejected_too_narrow
post_compression_directional_expansion_backtest_spec_rejected_optimizer_contamination
post_compression_directional_expansion_backtest_spec_rejected_closed_family_reslice
```
