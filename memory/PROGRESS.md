# Progress

This file is the always-read project snapshot. Keep it compact: current state,
latest verification, result paths, and a milestone index. Detailed evidence
belongs in focused `docs/` reviews, generated artifacts under `results/`, and
git history.

## Current State

- Scope is now offline Binance USDT-M futures range-strategy construction.
  The implemented CLI default remains BTCUSDT `5m`. Research is not stopped,
  but automatic reuse of failed premises is stopped: structured compression,
  breakout-retest/acceptance, and the BTCUSDT higher-timeframe nested
  range-rotation audit are exclusion evidence in their reviewed forms. The
  user-approved next posture is range-first, BTCUSDT-first strategy
  construction from scratch through documented offline backtesting and
  optimization stages. The first v1 construction implementation,
  `range_occupancy_rotation_v1`, has now run its bounded optimizer/backtester
  and failed with no selectable grid config. No fixed replay, walk-forward,
  strategy package, retune, symbol expansion, or live-adjacent work is
  authorized from this V1 grammar.
- Active market target is Binance USDT-M futures, not spot. Spot-generated
  audits/reviews are historical context only unless a futures rerun explicitly
  revalidates a specific conclusion.
- Default strategy remains `lab.EmptyStrategy`; trades remain `0` unless an
  explicit offline prototype flag is passed.
- Closed-candle decision semantics remain required.
- `cmd/rangelab` now enforces source identity before audits/backtests. The
  default `-csv` is the full-history Binance USDT-M futures file
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`; non-default CSVs
  require `-source-product`; spot CSVs require explicit comparison flags; every
  accepted run writes `source_manifest.json`.
- No live code, API keys, deploy scripts, grid, martingale, averaging down, or
  two-exchange execution is allowed.
- `memory/NEXT_CODEX_BRIEF.md` is the only canonical next-session prompt.
- The first minimal offline prototype for the futures-revalidated
  `hold_3_inside` + first `mid_touch` within `12` bars + event close-position
  bucket `mid_50` surface has been built and reviewed. The close-back
  boundary-target template failed P&L and is not promoted.
- Futures impulse absorption after abnormal OHLCV impulse candles has been
  implemented, run, and reviewed on Binance USDT-M futures data. It failed the
  non-trading review gate: continuation-first and quick continuation dominated
  midpoint reclaim-first across all period splits and horizons. No prototype or
  promotion is approved.
- Futures range candidate discovery has been implemented, run, and reviewed.
  It passed the discovery gate only for clean breakout continuation surfaces.
  The fixed-rule offline baseline backtest for the top non-duplicative `4h`
  up-breakout and `1h` all-side clean-breakout candidates has now been built
  and reviewed. Both candidates failed after costs; no optimization,
  portfolio-style stream, or automatic `15m` expansion is approved from that
  baseline.
- The local BTC/ETH/SOL range-universe discovery audit has been implemented,
  run, and reviewed. It passed the universe gate for structured compression
  expansion and breakout retest/acceptance rows. The fixed-rule structured
  compression universe baseline has now been built and reviewed. The `4h all
  h6` aggregate passed after costs and is approved for a bounded offline
  optimization/robustness brief only; BTCUSDT was weak, and the result depends
  on ETH/SOL strength. The `1h all h12` surface failed and is not promoted.
- The bounded `4h` structured-compression optimization selected
  `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`, the strategy spec froze
  it, and the fixed offline replay/backtest passed. The bounded offline
  walk-forward robustness pass has now run and is fragile: no candidate
  strategy package is approved, and the structured-compression ETH/SOL stream
  is exclusion evidence. The post-fragility pivot review authorized one
  bounded offline breakout-retest/acceptance baseline from existing
  range-universe discovery evidence; that baseline also failed after costs and
  is not promoted.
- The user-selected higher-timeframe nested range-rotation premise has now
  been implemented as a non-trading BTCUSDT audit and reviewed. Source and
  resample validation passed, but only `3` valid events appeared across the
  full sample, all upside, so the audit failed the no-baseline gate. No
  baseline backtest, optimizer, replay, walk-forward, source expansion, symbol
  expansion, or strategy package is approved from this premise.

## 2026-06-27

Futures range-first occupancy rotation V1 optimizer:

- Review doc:
  `docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md`.
- CLI flag:
  `-futures-range-first-occupancy-rotation-v1-optimization`.
- Result dir:
  `results/futures-range-first-occupancy-rotation-v1-optimization/`.
- Stop state:
  `range_first_strategy_v1_optimizer_failed_no_replay`.
- Source facts:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  - product/symbol/interval: Binance USDT-M futures `BTCUSDT` `5m`;
  - rows: `573,984`;
  - first/last open:
    `2021-01-01T00:00:00Z` to `2026-06-16T23:55:00Z`;
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`.
- Closed UTC resample facts:
  - `1h`: `47,832` rows, first open `2021-01-01T00:00:00Z`,
    last open `2026-06-16T23:00:00Z`, complete, accepted;
  - `15m`: `191,328` rows, first open `2021-01-01T00:00:00Z`,
    last open `2026-06-16T23:45:00Z`, complete, accepted.
- Grid and selection:
  - declared grid rows: `1,152`;
  - passing grid rows: `0`;
  - selected config: none;
  - top ranked row:
    `range_occupancy_rotation_v1_1h_l24_w020_ow8_occ070_rec25_t66_h12_sb000`;
  - top row did not pass because it had only `2` train trades, too few
    OOS/recent/full trades, negative OOS and full net P&L, weak PF, and side
    weakness.
- Fixed baseline:
  `range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005`.
  Common `summary.csv/json` and `trades.json` describe this baseline only.
- Baseline outcome:
  - `2021_2022_stress`: `8` trades, net P&L `-22.653070`, PF `0.464536`;
  - `2023_2024_oos`: `16` trades, net P&L `-38.248119`, PF `0.428362`;
  - `2025_2026_recent`: `19` trades, net P&L `-33.307557`, PF `0.504958`;
  - `full_2021_2026`: `43` trades, net P&L `-94.208745`, PF `0.466232`;
  - baseline rank: `769` of `1,152`.
- CSV line counts including headers:
  - `futures_range_first_occupancy_rotation_v1_baseline.csv`: `2`;
  - `futures_range_first_occupancy_rotation_v1_coverage.csv`: `3`;
  - `futures_range_first_occupancy_rotation_v1_grid.csv`: `1,153`;
  - `futures_range_first_occupancy_rotation_v1_rankings.csv`: `1,153`;
  - `futures_range_first_occupancy_rotation_v1_selection.csv`: `3`;
  - `futures_range_first_occupancy_rotation_v1_signals.csv`: `84`;
  - `futures_range_first_occupancy_rotation_v1_skips.csv`: `42,543`;
  - `futures_range_first_occupancy_rotation_v1_sources.csv`: `2`;
  - `futures_range_first_occupancy_rotation_v1_summary.csv`: `13,825`;
  - `futures_range_first_occupancy_rotation_v1_trades.csv`: `44`;
  - common `summary.csv`: `13`.
- Review outcome: source and resample validation passed, but fixed baseline
  failed all split gates after costs and no declared grid config passed the
  optimizer gates. No fixed replay spec is authorized.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a review-only stop requiring a
  materially different user-approved offline range-first premise before any
  further implementation.
- Verification commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-first-occupancy-rotation-v1-optimization -out-dir results/futures-range-first-occupancy-rotation-v1-optimization`
  - `wc -l results/futures-range-first-occupancy-rotation-v1-optimization/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`

## 2026-06-26

Futures range-first strategy construction v1 spec:

- Spec doc:
  `docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md`.
- Stop state:
  `range_first_strategy_v1_spec_ready_for_optimizer_implementation`.
- This was documentation-only. No strategy, optimizer, replay, walk-forward
  run, CLI flag, generated result directory, artifact writer, data download,
  live/paper/testnet path, exchange API, credential, deploy file, martingale,
  averaging down, or two-exchange logic was added.
- Selected v1 grammar:
  `range_occupancy_rotation_v1`.
- The grammar builds a rolling closed-candle range envelope, requires recent
  closes to cluster persistently in one outer part of the range, then enters
  only after a closed candle recaptures an interior line back toward range
  value. It is not a structured-compression, breakout-retest, clean-breakout,
  hold-inside/midline, impulse-absorption, or nested range-rotation reslice.
- Source contract remains BTCUSDT Binance USDT-M futures `5m`:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  `573,984` loaded candles; open-time coverage
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Expected V1 signal resamples:
  - `15m`: `191,328` rows, `2021-01-01T00:00:00Z` through
    `2026-06-16T23:45:00Z`;
  - `1h`: `47,832` rows, `2021-01-01T00:00:00Z` through
    `2026-06-16T23:00:00Z`.
- Declared future implementation flag:
  `-futures-range-first-occupancy-rotation-v1-optimization`.
- Declared future result dir:
  `results/futures-range-first-occupancy-rotation-v1-optimization/`.
- Declared grid:
  `1,152` configs plus the fixed baseline row
  `range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005`.
- Common outputs for the future optimizer are locked to the fixed baseline
  (`source_manifest.json`, `summary.csv/json`, `trades.json`); grid, ranking,
  selected-config, comparison, skipped-signal, and selected-trade rows belong
  only in V1-specific artifacts until a later fixed replay is approved.
- Added a durable decision limiting the next implementation to this exact
  BTCUSDT, range-first, `15m`/`1h`, declared-grid optimizer/backtester.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the bounded implementation brief
  for `range_occupancy_rotation_v1`.
- Verification commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`

