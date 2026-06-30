# Decisions

## Durable Constraints

- This project is offline research only.
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
- After completing a brief or milestone, Codex should run closeout checks and
  commit completed repo changes unless the user explicitly says not to commit.

## Current Research Decision

- The post-rotation premise failure pivot review in
  `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md` stopped at
  `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
  No automatic next BTCUSDT-only, candle-price-only range-premise audit is
  selected from the current state/router/premise evidence.
- The BTC regime plus ETH/SOL context scope review in
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md` approved
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md` only for a separate
  zero-trade audit brief-writing task. User then explicitly approved the audit
  implementation described in
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md`.
- The implemented audit in
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md` stopped
  at
  `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.
  Source validation, resampled coverage, zero-trade common outputs, and
  anti-leakage checks passed, but `0` BTC-regime-plus-ETH/SOL context cohorts
  passed the declared gates.
- The BTC regime plus ETH/SOL context path is closed in reviewed zero-trade
  form. Do not retune, reslice, rename, gate-relax, replay, walk-forward,
  promote BTC regime rows, treat ETH/SOL context rows as strategy authority, or
  convert this result into entries, exits, P&L backtests, optimizer grids,
  source downloads, paper/testnet/live paths, exchange API, credentials, deploy
  files, martingale, averaging down, or two-exchange logic.
- The only source scope used and approved for that audit was the already local
  Binance USDT-M futures `5m` BTCUSDT, ETHUSDT, and SOLUSDT files. BTCUSDT
  remains market-regime context and diagnostic-only authority. ETHUSDT/SOLUSDT
  remain failed zero-trade context authority candidates only, not strategy
  promotion.
- The derivatives market-data context source scope review in
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md` approved only a
  separate zero-trade source-audit brief-writing task and stopped at
  `derivatives_context_source_scope_review_approved_needs_zero_trade_source_audit_brief`.
  It found no approved durable derivatives market-data source rows in the lab's
  current local data scope. A later brief may reference adjacent
  `../binance-bot/research/` source-proof artifacts only as source/process
  evidence and may use existing BTC/ETH/SOL Binance USDT-M futures `5m` candle
  files only as alignment anchors.
- The derivatives context zero-trade source-audit brief in
  `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md` stopped
  at
  `derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval`.
  It selects Binance USDT-M futures mark/index/premium basis klines as the only
  first source family for a possible later source audit. Funding remains second,
  and aggregate trades remain parked as high-volume secondary source proof.
- The derivatives context source materialization plan in
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md` stopped at
  `derivatives_context_source_materialization_plan_ready_for_execution_approval`.
  User approved a docs-only offline materialization plan for Binance public Data
  Vision USDT-M futures `markPriceKlines`, `indexPriceKlines`, and optional
  `premiumIndexKlines` `5m` archive objects for `BTCUSDT`, `ETHUSDT`, and
  `SOLUSDT` under `../binance-bot/data/derivatives/`.
- Derivatives context implementation is not authorized by the brief or
  materialization plan alone. The next permitted derivatives step is only an
  explicit execution approval gate for source materialization. Passing
  materialization would create durable local/offline source files, but would
  still require separate explicit approval before the zero-trade source audit.
  Open interest, long/short ratios, liquidation/force-order history,
  order-book/depth, funding, aggregate trades, taker flow, unapproved source
  downloads, live probes, private endpoints, API keys, credentials,
  context-gain features, labels, cohorts, rankings, entries, exits, P&L
  backtests, replay, walk-forward, and promotion remain forbidden from the
  materialization and first source-audit scope.
