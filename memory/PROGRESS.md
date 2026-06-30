# Progress

This file is the always-read project snapshot. Keep it compact: current state,
latest verification, result paths, and a milestone index. Detailed evidence
belongs in focused `docs/` reviews, generated artifacts under `results/`, and
git history.

## Current State

- Scope is offline Binance USDT-M futures range-strategy research. The default
  CLI remains BTCUSDT `5m` with `lab.EmptyStrategy`; trades remain `0` unless an
  explicit offline research flag is passed.
- Active source contract remains:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed-candle decision semantics, next-bar-open entries, one-position max,
  stop-first ambiguity, cost accounting, split metrics, and source manifests
  remain core infrastructure.
- Research is not stopped, but automatic reuse of failed premises is stopped.
  Reviewed exclusion evidence includes structured compression, breakout-retest
  acceptance, clean breakout continuation, hold-inside/midline, impulse
  absorption, higher-timeframe nested range rotation, `range_occupancy_rotation_v1`,
  and range quality/session/failure-mode triage cohorts in their reviewed forms.
- The latest completed research doc is
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_REVIEW.md`.
  The user explicitly approved implementing the offline backtest for
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`, and it failed at
  `post_compression_directional_expansion_backtest_failed_no_usable_strategy`.
  The implementation reproduced the accepted source/resample contract and the
  representative raw candidate identity (`468` raw rows), executed `421` trades,
  and passed the trade-count gate, but failed gross edge, extra
  slippage-stress edge, stress profit factor, and drawdown gates. Full gross P&L
  was `208.560999`, engine net was `-129.258571`, extra slippage-stress net was
  `-227.226250`, full stress PF was `0.799666`, and full stress max drawdown was
  `0.289326`. The candidate is closed as no usable fixed strategy in this form;
  no adjacent-cell rescue, exit retune, derivatives-veto interaction, replay,
  walk-forward, paper/testnet/live path, or promotion is authorized.
- The prior completed research doc is
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_SPEC.md`.
  The user explicitly approved the docs-only offline backtest spec for
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`. It stopped at
  `post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval`
  and fixed one model: next-`15m`-open entry, stop at
  `entry_price - 1.0 * ATR(14)[d-1]`, target at
  `entry_price + 2.0 * ATR(14)[d-1]`, max hold `48` closed `15m` bars,
  one-position max, stop-first ambiguity, `1%` risk-at-stop sizing, `1x`
  notional cap, `0.0004` fee per side, and `0.000116` slippage per side. The
  later approved implementation consumed that gate and failed in reviewed form.
- The prior completed research doc is
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_STRATEGY_PREMISE_SPEC.md`.
  The user explicitly approved the docs-only strategy-premise spec for
  `btc_15m_post_compression_directional_expansion_v1`. It stopped at
  `post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval`
  and selected exactly one conservative representative candidate for a later
  docs-only offline backtest spec:
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`. The candidate is
  long-only, uses the `192`-bar lookback, bottom `20%` compression threshold,
  `0.2` prior-bar `ATR(14)` upside breakout, no volume filter, and the `48`
  closed `15m` bar evidence horizon. It does not authorize backtest
  implementation, P&L, replay, walk-forward, optimizer selection, derivatives
  veto interaction, or promotion. The later backtest spec fixed one risk/exit
  model and moved the line to implementation approval.
- The prior completed research doc is
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md`.
  The user explicitly approved implementing the zero-trade audit for
  `btc_15m_post_compression_directional_expansion_v1`, and it passed at
  `btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review`.
  The audit used only local BTCUSDT Binance USDT-M futures `5m` candles
  resampled to exact closed UTC `15m` bars, emitted `0` trades, and found a
  narrow long-only `48`-bar label-separation pocket at the `192`-bar lookback
  and bottom `20%` compression threshold across adjacent breakout/volume cells.
  It does not authorize a strategy, P&L, replay, walk-forward, or veto
  interaction. The later strategy-premise spec narrowed this to one
  representative backtest-spec candidate.
- The prior completed research doc is
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_PREMISE_SPEC.md`.
  The user explicitly approved the docs-only BTCUSDT `15m` post-compression
  directional expansion premise spec. It stopped at
  `independent_entry_premise_spec_ready_for_user_approval` and selected
  `btc_15m_post_compression_directional_expansion_v1` as the one independent
  BTCUSDT `15m` local-source premise family for zero-trade audit.
- The prior completed research doc is
  `docs/FUTURES_INDEPENDENT_ENTRY_PREMISE_AND_HYPOTHESIS_MAP.md`. The user
  explicitly approved implementing the combined docs-only hypothesis map and
  independent entry-premise spec plan. The review stopped at
  `independent_entry_premise_and_hypothesis_map_needs_user_scope_choice`
  because the current reviewed evidence does not expose exactly one eligible
  BTCUSDT `15m` local-source independent entry premise. The later
  post-compression premise spec supersedes its user-choice gate by selecting
  exactly one new BTCUSDT `15m` local-source event premise for approval.
- The prior completed research doc is
  `docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_INTEGRATION_SPEC.md`.
  The user explicitly approved the docs-only derivatives no-trade filter
  integration spec, and it stopped at
  `derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise`.
  The canonical veto `btc_15m_basis_discount_no_trade_veto_v1` is preserved as a
  future veto candidate only, but no implementation gate is selected because
  there is no independently approved entry premise for it to filter.
- The prior completed research doc is
  `docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_PREMISE_AUDIT_REVIEW.md`.
  The user explicitly approved implementing the zero-trade derivatives no-trade
  filter premise audit, and it passed at
  `derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec`.
  The audit (`-futures-derivatives-no-trade-filter-premise-audit`) reproduced
  all `5` selected BTCUSDT `15m` `h48` exact toxic rows, built a de-duplicated
  canonical veto union with `1,823` rows, reported `821` overlap rows and
  collateral damage, and produced `0` trades. This is a filter-integration
  premise only, not a strategy or P&L result.
- The prior completed doc is
  `docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md`. The user
  explicitly approved the docs-only strategy-premise spec, and it stopped at
  `derivatives_context_strategy_premise_spec_ready_for_user_approval`. The spec
  selected one later premise track: a BTCUSDT `15m` derivatives-context no-trade
  filter audit. It rejected the rotation-entry and two-track alternatives: the
  five toxic/no-trade cohorts were coherent enough for a zero-trade filter audit,
  while the single rotation candidate remained diagnostic only.
- The prior completed research doc is
  `docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md`. The user explicitly
  approved implementing the zero-trade derivatives context audit, and it passed
  at
  `derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec`.
  The audit (`-futures-derivatives-context-audit`) used only the `9` validated
  derivatives mark/index/premium `5m` CSVs plus the `3` candle anchors, kept the
  conservative one-`5m` lag and no-fill/no-interp/no-nearest-future policy, and
  produced `6` passing BTCUSDT `15m` separation cohorts (`5` no-trade/toxic,
  `1` rotation candidate). ETHUSDT and SOLUSDT produced `0` passing cohorts.
  Common outputs stayed zero-trade.
- The earlier completed research doc is
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md`. The user explicitly
  approved implementing the zero-trade derivatives context source audit, and it
  passed at
  `derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief`.
  The audit (`-futures-derivatives-context-source-audit`) validated the `9`
  materialized mark/index/premium `5m` source files, SHA-256-bound their
  provenance (hashes match the materialization manifests), and proved
  anti-lookahead alignment to the `5m` candle anchors under a conservative
  one-interval lag: all `6` required mark/index streams cleared the `0.99`
  coverage bar (min `0.994472`, index BTCUSDT), recorded missingness with
  `forward_filled_rows=0`, and produced `0` trades.
- The prior derivatives steps remain durable boundaries:
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md` passed at
  `derivatives_context_source_materialization_passed_ready_for_source_audit_approval`;
  `729` checksum-verified raw zips, `9` normalized CSVs, and `5` manifests live
  under `../binance-bot/data/derivatives/`.
- The next step is not automatic implementation. Neither the passing source
  audit, context-audit brief, context-audit review, strategy-premise spec,
  no-trade filter premise audit, no-trade filter integration spec,
  independent-entry map, post-compression premise spec, nor post-compression
  zero-trade audit authorizes entries, exits, P&L, replay, walk-forward,
  packaging, paper/testnet/live paths, exchange API work, credentials, deploy
  files, or promotion. The post-compression backtest implementation consumed the
  fixed backtest-spec approval gate and failed in reviewed form, so there is no
  selected next implementation. Any materially different premise, any
  derivatives-veto interaction, and any new backtest/replay/walk-forward path
  require separate explicit approval.
- The prior dependency docs are
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md` and
  `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.