Futures range-first strategy construction protocol:

- Protocol doc:
  `docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_PROTOCOL.md`.
- Stop state:
  `range_first_strategy_construction_protocol_ready_for_v1_spec`.
- This was documentation-only. No strategy, optimizer, replay, walk-forward
  run, CLI flag, generated result directory, data download, live/paper/testnet
  path, exchange API, credential, deploy file, martingale, averaging down, or
  two-exchange logic was added.
- User-approved posture:
  - research continues as a fresh range-derived strategy construction process;
  - scope remains range-first and BTCUSDT-first;
  - the first construction source remains the accepted Binance USDT-M futures
    BTCUSDT `5m` file
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  - from-scratch means defining a new strategy grammar and review ladder, not
    retuning failed structured-compression, breakout-retest, clean-breakout,
    hold-inside/midline, impulse absorption, or nested-rotation families.
- The protocol ladder is:
  strategy grammar spec -> baseline backtest -> bounded optimization -> fixed
  replay -> walk-forward robustness -> package review.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a spec-only brief for
  `docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md`; the next brief
  should define optimizer/backtester grammar, bounds, gates, artifacts, and
  stop states, but not implement code yet.
- Added a durable decision that from-scratch strategy construction is allowed
  only as offline, range-first, BTCUSDT-first work until a later approved brief
  expands scope.
- Verification commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`

Futures higher-timeframe nested range rotation audit:

- Review doc:
  `docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md`.
- CLI flag:
  `-futures-higher-tf-nested-range-rotation-audit`.
- Result dir:
  `results/futures-higher-tf-nested-range-rotation-audit/`.
- Stop state:
  `higher_tf_nested_range_rotation_audit_failed_no_baseline`.
- Common outputs stayed non-trading:
  `trades.json` contains `0` trades and common `summary.csv/json` are
  zero-trade compatibility outputs.
- Source facts:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  - product/symbol/interval: Binance USDT-M futures `BTCUSDT` `5m`;
  - rows: `573,984`;
  - first/last open:
    `2021-01-01T00:00:00Z` to `2026-06-16T23:55:00Z`;
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`.
- Closed UTC resample facts:
  - parent `4h`: `11,958` rows, first open `2021-01-01T00:00:00Z`,
    last open `2026-06-16T20:00:00Z`, complete, accepted;
  - child `1h`: `47,832` rows, first open `2021-01-01T00:00:00Z`,
    last open `2026-06-16T23:00:00Z`, complete, accepted.
- CSV line counts including headers:
  - `futures_higher_tf_nested_range_rotation_child_ranges.csv`: `283`;
  - `futures_higher_tf_nested_range_rotation_coverage.csv`: `3`;
  - `futures_higher_tf_nested_range_rotation_events.csv`: `20`;
  - `futures_higher_tf_nested_range_rotation_parent_ranges.csv`: `69`;
  - `futures_higher_tf_nested_range_rotation_sources.csv`: `2`;
  - `futures_higher_tf_nested_range_rotation_summary.csv`: `13`;
  - common `summary.csv`: `13`.
- Audit counts:
  - parent ranges: `68`;
  - child ranges: `282`, with `11` eligible, `262` no-valid-parent skips,
    `8` child-width-above-40pct-parent skips, and `1` child-not-inside-parent
    skip;
  - event rows: `19`, with `3` valid events and skips for
    `duplicate_child_event=9`, `event_beyond_parent_midpoint=6`, and
    `event_wrong_rotation_side=1`;
  - all valid events were `up`; no `down` event passed the nested geometry and
    first-break rules.
- Outcome facts:
  - `2021_2022_stress`: `1` valid event, adverse parent/child invalidation;
  - `2023_2024_oos`: `1` valid event, favorable midpoint and far quartile;
  - `2025_2026_recent`: `1` valid event, no resolution;
  - full sample: `3` valid events, favorable midpoint/far quartile rate
    `33.33%`, adverse parent/child invalidation rate `33.33%`,
    no-resolution rate `33.33%`, quick invalidation `0%`.
- Review outcome: source/resample validation passed, but the premise failed
  the `100` full-event gate, every `25` event split gate, and the side-balance
  gate. This is exclusion evidence for the frozen `4h` parent / `1h` child
  nested range-rotation audit, not permission to retune width/horizon/gates.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a review-only no-automatic
  implementation stop requiring a materially different user-approved offline
  range premise before more strategy work.
- Verification commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-higher-tf-nested-range-rotation-audit -out-dir results/futures-higher-tf-nested-range-rotation-audit`
  - `wc -l results/futures-higher-tf-nested-range-rotation-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`

Futures higher-timeframe nested range rotation premise spec:

- Review doc:
  `docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_PREMISE_SPEC.md`.
- Stop state:
  `higher_tf_nested_range_rotation_premise_ready_for_audit`.
- This was review-only. No code, CLI flag, generated result directory,
  strategy, optimizer, replay, walk-forward, data download, live/paper/testnet
  path, exchange API, credential, deploy file, martingale, averaging down, or
  two-exchange logic was added.
- Premise:
  - use the accepted BTCUSDT Binance USDT-M futures `5m` parent source;
  - derive closed UTC `4h` parent ranges and `1h` child ranges;
  - after a mature `1h` child range forms fully inside one half of a frozen
    mature `4h` parent range, measure whether the first closed `1h`
    displacement toward the parent interior reaches parent midpoint / far
    quartile before child-range invalidation.
- Source facts carried forward from the higher-timeframe source spec:
  - parent `5m` source:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  - loaded candles: `573,984`;
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`;
  - accepted manifest facts:
    `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`;
  - expected closed UTC resamples: `1h=47,832` rows through
    `2026-06-16T23:00:00Z`, `4h=11,958` rows through
    `2026-06-16T20:00:00Z`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` toward the bounded non-trading audit
  behind `-futures-higher-tf-nested-range-rotation-audit`.
- Added a durable decision limiting this premise to a non-trading audit until
  source/resample, candidate-count, split-stability, side-balance, and outcome
  gates pass.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Futures range universe breakout retest acceptance baseline:

- Review doc:
  `docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md`.
- Stop state:
  `breakout_retest_acceptance_baseline_failed_no_promotion`.
- Added explicit offline CLI flag:
  `-futures-range-universe-breakout-retest-acceptance-baseline-backtest`.
- This fixed-rule baseline selected only passing all-side
  `breakout_retest_acceptance` discovery rows and skipped duplicate
  timeframe/side variants:
  - selected `breakout_retest_acceptance_15m_all_h12` from discovery rank
    `22`;
  - selected `breakout_retest_acceptance_1h_all_h12` from discovery rank
    `28`.
- Result dir:
  `results/futures-range-universe-breakout-retest-acceptance-baseline-backtest/`.
- Outputs:
  - `futures_range_universe_breakout_retest_acceptance_baseline_sources.csv/json`
  - `futures_range_universe_breakout_retest_acceptance_baseline_coverage.csv/json`
  - `futures_range_universe_breakout_retest_acceptance_baseline_selection.csv/json`
  - `futures_range_universe_breakout_retest_acceptance_baseline_signals.csv/json`
  - `futures_range_universe_breakout_retest_acceptance_baseline_trades.csv/json`
  - `futures_range_universe_breakout_retest_acceptance_baseline_summary.csv/json`
  - common `source_manifest.json`, `summary.csv/json`, and `trades.json`
- Source facts:
  - each of `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` loaded `573,984` Binance
    USDT-M futures `5m` candles from `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`;
  - gaps / duplicates were `0` / `0` for every symbol;
  - zero-volume counts: `BTCUSDT=66`, `ETHUSDT=47`, `SOLUSDT=47`;
  - physical non-monotonic counts: `BTCUSDT=0`, `ETHUSDT=0`,
    `SOLUSDT=1`; SOLUSDT was sorted and accepted.
- Selected closed UTC resamples:
  - `15m`: `191,328` rows per symbol, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T23:45:00Z`;
  - `1h`: `47,832` rows per symbol, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T23:00:00Z`;
  - every selected coverage row had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`.
- CSV line counts including headers:
  - coverage `7`
  - selection `3`
  - signals `10,827`
  - sources `4`
  - summary `97`
  - trades `7,183`
  - common summary `13`
- Baseline outcomes:
  - `breakout_retest_acceptance_15m_all_h12`: aggregate full `5,825`
    trades, net P&L `-2329.18`, PF `0.6778`; all period splits were negative
    after costs, with `2023_2024_oos=-860.87` and
    `2025_2026_recent=-212.23`;
  - `breakout_retest_acceptance_1h_all_h12`: aggregate full `1,357` trades,
    net P&L `-604.03`, PF `0.8652`; all period splits were negative after
    costs, with `2023_2024_oos=-285.08` and
    `2025_2026_recent=-220.25`;
  - full-period symbol rows were negative for every symbol in both candidates
    (`15m`: BTC `-828.60`, ETH `-639.03`, SOL `-861.54`; `1h`: BTC
    `-361.79`, ETH `-5.18`, SOL `-237.06`).
- Review outcome:
  - source and resample validation passed;
  - enough trades existed, but both fixed candidates failed positive net,
    PF `1.2`, OOS/recent split, and symbol-transfer gates;
  - no optimization, robustness review, BTCUSDT promotion, paper/testnet/live
    path, exchange API, deployment, data download, martingale, averaging down,
    or two-exchange work is approved from this baseline.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a review-only stop brief requiring
  user input before any new implementation premise.
