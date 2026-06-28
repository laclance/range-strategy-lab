# Next Codex Brief: BTC Regime Plus ETH/SOL Context Scope Review

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md.
- Read docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md only for
  parked-direction ordering context.
- Read docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md only if exact
  failed-premise evidence is needed.
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- The accepted BTCUSDT futures 5m source remains:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
- Source facts:
  - loaded candles: 573,984;
  - coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - comparison_only=false, validation_status=accepted.
- The futures range state construction loop audit stopped at:
  range_state_construction_loop_audit_passed_needs_router_spec.
- The futures range context router audit stopped at:
  range_context_router_passed_needs_rotation_premise_spec.
- The router rotation premise audit stopped at:
  range_router_rotation_premise_audit_failed_no_premise.
- The post-rotation premise failure pivot review stopped at:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- The reviewed premise router_gated_boundary_reclaim_rotation_v1 is closed in
  reviewed form.
- Do not convert the 278 premise context segments, 97 boundary-reclaim events,
  or 1,299 tradable_rotation router rows into trades.
- No materially different BTCUSDT-only, candle-price-only range-premise audit is
  selected from current evidence.

User-approved lane:
- Review BTC regime plus ETH/SOL context first.
- Treat derivatives market-data context as the parked second candidate.
- Treat spread-range/pair-range as the parked third candidate.
- Treat volatility-aware exits as unavailable unless a future independent entry
  premise first shows gross edge before costs.

Goal:
Create a documentation-only approval review for BTC regime plus ETH/SOL context.
Do not implement an audit. The review should decide whether
docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md is materially different enough
and tightly scoped enough to become the next zero-trade audit brief.

Expected doc:
- Create docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md.
- Add it to README.md's docs index if created.
- Keep the review concise and decision-oriented.

The review must cover:
- why BTC-as-regime-context plus ETH/SOL local range context is materially
  different from another BTCUSDT-only price-range slice;
- why it is not a structured-compression rescue, ETH/SOL replay, BTC promotion,
  broad symbol-mining path, or portfolio construction path;
- the exact allowed source scope: already local Binance USDT-M futures 5m files
  for BTCUSDT, ETHUSDT, and SOLUSDT only;
- BTCUSDT role: market-regime context and diagnostic-only authority row;
- ETHUSDT/SOLUSDT role: possible authority rows only in a later zero-trade
  context audit, not strategy promotion;
- the minimum next audit question: whether BTC regime buckets improve separation
  of ETH/SOL usable, toxic, rotation, continuation, or no-trade range states;
- rejection criteria for the review: closed-family reslice, broad mining, source
  gap, hidden future-label input, or any move toward entries/backtests.

Boundaries:
- Do not add Go code, CLI flags, generated result directories, source downloads,
  strategy code, entries, exits, P&L strategy backtests, optimizer grids, replay,
  walk-forward logic, strategy packages, paper/testnet/live paths, exchange API,
  credentials, deploy files, broad mining, martingale, averaging down, or
  two-exchange logic.
- Do not retune, rename, relax gates for, or directly repackage the failed
  router-gated boundary-reclaim rotation premise.
- Do not rescue the structured-compression branch or reuse its ETH/SOL authority
  result as promotion evidence.
- Do not import old binance-bot strategy/scoring/live code.
- Do not use future labels as premise inputs.

Allowed outcomes:
- If approved, stop at:
  btc_regime_eth_sol_context_scope_review_approved_needs_zero_trade_audit_brief
  and refresh memory/NEXT_CODEX_BRIEF.md to a separate brief-writing task for a
  zero-trade audit spec or implementation plan that still requires explicit
  approval before code.
- If rejected, stop at:
  btc_regime_eth_sol_context_scope_review_rejected_closed_family_reslice
  or:
  btc_regime_eth_sol_context_scope_review_rejected_scope_gap
  and refresh memory/NEXT_CODEX_BRIEF.md to either no next audit or another user
  scope-choice prompt.

Verification for this documentation-only closeout:
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the review creates a durable boundary,
  no-go rule, or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next bounded step.
- Commit completed docs/memory updates after checks pass unless explicitly told
  not to commit.
```