- The derivatives context source materialization in
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md` was
  explicitly approved and executed. It passed at
  `derivatives_context_source_materialization_passed_ready_for_source_audit_approval`.
  Durable Binance public Data Vision USDT-M futures `markPriceKlines`,
  `indexPriceKlines`, and optional `premiumIndexKlines` `5m` source files for
  `BTCUSDT`/`ETHUSDT`/`SOLUSDT` now exist under
  `../binance-bot/data/derivatives/` (`729` checksum-verified raw zips, `9`
  normalized CSVs, `5` manifests). These files are outside this repo and are not
  tracked by Git; the generator was a one-shot offline scratchpad tool, not
  tracked lab code.
- Gap-handling rule established for these derivatives source files: the required
  mark/index `5m` public archives contain real, whole-day-aligned
  publication-outage gaps (required missing `9,820` of `3,443,904`, `0.285%`;
  trade-candle anchors had `gap_count=0` over the same era). The user decided to
  record gaps in the manifests (`gap_count`, `missing_interval_count`,
  no-imputation policy) and pass, deferring bounded-missingness gating to the
  later zero-trade source audit, rather than fail closed on gaps. Integrity
  faults (duplicate-conflicting rows, schema ambiguity, checksum mismatch, or any
  missing required object) remain fail-closed and would reject materialization.
- These materialized rows are durable candidate source inputs only. They are not
  approved context inputs and do not authorize the source-audit implementation,
  context features, labels, cohorts, rankings, entries, exits, P&L backtests,
  replay, walk-forward, or promotion. The zero-trade source audit over them needs
  separate explicit user approval.
- The derivatives context zero-trade source audit in
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md` was explicitly
  approved and implemented behind `-futures-derivatives-context-source-audit`. It
  passed at
  `derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief`.
  It validated the `9` materialized mark/index/premium `5m` source files,
  SHA-256-bound their provenance (hashes match the materialization manifests),
  and proved anti-lookahead alignment to the `5m` candle anchors. All `6`
  required mark/index streams cleared the `0.99` coverage bar under a
  conservative one-`5m`-interval lag (min `0.994472`, index BTCUSDT) with
  recorded missingness and no forward fill; the run produced `0` trades.
- Derivatives alignment finality rule established: publication lag for
  mark/index/premium klines is unproven, so any future use must apply the
  conservative one-native-interval lag
  (`source_close_time + 5m <= decision_candle_close_time`); no forward fill,
  interpolation, or nearest-future joins. The materialization gaps are
  bounded missingness to be surfaced, never silently filled.
- The passing source audit does not authorize context-gain implementation,
  labels, cohorts, rankings, entries, exits, P&L backtests, optimizer grids,
  replay, walk-forward, packaging, source downloads, paper/testnet/live paths,
  exchange API work, credentials, deploy files, or promotion. The next
  derivatives step is a separate, approval-gated zero-trade derivatives
  context-audit brief. Mark-minus-index basis derivation remains deferred to
  that later context stage and was intentionally not computed as an artifact in
  the source audit.
- The derivatives context zero-trade context-audit brief in
  `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md` stopped at
  `derivatives_context_zero_trade_context_audit_brief_ready_for_user_approval`.
  It defines the only later derivatives context-audit question: whether
  mark-minus-index basis, premium-index level, or basis-change buckets known at
  closed decision-candle time improve separation of BTCUSDT/ETHUSDT/SOLUSDT local
  range states (usable, toxic, rotation, continuation, no-trade) beyond the local
  price/volume state alone. The later audit must not test basis tradability and
  must not measure entry/exit/P&L. Approved later inputs are only the `9`
  validated derivatives CSVs plus the `3` candle anchors; carried-forward facts
  are SHA-256 provenance, the conservative one-`5m`-interval lag
  (`source_close_time + 5m <= decision_candle_close_time`), exact closed-interval
  joins, no forward fill/interpolation/nearest-future joins, and bounded recorded
  missingness with a required basis context coverage floor of `0.994472`.
- Durable boundary for that later audit: basis/premium is treated as an
  orthogonal source (perp-vs-index dislocation), materially different from the
  closed price-only and BTC-regime families. All three symbols may be local
  range-state authority candidates there, but the audit must reject itself
  (orthogonality gate plus
  `derivatives_context_zero_trade_context_audit_rejected_closed_family_rescue`)
  if basis buckets are collinear with or dominated by the local price/volume
  state, and it must not reopen, retune, rename, gate-relax, or promote any
  closed family. Forward labels stay in label/cohort/ranking/summary artifacts
  only, basis/premium context is built only from lagged closed source rows,
  missing context produces missingness/skip rows (never silent defaults), and
  `summary.json`/`summary.csv`/`trades.json` stay zero-trade compatible. At
  brief closeout, it authorized no implementation and required explicit user
  approval before the context-audit implementation.