- Added a durable decision blocking promotion and automatic retuning of this
  fixed-rule breakout-retest/acceptance baseline.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-breakout-retest-acceptance-baseline-backtest -out-dir results/futures-range-universe-breakout-retest-acceptance-baseline-backtest`
  - `wc -l results/futures-range-universe-breakout-retest-acceptance-baseline-backtest/*.csv`

Futures range universe post structured compression pivot review:

- Review doc:
  `docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md`.
- Stop state:
  `post_structured_compression_pivot_ready_for_breakout_retest_acceptance_baseline`.
- This was review-only. No strategy, optimizer, replay, result rerun, data
  download, live/paper/testnet path, exchange API, credential, deploy file,
  martingale, averaging down, or two-exchange logic was added.
- The structured-compression walk-forward result is now exclusion evidence for
  candidate packaging:
  - fold `wf_2021_2022_train__2023_2024_test` selected
    `sc4h_btc_diagnostic_eth_sol_cw2_h12_t0_75_sb0_10`, but selected test net
    P&L was `92.68` versus frozen test net P&L `229.02`;
  - fold `wf_2021_2024_train__2025_2026_test` selected no config because the
    combined train period had `97` authority trades, below the `100`
    multi-split train gate;
  - fold `wf_2023_2024_train__2025_2026_test` selected the exact frozen
    config and passed, but this was only one of three folds.
- Artifact sanity check, without rerunning walk-forward:
  `wc -l results/futures-range-universe-structured-compression-walk-forward-robustness/*.csv`
  matched the documented line counts: coverage `4`, folds `4`, grid `217`,
  rankings `649`, sources `4`, walk-forward summary `9,505`,
  walk-forward trades `35,881`, common summary `13`, total `46,277`.
- Source/resample facts carried forward from the walk-forward review:
  BTCUSDT, ETHUSDT, and SOLUSDT each had `573,984` accepted Binance USDT-M
  futures `5m` candles from `2021-01-01T00:00:00Z` through
  `2026-06-16T23:55:00Z`; closed UTC `4h` resamples had `11,958` rows per
  symbol through `2026-06-16T20:00:00Z`.
- Inventory outcome:
  - the `4h` ETH/SOL structured-compression stream is not package-ready;
  - the failed `1h` structured-compression surface stays closed;
  - boundary touch rejection, single-candle wick rejection, failed breakout
    re-entry, and mature balance persistence still have no baseline approval
    from current docs;
  - `breakout_retest_acceptance` remains the only open range-universe premise
    because the universe discovery review recorded `14` passing rows as
    secondary evidence.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` toward a bounded offline
  breakout-retest/acceptance fixed-rule baseline selected from existing
  range-universe discovery artifacts.
- Added a durable decision authorizing only that next automatic premise and
  blocking structured-compression rescue retunes.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended review doc, README, and memory
    changes before commit.

Futures range universe structured compression walk-forward robustness:

- Review doc:
  `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md`.
- Stop state:
  `structured_compression_walk_forward_fragile_needs_review`.
- Added explicit offline CLI flag:
  `-futures-range-universe-structured-compression-walk-forward-robustness`.
- This walk-forward pass reused the declared `4h` structured-compression grid
  only: confirmation window `2,3,4`, max hold `4,6,8,12`, target multiple
  `0.75,1.0,1.25`, stop buffer `0.0,0.10`, and symbol sets `BTC_ETH_SOL`,
  `ETH_SOL`, and `BTC_DIAGNOSTIC_ETH_SOL`. It did not add new grid dimensions,
  new symbols, new timeframes, new candidate families, live/paper/testnet
  wiring, exchange API use, credentials, deployment files, data downloads,
  martingale, averaging down, or two-exchange logic.
- Result dir:
  `results/futures-range-universe-structured-compression-walk-forward-robustness/`.
- Outputs:
  - `futures_range_universe_structured_compression_walk_forward_sources.csv/json`
  - `futures_range_universe_structured_compression_walk_forward_coverage.csv/json`
  - `futures_range_universe_structured_compression_walk_forward_grid.csv/json`
  - `futures_range_universe_structured_compression_walk_forward_folds.csv/json`
  - `futures_range_universe_structured_compression_walk_forward_trades.csv/json`
  - `futures_range_universe_structured_compression_walk_forward_summary.csv/json`
  - `futures_range_universe_structured_compression_walk_forward_rankings.csv/json`
  - common `source_manifest.json`, `summary.csv/json`, and `trades.json`
- Source facts:
  - each of `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` loaded `573,984` Binance
    USDT-M futures `5m` candles from `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`;
  - gaps / duplicates were `0` / `0` for every symbol;
  - zero-volume counts: `BTCUSDT=66`, `ETHUSDT=47`, `SOLUSDT=47`;
  - physical non-monotonic counts: `BTCUSDT=0`, `ETHUSDT=0`,
    `SOLUSDT=1`; SOLUSDT was sorted and accepted.
- Closed UTC `4h` resamples:
  - `11,958` rows per symbol, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T20:00:00Z`;
  - every coverage row had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`.
- CSV line counts including headers:
  - coverage `4`
  - folds `4`
  - grid `217`
  - rankings `649`
  - sources `4`
  - walk-forward summary `9,505`
  - walk-forward trades `35,881`
  - common summary `13`
- Normal common outputs stayed frozen to the ETH/SOL authority replay:
  `trades.json` has `129` trades, and `summary.csv/json` matches the replay
  authority stream. BTCUSDT appears only in walk-forward-specific diagnostic
  or historical comparison artifacts.
- Fold results:
  - `wf_2021_2022_train__2023_2024_test` selected
    `sc4h_btc_diagnostic_eth_sol_cw2_h12_t0_75_sb0_10`; selected test net P&L
    was `92.68` versus frozen test net P&L `229.02`, so the fold failed with
    `selected_test_net_worse_than_frozen`;
  - `wf_2021_2024_train__2025_2026_test` selected no config because the
    combined train period had `97` authority trades, below the `100` aggregate
    multi-split train gate; frozen test net P&L was `193.06`;
  - `wf_2023_2024_train__2025_2026_test` selected the exact frozen config and
    passed, with `32` test trades, net P&L `193.06`, PF `1.9121`.
- Review outcome:
  - source/resample validation passed;
  - the frozen config was present in every fold;
  - only one of three folds selected the exact frozen config and passed;
  - no fold required BTCUSDT authority to pass;
  - the result is fragile and does not authorize candidate packaging,
    retuning, BTCUSDT promotion, or live-adjacent work.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` toward a post-fragility hypothesis
  pivot brief rather than another structured-compression retune.
- Added a durable decision that the walk-forward result blocks candidate
  strategy packaging for this frozen ETH/SOL structured-compression stream.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-walk-forward-robustness -out-dir results/futures-range-universe-structured-compression-walk-forward-robustness`
  - `wc -l results/futures-range-universe-structured-compression-walk-forward-robustness/*.csv`

Futures range universe structured compression strategy replay:

- Review doc:
  `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_REPLAY_REVIEW.md`.
- Stop state:
  `structured_compression_strategy_replay_passed_needs_walk_forward_robustness_brief`.
- Added explicit offline CLI flag:
  `-futures-range-universe-structured-compression-strategy-replay`.
- This replay is fixed to
  `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`. It does not add a new
  grid, new symbols, new timeframes, new candidate families, scoring search,
  optimization, live/paper/testnet wiring, exchange API use, credentials,
  deployment files, data downloads, martingale, averaging down, or
  two-exchange logic.
- Result dir:
  `results/futures-range-universe-structured-compression-strategy-replay/`.
- Outputs:
  - `futures_range_universe_structured_compression_strategy_sources.csv/json`
  - `futures_range_universe_structured_compression_strategy_coverage.csv/json`
  - `futures_range_universe_structured_compression_strategy_signals.csv/json`
  - `futures_range_universe_structured_compression_strategy_trades.csv/json`
  - `futures_range_universe_structured_compression_strategy_summary.csv/json`
  - common `source_manifest.json`, `summary.csv/json`, and `trades.json`
- Source facts:
  - each of `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` loaded `573,984` Binance
    USDT-M futures `5m` candles from `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`;
  - gaps / duplicates were `0` / `0` for every symbol;
  - zero-volume counts: `BTCUSDT=66`, `ETHUSDT=47`, `SOLUSDT=47`;
  - physical non-monotonic counts: `BTCUSDT=0`, `ETHUSDT=0`,
    `SOLUSDT=1`; SOLUSDT was sorted and accepted.
- Closed UTC `4h` resamples:
  - `11,958` rows per symbol, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T20:00:00Z`;
  - every coverage row had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`.
- CSV line counts including headers:
  - coverage `4`
  - signals `186`
  - sources `4`
  - strategy summary `49`
  - strategy trades `185`
  - common summary `13`
- Replay result:
  - strategy-specific rows: `185` signals and `184` trades;
  - common authority outputs: `129` ETH/SOL trades only;
  - full authority gross P&L `641.05`, net P&L `573.87`, PF `1.8089`,
    max drawdown `9.82%`;
  - `2021_2022_stress`: `54` trades, net P&L `151.79`, PF `1.4867`;
  - `2023_2024_oos`: `43` trades, net P&L `229.02`, PF `2.2318`;
  - `2025_2026_recent`: `32` trades, net P&L `193.06`, PF `1.9121`;
  - long side: `69` trades, net P&L `385.79`, PF `2.0488`;
  - short side: `60` trades, net P&L `188.08`, PF `1.5506`.
- Symbol facts:
  - `BTCUSDT` diagnostic-only remained negative: `55` trades, net P&L
    `-100.67`, PF `0.6507`;
  - `ETHUSDT` authority full sample: `70` trades, net P&L `351.81`,
    PF `1.9044`, but recent split remained negative;
  - `SOLUSDT` authority full sample: `59` trades, net P&L `222.06`,
    PF `1.6930`, but stress split remained negative.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement a bounded offline
  walk-forward robustness pass for the structured-compression `4h` stream.
- Added a durable decision that the replay authorizes walk-forward robustness
  only; it still does not authorize live, paper/testnet, deploy, new symbols,
  new grid dimensions, or BTCUSDT promotion.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-strategy-replay -out-dir results/futures-range-universe-structured-compression-strategy-replay`
  - `wc -l results/futures-range-universe-structured-compression-strategy-replay/*.csv`

Futures range universe structured compression strategy spec:

- Spec doc:
  `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_SPEC.md`.
- Stop state:
  `structured_compression_strategy_spec_ready_for_offline_replay`.
- This was a docs/memory-only milestone. No code, CLI flags, source files,
  generated results, optimization grids, entries, exits, sizing changes,
  live/paper/testnet wiring, exchange API use, deployment files, data
  downloads, broad symbol mining, martingale, averaging down, or two-exchange
  logic were added.
- Frozen selected config:
  `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`.
- Authority symbols: `ETHUSDT,SOLUSDT`.
- Diagnostic-only symbol: `BTCUSDT`.
- Frozen rules:
  - closed UTC `4h` resampling from accepted local `5m` Binance USDT-M futures
    parents;
  - detector `p30_c12_bollinger_on_adx_off`;
  - first closed breakout within `24` closed `4h` bars after a completed
    mature range;
  - closed confirmation within `2` bars;
  - next `4h` bar open entry;
  - target `1.0` completed range width from slipped entry;
  - stop buffer `0.0` range width;
  - max hold `12` closed `4h` bars.
- Source facts carried forward:
  - each of `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` loaded `573,984` Binance
    USDT-M futures `5m` candles from `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`;
  - gaps / duplicates were `0` / `0` for every symbol;
  - zero-volume counts: `BTCUSDT=66`, `ETHUSDT=47`, `SOLUSDT=47`;
  - physical non-monotonic counts: `BTCUSDT=0`, `ETHUSDT=0`,
    `SOLUSDT=1`;
  - SOLUSDT may be sorted before downstream validation and accepted only when
    the sorted stream has no gaps, duplicates, or invalid OHLCV.
- Optimization facts preserved in the spec:
  - selected authority result had `129` trades, gross P&L `641.05`, net P&L
    `573.87`, PF `1.8089`, max drawdown `9.82%`, average net R `0.3465`;
  - stress/OOS/recent authority splits were positive after costs;
  - BTCUSDT diagnostic-only remained negative at `55` trades, net P&L
    `-100.67`, PF `0.6507`;
  - ETHUSDT recent and SOLUSDT stress split caveats remain explicit.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement the exact offline
  candidate strategy replay/backtest behind
  `-futures-range-universe-structured-compression-strategy-replay`.
- Added a durable decision that the strategy spec freezes the selected
  ETH/SOL authority configuration for replay; any replay mismatch, BTC
  promotion, or grid-style retuning must stop for review.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended docs/memory changes before
    commit.

Futures range universe structured compression optimization:

- Review doc:
  `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md`.
- Stop state:
  `structured_compression_optimization_passed_needs_strategy_spec`.
- Added explicit offline CLI flag:
  `-futures-range-universe-structured-compression-optimization`.
- This optimization is bounded to the passing `4h` structured-compression
  universe stream only. The failed `1h` structured-compression surface was not
  optimized. Default runs still use `lab.EmptyStrategy`; this run writes common
  `trades.json` and `summary.*` only for the selected top configuration.
- Source facts:
  - each of `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` loaded `573,984` Binance
    USDT-M futures `5m` candles from `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`;
  - gaps / duplicates were `0` / `0` for every symbol;
  - zero-volume counts: `BTCUSDT=66`, `ETHUSDT=47`, `SOLUSDT=47`;
  - physical non-monotonic counts: `BTCUSDT=0`, `ETHUSDT=0`,
    `SOLUSDT=1`;
  - the SOL source was sorted before downstream validation and accepted.
- Closed UTC `4h` resamples:
  - `11,958` rows per symbol, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T20:00:00Z`;
  - all coverage rows had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`.
- Result dir:
  `results/futures-range-universe-structured-compression-optimization/`.
- Outputs:
  - `futures_range_universe_structured_compression_optimization_sources.csv/json`
  - `futures_range_universe_structured_compression_optimization_coverage.csv/json`
  - `futures_range_universe_structured_compression_optimization_grid.csv/json`
  - `futures_range_universe_structured_compression_optimization_trades.csv/json`
  - `futures_range_universe_structured_compression_optimization_summary.csv/json`
  - `futures_range_universe_structured_compression_optimization_rankings.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - coverage `4`
  - grid `217`
  - rankings `217`
  - sources `4`
  - optimization summary `9,505`
  - optimization trades `35,881`
  - common summary `13`
- Grid/result facts:
  - grid size: `216` configurations;
  - passing configurations: `115`;
  - selected config:
    `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`;
  - selected authority symbols: `ETHUSDT,SOLUSDT`;
  - selected diagnostic symbol: `BTCUSDT`;
  - selected config settings: confirmation window `2`, max hold `12`,
    target multiple `1.0`, stop buffer `0.0`.
- Selected authority result:
  - `129` trades, gross P&L `641.05`, net P&L `573.87`, PF `1.8089`,
    max drawdown `9.82%`, average net R `0.3465`;
  - `2021_2022_stress`: `54` trades, net P&L `151.79`, PF `1.4867`;
  - `2023_2024_oos`: `43` trades, net P&L `229.02`, PF `2.2318`;
  - `2025_2026_recent`: `32` trades, net P&L `193.06`, PF `1.9121`;
  - long side: `69` trades, net P&L `385.79`, PF `2.0488`;
  - short side: `60` trades, net P&L `188.08`, PF `1.5506`.
- Symbol facts:
  - `BTCUSDT` diagnostic-only remained negative: `55` trades, net P&L
    `-100.67`, PF `0.6507`;
  - `ETHUSDT` authority full sample: `70` trades, net P&L `351.81`,
    PF `1.9044`, but recent split was negative after costs;
  - `SOLUSDT` authority full sample: `59` trades, net P&L `222.06`,
    PF `1.6930`, but stress split was negative after costs.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to create a first offline candidate
  strategy spec for the selected ETH/SOL `4h` structured-compression stream,
  with BTC diagnostic-only and no live-adjacent work.
- Added a durable decision that this optimization authorizes an offline
  ETH/SOL universe strategy spec only. It does not authorize BTC strategy
  promotion or live/paper/testnet/deploy work.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-optimization -out-dir results/futures-range-universe-structured-compression-optimization`
  - `wc -l results/futures-range-universe-structured-compression-optimization/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended implementation, doc, and memory
    changes before commit.

Futures range universe structured compression baseline backtest:

- Review doc:
  `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_BASELINE_REVIEW.md`.
- Stop state: `structured_compression_baseline_passed_needs_optimization_brief`.
- Added explicit offline CLI flag:
  `-futures-range-universe-structured-compression-baseline-backtest`.
- Default runs still use `lab.EmptyStrategy`; this baseline runs only when the
  explicit flag is passed and rejects spot/comparison sources through the
  futures source guard. Universe sources remain limited to local BTCUSDT,
  ETHUSDT, and SOLUSDT Binance USDT-M futures files.
- Normal source manifest:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - product: Binance USDT-M futures
  - symbol / interval: `BTCUSDT` / `5m`
  - row count: `573,984`
  - first / last open: `2021-01-01T00:00:00Z` /
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Universe source facts:
  - each of `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` loaded `573,984` candles from
    `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  - gaps / duplicates were `0` / `0` for every symbol;
  - zero-volume counts: `BTCUSDT=66`, `ETHUSDT=47`, `SOLUSDT=47`;
  - physical non-monotonic counts: `BTCUSDT=0`, `ETHUSDT=0`,
    `SOLUSDT=1`;
  - the SOL source was sorted before downstream validation and accepted.
- Closed UTC resamples used:
  - `1h`: `47,832` rows per symbol, last open `2026-06-16T23:00:00Z`
  - `4h`: `11,958` rows per symbol, last open `2026-06-16T20:00:00Z`
  - all coverage rows had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`.
- Result dir:
  `results/futures-range-universe-structured-compression-baseline-backtest/`.
- Outputs:
  - `futures_range_universe_structured_compression_baseline_sources.csv/json`
  - `futures_range_universe_structured_compression_baseline_coverage.csv/json`
  - `futures_range_universe_structured_compression_baseline_signals.csv/json`
  - `futures_range_universe_structured_compression_baseline_trades.csv/json`
  - `futures_range_universe_structured_compression_baseline_summary.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - coverage `7`
  - signals `926`
  - sources `4`
  - baseline summary `97`
  - baseline trades `913`
  - common summary `13`
- Baseline result:
  - `925` signal rows, `912` executed trades, `6` skipped signals;
  - `structured_compression_4h_all_h6`: `186` signals, `185` trades,
    gross P&L `395.68`, net P&L `302.92`, PF `1.3598`, max drawdown
    `12.92%`;
  - `structured_compression_1h_all_h12`: `739` signals, `727` trades,
    gross P&L `239.01`, net P&L `-200.13`, PF `0.9362`, max drawdown
    `35.46%`;
  - common combined compatibility view: `912` trades, gross P&L `634.69`,
    net P&L `102.79`, PF `1.0258`.
- Important verdict facts:
  - the `4h` aggregate passed after costs and in `2023_2024_oos` plus
    `2025_2026_recent`;
  - the `4h` aggregate `2021_2022_stress` split was slightly negative after
    costs at net P&L `-4.82`;
  - `BTCUSDT` was negative on the `4h` surface: `56` trades, net P&L
    `-92.27`, PF `0.6404`;
  - `ETHUSDT` and `SOLUSDT` carried the `4h` result: net P&L `243.61` and
    `151.58`;
  - the `1h` surface failed and is not promoted.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement a bounded offline
  optimization/robustness pass for `4h structured_compression_4h_all_h6` only,
  with BTC weakness and symbol inclusion constraints explicit.
- Added a durable decision that the `4h` aggregate is authorized for bounded
  offline optimization only, while `1h` is closed from this baseline.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-baseline-backtest -out-dir results/futures-range-universe-structured-compression-baseline-backtest`
  - `wc -l results/futures-range-universe-structured-compression-baseline-backtest/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended implementation, doc, and memory
    changes before commit.

Futures range universe discovery audit:

- Review doc: `docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md`.
- Stop state: `range_universe_audit_ready_for_baseline_backtest`.
- Added explicit non-trading CLI flag:
  `-futures-range-universe-discovery-audit`.
- Default runs still use `lab.EmptyStrategy`; this audit writes source,
  coverage, candidate, summary, ranking, and stability artifacts only and
  produced `0` trades.
- V1 local universe:
  - `BTCUSDT`:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - `ETHUSDT`:
    `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`
  - `SOLUSDT`:
    `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`
- Source facts:
  - each source loaded `573,984` candles from `2021-01-01T00:00:00Z`
    through `2026-06-16T23:55:00Z`;
  - gaps / duplicates were `0` / `0` for every symbol;
  - zero-volume counts: `BTCUSDT=66`, `ETHUSDT=47`, `SOLUSDT=47`;
  - physical non-monotonic counts: `BTCUSDT=0`, `ETHUSDT=0`,
    `SOLUSDT=1`;
  - the SOL source was explicitly sorted before downstream validation, then
    accepted as monotonic with no gaps or duplicates.
- Resample coverage per symbol:
  - `5m`: `573,984` rows, last open `2026-06-16T23:55:00Z`
  - `15m`: `191,328` rows, last open `2026-06-16T23:45:00Z`
  - `1h`: `47,832` rows, last open `2026-06-16T23:00:00Z`
  - `4h`: `11,958` rows, last open `2026-06-16T20:00:00Z`
  - every resample had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`.
- Result dir:
  `results/futures-range-universe-discovery-audit/`.
- Outputs:
  - `futures_range_universe_sources.csv/json`
  - `futures_range_universe_coverage.csv/json`
  - `futures_range_universe_candidates.csv/json`
  - `futures_range_universe_summary.csv/json`
  - `futures_range_universe_rankings.csv/json`
  - `futures_range_universe_stability.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - candidates `415,954`
  - coverage `13`
  - rankings `148`
  - sources `4`
  - stability `442`
  - universe summary `1,765`
  - common summary `13`
- Audit result:
  - `415,953` candidate rows;
  - candidate rows by symbol: `BTCUSDT=134,379`, `ETHUSDT=138,324`,
    `SOLUSDT=143,250`;
  - `147` ranked universe surfaces;
  - `35` passing rows: `21` structured compression expansion and `14`
    breakout retest/acceptance rows.
- Top next baseline surfaces:
  - `4h structured_compression_expansion all h6`
  - `1h structured_compression_expansion all h12`
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement a fixed-rule offline
  structured-compression universe baseline backtest for those two surfaces.
- Added a durable decision that only those two top structured-compression
  surfaces are authorized for the next baseline from this audit; no
  optimization or live-adjacent path is approved.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-discovery-audit -out-dir results/futures-range-universe-discovery-audit`
  - `wc -l results/futures-range-universe-discovery-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended implementation, doc, and memory
    changes before commit.

Futures range universe discovery spec:

- Spec doc: `docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_SPEC.md`.
- Stop state: `range_universe_spec_ready_for_audit_implementation`.
- Outcome: the project explicitly broadens from BTCUSDT-only to a local
  Binance USDT-M futures range universe for discovery. The v1 universe is
  limited to local BTCUSDT, ETHUSDT, and SOLUSDT `5m` CSVs under
  `../binance-bot/data/`.
- Motivation: the BTCUSDT clean-breakout baseline had adequate trade counts
  but failed after costs, so the next useful step is source-validated
  cross-symbol range discovery rather than retuning the failed immediate
  breakout template.
- Local source quick checks:
  - `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`: `573,985`
    CSV lines including header, `573,984` rows, min open `1609459200000`,
    max open `1781654100000`, non-monotonic physical row count `0`
  - `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`: `573,985`
    CSV lines including header, `573,984` rows, min open `1609459200000`,
    max open `1781654100000`, non-monotonic physical row count `0`
  - `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`: `573,985`
    CSV lines including header, `573,984` rows, min open `1609459200000`,
    max open `1781654100000`, non-monotonic physical row count `1`
- The SOL caveat is recorded in the spec: the future implementation must
  validate and sort/fail explicitly instead of trusting filename coverage.
- Future candidate families:
  - breakout retest / acceptance after completed mature range;
  - boundary touch / rejection;
  - failed breakout re-entry;
  - mature balance rotation / persistence;
  - compression-to-expansion only if materially reframed from the failed
    immediate clean-breakout baseline.
- Required future outputs are scoped under
  `results/futures-range-universe-discovery-audit/` and remain zero-trade
  common outputs until a later baseline backtest is approved.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement
  `-futures-range-universe-discovery-audit` as a non-trading source/universe
  validation plus discovery audit.
- This was docs/memory-only: no code, CLI flags, generated results, entries,
  exits, scoring, sizing, optimization, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, or sibling
  repo mutation was added.
- Verification passed:
  - `wc -l ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`
  - `awk -F, 'NR>1{if(NR==2||$1<min)min=$1;if($1>max)max=$1;if(NR>2 && $1<=prev) nonmono++; prev=$1} END{print FILENAME, NR-1, min, max, nonmono+0}' <source>`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures clean breakout baseline backtest:

- Review doc: `docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md`.
- Stop state: `clean_breakout_baseline_failed_no_promotion`.
- Added explicit offline CLI flag:
  `-futures-clean-breakout-baseline-backtest`.
- Default runs still use `lab.EmptyStrategy`; this baseline runs only when the
  explicit flag is passed and rejects spot/comparison sources.
- First baseline evaluated candidates independently, per user direction. The
  user-approved portfolio-stream routing rule was checked and did not trigger:
  neither candidate was near-viable after costs.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT `5m`.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Resample coverage used by the baseline:
  - `1h`: `47,832` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:00:00Z`
  - `4h`: `11,958` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T20:00:00Z`
- Result dir:
  `results/futures-clean-breakout-baseline-backtest/`.
- Outputs:
  - `futures_clean_breakout_baseline_signals.csv/json`
  - `futures_clean_breakout_baseline_trades.csv/json`
  - `futures_clean_breakout_baseline_summary.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - signals `4,849`
  - trades `1,186`
  - baseline summary `25`
  - common summary `13`
- Baseline result:
  - `4,848` signal rows
  - `1,185` executed trades
  - `clean_breakout_4h_up_h12`: `496` signals, `121` trades, gross P&L
    `2.55`, net P&L `-69.29`, PF `0.8386`
  - `clean_breakout_1h_all_h12`: `4,352` signals, `47` skipped signals,
    `1,064` trades, gross P&L `272.63`, net P&L `-320.25`, PF `0.8846`
  - aggregate compatibility summary: `1,185` trades, gross P&L `275.18`,
    net P&L `-389.54`, PF `0.8784`
- Failure facts:
  - both candidates had negative full net P&L after costs;
  - both had full PF below `1.2`;
  - both were negative in `2023_2024_oos` and `2025_2026_recent`;
  - the `1h` all-side candidate was negative on both long and short sides.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to request a new futures premise or
  explicit scope choice before another backtest; do not optimize this clean
  breakout baseline.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-clean-breakout-baseline-backtest -out-dir results/futures-clean-breakout-baseline-backtest`
  - `wc -l results/futures-clean-breakout-baseline-backtest/*.csv`

Futures range candidate discovery audit:

- Review doc: `docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md`.
- Stop state: `range_discovery_audit_ready`.
- Added explicit non-trading CLI flag:
  `-futures-range-candidate-discovery-audit`.
- Default runs still use `lab.EmptyStrategy`; this audit writes labels,
  coverage, rankings, summaries, and stability rows only and produced `0`
  trades.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT `5m`.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Resample coverage:
  - `5m`: `573,984` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:55:00Z`
  - `15m`: `191,328` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:45:00Z`
  - `1h`: `47,832` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:00:00Z`
  - `4h`: `11,958` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T20:00:00Z`
  - every timeframe had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`
- Result dir:
  `results/futures-range-candidate-discovery-audit/`.
- Outputs:
  - `futures_range_discovery_coverage.csv/json`
  - `futures_range_discovery_candidates.csv/json`
  - `futures_range_discovery_summary.csv/json`
  - `futures_range_discovery_rankings.csv/json`
  - `futures_range_discovery_stability.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - candidates `144,637`
  - coverage `5`
  - rankings `121`
  - stability `121`
  - discovery summary `481`
  - common summary `13`
- Audit result:
  - `144,636` event/horizon rows
  - `48,212` distinct events
  - `120` ranked surfaces
  - `24` passing rows, all from `clean_breakout_continuation`
  - passing rows by timeframe: `15m=9`, `1h=9`, `4h=6`
- Top non-duplicative candidates for the next baseline:
  - `4h clean_breakout_continuation up h12`: full count `496`, weakest split
    count `112`, full favorable rate `86.69%`, weakest favorable rate
    `82.38%`, weakest cost buffer `4.2910%`
  - `1h clean_breakout_continuation all h12`: full count `4,352`, weakest
    split count `1,213`, full favorable rate `90.37%`, weakest favorable rate
    `89.01%`, weakest cost buffer `2.5907%`
- Rejected for this baseline: boundary touch rejection, single-candle wick
  rejection, failed breakout re-entry, and mature balance persistence.
- Added README docs-index entry for the new review doc.
- Added a durable decision that only clean breakout continuation is authorized
  for the next fixed-rule baseline; the other discovery families are not.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement
  `-futures-clean-breakout-baseline-backtest`.
- No entries, exits, scoring, sizing, optimization, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, spot
  comparison, symbol expansion, or sibling repo mutation was added by this
  audit.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-range-candidate-discovery-audit -out-dir results/futures-range-candidate-discovery-audit`
  - `wc -l results/futures-range-candidate-discovery-audit/*.csv`

Futures range candidate discovery spec:

- Spec doc: `docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_SPEC.md`.
- Stop state: `range_discovery_spec_ready_for_audit_implementation`.
- Outcome: the repo scope is corrected from overly narrow range-only gating to
  range-first broad discovery. `5m`, buy/sell touch, single-candle reactions,
  boundary rejection, failed-break re-entry, and breakout continuation remain
  open when materially reframed and tied to a fast baseline-backtest gate.
- Current lab authority remains Binance USDT-M futures BTCUSDT `5m`:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - required CLI product: `-source-product binance-usdm-futures`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Discovery menu for the next audit:
  - mature balance persistence on `1h`/`4h` with optional `15m` cross-check;
  - boundary touch rejection on `5m`/`15m`/`1h`;
  - single-candle wick rejection on `5m`/`15m`/`1h`;
  - failed breakout re-entry on `5m`/`15m`/`1h`;
  - clean breakout continuation after mature range compression on
    `15m`/`1h`/`4h`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement
  `-futures-range-candidate-discovery-audit` as a non-trading discovery audit.
  If one or two candidates clear the gate, the following brief should be a
  baseline offline backtest, not another review-only loop.
- Added README docs-index entry for the new spec.
- Added a durable decision that closed families are not permanently banned;
  exact failed templates are banned, but materially reframed range-first ideas
  may compete in a broad futures discovery funnel.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, spot
  comparison, symbol expansion, or sibling repo mutation was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures higher-timeframe range source spec:

- Spec doc: `docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md`.
- Stop state: `higher_tf_range_source_spec_no_viable_range_premise`.
- Outcome: the parent source and closed UTC resampling contract are specified,
  but no higher-timeframe audit is ready until the user supplies a materially
  different range premise, closed-candle candidate event, and falsification
  rule.
- Current lab authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - required CLI product: `-source-product binance-usdm-futures`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Candidate higher-timeframe bars from the current parent:
  - `15m`: `3` child bars, expected complete bars `191,328`, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T23:45:00Z`
  - `1h`: `12` child bars, expected complete bars `47,832`, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T23:00:00Z`
  - `4h`: `48` child bars, expected complete bars `11,958`, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T20:00:00Z`
- Resampling contract: UTC bucket start as higher-timeframe open; complete
  buckets only; open first child open, high max child high, low min child low,
  close last child close, volume sum; reject partial, missing-child,
  duplicate, synthetic, forward-filled, spot/comparison, non-UTC, or
  future-looking bars.
- Added README docs-index entry for the new source spec.
- Added a durable decision that higher-timeframe range work must derive from
  the accepted BTCUSDT futures 5m parent source until an explicit source/scope
  change.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` into a premise-request brief, not an
  audit brief.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, or sibling
  repo mutation was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures scope pivot review, range-only:

- Review doc: `docs/FUTURES_SCOPE_PIVOT_REVIEW.md`.
- Stop state: `range_scope_pivot_ready_for_higher_timeframe_source_spec`.
- Decision: pause BTCUSDT 5m range continuation work. The next range-only
  route is a BTCUSDT higher-timeframe futures source/premise spec, not an audit
  or prototype.
- Added README docs-index entry for the new review doc.
- Added a durable decision that BTCUSDT higher-timeframe range-source work is
  the next authorized lane; BTC/ETH expansion remains deferred until that
  source/premise review is complete or the user explicitly changes scope.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to create
  `docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md`.
- Current lab authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Sibling repo facts used as context only:
  - `~/binance-bot` BTCUSDT 5m range sleeve failed on legacy spot data; caution
    only, not futures authority.
  - `~/binance-bot` cross-exchange spread and broad selector work remain
    exclusion evidence for this lab.
  - `~/binance-bot` daily ATR contraction is non-range trend/volatility context
    and cannot promote range work.
  - `~/crypto-trading-bot` BTC/ETH USD-M source-contract discipline is process
    context only.
  - `~/crypto-trading-bot` BTC/ETH 4h range-reversal/re-entry is closed
    exclusion evidence, not a port target.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, or sibling repo mutation was
  added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures scope pivot review spec, range-only:

- Spec doc: `docs/FUTURES_SCOPE_PIVOT_REVIEW_SPEC.md`.
- Stop state: `range_scope_pivot_spec_ready`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to make the next task a docs-only
  `docs/FUTURES_SCOPE_PIVOT_REVIEW.md`, focused on range strategy scope only.
- Added README docs-index entry for the new spec.
- Added a durable decision: a scope pivot remains range-strategy-only unless
  the user explicitly changes the project objective. Higher-timeframe,
  multi-symbol, or sibling-repo context may be used only as range-source,
  range-premise, process, or exclusion evidence.
- Current lab authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Sibling repo facts recorded for next-review context only:
  - `~/binance-bot` range seed: Binance-only BTCUSDT 5m range sleeve on
    legacy spot data failed even gross; use only as caution.
  - `~/binance-bot` cross-exchange spread seed failed local gates and remains
    out of scope because this lab forbids two-exchange execution.
  - `~/binance-bot` multi-pair selector/alt overlay work is no-go evidence for
    broad symbol selection.
  - `~/binance-bot` daily ATR contraction is positive-looking but non-range
    trend/volatility work, so it is outside this range-only pivot.
  - `~/crypto-trading-bot` BTC/ETH 4h range-reversal/re-entry workbench failed
    as a frozen family; use only as exclusion evidence.
  - `~/crypto-trading-bot` BTC/ETH USD-M source-contract discipline may be
    useful if the review selects higher-timeframe or BTC/ETH range-source
    specification.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, or two-exchange logic was added.
- Sibling repo status note:
  - `~/binance-bot` was clean during inspection.
  - `~/crypto-trading-bot` had pre-existing dirty docs/memory/research files
    and was treated as read-only.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

## 2026-06-25

Futures impulse absorption audit:

- Review doc: `docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md`.
- Stop state: `impulse_absorption_no_viable_edge`.
- Added explicit non-trading CLI flag:
  `-futures-impulse-absorption-audit`.
- Default runs still use `lab.EmptyStrategy`; this audit writes labels and
  summaries only and produced `0` trades.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT 5m.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Result dir:
  `results/futures-impulse-absorption-audit/`.
- Outputs:
  - `futures_impulse_absorption_candidates.csv/json`
  - `futures_impulse_absorption_summary.csv/json`
  - `futures_impulse_absorption_stability.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - candidates `13,737`
  - summary `949`
  - stability `241`
  - common summary `13`
- Audit result:
  - `3,434` distinct impulse events
  - `13,736` event/horizon candidate rows
  - `1,838` up impulse events and `1,596` down impulse events
  - no missing future windows
  - all period splits have at least `1,003` all-direction/all-bucket events
    at every horizon
- Gate failure:
  - midpoint reclaim-first never beats continuation-first in any period split
    at horizons `3`, `6`, `12`, or `24`
  - all-direction/all-bucket continuation-first rates are roughly
    `57.74%` to `59.58%`
  - reclaim-first rates are only roughly `24.13%` to `25.43%`
  - maximum quick-continuation rate is `80.76%`
  - minimum reclaim-minus-continuation margin stays negative at every horizon
- No entry, exit, scoring, sizing, strategy replacement, paper/testnet/live
  wiring, exchange API use, deploy files, grid, martingale, averaging down, or
  two-exchange logic was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-impulse-absorption-audit -out-dir results/futures-impulse-absorption-audit`
  - `wc -l results/futures-impulse-absorption-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Futures impulse absorption next-brief refresh:

- Stop state: `impulse_absorption_audit_ready`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` from the open-ended
  new-hypothesis prompt into a scoped non-trading audit brief for
  `-futures-impulse-absorption-audit`.
- Chosen premise: after a liquidation-like abnormal Binance USDT-M futures 5m
  impulse candle, BTCUSDT may show absorption by reclaiming the event candle
  midpoint before extending beyond the event extreme.
- Candidate event definition for the next audit:
  - closed BTCUSDT 5m futures candle after a `30` day prior rolling warmup
  - true-range percentile rank at least `p99`
  - volume percentile rank at least `p95`
  - close position `<=0.25` for down impulses or `>=0.75` for up impulses
- Required next-audit source remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - expected accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- This was docs/memory-only handoff work: no code, CLI flags, audits, result
  directories, entries, exits, scoring, sizing, paper/testnet/live wiring,
  exchange API use, deploy files, grid, martingale, averaging down, or
  two-exchange logic was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` before commit showed only intended memory changes.

Futures hypothesis pivot inventory:

- Inventory doc: `docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md`.
- Stop state: `pivot_inventory_needs_user_hypothesis`.
- Outcome: reviewed material is now an exclusion map plus reusable
  infrastructure, not a new strategy queue.
- Source authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Classification:
  - SR timing and compression breakout are closed legacy spot-only evidence.
  - Range durability is diagnostic/infrastructure, not an entry surface.
  - Detector durability/context refinement is reusable infrastructure; futures
    review revalidated context shape only, not entry promotion.
  - Hold-inside directional edge is closed.
  - Hold-inside midline transition/reaction is diagnostic after the failed
    prototype; do not broaden into `hold_6_inside`, `mid_close_across`, side
    cohorts, or old spot authority.
  - Futures midline touch prototype is a futures-authoritative failure.
- No code, CLI flags, audits, result directories, entry/exit logic, scoring,
  sizing, paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, or two-exchange logic was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` before commit showed only intended doc/memory changes.

Minimal futures midline touch prototype:

- Review doc: `docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md`.
- Stop state: `prototype_failed_no_promotion`.
- Added explicit offline CLI flag:
  `-hold-inside-midline-touch-prototype`.
- Default runs still use `lab.EmptyStrategy`; the prototype runs only when the
  flag is passed.
- Prototype source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT 5m.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Result dir:
  `results/futures-hold-inside-midline-touch-prototype/`.
- Prototype outputs:
  - `hold_inside_midline_touch_prototype_signals.csv/json`
  - `hold_inside_midline_touch_prototype_trades.csv/json`
  - `hold_inside_midline_touch_prototype_summary.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - signals `533`
  - prototype trades `532`
  - prototype summary `13`
  - common summary `13`
- Run output:
  - `signal_rows=532`
  - `trades=531`
  - `summary_rows=12`
  - one exact-mid skip
  - exit reasons: `140` stop losses, `82` take profits, `309` time stops
- Full-sample result:
  - `531` trades, win rate `29.94%`
  - gross P&L `-95.54`, net P&L `-418.99`
  - profit factor `0.3409`, average net R `-0.4276`
  - max drawdown `42.11%`
- Period splits all failed:
  - `2021_2022_stress`: net P&L `-226.03`, PF `0.3824`
  - `2023_2024_oos`: net P&L `-116.58`, PF `0.3165`
  - `2025_2026_recent`: net P&L `-76.38`, PF `0.2298`
- Side splits both failed:
  - long: `248` trades, net P&L `-234.50`, PF `0.2805`
  - short: `283` trades, net P&L `-184.49`, PF `0.4045`
- Added engine entry-geometry guard so next-bar entries are skipped when the
  stop/target would be on the wrong side of the actual entry.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -hold-inside-midline-touch-prototype -out-dir results/futures-hold-inside-midline-touch-prototype`
  - `wc -l results/futures-hold-inside-midline-touch-prototype/*.csv`

Futures data impact review:

- Review doc: `docs/FUTURES_DATA_IMPACT_REVIEW.md`.
- Stop state: `futures_reaction_gate_passed_needs_minimal_entry_brief`.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT 5m.
- Source manifests in all futures result dirs:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Result dirs:
  - `results/futures-detector-context-refinement-audit/`
  - `results/futures-hold-inside-midline-transition-audit/`
  - `results/futures-hold-inside-midline-reaction-audit/`
- Audit sizes:
  - detector context refinement: `117,848` candidate rows, `640` summary rows,
    `160` stability rows
  - hold-inside midline transition: `8,600` candidate rows, `672` summary rows,
    `168` stability rows
  - hold-inside midline reaction: `10,172` candidate rows, `24` funnel rows,
    `1,240` summary rows, `336` stability rows
- Futures reaction gate:
  - `hold_3_inside + mid_touch`: weakest split `129` event candidates,
    `55.60%` minimum event rate, `44.40%` maximum missing-event rate
  - h6 all-bucket: weakest split `129` candidates, `56.59%` minimum
    close-back, `48.06%` minimum mid-rejection before boundary, `29.46%`
    maximum boundary-before-rejection, `18.60%` maximum quick invalidation,
    `22.48%` maximum trend
  - h6 `mid_50`: weakest split `114` candidates, `57.89%` minimum close-back,
    `52.63%` minimum mid-rejection before boundary, `23.68%` maximum
    boundary-before-rejection, `13.16%` maximum quick invalidation, `21.93%`
    maximum trend
- Revalidated for the next task: first minimal offline prototype around
  `hold_3_inside` + first `mid_touch` within `12` bars + event close-position
  bucket `mid_50`.
- Still diagnostic only: old spot approval as evidence, `hold_6_inside`,
  `mid_close_across`, side-specific cohorts, and `hold_3_inside_mid_50`.
- No entry, exit, scoring, sizing, strategy replacement, paper, testnet, live,
  exchange-key, deploy, grid, martingale, averaging-down, or two-exchange work
  was added. The runs reported `strategy=empty trades=0`.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -detector-context-refinement-audit -out-dir results/futures-detector-context-refinement-audit`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -hold-inside-midline-transition-audit -out-dir results/futures-hold-inside-midline-transition-audit`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -hold-inside-midline-reaction-audit -out-dir results/futures-hold-inside-midline-reaction-audit`
  - `wc -l results/futures-*/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Futures source guard and next-brief refresh:

- Added CLI source enforcement:
  - default CSV:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - valid `-source-product` values: `binance-usdm-futures` and `binance-spot`
  - non-default CSV paths must pass `-source-product`
  - spot paths require `-source-product binance-spot` plus
    `-allow-spot-comparison`
  - accepted runs write `source_manifest.json`
- Added Go-native source validation inspired by `crypto-trading-bot` source
  contract discipline, without importing its Python helpers or strategy
  evidence.
- Real futures smoke manifest:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - market type: Binance USDT-M futures
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - close-time end from smoke output: `2026-06-16T23:59:59Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Updated `README.md`, `docs/VERIFICATION.md`, `memory/DECISIONS.md`, and
  `memory/NEXT_CODEX_BRIEF.md`. The next brief now starts from the full-history
  futures CSV and asks for `docs/FUTURES_DATA_IMPACT_REVIEW.md` plus the three
  paused futures audit reruns.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -out-dir results/source-guard-smoke`
  - `wc -l results/source-guard-smoke/*.json results/source-guard-smoke/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

## 2026-06-16

Futures data-source correction:

- User clarified the real trading target is Binance futures rather than spot.
- Previous audit/review outputs used the spot path
  `../binance-bot/data/btcusdt_spot_5m_2021_2026.csv`; treat all promotion
  implications from those results as suspended until futures revalidation.
- Local sibling data currently visible:
  - `../binance-bot/data/btcusdt_spot_5m_2021_2026.csv`: `573,697` CSV lines
    including header, spanning `2021-01-01T00:00:00Z` through
    `2026-06-15T23:59:59Z`
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-13_2026-06-15.csv`:
    `865` CSV lines including header
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-14_2026-06-15.csv`:
    `577` CSV lines including header
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-15.csv`: `289` CSV
    lines including header
