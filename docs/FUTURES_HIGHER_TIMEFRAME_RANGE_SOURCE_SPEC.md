# Futures Higher-Timeframe Range Source Spec

Date: 2026-06-26

## Verdict

Stop state:
`higher_tf_range_source_spec_no_viable_range_premise`.

The BTCUSDT higher-timeframe source contract is now specified for future
range-only work, but no non-trading audit is ready yet. The next audit still
needs a materially different higher-timeframe range premise from the user,
including a closed-candle candidate event and a falsification rule.

This milestone does not add code, CLI flags, audits, generated result
directories, entries, exits, scoring, sizing, strategy replacement, source
downloads, sibling repo mutation, paper/testnet/live wiring, exchange API use,
deploy files, grid, martingale, averaging down, or two-exchange logic.

## Parent Source

The only accepted parent source for this higher-timeframe range lane is the
current Binance USDT-M futures BTCUSDT 5m CSV:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Required CLI product | `-source-product binance-usdm-futures` |
| Venue / product | `Binance` / `Binance USDT-M futures` |
| Symbol / interval | `BTCUSDT` / `5m` |
| CSV lines including header | `573,985` |
| Loaded candles | `573,984` |
| First parent open | `2021-01-01T00:00:00Z` |
| Last parent open | `2026-06-16T23:55:00Z` |
| Manifest status | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

Spot sources and comparison-only manifests are not valid parents for this lane.
Any source change is a research break and must be reviewed before a
higher-timeframe audit can use it.

## Candidate Bars

The approved candidate intervals are `15m`, `1h`, and `4h`. They are source
shapes for future range audits, not strategy approvals.

| Interval | Expected 5m Children | Expected Complete Bars From Current Parent | First HTF Open | Last HTF Open |
| --- | ---: | ---: | --- | --- |
| `15m` | `3` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` |
| `1h` | `12` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` |
| `4h` | `48` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` |

These counts are deterministic from the accepted parent coverage because the
first and last parent candles align to complete UTC buckets. A future audit
must still write or document its generated coverage and row counts before any
result is trusted.

## Closed UTC Resampling Contract

Parent bars are `5m` open-time candles. A higher-timeframe open time is the UTC
bucket start, and the bucket contains the expected child opens in
`[bucket_start, bucket_start + interval)`.

Aggregation rules:

- `open`: first child open;
- `high`: maximum child high;
- `low`: minimum child low;
- `close`: last child close;
- `volume`: sum of child volume.

A higher-timeframe bar is accepted only when every expected child `5m` open
exists exactly once. The final bucket must be dropped if it is partial. No
forward-filled bar, synthetic bar, duplicate timestamp, child-gap repair, local
timezone bucket, or future-looking close is allowed.

Closed-candle finality:

- a `15m` bar is knowable only after its third `5m` child closes;
- a `1h` bar is knowable only after its twelfth `5m` child closes;
- a `4h` bar is knowable only after its forty-eighth `5m` child closes;
- any later audit must evaluate signals from completed higher-timeframe bars
  only.

## Source Acceptance Gates

Before any higher-timeframe audit can run, all of these must be true:

- parent manifest is accepted, Binance USDT-M futures, `BTCUSDT`, `5m`, and
  `comparison_only=false`;
- parent coverage and row count match the accepted source facts above, or a
  source-change review explicitly replaces them;
- generated `15m`, `1h`, and/or `4h` coverage, first open, last open, row
  count, gap count, duplicate count, and child completeness are documented;
- no missing child opens exist inside accepted higher-timeframe buckets;
- timestamps remain UTC open times;
- the generated bars preserve closed-candle finality.

If any gate fails, the correct stop state is
`higher_tf_range_source_spec_blocked_by_source_or_resampling_gap`.

## Premise Acceptance Gates

The source contract alone does not justify an audit. The next premise must:

- stay BTCUSDT-only and range-only;
- use one or more of `15m`, `1h`, or `4h`;
- be materially different from failed BTCUSDT 5m SR timing, compression
  breakout, hold-inside/midline, and impulse surfaces;
- define the closed-candle observable that creates a candidate event;
- define the outcome label and falsification rule before any entry prototype;
- explain why the higher-timeframe source shape should reveal behavior that
  the failed `5m` families did not;
- avoid entries, exits, sizing, scoring, paper/testnet/live wiring, API keys,
  deploy work, grid, martingale, averaging down, and two-exchange logic.

If the premise is merely a renamed `5m` family resliced onto higher-timeframe
bars, the correct stop state is
`higher_tf_range_source_spec_rejected_as_5m_reslice`.

## Next Required Input

The next task should not start an audit automatically. It should ask for a
materially different higher-timeframe range premise:

1. Which interval or interval set should be used: `15m`, `1h`, `4h`, or a
   fixed comparison among them?
2. What range behavior should be audited?
3. Why should that behavior differ from the closed `5m` families?
4. Which closed-candle observable defines each candidate event?
5. What outcome would falsify the premise before any entry prototype?

Once that premise is explicit, a later brief may specify a non-trading audit
only. Until then, keep `lab.EmptyStrategy` and do not add strategy behavior.