- The router produced `58` rules from passing state-audit rankings, assigned
  `13,546` `no_trade`, `1,299` `tradable_rotation`, `0`
  `trend_continuation`, and `14,939` `diagnostic_only` rows, and found `3`
  passing router cohorts: `2` no-trade and `1` rotation. The follow-up router
  rotation premise audit collapsed that rotation context to `278` segments and
  `97` valid boundary-reclaim events, found `0` passing premise cohorts, and
  closed `router_gated_boundary_reclaim_rotation_v1` in reviewed form. It does
  not authorize a strategy, entry, exit, optimizer, replay, walk-forward, symbol
  expansion, source expansion, live/paper/testnet path, exchange API, deploy
  file, martingale, averaging down, or two-exchange logic.
- The post-rotation premise failure pivot review stopped with
  `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
  No automatic BTCUSDT-only price-only audit is selected.
- User explicitly approved the BTC regime plus ETH/SOL zero-trade context audit
  described in
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md`.
  Implementation used only the approved local Binance USDT-M futures `5m`
  BTCUSDT, ETHUSDT, and SOLUSDT files, produced `0` passing context cohorts,
  and closed the path in reviewed zero-trade form. BTCUSDT remains
  market-regime context and diagnostic-only authority; ETHUSDT/SOLUSDT remain
  failed context authority candidates only. This does not authorize entries,
  exits, P&L backtests, optimizer grids, replay, walk-forward,
  paper/testnet/live paths, exchange API, credentials, deploy files, broad
  mining, martingale, averaging down, or two-exchange logic.
- Parked future directions remain documented but not implementation-ready:
  spread-range/pair-range is parked behind source/engine complexity, broad
  derivatives source expansion remains stopped because the current
  strategy-premise spec creates no source-expansion need, and volatility-aware
  exits remain rejected until a new independent entry premise first shows gross
  edge before costs.
- `memory/NEXT_CODEX_BRIEF.md` is the canonical next-session prompt.

## 2026-06-30

BTCUSDT `15m` post-compression directional expansion fixed backtest:

- Added implementation and review doc:
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_REVIEW.md`.
- Added CLI flag:
  `-futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest`.
- Output path:
  `results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/`.
- Stop state:
  `post_compression_directional_expansion_backtest_failed_no_usable_strategy`.
- User explicitly approved the offline backtest implementation for
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`.
- Source reproduced:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`; Binance USDT-M
  futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`; `gap_count=0`;
  `duplicate_count=0`; `zero_volume_count=66`; `comparison_only=false`;
  `validation_status=accepted`.
- Exact `15m` resample reproduced: `191,328` rows;
  `2021-01-01T00:00:00Z` first open; `2026-06-16T23:45:00Z` last open;
  `2026-06-16T23:59:59Z` last close; `3` expected child bars;
  `0` missing child opens; validation accepted.
- Candidate identity reproduced: expected `468` raw representative-cell signal
  rows before one-position filtering, got `468`.
- Executed `421` trades after one-position filtering; primary split trade
  counts were `152` (`2021_2022_stress`), `146` (`2023_2024_oos`), and `123`
  (`2025_2026_recent`).
- Full-sample economics: gross P&L `208.560999`, engine net `-129.258571`,
  extra slippage-stress net `-227.226250`, stress PF `0.799666`, stress max
  drawdown `0.289326`, win rate `0.368171`.
- Split economics: `2021_2022_stress` stress net `20.837617` and PF
  `1.042668`; `2023_2024_oos` stress net `-115.553702` and PF `0.684405`;
  `2025_2026_recent` gross P&L `-15.799742`, stress net `-132.510165`, and PF
  `0.526275`.
- Falsification: source/resample, candidate identity, leakage, trade count,
  robustness, optimizer-contamination, closed-family, and derivatives-veto
  gates passed; gross edge, extra slippage-stress edge, stress PF, and drawdown
  gates failed.
- CSV artifact line counts including headers: signals `469`, trades `422`,
  strategy summary `13`, cost stress `13`, skips `7`, sources `2`, resample
  coverage `2`, common summary `13`, total `941`.
- Durable outcome: this exact fixed candidate is closed as no usable strategy in
  this form; no adjacent-cell rescue, exit retune, derivatives-veto interaction,
  replay, walk-forward, or promotion is authorized.
- Commands run:
  - `/usr/local/go/bin/gofmt -w cmd/rangelab/main.go cmd/rangelab/main_test.go internal/lab/futures_btc_15m_post_compression_l192_q20_m020_none_long_h48_backtest.go internal/lab/futures_btc_15m_post_compression_l192_q20_m020_none_long_h48_backtest_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest`
  - `wc -l results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
  - `git diff --cached --check`
- Verification outcomes: `gofmt` completed; full Go tests passed; the backtest
  rerun reproduced `468` signals, `421` trades, and the failed stop state;
  `wc -l` totaled `941` CSV lines; reference scan found the canonical
  `memory/NEXT_CODEX_BRIEF.md` references plus historical/checklist mentions
  only; `git diff --check` passed.

BTCUSDT `15m` post-compression directional expansion backtest spec:

- Added docs-only spec:
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_SPEC.md`.
- Stop state:
  `post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval`.
- User explicitly approved the docs-only offline backtest spec for
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`.
- Decision: request a later offline backtest implementation approval gate for
  exactly one candidate stream: BTCUSDT closed UTC `15m`, long only, prior
  `192`-bar range, bottom `20%` compression against the prior `1,920` valid
  range-width observations, decision close above the prior range by `0.2`
  prior-bar `ATR(14)`, no volume filter, and next-`15m`-open entry timing.
- Fixed risk/exit model for the later implementation: one open position max;
  stop-first ambiguity; stop at `entry_price - 1.0 * ATR(14)[d-1]`; target at
  `entry_price + 2.0 * ATR(14)[d-1]`; max hold `48` closed `15m` bars;
  `1%` risk-at-stop sizing; `1x` notional cap; start balance `1000`; fee
  `0.0004` per side; slippage `0.000116` per side.
- Future implementation must reproduce the representative zero-trade audit cell
  count (`l192_q20_m020_none` long `h48` full rows = `468`) before
  one-position filtering, report both current engine net and an extra
  slippage-stress net, and make pass/fail decisions on the stress view.
- Future pass gates include source/resample reproduction, candidate identity,
  no leakage, minimum executed trade counts (`120` full, `25` in each primary
  split), positive gross and extra slippage-stress net in full plus
  non-negative `2023_2024_oos` and `2025_2026_recent`, full stress PF at least
  `1.20`, OOS/recent stress PF at least `1.05`, drawdown bounds, and no
  optimizer/closed-family/veto contamination.
- The canonical derivatives veto
  `btc_15m_basis_discount_no_trade_veto_v1` remains parked as future
  skip/retain evidence only and may not shape entries, exits, side, ranking,
  scoring, P&L, pass/fail decisions, replay, walk-forward, or promotion.
- No Go code, CLI flag, generated result directory, audit run, backtest run,
  source download, network request, data write, P&L artifact, replay,
  walk-forward, derivatives veto interaction, or strategy promotion was made or
  authorized.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the offline backtest implementation
  approval gate.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
  - `git diff --cached --check`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references plus historical/checklist mentions
  only; `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes; staged diff check passed.

BTCUSDT `15m` post-compression directional expansion strategy-premise spec:

- Added docs-only spec:
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_STRATEGY_PREMISE_SPEC.md`.
- Stop state:
  `post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval`.
- User explicitly approved the docs-only strategy-premise spec for
  `btc_15m_post_compression_directional_expansion_v1`.
- Decision: request a later docs-only offline backtest spec for exactly one
  representative candidate from the passing zero-trade pocket:
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`.
- Candidate boundary: long side only, closed UTC `15m` BTCUSDT futures
  decisions, `192` prior closed `15m` bars for the compression range, bottom
  `20%` compression threshold from the prior `1,920` valid range-width
  observations, close above the prior range by `0.2` prior-bar `ATR(14)`, no
  volume confirmation, next-`15m`-open timing for any later approved backtest,
  and the `48` closed `15m` bar zero-trade evidence horizon.
- The adjacent passing cells from the zero-trade audit remain supporting
  robustness evidence only. They may not reopen the `81`-cell grid, rank P&L,
  rescue a failed representative-cell backtest, or authorize optimizer
  selection.
- The later docs-only backtest spec must choose one fixed risk/exit model
  before any implementation. If it cannot do so without user preference, it
  must stop at a user-choice gate.
- The canonical derivatives veto
  `btc_15m_basis_discount_no_trade_veto_v1` remains parked as future
  skip/retain evidence only and may not shape entries, exits, side, ranking,
  scoring, P&L, pass/fail decisions, replay, walk-forward, or promotion.
- No Go code, CLI flag, generated result directory, audit run, backtest run,
  source download, network request, data write, P&L artifact, replay,
  walk-forward, derivatives veto interaction, or strategy promotion was made or
  authorized.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the docs-only offline backtest-spec
  approval gate.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
  - `git diff --cached --check`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references plus historical/checklist mentions
  only; `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes; staged diff check passed.

BTCUSDT `15m` post-compression directional expansion zero-trade audit:

- Added implementation review:
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md`.
- Stop state:
  `btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review`.
