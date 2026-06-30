# Futures BTCUSDT 15m Post-Compression Directional Expansion Premise Spec

Date: 2026-06-30

## Verdict

Stop state:
`independent_entry_premise_spec_ready_for_user_approval`.

This documentation-only milestone selects one independent BTCUSDT `15m`
local-source premise family for a later zero-trade audit:

```text
btc_15m_post_compression_directional_expansion_v1
```

The premise asks whether a closed BTCUSDT `15m` candle that expands out of a
recent, low-percentile local range has directional label separation beyond the
unconditional eligible `15m` baseline. The later audit may test only the bounded
parameter grid declared here. It may not pick a P&L winner, simulate trades,
apply the derivatives veto, replay, walk forward, or promote a strategy.

This spec authorizes no Go code, CLI flag, generated result directory, audit
run, source download, source materialization, data write, entry, exit, P&L
backtest, optimizer grid, replay, walk-forward, portfolio construction,
paper/testnet/live path, exchange API, credential, deploy file, martingale,
averaging down, two-exchange logic, derivatives veto interaction, closed-family
rescue, or strategy promotion.

## Approved Premise

| Field | Decision |
| --- | --- |
| Premise id | `btc_15m_post_compression_directional_expansion_v1` |
| Route | New BTCUSDT `15m` local-source premise |
| Source | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Market | Binance USDT-M futures `BTCUSDT` |
| Parent interval | `5m` |
| Decision interval | closed UTC `15m` bars resampled from exact `5m` children |
| Candidate family | post-compression directional expansion |
| Side/timing | long after upside expansion at next `15m` bar open; short after downside expansion at next `15m` bar open |
| Audit type | later zero-trade label-separation audit only |
| Derivatives veto | excluded from entry inputs, candidate rows, side, ranking, and pass/fail scoring |