- The derivatives context zero-trade audit in
  `docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md` was explicitly approved and
  implemented behind `-futures-derivatives-context-audit`. It passed at
  `derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec`.
  It used only the `9` validated derivatives mark/index/premium `5m` CSVs plus
  the `3` candle anchors, enforced the conservative one-`5m` source lag, kept
  forward labels out of feature/state/gating inputs, recorded missing context as
  skips with no fill/interpolation/nearest-future joins, and produced `0`
  trades.
- Durable boundary from that pass: only the `6` passing BTCUSDT `15m`
  derivatives-context cohorts are candidates for a later strategy-premise spec
  (`5` no-trade/toxic, `1` rotation candidate). ETHUSDT and SOLUSDT produced
  `0` passing cohorts and are not promoted. The pass means basis/premium context
  showed separation beyond the local price/volume state under the declared
  gates; it does not authorize entries, exits, P&L backtests, optimizers, replay,
  walk-forward, source expansion, paper/testnet/live paths, exchange API work,
  credentials, deploy files, or strategy promotion.
- The strategy-premise decision from the passed derivatives context audit is no
  longer open: the later track is the BTCUSDT `15m` no-trade filter premise
  only. The rotation-entry, two-track, and no-further-work alternatives were
  considered in the docs-only spec; rotation remains diagnostic only.
- The derivatives context strategy-premise spec in
  `docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md` was explicitly
  approved as a docs-only task and stopped at
  `derivatives_context_strategy_premise_spec_ready_for_user_approval`. It selects
  only a BTCUSDT `15m` derivatives-context no-trade filter premise for a later
  zero-trade audit. The `5` toxic/no-trade cohorts are candidates for exact and
  canonical no-trade veto testing; the single rotation candidate remains
  diagnostic only and does not authorize a rotation-entry premise or a second
  implementation track.
- Durable boundary from the strategy-premise spec: a later no-trade filter audit
  may use only local BTCUSDT Binance USDT-M futures candles plus the already
  validated BTCUSDT mark/index/premium `5m` rows needed by the selected premise,
  with confirmed closed-candle decisions, the conservative one-`5m` source lag
  (`source_close_time + 5m <= decision_candle_close_time`), no forward
  fill/interpolation/nearest-future joins, recorded missingness/skips, and
  forward labels only as evaluation metadata. It must stay zero-trade unless a
  later explicit spec changes scope.
- The strategy-premise spec does not authorize Go code, CLI flags, generated
  results, audit runs, source downloads, data writes under
  `../binance-bot/data/derivatives/`, entries, exits, P&L backtests, optimizer
  grids, replay, walk-forward, source/symbol expansion, paper/testnet/live paths,
  exchange API work, credentials, deploy files, strategy promotion, martingale,
  averaging down, two-exchange logic, or closed-family rescue.
- The derivatives no-trade filter premise audit in
  `docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_PREMISE_AUDIT_REVIEW.md`
  was explicitly approved and implemented behind
  `-futures-derivatives-no-trade-filter-premise-audit`. It passed at
  `derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec`.
  It reproduced all `5` BTCUSDT `15m` `h48` exact toxic rows, built a
  de-duplicated canonical veto union (`1,823` rows, `1,241` no-trade toxic,
  full toxic rate `0.680746`, min split toxic rate `0.665485`), reported
  overlaps and collateral damage, and produced `0` trades.
- Durable boundary from the no-trade filter premise audit: the only next
  authorized direction is a separate docs-only filter integration spec. The
  passing veto premise may not be treated as an entry signal, "trade the
  opposite" rule, basis-tradability claim, P&L result, replay/walk-forward
  result, rotation-entry rescue, source/symbol expansion, or strategy promotion.
  The diagnostic rotation row remains diagnostic only, and no closed family is
  reopened.