- User explicitly approved implementing the zero-trade audit for
  `btc_15m_post_compression_directional_expansion_v1`.
- Implemented CLI flag:
  `-futures-btc-15m-post-compression-directional-expansion-audit`, defaulting
  to
  `results/futures-btc-15m-post-compression-directional-expansion-audit/`.
- Source facts reproduced: `573,984` BTCUSDT Binance USDT-M futures `5m`
  candles from `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`,
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Exact closed UTC `15m` resample facts: `191,328` rows, first open
  `2021-01-01T00:00:00Z`, last open `2026-06-16T23:45:00Z`, last close
  `2026-06-16T23:59:59Z`, `3` expected child bars, `0` missing child opens,
  validation accepted.
- Result path:
  `results/futures-btc-15m-post-compression-directional-expansion-audit/`.
  Audit artifacts reported `parameter_cells=81`, `candidate_rows=386,694`,
  `dedup_events=4,677`, `baseline_rows=24`, `split_summary_rows=1,944`,
  `adjacency_rows=486`, `missingness_rows=4`, `passing_cells=9`,
  `adjacent_pass_clusters=9`, and common outputs with `trades=0`.
- Falsification gates passed: source/resample, leakage, full de-duplicated
  candidate size (`4,677` versus required `300`), split size (minimum primary
  split `1,484` versus required `50`), baseline separation, adjacent-cell
  cluster, split stability, closed-family protection, derivatives-veto
  contamination, and zero-trade common outputs.
- Passing evidence is narrow and diagnostic only: long-side `48`-bar labels at
  lookback `192`, compression threshold bottom `20%`, all breakout thresholds
  `0.1`/`0.2`/`0.3` prior-bar `ATR(14)`, and all volume modes. No short-side,
  `16`-bar, `32`-bar, lookback `48`/`96`, or compression threshold `30%`/`40%`
  surface passed the full gate.
- Missingness was skipped, never filled: audit warmup `2,112`, missing forward
  label `27`, missing max-horizon future `48`, and missing volume reference
  `780`.
- Added Go tests for event construction, label anchoring/symmetry,
  de-duplication, baseline comparison, adjacency, stop-state precedence, CLI
  artifact creation, zero-trade common outputs, spot rejection, and flag
  conflicts.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the docs-only strategy-premise spec
  approval gate. No trading implementation is authorized.
- Commands run:
  - `gofmt -w cmd/rangelab/main.go cmd/rangelab/main_test.go internal/lab/futures_btc_15m_post_compression_directional_expansion_audit.go internal/lab/futures_btc_15m_post_compression_directional_expansion_audit_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-15m-post-compression-directional-expansion-audit`
  - `wc -l results/futures-btc-15m-post-compression-directional-expansion-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
  - `git diff --cached --check`
- Verification outcomes: all package tests passed; audit run passed with `0`
  trades and the passing stop state; CSV line count totaled `393,934`
  including headers; reference scan found canonical `memory/NEXT_CODEX_BRIEF.md`
  references plus historical/checklist mentions only; `git diff --check`
  passed; pre-staged `git status --short` showed only intended code, docs, and
  memory changes; staged diff check passed.

BTCUSDT `15m` post-compression directional expansion premise spec:

- Added docs-only spec:
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_PREMISE_SPEC.md`.
- Stop state:
  `independent_entry_premise_spec_ready_for_user_approval`.
- User explicitly approved the docs-only premise spec plan after supplying the
  concrete route:
  `btc_15m_post_compression_directional_expansion_v1`.
- Decision: select one independent BTCUSDT `15m` local-source premise family for
  a later zero-trade audit. Candidate rows are closed-candle post-compression
  directional expansions over local BTCUSDT futures `15m` resamples only.
- The bounded parameter grid is predeclared but not run: compression lookback
  `48`/`96`/`192`; compression threshold bottom `20%`/`30%`/`40%` of prior
  `1,920` closed `15m` range-width observations; breakout beyond prior range by
  `0.1`/`0.2`/`0.3` prior-bar `ATR(14)`; volume confirmation `none`, above
  prior `96`-bar median, or above prior `96`-bar `60%` percentile.
- The later audit labels are zero-trade only: intended-side forward close
  return, favorable excursion, adverse excursion, and
  favorable-greater-than-adverse rate over `16`, `32`, and `48` `15m` bars,
  compared against the unconditional eligible `15m` baseline by split, side, and
  horizon.
- Falsification gates include source/resample failure, leakage, fewer than
  `300` de-duplicated `(decision_close, side)` candidates, fewer than `50`
  candidates in any primary period split, no baseline separation, only one
  isolated passing parameter cell, split instability, closed-family reslice, or
  derivatives-veto contamination.
- The canonical derivatives veto
  `btc_15m_basis_discount_no_trade_veto_v1` remains parked as future
  skip/retain evidence only. It cannot shape the entry premise, create entries,
  choose side, score P&L, or be tested until a separate later interaction audit
  is approved after an independent entry audit exists.
- The spec authorizes no Go code, CLI flag, generated result directory, audit
  run, source download, source materialization, data write, entry, exit, P&L
  backtest, optimizer grid, replay, walk-forward, portfolio construction,
  paper/testnet/live path, exchange API, credential, deploy file, martingale,
  averaging down, two-exchange logic, closed-family rescue, veto interaction, or
  strategy promotion.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the zero-trade audit approval gate.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
  - `git diff --cached --check`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references plus historical/checklist mentions
  only; `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes; staged diff check passed.

Next Codex brief gate hardening:

- User asked to apply the recommendations from the read-only review of
  `memory/NEXT_CODEX_BRIEF.md`.
- Refined the canonical next-session brief without changing the active stop
  state:
  `independent_entry_premise_and_hypothesis_map_needs_user_scope_choice`.
- The brief now requires one complete route before work starts, treats a route
  name alone as insufficient, and keeps the worktree unchanged when required
  route details are missing.
- The brief now keeps a new BTCUSDT `15m` local-source premise docs-first,
  points detailed veto metrics back to the integration spec instead of
  foregrounding them in the user-choice gate, defaults new source-family work to
  source-scope first unless the source is already local/provenance-bound/
  coverage-known, and adds `git diff --cached --check` to staged closeout.
- No Go code, CLI flag, generated result directory, audit run, source download,
  data write, entry, exit, P&L backtest, replay, walk-forward, veto interaction,
  strategy promotion, README index change, or durable `memory/DECISIONS.md`
  boundary change was made.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
  - `git diff --cached --check`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references plus historical/checklist mentions
  only; `git diff --check` passed; pre-commit `git status --short` showed only
  intended memory changes; staged diff check passed.

Independent entry-premise and hypothesis map:

- Added docs-only review:
  `docs/FUTURES_INDEPENDENT_ENTRY_PREMISE_AND_HYPOTHESIS_MAP.md`.
- Stop state:
  `independent_entry_premise_and_hypothesis_map_needs_user_scope_choice`.
- User explicitly approved implementing the combined hypothesis-map and
  independent-entry-premise spec plan.
- Decision: no single BTCUSDT `15m` local-source independent entry premise was
  selected from the current reviewed evidence. The reviewed candidates either
  collapse into closed families, remain filter/context evidence only, or require
  a fresh user-supplied event premise or scope change.
- The canonical derivatives veto
  `btc_15m_basis_discount_no_trade_veto_v1` remains parked as future
  skip/retain evidence only. It cannot shape the entry premise, create entries,
  choose side, score P&L, or be tested until a separate independently approved
  candidate-entry stream exists.
- The next gate is a user scope/premise choice: a new BTCUSDT `15m`
  local-source candidate event, a higher-timeframe premise, spread-range/source
  scope, another explicitly approved source family, or no further audit.
- The review authorizes no Go code, CLI flag, generated result directory, audit
  run, source download, source materialization, data write, entry, exit, P&L
  backtest, optimizer grid, replay, walk-forward, portfolio construction,
  paper/testnet/live path, exchange API, credential, deploy file, martingale,
  averaging down, two-exchange logic, closed-family rescue, veto integration, or
  strategy promotion.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the user-choice gate.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references plus historical/checklist mentions
  only; `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes.

## 2026-06-29

Derivatives context no-trade filter integration spec:

- Added docs-only spec:
  `docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_INTEGRATION_SPEC.md`.
- Stop state:
  `derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise`.
- User explicitly approved the docs-only derivatives context no-trade filter
  integration spec.
- Decision: preserve `btc_15m_basis_discount_no_trade_veto_v1` as a future veto
  candidate only, but defer integration because no independently approved entry
  premise exists. The passing no-trade filter premise audit stays useful as
  veto evidence, not as an entry signal, basis-tradability claim, or P&L result.
- No implementation gate is selected. A future interaction audit, if ever
  approved, must first name an independent entry premise and may only annotate
  candidate rows as skipped or retained; it may not create entries, alter entry
  logic, simulate fills, score P&L, optimize, replay, walk forward, promote a
  strategy, or reopen closed families.
