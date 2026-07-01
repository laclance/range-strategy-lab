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
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`; Binance USDT-M
  futures `BTCUSDT` `5m`; `573,984` loaded candles; `2021-01-01T00:00:00Z`
  through `2026-06-16T23:55:00Z`; `gap_count=0`, `duplicate_count=0`,
  `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted`.
- Closed-candle decision semantics, next-bar-open entries, one-position max,
  stop-first ambiguity, cost accounting, split metrics, and source manifests
  remain core infrastructure.
- Research is not stopped, but failed fixed baselines and closed premises must
  not be silently retuned, rescued, renamed, or promoted.
- `memory/NEXT_CODEX_BRIEF.md` is the canonical next-session prompt.

## Latest Milestone

Strategy class pivot assessment:

- Review doc:
  `docs/STRATEGY_CLASS_PIVOT_ASSESSMENT.md`.
- Stop state:
  `strategy_class_pivot_assessment_recommends_trend_pullback`.
- Result: docs-only scoring ranked trend-pullback continuation first (`26`),
  volatility expansion / breakout continuation second (`25`), session-based
  opening-range expansion third (`25`), liquidity sweep plus reclaim fourth
  (`19`), and cross-asset or regime filter before entries fifth (`14`).
- Durable result: the current range-reversion / midpoint / edge-fade /
  previous-day range / bounded range-optimization path remains closed for now.
  The recommended next lane is a docs-only trend-pullback continuation candidate
  packet after operator approval.
- Exact next allowed artifact after approval:
  `docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_CANDIDATE_PACKET.md`.
- No backtest, optimizer, generated result, source mutation, exchange API,
  credential, deployment, paper/testnet/live flow, martingale, averaging down,
  two-exchange logic, or promotion was authorized or added.

## 2026-07-01 Milestone Index

Strategy class pivot assessment:

- Added docs-only assessment:
  `docs/STRATEGY_CLASS_PIVOT_ASSESSMENT.md`.
- Recommended trend-pullback continuation as the cleanest materially different
  next research class for a future backtest-first candidate packet.
- Updated `memory/NEXT_CODEX_BRIEF.md` to point only to the next bounded
  docs-only candidate packet after operator approval.

## 2026-06-30 Milestone Index

Range optimization workbench verified no-candidate result:

- Review doc:
  `docs/RANGE_OPTIMIZATION_WORKBENCH_IMPLEMENTATION_REVIEW.md`.
- CLI flag:
  `-range-optimization-workbench-v1`.
- Stop state:
  `range_optimization_workbench_failed_no_candidate`.
- Immutable run path:
  `results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/`.
- Result: source/resample passed; `112` trials ran; `112` trials were rejected;
  `0` passing candidates; selected candidate is empty.
- Artifact counts: coverage `2`, rejected candidates `113`, source contract `2`,
  top candidates `1`, trial results `113`, trial summary `1345`, total CSV lines
  `1576`.
- Durable result: no candidate is selected for locked fixed validation. Optimizer
  output from this run cannot authorize paper/testnet/live trading or promotion.
- Further search, if desired, requires a separately approved spec revision or a
  materially different research lane with explicit search-space changes and
  guardrails.

Range optimization workbench implementation:

- Added offline implementation review:
  `docs/RANGE_OPTIMIZATION_WORKBENCH_IMPLEMENTATION_REVIEW.md`.
- Added CLI flag:
  `-range-optimization-workbench-v1`.
- Earlier stop state before local run:
  `range_optimization_workbench_implementation_added_needs_local_run`.
- The implementation consumes the approved docs-only spec:
  `docs/RANGE_OPTIMIZATION_WORKBENCH_SPEC.md`.
- Workbench runs must write to immutable run directories under
  `results/range-optimization-workbench-v1/runs/<run_id>/`; normal verification
  must not delete the canonical workbench results parent.

Range optimization workbench spec:

- Added docs-only spec:
  `docs/RANGE_OPTIMIZATION_WORKBENCH_SPEC.md`.
- Stop state:
  `range_optimization_workbench_spec_ready_for_implementation_approval`.
- This spec exists because all three candidates from
  `docs/BACKTEST_FIRST_CANDIDATE_PACKET.md` failed as fixed baselines.
- The workbench is a controlled offline discovery lane for bounded combinations
  of failed range-family components plus allowed decision-time context. It does
  not authorize paper/testnet/live trading, exchange API work, credentials,
  deploy files, martingale, averaging down, two-exchange logic, or promotion.

Backtest-first candidate packet:

- Added docs-only candidate packet:
  `docs/BACKTEST_FIRST_CANDIDATE_PACKET.md`.
- Selected simple fixed baselines to move faster from docs into actual P&L tests.
- All selected fixed baselines later failed; no retune/rescue is authorized from
  this packet.

BTCUSDT `5m` rolling value-area reversion fixed baseline:

- Review doc:
  `docs/BACKTEST_FIRST_BTC_5M_ROLLING_VALUE_AREA_REVERSION_IMPLEMENTATION_REVIEW.md`.
- Candidate: `btc_5m_rolling_value_area_reversion_v1`.
- Stop state: `btc_5m_rolling_value_area_reversion_backtest_failed_no_usable_strategy`.
- Result: `20,144` trades; full gross P&L `-489.49542250301454`; full net P&L
  `-999.9999999684289`; full PF `0.6949768102636272`; full max drawdown
  `0.9999999999689176`.
- Failed gates: gross edge, net edge, drawdown.
- Closed boundary: do not rescue with alternate VWAP windows, outer-zone
  thresholds, target changes, time-stop changes, side selection, filters,
  derivatives-veto interaction, replay, walk-forward, or optimizer grids.

BTCUSDT `15m` previous-day range reversion fixed baseline:

- Review doc:
  `docs/BACKTEST_FIRST_BTC_15M_PREVIOUS_DAY_RANGE_REVERSION_IMPLEMENTATION_REVIEW.md`.
- Candidate: `btc_15m_previous_day_range_reversion_v1`.
- Stop state: `btc_15m_previous_day_range_reversion_backtest_failed_no_usable_strategy`.
- Result: `1,956` trades; full gross P&L `-506.1461042753044`; full net P&L
  `-935.263471879938`; full PF `0.5772171917112771`; full max drawdown
  `0.9388204368154983`.
- Failed gates: gross edge, net edge, drawdown.
- Closed boundary: do not rescue with alternate UTC sessions, previous 2/3 day
  windows, changed outer-decile thresholds, derivatives context,
  calendar/session mining, replay, walk-forward, or optimizer grids.

BTCUSDT `15m` range-edge exhaustion fade fixed baseline:

- Review doc:
  `docs/BACKTEST_FIRST_BTC_15M_RANGE_EDGE_EXHAUSTION_FADE_IMPLEMENTATION_REVIEW.md`.
- Candidate: `btc_15m_range_edge_exhaustion_fade_v1`.
- Stop state: `btc_15m_range_edge_exhaustion_fade_backtest_failed_no_usable_strategy`.
- Result: `156` trades; full gross P&L `-154.40528599997904`; full net P&L
  `-261.59525647142874`; full PF `0.48125879295748447`; full max drawdown
  `0.28473381700333156`.
- Failed gates: gross edge, net edge, drawdown.
- Closed boundary: do not rescue with alternate range windows, progress
  thresholds, edge zones, midpoint variants, added volume filters, derivatives
  context, replay, walk-forward, or optimizer grids.

BTCUSDT `15m` post-compression directional expansion fixed backtest:

- Review doc:
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_REVIEW.md`.
- Candidate: `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`.
- Stop state: `post_compression_directional_expansion_backtest_failed_no_usable_strategy`.
- Result: reproduced `468` raw candidate rows, executed `421` trades, passed the
  trade-count gate, but failed gross edge, extra slippage-stress edge, stress PF,
  and drawdown gates. Full gross P&L was `208.560999`, engine net was
  `-129.258571`, extra slippage-stress net was `-227.226250`, full stress PF was
  `0.799666`, and full stress max drawdown was `0.289326`.
- Closed boundary: no adjacent-cell rescue, exit retune, derivatives-veto
  interaction, replay, walk-forward, paper/testnet/live path, or promotion is
  authorized.

Derivatives context/no-trade line:

- Durable source materialization produced local/offline Binance USDT-M futures
  mark/index/premium `5m` files under `../binance-bot/data/derivatives/`.
- Derivatives source audit passed and proved conservative one-`5m` lag alignment.
- Derivatives context audit passed as zero-trade context separation only.
- No-trade filter premise audit passed as a veto-premise artifact only.
- Integration spec deferred the canonical veto until a separately approved entry
  premise exists. The veto may only annotate an approved candidate stream as
  skipped/retained; it may not create entries, choose side, rank, score P&L,
  replay, walk forward, optimize, promote, expand sources/symbols, or reopen
  closed families.

Older reviewed range/premise work:

- Historical reviewed docs remain authoritative in `docs/` and git history.
- The always-read memory no longer repeats every older review in full; use the
  docs index and focused review files when older evidence is needed.