The active source facts carried into the later audit are the current accepted
BTCUSDT futures source contract: `573,984` loaded `5m` candles,
`2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`, `gap_count=0`,
`duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, and
`validation_status=accepted`.

## Candidate Event Definition

For each closed UTC `15m` decision candle `d`, compute the event with only data
known at or before `d` close.

For each compression lookback `L` in `{48, 96, 192}`:

1. Define the prior range over closed `15m` bars `[d-L, d-1]`.
2. `range_high_L(d)` is the maximum high in that prior window.
3. `range_low_L(d)` is the minimum low in that prior window.
4. `range_width_pct_L(d)` is
   `(range_high_L(d) - range_low_L(d)) / close[d-1]`.
5. The decision row is skipped for this `L` if the full prior window, prior
   close, prior range, or prior close price is missing or invalid.

For each compression threshold `Q` in `{0.20, 0.30, 0.40}`:

1. Build the rolling reference set from the previous `1,920` valid
   `range_width_pct_L` observations, ending at `d-1`.
2. The current decision value is excluded from its own reference set.
3. `compressed_L_Q(d)` is true when `range_width_pct_L(d)` is at or below the
   `Q` percentile of that prior reference set.
4. The row is skipped if fewer than `1,920` prior valid reference values exist.

For each breakout threshold `M` in `{0.1, 0.2, 0.3}`:

1. Compute Wilder-style `ATR(14)` on the resampled `15m` candles using the
   repo's existing `ATR(candles, 14)` convention.
2. Use only the prior-bar ATR: `atr14[d-1]`.
3. Upside expansion is true when
   `close[d] >= range_high_L(d) + M * atr14[d-1]`.
4. Downside expansion is true when
   `close[d] <= range_low_L(d) - M * atr14[d-1]`.
5. The row is skipped if prior ATR is missing, non-finite, or non-positive.

For each volume confirmation mode:

| Mode | Pass rule |
| --- | --- |
| `none` | no volume filter |
| `above_prior_96_median` | `volume[d]` is strictly greater than the median of volumes `[d-96, d-1]` |
| `above_prior_96_p60` | `volume[d]` is strictly greater than the 60th percentile of volumes `[d-96, d-1]` |

The current decision volume is known only at `d` close. It may be tested against
prior thresholds, but it must not be included in the threshold reference set.
Rows with fewer than `96` valid prior volume values are skipped for the two
volume-confirmed modes.

A candidate event exists for a parameter cell when compression is true, one
directional expansion is true, and the volume mode passes. The side is:

| Expansion | Side label | Timing label |
| --- | --- | --- |
| Upside | `long` | next `15m` bar open |
| Downside | `short` | next `15m` bar open |

The timing label is diagnostic only for this premise stage. The later audit may
use `open[d+1]` as a label anchor, but it must not simulate a fill, apply
fees/slippage, model stop/target exits, or write trades.

## Parameter Grid

The later audit may evaluate only this bounded grid:

| Dimension | Values |
| --- | --- |
| Compression lookback | `48`, `96`, `192` closed `15m` bars |
| Compression threshold | bottom `20%`, `30%`, `40%` of rolling prior range width |
| Percentile reference | prior `1,920` closed `15m` bars |
| Breakout location | close beyond prior range by `0.1`, `0.2`, `0.3` prior-bar `ATR(14)` |
| Volume confirmation | `none`, `above_prior_96_median`, `above_prior_96_p60` |

This creates `81` parameter cells before side. The grid is predeclared and
closed. A later audit may summarize the cells, but it may not add dimensions,
relax thresholds after result inspection, or rank cells by P&L.

Candidate viability counts must use de-duplicated `(decision_close, side)` rows
so repeated hits across parameter cells cannot inflate evidence.

## Later Zero-Trade Labels

Forward labels are evaluation metadata only. They may not be used in candidate
construction, filtering, side selection, or ranking before the label step.

For each candidate row and each horizon `H` in `{16, 32, 48}` closed `15m`
bars after the decision candle:

| Label | Definition |
| --- | --- |
| `label_anchor_open` | `open[d+1]`, used only as a diagnostic anchor |
| `intended_side_forward_close_return_bp` | side-adjusted return from `open[d+1]` to `close[d+H]`, in basis points |
| `intended_side_favorable_bp` | side-adjusted maximum favorable excursion from `open[d+1]` through bars `[d+1, d+H]`, in basis points |
| `intended_side_adverse_bp` | side-adjusted maximum adverse excursion from `open[d+1]` through bars `[d+1, d+H]`, in basis points |
| `favorable_gt_adverse` | true when favorable excursion is greater than adverse excursion |
| `favorable_minus_adverse_bp` | favorable excursion minus adverse excursion |

Rows without `d+1`, `d+H`, or a full forward high/low path are skipped and
counted in missingness. The later audit must also emit an unconditional eligible
`15m` baseline by split, side, and horizon using the same warmup and label
availability requirements.

## Falsification Gates

The later zero-trade audit must reject the premise if any required gate fails:

| Gate | Rejection condition |
| --- | --- |
| Source | source validation, `15m` resampling, or exact child coverage fails |
| Leakage | any feature uses data after decision candle close, or any threshold includes the current value when a prior reference is required |
| Candidate size | fewer than `300` de-duplicated `(decision_close, side)` candidates in the full sample |
| Split size | any primary period split has fewer than `50` de-duplicated candidates |
| Baseline separation | no side/horizon/cell cluster separates beyond the unconditional eligible `15m` baseline |
| Isolated-cell pass | only one isolated parameter cell passes without an adjacent passing cell |
| Split stability | full-sample separation is not supported in every primary period split |
| Closed-family boundary | passing evidence is only a retuned version of a reviewed closed family |
| Veto contamination | derivatives veto rows, basis/premium facts, or no-trade filter labels shape entries, side, candidate selection, or pass/fail scoring |

Primary period splits are `2021_2022_stress`, `2023_2024_oos`, and
`2025_2026_recent`, assigned by decision candle close time.

A side/horizon/cell separates beyond baseline only when all three metrics are
strictly better than the same-side unconditional eligible baseline in the full
sample and in every primary period split:

- mean `intended_side_forward_close_return_bp`;
- mean `favorable_minus_adverse_bp`;
- `favorable_gt_adverse` rate.

A passing cell is not enough by itself. The family can pass only if at least
two adjacent cells pass for the same side and horizon. Adjacent means the cells
differ by one ordered step in exactly one parameter dimension and match on all
other parameter dimensions. Ordered volume modes are `none`,
`above_prior_96_median`, then `above_prior_96_p60`.

## Closed-Family Separation

This premise intentionally shares the word "compression" with older work, but
it is not allowed to reopen the closed compression-breakout or structured-
compression families.

The material differences are:

- it is BTCUSDT-only and `15m`, not ETH/SOL-authoritative `4h` universe work;
- it uses a simple local prior-range percentile and ATR-displaced decision
  close, not detector `RawActive`/`Active` episodes;
- it has no max breakout-delay window after a detector episode ends;
- it has no target, stop, hold, confirmation window, cost model, P&L ranking,
  replay, walk-forward, or symbol-transfer step;
- it does not use router rows, occupancy rotation, midline/hold-inside states,
  BTC regime context, ETH/SOL context, derivatives basis/premium, or the
  no-trade veto.

The later audit must still reject itself as
`independent_entry_premise_spec_rejected_closed_family_reslice` if the only
passing evidence depends on detector episodes, old compression breakout fields,
structured-compression parameters, clean-breakout continuation, router
rotation, occupancy rotation, midline/hold-inside, higher-timeframe nested
rotation, BTC/ETH/SOL context, or the derivatives veto.

## Expected Later Audit Outputs

If the user later approves the zero-trade audit implementation, outputs should
remain compact and inspectable under a new ignored `results/` directory.
Expected artifacts:

- `source_manifest.json`;
- `summary.json`, `summary.csv`, and `trades.json` with `trades=0`;
- `btc_15m_post_compression_directional_expansion_sources.csv`;
- `btc_15m_post_compression_directional_expansion_resample_coverage.csv`;
- `btc_15m_post_compression_directional_expansion_parameter_cells.csv`;
- `btc_15m_post_compression_directional_expansion_candidates.csv`;
- `btc_15m_post_compression_directional_expansion_dedup_events.csv`;
- `btc_15m_post_compression_directional_expansion_baseline.csv`;
- `btc_15m_post_compression_directional_expansion_split_summary.csv`;
- `btc_15m_post_compression_directional_expansion_adjacency.csv`;
- `btc_15m_post_compression_directional_expansion_missingness.csv`;
- `btc_15m_post_compression_directional_expansion_falsification.json`.

These artifacts may count, label, group, and compare candidate rows. They must
not score strategy P&L or write simulated trades.

## Next Gate

The next step is explicit user approval for a zero-trade audit implementation
of `btc_15m_post_compression_directional_expansion_v1`.

Until that approval is given, do not add code, CLI flags, generated outputs, or
audit runs for this premise.

## Stop States

Used:

```text
independent_entry_premise_spec_ready_for_user_approval
```

Reserved for the later audit or rejection path:

```text
btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review
btc_15m_post_compression_directional_expansion_zero_trade_audit_failed_no_usable_entry_premise
independent_entry_premise_spec_rejected_closed_family_reslice
```