- The spec authorizes no Go code, CLI flag, generated result directory, audit
  run, source download, network request, source materialization, data write
  under `../binance-bot/data/derivatives/`, entry, exit, P&L backtest, optimizer
  grid, replay, walk-forward, portfolio construction, paper/testnet/live path,
  exchange API, credential, deploy file, martingale, averaging down,
  two-exchange logic, closed-family rescue, or strategy promotion.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to record that no derivatives filter
  integration implementation is selected until an independent entry premise is
  supplied and explicitly approved.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes.

Derivatives context no-trade filter premise audit implementation:

- Added implementation review:
  `docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_PREMISE_AUDIT_REVIEW.md`.
- Stop state:
  `derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec`.
- User explicitly approved implementing the zero-trade derivatives no-trade
  filter premise audit selected by
  `docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md`.
- Added CLI flag `-futures-derivatives-no-trade-filter-premise-audit` (default
  out-dir `results/futures-derivatives-no-trade-filter-premise-audit`) plus
  audit engine/tests in
  `internal/lab/futures_derivatives_no_trade_filter_premise_audit.go`,
  `internal/lab/futures_derivatives_no_trade_filter_premise_audit_test.go`, and
  `cmd/rangelab/main_test.go`.
- Inputs stayed within the approved BTCUSDT source scope:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_BTCUSDT_2021_2026.csv`,
  `../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_BTCUSDT_2021_2026.csv`,
  and
  `../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_BTCUSDT_2021_2026.csv`.
- Source facts reproduced: BTCUSDT candle anchor `573,984` rows, `gap_count=0`,
  `duplicate_count=0`, `zero_volume_count=66`; mark rows `571,675`, gaps `6`,
  SHA-256
  `424c05ca880a31270eea1286d6cdd96ac1132d848e8f5e9d6b3b7177bb7c2858`;
  index rows `570,812`, gaps `8`, SHA-256
  `7ba5a375311e0324dab38f18c2a7137376b619ce63cfb01a17e5684c58390aca`;
  premium rows `571,959`, gaps `7`, SHA-256
  `094e610617f812f032e6a68b3ae6186b20359592415cdb678845c5d287ec298c`.
- Anti-lookahead model carried forward the conservative one-`5m` lag
  (`source_close_time + 5m <= decision_candle_close_time`), exact closed source
  rows, no fill/interpolation/nearest-future joins, and forward labels only as
  evaluation metadata. Missing lagged context was skipped and counted
  (`311` missing BTCUSDT `15m` basis-context rows out of `24,067` state rows).
- Result counts:
  `source_rows=4`, `coverage_rows=7`, `filter_definition_rows=5`,
  `exact_candidate_rows=20`, `canonical_union_rows=4`, `overlap_rows=40`,
  `veto_candidate_rows=1823`, `collateral_rows=37`, `missingness_rows=4`,
  `exact_candidates_passed=5`, `canonical_union_passed=true`, `trades=0`.
- Exact candidates reproduced the five selected toxic rows (`515`, `622`,
  `356`, `613`, and `538` full-sample rows; weakest splits `110`, `142`, `62`,
  `124`, and `115` respectively). The diagnostic rotation row was not selected
  and not converted into an entry premise.
- Canonical union:
  `btc_15m_basis_discount_no_trade_veto_v1`; full rows `1,823`; overlap rows
  `821`; nested trend-down premium overlap rows `515`; no-trade toxic rows
  `1,241`; full toxic rate `0.680746`; min split toxic rate `0.665485`; weakest
  split rows `387`; full toxic improvement versus local-only baseline
  `0.046269`. Collateral was reported (`311` rotation-useful and `271`
  continuation-useful full-sample rows blocked).
- Common outputs stayed zero-trade compatible:
  `summary.csv`/`summary.json` report `0` trades and `trades.json` is empty.
  This pass authorizes only a later docs-only filter integration spec, not an
  entry, exit, P&L backtest, optimizer, replay, walk-forward, paper/testnet/live
  path, exchange API, credential, deploy file, strategy promotion, martingale,
  averaging down, two-exchange logic, or closed-family rescue.
- Commands run:
  - `gofmt -w internal/lab/futures_derivatives_no_trade_filter_premise_audit.go internal/lab/futures_derivatives_no_trade_filter_premise_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-no-trade-filter-premise-audit -out-dir results/futures-derivatives-no-trade-filter-premise-audit`
  - `wc -l results/futures-derivatives-no-trade-filter-premise-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: all package tests passed; the audit produced the
  passing stop state above; generated CSV line count total was `1,971`
  including headers; reference scan found canonical `memory/NEXT_CODEX_BRIEF.md`
  references and historical/checklist mentions only; `git diff --check` passed;
  pre-commit `git status --short` showed only intended code, docs, and memory
  changes.

Derivatives context strategy-premise spec:

- Added docs-only spec:
  `docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md`.
- Stop state:
  `derivatives_context_strategy_premise_spec_ready_for_user_approval`.
- User explicitly approved writing the docs-only strategy-premise spec from the
  `6` passing BTCUSDT `15m` derivatives-context cohorts.
- Decision: select only a BTCUSDT `15m` no-trade filter premise for a later
  zero-trade audit. The `5` toxic/no-trade cohorts are coherent enough to test as
  exact/canonical veto candidates; the single rotation candidate is preserved as
  diagnostic evidence only and does not justify a rotation-entry premise or a
  second implementation track.
- Later allowed implementation, only after separate explicit user approval, is a
  zero-trade no-trade filter premise audit over local BTCUSDT Binance USDT-M
  futures candles and the already validated BTCUSDT mark/index/premium rows. It
  must keep the conservative one-`5m` lag
  (`source_close_time + 5m <= decision_candle_close_time`), no forward
  fill/interpolation/nearest-future joins, recorded missingness/skips, and
  forward labels only as evaluation metadata.
- The spec authorizes no Go code, CLI flag, generated result directory, audit
  run, source download, data write under `../binance-bot/data/derivatives/`,
  entry, exit, P&L backtest, optimizer, replay, walk-forward, paper/testnet/live
  path, exchange API, credential, deploy file, martingale, averaging down,
  two-exchange logic, or closed-family rescue.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the no-trade filter premise audit
  implementation approval gate.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes.

Derivatives context zero-trade audit implementation:

- Added implementation review:
  `docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md`.
- Stop state:
  `derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec`.
- User explicitly approved implementing the zero-trade derivatives context audit
  from
  `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md`.
- Added CLI flag `-futures-derivatives-context-audit` (default out-dir
  `results/futures-derivatives-context-audit`) plus audit engine/tests in
  `internal/lab/futures_derivatives_context_audit.go`,
  `internal/lab/futures_derivatives_context_audit_test.go`, and
  `cmd/rangelab/main_test.go`.
- Inputs stayed exactly within scope: the `9` validated derivatives
  mark/index/premium `5m` source CSVs under
  `../binance-bot/data/derivatives/` plus the `3` Binance USDT-M futures
  BTCUSDT/ETHUSDT/SOLUSDT `5m` candle anchors. No source downloads, network
  requests, data writes under `../binance-bot/data/derivatives/`, or new source
  families were used.
- Anti-lookahead model carried forward the source audit's conservative
  one-`5m` lag (`source_close_time + 5m <= decision_candle_close_time`), exact
  closed-interval joins, no forward fill/interpolation/nearest-future joins, and
  missing context as recorded skip rows. Forward labels remained metadata only
  in label/cohort/ranking/summary artifacts.
- Results: `source_rows=12`, `coverage_rows=18`,
  `basis_feature_rows=83,004`, `local_state_rows=83,640`,
  `label_rows=249,012`, `cohort_rows=512,190`, `ranking_rows=181,827`,
  `missingness_rows=36`, `passing_cohorts=6`, `trades=0`.
- Passing cohorts were narrow: all `6` are BTCUSDT `15m`; `5` are
  no-trade/toxic separation cohorts and `1` is a rotation candidate. ETHUSDT and
  SOLUSDT had `0` passing cohorts. Rows ranked after the passing set mostly
  failed inadequate-count gates.
- Coverage/missingness: required source lag coverage stayed at or above the
  rounded `0.994472` floor; derived local-state subset coverage was recorded
  separately (minimum `0.982930`, BTCUSDT `4h`) and missing rows were skipped,
  never filled or encoded as default context.
- Common outputs stayed zero-trade compatible (`summary.csv`/`summary.json`
  report `0` trades; `trades.json` contains no trades). The audit does not
  authorize entries, exits, P&L backtests, optimizers, replay, walk-forward,
  source expansion, paper/testnet/live paths, exchange API work, credentials,
  deploy files, or promotion.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the strategy-premise spec approval
  gate.
