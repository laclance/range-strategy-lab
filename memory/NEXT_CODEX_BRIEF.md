# Next Codex Brief: Post BTC Regime ETH/SOL Context Audit Pivot Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md only
  for the approved audit boundary that has now been executed.
- Read docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md only
  for the BTCUSDT price-only stop boundary.
- Read docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md only if
  the user asks for a next-lane recommendation.
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- The post-rotation premise failure pivot stopped at:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- The reviewed premise router_gated_boundary_reclaim_rotation_v1 is closed in
  reviewed form.
- Do not convert the 278 premise context segments, 97 boundary-reclaim events,
  or 1,299 tradable_rotation router rows into trades.
- The BTC regime plus ETH/SOL zero-trade context audit was explicitly approved,
  implemented, and reviewed at:
  docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md.
- That audit stopped at:
  btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context.
- It produced 0 passing context cohorts. Source validation, coverage,
  anti-leakage checks, and common zero-trade outputs passed, but BTC regime
  context did not add durable separation beyond ETH/SOL local-only baselines.

Hard boundary:
- Do not retune, reslice, rename, gate-relax, replay, walk-forward, or promote
  the BTC regime plus ETH/SOL context audit.
- Do not treat ETHUSDT/SOLUSDT context rows as strategy authority.
- Do not promote BTCUSDT regime rows.
- Do not add entries, exits, P&L strategy backtests, optimizer grids, replay,
  walk-forward logic, strategy packages, paper/testnet/live paths, exchange API,
  credentials, deploy files, source downloads, broad mining, martingale,
  averaging down, or two-exchange logic.
- Do not rescue the structured-compression branch or reuse its ETH/SOL
  authority result as promotion evidence.
- Do not import old binance-bot strategy/scoring/live code.

Next-task gate:
- If the user only asks "what next" or asks for a recommendation, answer
  read-only from the current review/memory: the BTC regime plus ETH/SOL audit
  failed and the next step is choosing a materially different scope.
- If the user asks for a docs-only scope review, create or update only a focused
  scope-review doc and memory. Do not add Go code or generated result
  directories.
- If the user explicitly approves a materially different implementation, verify
  the new scope against memory/DECISIONS.md first. Implementation is allowed
  only if the request is not a retune or direct extension of a closed failed
  family.

Recommended next scope choices, still requiring explicit user approval:
1. Derivatives market-data context source review: docs-only first, because it
   needs explicit source/alignment approval and must remain market-data context
   only.
2. Spread-range/pair-range source and engine review: docs-only first, because
   it needs separate multi-leg source and engine scope before any P&L work.
3. A new independent entry premise only if it is materially different from
   range_occupancy_rotation_v1, router-gated boundary reclaim, BTC regime plus
   ETH/SOL context, structured compression, breakout-retest/acceptance, clean
   breakout continuation, hold-inside/midline, impulse absorption, and
   higher-timeframe nested range rotation.

Still blocked:
- Volatility-aware exits remain unavailable until a future independent entry
  premise first shows gross edge before costs.
- BTCUSDT-only price-range audits remain stopped by the post-rotation premise
  failure pivot unless the user approves a materially different non-price-only
  premise.
- BTC regime plus ETH/SOL context retunes remain stopped by the zero-trade audit
  review.

Expected stop state for a read-only or docs-only next-scope closeout:
btc_regime_eth_sol_context_zero_trade_audit_failed_waiting_for_new_scope_choice.

Verification for a docs-only closeout:
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if a durable boundary or permission rule
  changes.
- Refresh memory/NEXT_CODEX_BRIEF.md if the next approved scope changes.
- Commit completed docs/memory/code updates after checks pass unless explicitly
  told not to commit.
```