- Updated memory and docs so the next brief is an impact assessment, not entry
  implementation. If full-history futures data is unavailable, the next verdict
  should be "data gap first," not a prototype.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Memory context-budget wording adjustment:

- Updated `AGENTS.md`, `memory/README.md`, and `memory/DECISIONS.md` so the
  size guidance for all always-read memory files is a soft `300-350` line
  judgment band, not a hard threshold. Compact or split memory when it feels
  bulky or repetitive, not merely because one file crosses `300` lines.

## 2026-06-15

Hold-inside midline reaction review milestone:

- Review doc: `docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md`.
- Inputs: `results/hold-inside-midline-reaction-audit/`.
- Original spot-data verdict: one first minimal offline entry prototype was
  justified, but only for `hold_3_inside` + first `mid_touch` within `12` bars
  + event close-position bucket `mid_50`.
- Current status after futures correction: suspended pending futures rerun and
  impact review.
- Not promoted: live use, strategy promotion, broad detector-family entries,
  `hold_6_inside`, `mid_close_across`, side-specific cohorts, and
  `hold_3_inside_mid_50`.
- Key evidence:
  - funnel pass: `hold_3_inside` + `mid_touch` has weakest-split event rate
    `52.25%` with `116` event candidates
  - `hold_3_inside` + `mid_touch`, h6 all-bucket: weakest split `116`
    candidates, `55.17%` minimum close-back, `45.69%` minimum mid-rejection
    before boundary, `28.25%` maximum boundary-before-rejection, `18.39%`
    maximum quick invalidation, `22.22%` maximum trend
  - `hold_3_inside` + `mid_touch` + event-close `mid_50`, h6: weakest split
    `104` candidates, `58.65%` minimum close-back, `50.96%` minimum
    mid-rejection before boundary, `22.12%` maximum boundary-before-rejection,
    `11.54%` maximum quick invalidation, `19.23%` maximum trend