- The derivatives no-trade filter integration spec in
  `docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_INTEGRATION_SPEC.md` was
  explicitly approved as a docs-only task and stopped at
  `derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise`.
  The canonical veto `btc_15m_basis_discount_no_trade_veto_v1` is preserved as a
  future veto candidate only, but no implementation gate is selected because no
  independent entry premise currently exists for it to filter.
- Durable boundary from the integration spec: a future interaction audit may be
  considered only after a separate independent entry premise exists and is
  explicitly approved for interaction testing. The veto may only annotate an
  already-approved candidate stream as skipped or retained; it may not create
  entries, change entry logic, choose side, act as an exit, rank trades, score
  P&L, replay, walk forward, optimize, promote a strategy, expand sources or
  symbols, or reopen closed families.
- The independent entry-premise and hypothesis map in
  `docs/FUTURES_INDEPENDENT_ENTRY_PREMISE_AND_HYPOTHESIS_MAP.md` was explicitly
  approved as a docs-only task and stopped at
  `independent_entry_premise_and_hypothesis_map_needs_user_scope_choice`. No
  single BTCUSDT `15m` local-source independent entry premise is selected from
  the current reviewed evidence.
- Durable boundary from the independent-entry map: future work must first choose
  exactly one route: a new BTCUSDT `15m` closed-candle local-source event
  premise, a higher-timeframe premise, spread-range/source scope, another
  explicitly approved source family, or no further audit. The derivatives veto
  remains parked as future skip/retain evidence only and may not shape candidate
  rows, create entries, choose side, rank, score P&L, replay, walk forward,
  optimize, promote a strategy, or reopen closed families.
- The BTCUSDT `15m` post-compression directional expansion premise spec in
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_PREMISE_SPEC.md`
  was explicitly approved as a docs-only task and stopped at
  `independent_entry_premise_spec_ready_for_user_approval`. It supersedes the
  independent-entry map's user-choice gate by selecting exactly one new
  BTCUSDT `15m` local-source premise family for a later zero-trade audit:
  `btc_15m_post_compression_directional_expansion_v1`.
- Durable boundary from the post-compression premise spec: the later audit may
  test only the predeclared local-source parameter family: compression lookbacks
  `48`/`96`/`192`, compression thresholds bottom `20%`/`30%`/`40%` of prior
  `1,920` closed `15m` range-width observations, breakout beyond prior range by
  `0.1`/`0.2`/`0.3` prior-bar `ATR(14)`, and volume confirmation `none`, above
  prior `96`-bar median, or above prior `96`-bar `60%` percentile. Candidate
  viability must use de-duplicated `(decision_close, side)` rows; passing
  evidence must separate from the unconditional eligible `15m` baseline across
  splits and cannot be a single isolated parameter cell. The audit must reject
  itself if it becomes a closed-family reslice or uses derivatives veto facts to
  shape candidate rows, side, scoring, or pass/fail decisions.
- The BTCUSDT `15m` post-compression directional expansion zero-trade audit in
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md`
  was explicitly approved for implementation and passed at
  `btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review`.
  It produced `0` trades and found only a narrow long-side `48`-bar diagnostic
  pocket at lookback `192` and bottom `20%` compression, across adjacent
  breakout/volume cells. No short-side, `16`-bar, `32`-bar, lookback `48`/`96`,
  or compression threshold `30%`/`40%` surface passed the full gate.
- Durable boundary from the post-compression audit: this pass is zero-trade
  label-separation evidence only. It authorizes no entries, exits, P&L,
  optimizer selection, replay, walk-forward, derivatives veto interaction, or
  promotion. A later docs-only strategy-premise spec must decide whether the
  narrow long `48`-bar pocket justifies requesting a separate offline backtest
  spec, or whether the line should stop.