- Commands run:
  - `gofmt -w internal/lab/futures_derivatives_context_audit.go internal/lab/futures_derivatives_context_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-audit -out-dir results/futures-derivatives-context-audit`
  - `wc -l results/futures-derivatives-context-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; the audit reproduced the passing stop
  state and counts above; generated CSV line counts totaled `1,109,870`
  including headers; reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed.

Derivatives context zero-trade context-audit brief:

- Added docs-only brief:
  `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md`.
- Stop state:
  `derivatives_context_zero_trade_context_audit_brief_ready_for_user_approval`.
- Converts the passed source audit into a decision-complete plan for a later
  zero-trade derivatives context audit. The later audit question is
  separation-only: whether mark-minus-index basis, premium-index level, or
  basis-change buckets known at closed decision-candle time improve separation of
  BTCUSDT/ETHUSDT/SOLUSDT local range states (usable, toxic, rotation,
  continuation, no-trade) beyond the local price/volume state alone. It must not
  test basis tradability and must not measure entry/exit/P&L.
- Approved later inputs are only the `9` validated derivatives CSVs plus the `3`
  candle anchors. Carried-forward facts: SHA-256 provenance, conservative
  one-`5m`-interval lag (`source_close_time + 5m <= decision_candle_close_time`),
  exact closed-interval joins, no forward fill/interpolation/nearest-future
  joins, bounded recorded missingness (required basis context coverage floor
  `0.994472`, index BTCUSDT).
- Required anti-leakage rules: forward labels only in label/cohort/ranking/
  summary artifacts; basis/premium context built only from lagged source rows;
  missing context produces missingness/skip rows, never silent defaults; common
  outputs (`summary.json`, `summary.csv`, `trades.json`) stay zero-trade
  compatible.
- Material-difference guard: basis/premium is an orthogonal source (perp-vs-index
  dislocation), not a reslice of closed price-only or BTC-regime families; the
  later audit must reject itself (orthogonality gate plus
  `...rejected_closed_family_rescue`) if basis buckets are collinear with the
  local price/volume state. All three symbols may be local range-state authority
  candidates here, but no closed family is reopened and no symbol is promoted.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the context-audit implementation
  approval gate.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit `git status --short` showed only intended
  docs and memory changes.

Derivatives context zero-trade source audit implementation:

- Added `-futures-derivatives-context-source-audit` (default out-dir
  `results/futures-derivatives-context-source-audit`).
- Audit engine `internal/lab/futures_derivatives_context_source_audit.go` plus
  tests `internal/lab/futures_derivatives_context_source_audit_test.go` and a CLI
  flag test in `cmd/rangelab/main_test.go`.
- Review doc: `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md`.
- Stop state:
  `derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief`.
- User explicitly approved implementing the source audit from
  `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md`.
- Inputs: the `9` durable materialized mark/index/premium `5m` source CSVs under
  `../binance-bot/data/derivatives/` (read-only; recomputed file SHA-256 matched
  the materialization manifests, e.g. mark BTCUSDT
  `424c05ca...c2858`, index BTCUSDT `7ba5a375...390aca`) plus the three
  `573,984`-row candle anchors used only for alignment.
- Anti-lookahead model: publication lag unproven, so a conservative one-`5m`-
  interval lag (`source_close_time + 5m <= decision_candle_close_time`); no
  forward fill (`MaxExtraStalenessIntervals=0`), no interpolation, no
  nearest-future joins; every alignment row recorded `uses_future_rows=false`
  and `exact_closed_interval_join=true`.
- Results: `source_rows=9`, `anchor_rows=3`, `coverage_rows=9`,
  `alignment_rows=9`, `lag_rows=9`, `missingness_rows=9`, `provenance_rows=9`,
  `skip_rows=18`, `required_streams=6`, `aligned_required_streams=6`. All `9`
  streams validated with `0` duplicate-conflict, `0` close-time violations,
  `0` non-monotonic rows. Required-stream min lag coverage `0.994472`
  (index BTCUSDT) cleared the `0.99` bar; missingness recorded with
  `forward_filled_rows=0`. `trades.json` had no trades.
- Compliance: wrote exactly the brief's nine artifact families (CSV+JSON) plus
  zero-trade common outputs; deferred mark-minus-index basis derivation to a
  later context-audit stage (no basis/feature/cohort/label/ranking artifact).
- Commands run:
  - `gofmt -w internal/lab/futures_derivatives_context_source_audit.go internal/lab/futures_derivatives_context_source_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-source-audit -out-dir results/futures-derivatives-context-source-audit`
  - `wc -l results/futures-derivatives-context-source-audit/*.csv`
  - `rg --files ../binance-bot/data | rg -i "(mark[_-]?price|index[_-]?price|premium[_-]?index|premiumIndex|basis)"`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; the audit reproduced the counts above and
  the passing stop state; CSV artifacts totaled `95` lines plus the common
  `summary.csv`; the `rg --files` basis check returned empty because the adjacent
  `binance-bot` repo Git-ignores its `data/` tree (`find` confirms all nine
  durable source CSVs are present and were read); `git diff --check` passed;
  pre-commit `git status --short` showed only intended code, docs, and memory
  changes (generated `results/` are Git-ignored; derivatives data stays outside
  this repo).

## 2026-06-28

Derivatives context source materialization execution:

- Review doc:
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md`.
- Stop state:
  `derivatives_context_source_materialization_passed_ready_for_source_audit_approval`.
- User explicitly approved executing the materialization plan. Execution used a
  one-shot offline Go generator (run from session scratchpad, not tracked in this
  repo) that downloaded deterministic public Data Vision object URLs, verified
  each `.zip` against its published `.CHECKSUM` SHA-256, wrote raw zips,
  normalized CSVs, and manifests only under `../binance-bot/data/derivatives/`.
- Scope executed exactly as approved: owner `Binance Data Vision`, product
  Binance USDT-M futures, families `markPriceKlines`/`indexPriceKlines`
  (required) and `premiumIndexKlines` (optional cross-check), symbols
  `BTCUSDT`/`ETHUSDT`/`SOLUSDT`, interval `5m`, era
  `2021-01-01T00:00:00Z`..`2026-06-16T23:55:00Z`, object set monthly
  `2021-01`..`2026-05` plus daily `2026-06-01`..`2026-06-16` = `81` objects per
  family/symbol, `729` total.
- Object outcomes: `objects_ok=729`, `objects_missing=0`, `objects_error=0`; all
  `729` checksum-verified (`validation_status=accepted`). Total compressed bytes
  `125,508,895`. Family bytes: `markPriceKlines` `42,833,309`,
  `indexPriceKlines` `47,233,885`, `premiumIndexKlines` `35,441,701`. Required
  mark+index `486` objects / `90,067,194` bytes reproduced the planning
  inventory exactly.
- Normalized streams (`9`) all span the full era with `0` duplicate-conflict,
  `0` parse-error, `0` out-of-range, `0` non-monotonic rows. Required mark/index
  missing intervals total `9,820` of `3,443,904` (`0.285%`); rows/gaps:
  mark BTC `571,675`/`6`, ETH `573,402`/`4`, SOL `571,963`/`5`;
  index BTC `570,812`/`8`, ETH `573,116`/`4`, SOL `573,116`/`4`;
  premium BTC `571,959`/`7`, ETH `571,960`/`6`, SOL `572,248`/`5`.
- Gap decision: required mark/index `5m` public archives have real, whole-day
  aligned publication-outage gaps (e.g. `2021-06-30→07-02`, `2021-07-23→07-28`,
  `2022-10-01→10-03`, `2023-02-23→02-25`, `2023-11-10 ~04:00Z`), recurring across
  symbols/families; trade-candle anchors had `gap_count=0` over the same era.
  User chose to record gaps and pass (no imputation; bounded-missingness gating
  deferred to the later source audit) rather than fail closed. Integrity faults
  (duplicate-conflict, schema, checksum, missing required object) remained
  fail-closed; none occurred.
- Durable layout: `743` files total under `../binance-bot/data/derivatives/`
  (`729` raw zips, `9` normalized CSVs, `5` manifests). Normalized schema:
  `open_time,open,high,low,close,close_time,source_object_id`.
- Commands run:
  - `curl -sS -I ...markPriceKlines/BTCUSDT/5m/BTCUSDT-5m-2021-01.zip` and
    `.CHECKSUM` (connectivity + header/schema probe).
  - offline Go generator (download + checksum-verify + normalize + manifests).
  - `find ../binance-bot/data/derivatives -type f | sort`
  - `find ../binance-bot/data/derivatives/raw -name '*.zip' | wc -l`
  - `wc -l ../binance-bot/data/derivatives/*.csv`
  - `wc -l ../binance-bot/data/derivatives/manifests/*.csv`
  - `sha256sum ../binance-bot/data/derivatives/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: `743` files present; `729` raw zips; normalized CSV
  lines totaled `5,150,260` (rows plus `9` headers); `objects.csv` `730` lines,
  `files.csv` `10` lines; `9` normalized SHA-256 hashes recorded in the review
  doc; reference scan found only canonical `memory/NEXT_CODEX_BRIEF.md`
  references; `git diff --check` passed; pre-commit `git status --short` showed
  only intended `docs/` and `memory/` changes (generated derivatives data is
  outside this repo and not committed).

Derivatives context source materialization plan:

- Added docs-only materialization plan:
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md`.
- Stop state:
  `derivatives_context_source_materialization_plan_ready_for_execution_approval`.
- User explicitly approved an offline materialization plan for Binance public
  Data Vision mark-price, index-price, and optional premium-index `5m` klines
  for `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` under
  `../binance-bot/data/derivatives/`.
- The plan does not execute downloads, source parsing, normalization, durable
  data writes, source-audit implementation, context features, labels, cohorts,
  rankings, entries, exits, P&L backtests, replay, walk-forward, live probes,
  private endpoints, API keys, credentials, deploy files, or strategy
  promotion.
- Required future materialization families are Binance USDT-M futures
  `markPriceKlines` and `indexPriceKlines`. Optional cross-check family is
  `premiumIndexKlines`; missing optional premium-index rows may be recorded as
  an optional-family gap, but missing required mark/index objects must reject
  materialization.
- Approved object shape is `5m`, `BTCUSDT`/`ETHUSDT`/`SOLUSDT`, monthly
  objects from `2021-01` through `2026-05`, plus daily tail objects from
  `2026-06-01` through `2026-06-16`, matching the existing candle-anchor era.
- Adjacent local source inventory facts from
  `../binance-bot/research/2026-06-18_futures_perp_basis_reversion/event_study/source_inventory.csv`:
  required mark/index scope has `486` archive objects and `90,067,194`
  compressed bytes; `markPriceKlines` has `243` objects and `42,833,309`
  bytes; `indexPriceKlines` has `243` objects and `47,233,885` bytes.
  Optional `premiumIndexKlines` is estimated, not proven in this lab, at another
  same-shaped `243` objects and about `45` MB compressed.
- Target durable raw, normalized CSV, and manifest layouts are defined under
  `../binance-bot/data/derivatives/`. Future execution must record object URLs,
  byte counts, SHA-256 hashes, ETag/Last-Modified/Content-Length when returned,
  row counts, gaps, duplicates, parse errors, timestamp semantics, finality, and
  validation status.
- Future execution stop states include source gap, checksum/schema gap,
  unapproved source path, unapproved live/private path, and passed materialized
  states ready only for separate source-audit approval. Passing materialization
  will not authorize the source audit or any context/strategy work by itself.
- Local inventory check found no existing durable files under
  `../binance-bot/data/derivatives/`.
- Commands run:
  - `awk -F, 'NR>1 {count[$1]++; bytes[$1]+=$10; rows[$1]+=$17; total_count++; total_bytes+=$10; total_rows+=$17} END {for (f in count) printf "%s count=%d bytes=%d rows_total=%d\n", f, count[f], bytes[f], rows[f]; printf "total count=%d bytes=%d rows_total=%d\n", total_count, total_bytes, total_rows}' ../binance-bot/research/2026-06-18_futures_perp_basis_reversion/event_study/source_inventory.csv`
  - `awk -F, 'NR>1 {key=$1 "/" $3 "/" $4; count[key]++; bytes[key]+=$10; rows[key]+=$17} END {for (k in count) printf "%s count=%d bytes=%d rows_total=%d\n", k, count[k], bytes[k], rows[k]}' ../binance-bot/research/2026-06-18_futures_perp_basis_reversion/event_study/source_inventory.csv | sort`
  - `find ../binance-bot/data -maxdepth 4 -type f | sort | sed -n '1,160p'`
  - `find ../binance-bot/data/derivatives -maxdepth 5 -type f 2>/dev/null | sort | sed -n '1,120p'`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: adjacent inventory aggregation reproduced the required
  mark/index `486` object and `90,067,194` byte total; local durable data
  inventory showed candle CSVs only and no existing derivatives files under
  `../binance-bot/data/derivatives/`; reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit status showed only intended docs and
  memory changes.

Derivatives context zero-trade source audit brief:

- Added docs-only brief:
  `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md`.
- Stop state:
  `derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval`.
- The brief selects Binance USDT-M futures mark-price, index-price, or
  premium-index klines as the only first derivatives source family for a later
  source audit. Funding remains second; aggregate trades remain high-volume
  secondary and parked.
- The brief is source/alignment only. It does not approve implementation,
  source downloads, source materialization, source parsing, context-gain
  features, labels, cohorts, rankings, entries, exits, P&L backtests,
  optimizer grids, replay, walk-forward, paper/testnet/live paths, exchange
  APIs, credentials, deploy files, broad mining, martingale, averaging down, or
  two-exchange logic.
- No derivatives market-data source rows are approved as direct lab inputs
  today. A later implementation, only after explicit user approval, must prove
  durable local/offline mark/index/premium files under `../binance-bot/data/`
  or a documented subdirectory, or stop at source gap.
- Existing BTC/ETH/SOL Binance USDT-M futures `5m` candle files remain
  alignment anchors only:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Candle anchor facts preserved: each has `573,984` loaded candles from
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`; sorted streams had
  `gap_count=0` and `duplicate_count=0`; zero-volume counts were BTC `66`, ETH
  `47`, SOL `47`; SOL had one physical non-monotonic row and was accepted only
  after sorting.
- Required future source-audit rules include explicit timestamp semantics,
  publication/finality lag, checksum or provenance identifier, bounded
  missingness/staleness, and anti-lookahead joins where
  `source_close_time + publication_lag <= decision_candle_close_time`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to the source audit approval gate.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes.

Derivatives market-data context source scope review:

- Added docs-only scope review:
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md`.
- Stop state:
  `derivatives_context_source_scope_review_approved_needs_zero_trade_source_audit_brief`.
- Decision: a later zero-trade derivatives source-audit brief is justified only
  for source and alignment approval, not context-gain implementation.
- The review found no approved durable derivatives market-data rows in the
  lab's current local data scope. `../binance-bot/data/` contained durable
  candle CSVs only, and `../binance-bot/data/raw/` had no files in the reviewed
  search depth.
- Approved local/offline inputs for the later brief are only alignment anchors:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Adjacent `../binance-bot/research/` source-proof artifacts may be referenced
  only as process/source evidence, not as lab input data or strategy evidence.
  Candidate source families for the later brief are mark/index/premium basis
  first, funding second, and aggregate trades only as a high-volume secondary
  source-proof candidate.
- Rejected or blocked from current scope: `/tmp` caches as durable inputs,
  source downloads, live probes, private endpoints, exchange API keys, open
  interest from the current evidence set, long/short ratios without full-era
  archive proof, liquidation/force-order history from current evidence,
  order-book/depth without separate historical archive proof, entries, exits,
  P&L backtests, replay, walk-forward, optimizer grids, and promotion.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a separate zero-trade source-audit
  brief-writing task that stops before code at
  `derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval`.
- Commands run:
  - `rg --files ../binance-bot/data | rg -i "(funding|fund|open.?interest|basis|premium|mark|index|taker|long.?short|book|depth|liquid|oi|agg.?trade|trade)"`
  - `rg --files ../binance-bot/data | rg -i "(btcusdt|ethusdt|solusdt).*futures.*(5m|1h|15m|4h|um)|futures_um"`
  - `rg --files . | rg -i "(funding|fund|open.?interest|basis|premium|mark|index|taker|long.?short|book|depth|liquid|oi|agg.?trade)"`
  - `rg -n "funding|open interest|basis|premium|taker|long/short|order-book|order book|derivatives" docs memory README.md AGENTS.md`
  - `rg --files ../binance-bot | rg -i "(funding|open[_-]?interest|\\boi\\b|oi_|basis|premium[_-]?index|premiumIndex|mark[_-]?price|index[_-]?price|taker|long[_-]?short|longshort|depth|order[_-]?book|book[_-]?ticker|liquidat|agg[_-]?trade)"`
  - `rg --files ../binance-bot/data`
  - `ls -lah ../binance-bot/data`
  - `find ../binance-bot/data/raw -maxdepth 3 -type f`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: inventory checks found durable local candle CSVs but
  no approved durable derivatives market-data source rows under
  `../binance-bot/data/`; adjacent source-proof artifacts exist only as
  references; final brief-reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit status showed only intended docs and
  memory changes.

BTC regime plus ETH/SOL zero-trade audit implementation:

- Added zero-trade audit implementation behind
  `-futures-btc-regime-eth-sol-context-audit`.
- Review doc:
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md`.
- Result directory:
  `results/futures-btc-regime-eth-sol-context-audit/`.
- Stop state:
  `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.
- The audit used only the approved local Binance USDT-M futures `5m` files:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Source facts reproduced: each file had `573,984` loaded candles from
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`; all accepted sorted
  streams had `gap_count=0` and `duplicate_count=0`; zero-volume counts were
  BTC `66`, ETH `47`, SOL `47`; SOL had one physical non-monotonic row and was
  accepted only after sorting.
- Resampled coverage passed for BTC/ETH/SOL `15m`, `1h`, and `4h`: rows were
  `191,328`, `47,832`, and `11,958` respectively for each symbol, with
  `gap_count=0`, `duplicate_count=0`, `missing_child_open_count=0`, and
  `complete=true`.
- Generated artifact counts: BTC state rows `29,784`, ETH/SOL local state rows
  `30,717`, relative-strength rows `30,717`, label rows `92,151`, cohort rows
  `359,055`, ranking rows `160,983`, and passing cohorts `0`.
- Main ranking failure inventory: `142,615` rows failed
  `btc_context_improvement_gate_failed`, `16,095` failed
  `single_split_contribution_above_gate`, `1,627` failed
  `route_rate_gate_failed`, `419` failed `missing_period_split`, and `227`
  failed `inadequate_cohort_count`.
- Anti-leakage and boundary outcomes: BTCUSDT stayed diagnostic market-regime
  context only; ETHUSDT/SOLUSDT stayed zero-trade context authority candidates
  only; forward labels were not used as state/context/gating inputs; common
  outputs remained zero-trade compatible; no entries, exits, P&L backtests,
  replay, walk-forward, source downloads, or strategy promotion were added.
- Commands run:
  - `gofmt -w internal/lab/futures_btc_regime_eth_sol_context_audit.go internal/lab/futures_btc_regime_eth_sol_context_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-regime-eth-sol-context-audit -out-dir results/futures-btc-regime-eth-sol-context-audit`
  - `wc -l results/futures-btc-regime-eth-sol-context-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced `source_rows=3`,
  `coverage_rows=9`, `btc_state_rows=29784`, `local_state_rows=30717`,
  `relative_strength_rows=30717`, `label_rows=92151`,
  `cohort_rows=359055`, `ranking_rows=160983`, `passing_cohorts=0`, and stop
  state `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`;
  `wc -l` over CSV artifacts plus common `summary.csv` totaled `703,514`
  lines; reference scan found canonical `memory/NEXT_CODEX_BRIEF.md`
  references and checklist mentions only; `git diff --check` passed; pre-commit
  status showed only intended code, docs, and memory changes.

BTC regime plus ETH/SOL zero-trade audit brief:

- Added docs-only brief:
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md`.
- Stop state:
  `btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval`.
- The brief defines the later audit question: whether BTCUSDT regime buckets
  known at closed decision-candle time improve separation of ETHUSDT/SOLUSDT
  usable, toxic, rotation, continuation, or no-trade local range states.
- BTCUSDT role remains market-regime context and diagnostic-only authority row;
  ETHUSDT/SOLUSDT may become possible authority rows only inside a zero-trade
  context audit, not strategy promotion.
- Allowed source scope remains only the already local Binance USDT-M futures
  `5m` files:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Source facts preserved from prior validation: each file has `573,984` loaded
  candles from `2021-01-01T00:00:00Z` through
  `2026-06-16T23:55:00Z`; all sorted streams had `gap_count=0` and
  `duplicate_count=0`; zero-volume counts were BTC `66`, ETH `47`, SOL `47`;
  SOL had one physical non-monotonic row and was accepted only after sorting.
- Required anti-leakage rule: forward labels may appear only in label, cohort,
  ranking, and summary artifacts, never premise, state-ID, router, gating, or
  feature-bucket inputs.
- Required common-output rule: `summary.json`, `summary.csv`, and `trades.json`
  must remain zero-trade compatible for the later audit.
- Rejection criteria include closed-family reslice, broad mining, source gap,
  hidden future-label input, structured-compression rescue, ETH/SOL replay,
  BTC promotion, entries, exits, P&L backtests, optimizer grids, replay, or
  walk-forward.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to wait for explicit user approval
  before any implementation.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes.

BTC regime plus ETH/SOL context scope review:

- Added docs-only review:
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md`.
- Approved `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md` as materially
  different enough and tightly scoped enough for a separate zero-trade audit
  brief-writing task, not audit implementation.
- Stop state:
  `btc_regime_eth_sol_context_scope_review_approved_needs_zero_trade_audit_brief`.
- Allowed source scope for the later brief is only the already local Binance
  USDT-M futures `5m` files:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Source facts preserved from prior validation: each file has `573,984` loaded
  candles from `2021-01-01T00:00:00Z` through
  `2026-06-16T23:55:00Z`; all sorted streams had `gap_count=0` and
  `duplicate_count=0`; zero-volume counts were BTC `66`, ETH `47`, SOL `47`;
  SOL had one physical non-monotonic row and was accepted only after sorting.
- BTCUSDT role is market-regime context and diagnostic-only authority row.
  ETHUSDT/SOLUSDT role is possible authority rows only in a later zero-trade
  context audit, not strategy promotion.
- Minimum next audit question: whether BTC regime buckets improve separation of
  ETH/SOL usable, toxic, rotation, continuation, or no-trade range states.
- Rejection criteria preserved: closed-family reslice, broad mining, source
  gap, hidden future-label input, or any move toward entries/backtests.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a separate zero-trade audit brief
  task that stops before code at
  `btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval`.
- Commands run:
  - `wc -l ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: source files existed with `573,985` CSV lines each
  including headers; reference scan found canonical `memory/NEXT_CODEX_BRIEF.md`
  references and checklist mentions only; `git diff --check` passed;
  pre-commit `git status --short` showed only intended docs and memory changes.

BTC/ETH/SOL context review handoff selection:

- Reviewed `memory/NEXT_CODEX_BRIEF.md` valid choices and selected BTC regime
  plus ETH/SOL context as the next bounded lane to review.
- The next task is a documentation-only scope approval review, not audit
  implementation.
- Recommended ranking preserved in memory:
  1. BTC regime plus ETH/SOL context first;
  2. derivatives market-data context second, pending source/alignment approval;
  3. spread-range source/engine work third, pending engine/source scope;
  4. volatility-aware exits only after a future independent entry premise shows
     gross edge before costs.
- The next brief should create or update a docs-only review that decides whether
  the existing parked BTC/ETH/SOL context spec can become a zero-trade audit
  brief. It must preserve the failed BTCUSDT price-only verdict and avoid
  structured-compression rescue.
- Stop state:
  `btc_regime_eth_sol_context_scope_review_selected_for_next_brief`.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; `git status --short` showed only intended changes
  to `memory/DECISIONS.md`, `memory/NEXT_CODEX_BRIEF.md`, and
  `memory/PROGRESS.md`.

## 2026-06-27

Future range-strategy directions and specs:

- Added research map:
  `docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md`.
- Added next implementation-ready spec:
  `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md`.
- Added parked follow-up specs:
  - `docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md`;
  - `docs/FUTURES_VOLATILITY_AWARE_EXIT_MODEL_SPEC.md`;
  - `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md`;
  - `docs/FUTURES_SPREAD_RANGE_STRATEGY_SPEC.md`;
  - `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md`.
- Updated `README.md` docs index so future sessions can discover the new docs.
- Updated memory and next Codex brief to make the non-trading range-state audit
  the next concrete task.
- Stop state:
  `range_state_construction_loop_spec_ready_for_audit_implementation`.
- This was a connector-only documentation/memory update. No local Go tests,
  `rg`, `git diff --check`, or `git status` were run in this chat. The next
  Codex brief requires those checks during implementation closeout.

Futures range-state construction loop audit implementation:

- Added `-futures-range-state-construction-loop-audit`.
- Added audit implementation and tests for source/coverage gates, closed-candle
  state construction, forward-label separation, route cohort gates, stop-state
  selection, CLI artifacts, spot rejection, and conflicting flag rejection.
- Result directory:
  `results/futures-range-state-construction-loop-audit/`.
- Source facts: `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed UTC resample rows: `15m=191,328`, `1h=47,832`, `4h=11,958`; all
  accepted with `gap_count=0`, `duplicate_count=0`,
  `missing_child_open_count=0`, and complete buckets.
- Generated artifact counts: states `29,784`, labels `89,352`, cohorts
  `68,796`, rankings `16,335`, passing cohorts `58`.
- Passing route counts: `no_trade_toxic=52`,
  `tradable_rotation_candidate=6`, `trend_continuation_candidate=0`.
- Common outputs stayed zero-trade compatible; `trades.json` contains no trades.
- Review doc:
  `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md`.
- Stop state:
  `range_state_construction_loop_audit_passed_needs_router_spec`.
- Commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-state-construction-loop-audit -out-dir results/futures-range-state-construction-loop-audit`
  - `wc -l results/futures-range-state-construction-loop-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced `state_rows=29784`,
  `label_rows=89352`, `cohort_rows=68796`, `ranking_rows=16335`,
  `passing_cohorts=58`, and stop state
  `range_state_construction_loop_audit_passed_needs_router_spec`;
  CSV line-count check totaled `204,345`; `rg` found only canonical
  `memory/NEXT_CODEX_BRIEF.md` references; `git diff --check` passed; pre-commit
  status showed only intended code, docs, and memory changes.
- That next step was completed by the router audit below; the state audit still
  does not authorize an entry strategy by itself.

Futures range context router audit implementation:

- Added `-futures-range-context-router-audit`.
- Added audit implementation and tests for source/coverage inheritance, state
  rollup rule matching, forward-label separation, conflicting-rule preservation,
  router cohort gates, stop-state selection, CLI artifacts, spot rejection, and
  conflicting flag rejection.
- Result directory:
  `results/futures-range-context-router-audit/`.
- Source facts: `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed UTC resample rows remained `15m=191,328`, `1h=47,832`,
  `4h=11,958`; all accepted with `gap_count=0`, `duplicate_count=0`,
  `missing_child_open_count=0`, and complete buckets.
- Generated artifact counts: rules `58`, router rows `29,784`, cohorts `84`,
  rankings `12`, passing cohorts `3`.
- Router row counts: `no_trade=13,546`, `tradable_rotation=1,299`,
  `trend_continuation=0`, `diagnostic_only=14,939`, conflicts `0`.
- Passing router cohorts: `no_trade=2`, `tradable_rotation=1`,
  `trend_continuation=0`.
- Common outputs stayed zero-trade compatible; `trades.json` contains no trades.
- Review doc:
  `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.
- Stop state:
  `range_context_router_passed_needs_rotation_premise_spec`.
- Commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-context-router-audit -out-dir results/futures-range-context-router-audit`
  - `wc -l results/futures-range-context-router-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced `rule_rows=58`,
  `router_rows=29784`, `cohort_rows=84`, `ranking_rows=12`,
  `passing_cohorts=3`, and stop state
  `range_context_router_passed_needs_rotation_premise_spec`; CSV line-count
  check totaled `30,024`; `rg` found canonical `memory/NEXT_CODEX_BRIEF.md`
  references and checklist mentions only; `git diff --check` passed;
  pre-commit status showed only intended code, docs, and memory changes.
- That next step was completed by the router rotation premise spec below; the
  router result still does not authorize entries, exits, backtests, optimizers,
  replay, or walk-forward.

Futures range router rotation premise spec:

- Added docs-only spec:
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md`.
- The spec starts from the passed router cohort
  `range_context_router_v1|15m|h24|tradable_rotation` and names the non-trading
  event premise `router_gated_boundary_reclaim_rotation_v1`.
- The router remains closed-candle context only, not an entry signal. The
  `1,299` `tradable_rotation` rows must not be converted directly into trades.
- The next permitted implementation is only a zero-trade non-trading audit
  behind `-futures-range-router-rotation-premise-audit`, with default output
  directory `results/futures-range-router-rotation-premise-audit/`.
- The spec requires context segments from `15m` `h24` `tradable_rotation`
  matched rules, frozen range bounds at the segment start, a subsequent closed
  in-range boundary-reclaim event, and labels beginning strictly after the event
  candle.
- Required next-audit stop states include
  `range_router_rotation_premise_audit_source_router_gap`,
  `range_router_rotation_premise_audit_rejected_closed_family_reslice`,
  `range_router_rotation_premise_audit_no_eligible_events`,
  `range_router_rotation_premise_audit_failed_no_premise`,
  `range_router_rotation_premise_audit_passed_needs_non_trading_trigger_audit`,
  and
  `range_router_rotation_premise_audit_rejected_as_strategy_backtest_request`.
- Stop state:
  `range_router_rotation_premise_spec_ready_for_non_trading_audit`.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: brief-reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit status showed only intended docs and
  memory changes.
- Current next step: implement only the non-trading audit from the premise spec;
  do not add entries, exits, P&L backtests, optimizers, replay, walk-forward,
  paper/testnet/live paths, source expansion, or symbol expansion.

Futures range router rotation premise audit implementation:

- Added `-futures-range-router-rotation-premise-audit`.
- Added audit implementation and tests for source/router dependency checks,
  router context collapse, closed-candle event construction, outcome labeling,
  gate failures, stop-state selection, CLI artifacts, spot rejection, default
  output behavior, and conflicting flag rejection.
- Result directory:
  `results/futures-range-router-rotation-premise-audit/`.
- Source facts: `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed UTC `15m` resample rows: `191,328`, from
  `2021-01-01T00:00:00Z` through `2026-06-16T23:45:00Z`, accepted with
  `gap_count=0`, `duplicate_count=0`, `missing_child_open_count=0`, and
  `complete=true`.
- Router dependency passed with stop state
  `range_context_router_passed_needs_rotation_premise_spec`; the required
  cohort was `range_context_router_v1|15m|h24|tradable_rotation`.
- Generated artifact counts: context segments `278`, events `97`, outcomes
  `97`, cohorts `12`, rankings `3`, passing cohorts `0`.
- Full-period event inventory: lower events `43`, upper events `54`, midline
  outcomes `71`, hard adverse outcomes `22`, chop/no-resolution outcomes `4`.
- Top failure reasons:
  `inadequate_event_count,inadequate_split_event_count,single_split_contribution_above_gate,behavior_gate_failed`.
- Common outputs stayed zero-trade compatible; `summary.csv` has `0` trades in
  every split/side row.
- Review doc:
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md`.
- Stop state:
  `range_router_rotation_premise_audit_failed_no_premise`.
- Commands run:
  - `gofmt -w internal/lab/futures_range_context_router_audit.go internal/lab/futures_range_router_rotation_premise_audit.go internal/lab/futures_range_router_rotation_premise_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-router-rotation-premise-audit -out-dir results/futures-range-router-rotation-premise-audit`
  - `wc -l results/futures-range-router-rotation-premise-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced
  `context_segments=278`, `events=97`, `outcomes=97`, `cohort_rows=12`,
  `ranking_rows=3`, `passing_cohorts=0`, and stop state
  `range_router_rotation_premise_audit_failed_no_premise`; CSV line-count check
  totaled `529`; `rg` found canonical `memory/NEXT_CODEX_BRIEF.md` references
  and checklist mentions only; `git diff --check` passed; pre-commit status
  showed only intended code, docs, and memory changes.
- Current next step: do not build entries, exits, P&L backtests, optimizers,
  replay, walk-forward, or trigger audits from this premise. Any follow-up must
  be a materially different non-trading premise or context audit.

Recent failed premise evidence to preserve:

- `docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md`: source/resampling passed,
  but no range-quality, UTC-session, or quality-plus-session cohort passed the
  declared gates; stop state `range_context_triage_failed_no_strategy_premise`.
- `docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md`: fixed
  baseline lost after costs and `0` of `1,152` grid rows passed selection gates;
  stop state `range_first_strategy_v1_optimizer_failed_no_replay`.
- `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md`:
  ETH/SOL authority stream stayed positive in the frozen replay but failed
  walk-forward robustness; no package approved.
- `docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md`:
  selected `15m` and `1h` breakout-retest/acceptance baselines failed after
  costs across BTCUSDT, ETHUSDT, and SOLUSDT.
- `docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md`: only
  `3` valid events across the full BTCUSDT sample, so no baseline was approved.

## 2026-06-28

Futures range post-rotation premise failure pivot review:

- Added docs-only review:
  `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md`.
- Preserved the failed audit verdict in
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md` and kept
  `router_gated_boundary_reclaim_rotation_v1` closed in reviewed form.
- Source facts remain:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Dependency facts preserved: state audit stopped at
  `range_state_construction_loop_audit_passed_needs_router_spec`; router audit
  stopped at `range_context_router_passed_needs_rotation_premise_spec`; premise
  audit stopped at `range_router_rotation_premise_audit_failed_no_premise`.
- Premise audit facts preserved: context segments `278`, events `97`, outcomes
  `97`, cohort rows `12`, ranking rows `3`, passing cohorts `0`, lower events
  `43`, upper events `54`, midline outcomes `71`, hard-adverse outcomes `22`,
  and chop/no-resolution outcomes `4`.
- Pivot decision: no materially different BTCUSDT-only, candle-price-only
  range-premise audit is worth specifying next from current evidence. Do not
  convert the `278` segments, `97` events, or `1,299` `tradable_rotation`
  router rows into trades.
- Stop state:
  `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
- Current next step: no automatic implementation. A future session should first
  obtain an explicit user scope choice before BTC/ETH/SOL context, derivatives
  context, spread-range source/engine work, or a no-further-audit stop is
  pursued.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: brief-reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit status showed only intended docs and
  memory changes.

## Milestone Index

Use `README.md` as the full docs index. The most relevant current docs are:

1. `docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md`.
2. `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md`.
3. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md`.
4. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md`.
5. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md`.
6. `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md`.
7. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md`.
8. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md`.
9. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md`.
10. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md`.
11. `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md`.
12. `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md`.
13. `docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md`.
14. `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.
15. `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md`.
16. `docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md`.
17. `docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md`.
18. `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md`.
19. `docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md`.
20. `memory/NEXT_CODEX_BRIEF.md`.

Historical details remain in the focused docs and git history rather than this
always-read memory file.