- Updated `README.md`, `memory/DECISIONS.md`, and
  `memory/NEXT_CODEX_BRIEF.md` so the next task is the bounded offline
  prototype.
- Verification passed:
  - `wc -l results/hold-inside-midline-reaction-audit/*.csv`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Hold-inside midline reaction audit milestone:

- Added CLI flag `-hold-inside-midline-reaction-audit`.
- Result directory: `results/hold-inside-midline-reaction-audit/`.
- Outputs:
  - `hold_inside_midline_reaction_candidates.csv/json`
  - `hold_inside_midline_reaction_funnel_summary.csv/json`
  - `hold_inside_midline_reaction_summary.csv/json`
  - `hold_inside_midline_reaction_stability.csv/json`
- Audit size:
  - profiles: `1`
  - rules: `3`
  - event types: `2`
  - candidate rows: `9,080`
  - funnel rows: `24`
  - summary rows: `1,296`
  - stability rows: `352`
  - CSV lines including header: `9,081` / `25` / `1,297` / `353`
  - horizons: `1`, `3`, `6`, `12`
  - max midline event delay: `12` bars
  - quick invalidation window: `3` bars
- Scope:
  - profile `p30_c12_bollinger_on_adx_off`
  - rules `hold_3_inside`, `hold_6_inside`, and diagnostic
    `hold_3_inside_mid_50`
  - event types `mid_touch` and `mid_close_across`
  - event candle is the reindexed decision candle; labels start at
    `event_index + 1`
  - no entries, exits, scoring, sizing, paper side, favorable/adverse fields,
    strategy replacement, or live wiring