- The BTCUSDT `15m` post-compression directional expansion strategy-premise spec
  in
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_STRATEGY_PREMISE_SPEC.md`
  was explicitly approved as a docs-only task and stopped at
  `post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval`.
  It selects exactly one conservative representative candidate for a later
  docs-only offline backtest spec:
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`.
- Durable boundary from the strategy-premise spec: any later backtest spec may
  use only this candidate stream: BTCUSDT Binance USDT-M futures `15m` closed
  decision candles, long side only, prior `192`-bar local range, bottom `20%`
  compression threshold against the prior `1,920` valid range-width
  observations, decision close above the prior range by `0.2` prior-bar
  `ATR(14)`, no volume confirmation, and next-`15m`-open timing. The adjacent
  passing zero-trade cells are supporting robustness evidence only; they may not
  become a P&L optimizer, rescue a failed representative-cell backtest, or
  reopen the full `81`-cell grid. The later docs-only backtest spec must choose
  one fixed risk/exit model before implementation or stop at a user-choice
  gate. The derivatives veto remains parked and may not shape entries, exits,
  side, ranking, scoring, P&L, pass/fail decisions, replay, walk-forward, or
  promotion.
- The BTCUSDT `15m` post-compression directional expansion backtest spec in
  `docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_SPEC.md`
  was explicitly approved as a docs-only task and stopped at
  `post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval`.
  It defines exactly one later offline implementation candidate:
  `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1`.
- Durable boundary from the backtest spec: the later implementation may add
  only one offline backtest flag for this candidate stream, with closed UTC
  `15m` BTCUSDT futures decisions, long side only, prior `192`-bar range,
  bottom `20%` compression against the prior `1,920` valid range-width
  observations, decision close above prior range by `0.2` prior-bar `ATR(14)`,
  no volume filter, and next-`15m`-open entry timing. The fixed risk/exit model
  is stop at `entry_price - 1.0 * ATR(14)[d-1]`, target at
  `entry_price + 2.0 * ATR(14)[d-1]`, max hold `48` closed `15m` bars,
  one-position max, stop-first ambiguity, `1%` risk-at-stop sizing, `1x`
  notional cap, start balance `1000`, fee `0.0004` per side, and slippage
  `0.000116` per side. The implementation must reproduce `468` representative
  raw candidate rows before one-position filtering, report both current-engine
  net and an extra slippage-stress net, and use the stress view for pass/fail.
  It may not add stop/target/hold/cell/volume/side/veto grids,
  adjacent-cell P&L rescue, replay, walk-forward, derivatives veto interaction,
  or promotion.
- Spread-range source/engine work remains parked; it does not authorize
  implementation from current state. Volatility-aware exits remain unavailable
  until a future independent entry premise first shows gross edge before costs.
- The futures range router rotation premise audit in
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md` implemented the
  zero-trade audit behind
  `-futures-range-router-rotation-premise-audit` and stopped at
  `range_router_rotation_premise_audit_failed_no_premise`.
- The reviewed premise `router_gated_boundary_reclaim_rotation_v1` is closed in
  its reviewed form. Do not convert its `278` context segments, `97`
  boundary-reclaim events, or the router's `1,299` `tradable_rotation` rows
  into trades.
- This failed premise audit does not authorize a non-trading trigger audit,
  entries, exits, P&L backtests, optimizer grids, fixed replay, walk-forward,
  packaging, source expansion, symbol expansion, live-adjacent work, or
  closed-family retuning.
- Any follow-up must be a materially different non-trading premise or context
  audit, not a retune, rename, gate relaxation, or direct strategy conversion of
  this boundary-reclaim surface.
- The futures range router rotation premise spec in
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md` is now historical
  dependency context for the failed audit above.
- The futures range context router audit in
  `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md` passed as a non-trading
  route-selection milestone and stopped at
  `range_context_router_passed_needs_rotation_premise_spec`.
- The router found `2` passing no-trade cohorts and `1` passing
  `tradable_rotation` cohort, with `0` passing `trend_continuation` cohorts.
  This authorized only the materially new rotation premise spec that has now
  failed audit review.
- The router audit does not authorize entries, exits, P&L backtests, optimizer
  grids, fixed replay, walk-forward, packaging, source expansion, symbol
  expansion, live-adjacent work, or closed-family retuning.