- Last run loaded `569,451` candles through `2026-06-01T23:59:59Z` and printed
  `strategy=empty trades=0`.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv -hold-inside-midline-reaction-audit -out-dir results/hold-inside-midline-reaction-audit`
  - `wc -l results/hold-inside-midline-reaction-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Context budget cleanup milestone:

- Added a durable context-budget rule to `AGENTS.md`, `memory/README.md`, and
  `memory/DECISIONS.md`.
- Compacted `memory/PROGRESS.md` from transcript-style history into this
  rolling snapshot plus milestone index.
- Updated `README.md` and `memory/NEXT_CODEX_BRIEF.md` so future sessions use
  docs as a task-scoped index instead of reading every historical doc by
  default.
- Verification:
  - `wc -l AGENTS.md README.md memory/*.md docs/*.md`: `memory/PROGRESS.md`
    is now `199` lines; counted project context files total `1,999` lines
    versus `3,506` before compaction.
  - `rg -n "docs/\\*|NEXT_CODEX_BRIEF|CODEX_BRIEF" README.md docs memory AGENTS.md`:
    no handoff requires reading `docs/*.md`; remaining `docs/*.md` mention is
    the required-file inventory in `docs/VERIFICATION.md`.
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`:
    passed.
  - `git diff --check`: passed.

Hold-inside midline transition review milestone:

- Review doc: `docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md`.
- Inputs: `results/hold-inside-midline-transition-audit/`; `7,988`
  candidate rows, `720` summary rows, `192` stability rows.
- Verdict: not entry-ready; broad `hold_3_inside`/`hold_6_inside` rows show
  split-stable midline touch/close-across behavior by `12` bars, but current
  midline labels are not entry context or strategy inputs.

Hold-inside midline transition audit milestone:

- Commit: `277ce49 Add hold-inside midline transition audit`.
- Added `-hold-inside-midline-transition-audit`; result dir
  `results/hold-inside-midline-transition-audit/`; `7,988` candidate rows,
  `720` summary rows, `192` stability rows; `strategy=empty trades=0`.

Hold-inside directional edge review milestone:

- Commit: `0601952 Review hold-inside directional edge audit`.
- Review doc: `docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md`.
- Inputs: `results/hold-inside-directional-edge-audit/`.
- Verdict: not entry-ready. `hold_3_inside` and `hold_6_inside` remain useful
  range-survival context but do not show a split-stable directional edge toward
  frozen range high or low.
- Durable rule added to `memory/DECISIONS.md`: do not promote
  `paper_side=toward_high` or `paper_side=toward_low` into entry context.
- Verification passed:
  - `wc -l results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `git diff --check`