- Any rotation premise spec must explain why it is materially different from
  `range_occupancy_rotation_v1`, hold-inside/midline, breakout-retest/
  acceptance, clean breakout continuation, structured compression, impulse
  absorption, and higher-timeframe nested range rotation.

## Parked Future Direction Decisions

The following specs are parked and not implementation-ready from current state:

- `docs/FUTURES_VOLATILITY_AWARE_EXIT_MODEL_SPEC.md`: may start only after a
  materially new entry template first shows gross edge before costs.
- `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md`: scope-approved, converted
  into
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md`, then
  implemented in
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md`; the
  audit failed with no usable context and is closed in reviewed zero-trade form.
- `docs/FUTURES_SPREAD_RANGE_STRATEGY_SPEC.md`: may start only with explicit
  engine/source approval; spread trading requires a separate multi-leg engine
  spec before any P&L strategy work.
- `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md`: may start only
  through the docs-only source scope and source-audit brief boundary in
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md` and
  `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md`; the
  source materialization, source audit, context-audit brief, and context audit
  have since executed in reviewed form. Further derivatives source expansion is
  not an automatic next step; the current strategy-premise spec creates no
  source-expansion need. It is market-data context only and does not permit API
  keys, private endpoints, live/paper/testnet, exchange order paths, entries,
  exits, P&L backtests, replay, walk-forward, or promotion.

## Exclusion Decisions

Reviewed failed or fragile families must not be retuned, renamed, gate-relaxed,
or promoted from their reviewed forms:

- `router_gated_boundary_reclaim_rotation_v1`;
- structured compression, including the fragile ETH/SOL authority stream;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline;
- impulse absorption after abnormal OHLCV candles;
- higher-timeframe nested range rotation;
- `range_occupancy_rotation_v1`;
- range quality, UTC session, and failure-mode triage cohorts by themselves;
- legacy spot-only SR rejection, confirmation, false-break, and compression
  promotion evidence.

These branches remain useful as exclusion evidence and infrastructure only.
Reusable infrastructure includes source guards, closed UTC resampling,
closed-candle semantics, event labels as labels only, split metrics, artifact
patterns, and detector/range episode helpers as feature extraction only.

## Historical Reviewed Decisions

- The futures range-context triage audit passed source/resampling but failed to
  find a gated strategy premise. Do not create a strategy, baseline, optimizer,
  replay, walk-forward, package, retune, source expansion, symbol expansion, or
  live-adjacent path from that audit.
- The futures range-state construction loop audit passed as a mixed router
  discovery milestone, not a strategy premise. It identified both toxic filters
  and rotation candidates, so do not create an entry, backtest, optimizer,
  replay, walk-forward, package, source expansion, symbol expansion, or
  live-adjacent path directly from it.
- The futures range context router audit passed as a rotation-premise-spec
  milestone, not a strategy premise. Do not create an entry, backtest,
  optimizer, replay, walk-forward, package, source expansion, symbol expansion,
  or live-adjacent path directly from it.
- The futures range-first occupancy rotation V1 optimizer evaluated the declared
  `1,152` row grid, the fixed baseline lost after costs, and `0` grid rows
  passed. Do not create a fixed replay or retune the grammar.
- The structured-compression universe stream selected ETH/SOL authority with BTC
  diagnostic-only, but walk-forward robustness was fragile. Do not package,
  retune, promote BTC, or broaden the grid from that result.
- The breakout-retest/acceptance baseline failed after costs on BTCUSDT,
  ETHUSDT, and SOLUSDT. Do not optimize or robustness-review that branch.
- The higher-timeframe nested range rotation audit produced only `3` valid
  events across the full BTCUSDT sample. Do not build a baseline from it.
- The impulse absorption audit found continuation-dominant behavior after
  abnormal OHLCV impulse candles. Do not convert that surface into an entry.

## Helper Module Decisions

External helper modules may be used for feature extraction and audit outputs
only. Strategy hypotheses, entries, exits, scoring, sizing, and backtest behavior
stay inside this lab.

Pinned research helper modules:

- `github.com/laclance/go-sr v1.0.0`
- `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
- `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`