Hold-inside directional edge audit milestone:

- Commit: `00b8844 Add hold-inside directional edge audit`.
- Added CLI flag `-hold-inside-directional-edge-audit`.
- Result directory: `results/hold-inside-directional-edge-audit/`.
- Audit size: `15,976` candidate rows, `624` summary rows, `168` stability
  rows; CSV lines including header `15,977` / `625` / `169`.
- Scope stayed non-trading; `strategy=empty trades=0`.
- Verification passed with `go test ./...`, audit run, `wc -l`, handoff
  reference check, and `git diff --check`.

Detector context refinement review milestone:

- Commit: `a4698bf Review detector context refinement audit`.
- Review doc: `docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md`.
- Inputs: `results/detector-context-refinement-audit/`.
- Verdict: delayed `hold_3_inside` and `hold_6_inside` are the leading context
  refinement and materially reduce quick invalidation and trend leakage, but
  are not promoted to entry context because the gain is conditioned survival
  and labels are regime-durability outcomes, not P&L.
- Verification passed with CSV line counts `113,825` / `641` / `161`,
  `go test ./...`, and `git diff --check`.

Detector context refinement audit milestone:

- Commit: `b980095 Add detector context refinement audit`.
- Added CLI flag `-detector-context-refinement-audit`.
- Result directory: `results/detector-context-refinement-audit/`.
- Audit size: `113,824` candidate rows, `640` summary rows, `160` stability
  rows; horizons `1`, `3`, `6`, `12`.
- Scope stayed non-trading; `strategy=empty trades=0`.

Detector durability sweep review milestone:

- Commit: `7bd2852 Review detector durability sweep`.
- Review doc: `docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md`.
- Inputs: `results/detector-durability-sweep/`.
- Verdict: no current `DefaultDetectorSweepProfiles` profile is approved as
  future entry context; `p30_c12_bollinger_on_adx_on` is diagnostic only.
- Verification passed with CSV line counts, `go test ./...`, and
  `git diff --check`.

Detector durability sweep milestone:

- Commit: `d70236b Add detector durability sweep`.
- Added CLI flag `-detector-durability-sweep`.
- Result directory: `results/detector-durability-sweep/`.
- Audit size: `304` broad rows, `9,088` slice rows, `76` stability rows.
- Scope stayed non-trading; `strategy=empty trades=0`.

Range regime durability review milestone:

- Commit: `bc81d9d Review range regime durability`.
- Review doc: `docs/RANGE_REGIME_DURABILITY_REVIEW.md`.
- Inputs: `results/range-regime-durability-audit/`.
- Verdict: current balanced detector regimes are not durable enough as context
  for future entry hypotheses. Refine detector/context before trigger work.
- Verification passed with CSV line counts, `go test ./...`, and
  `git diff --check`.

## Historical Milestone Index

- `7e1494b Add range regime durability audit`: added
  `-range-regime-durability-audit`; result dir
  `results/range-regime-durability-audit/`; `11,984` episode rows and `452`
  summary rows; `strategy=empty trades=0`.
- `b92207a Review compression breakout audit`: added
  `docs/COMPRESSION_BREAKOUT_REVIEW.md`; verdict not entry-ready; result dir
  `results/compression-breakout-audit/`.
- `be1e8d5 Add compression breakout audit mode`: added
  `-compression-breakout-audit`; `5,096` candidate rows and `24` summary rows;
  `strategy=empty trades=0`.
- `0256301 Review SR false-break reclaim timing`: added
  `docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md`; verdict not entry-ready;
  result dir `results/sr-false-break-reclaim-timing-audit/`.
- `1264a2f Add false-break reclaim timing audit`: added
  `-sr-false-break-reclaim-timing-audit`; `17,652` candidate rows and `24`
  summary rows; `strategy=empty trades=0`.
- `4e1a682 Add SR confirmation timing audit and review`: added
  `-sr-confirmation-timing-audit` and `docs/SR_CONFIRMATION_TIMING_REVIEW.md`;
  `9,692` candidate rows and `72` summary rows; verdict not entry-ready.
- `ed82756 Add SR rejection timing audit and review`: added
  `-sr-rejection-timing-audit` and `docs/SR_REJECTION_TIMING_REVIEW.md`; `968`
  candidate rows and `24` summary rows; verdict not entry-ready.
- `f60a8d9 Add entry readiness review gate`: added
  `docs/ENTRY_READINESS_REVIEW.md`; `internal/lab` coverage reached `99.8%`;
  first entries remained blocked pending non-trading timing evidence.
- `04edccc Use memory brief as canonical handoff`: removed duplicate root
  `CODEX_BRIEF.md`; `memory/NEXT_CODEX_BRIEF.md` became canonical.
- `1e26695 Add SR boundary inspection mode`: added `-sr-boundary-inspect`;
  `281,080` boundary events inspected and `192` comparison rows; `strategy=empty
  trades=0`.
- `9833525 Add SR boundary audit mode`: added `-sr-boundary-audit`; `281,080`
  boundary event rows and `192` boundary quality rows.
- `bdf1398 Add mods and next plan;`: pinned helper modules and documented
  adapter-only research helper boundaries in `docs/RESEARCH_HELPERS.md`.
- `6483e8c Add detector sweep audit mode`: added `-detector-sweep`; `19`
  profiles and `76` sweep rows; balanced baseline
  `p30_c12_bollinger_on_adx_off`.
- `e902a9e Initialize range strategy lab`: initialized standalone lab,
  tracked memory, and empty-strategy smoke path.

## Standard Closeout Checks

Use these unless the task has a narrower verifier:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```
