# Next Codex Brief: BTC Regime Plus ETH/SOL Zero-Trade Audit Brief

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md.
- Read docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md only
  for the BTCUSDT price-only stop boundary.
- Read docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md only for BTC/ETH/SOL
  source-validation facts if needed.
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
- The BTC regime plus ETH/SOL scope review stopped at:
  btc_regime_eth_sol_context_scope_review_approved_needs_zero_trade_audit_brief.
- That review approved only a separate zero-trade audit brief-writing task, not
  audit implementation.

Allowed source scope for the future brief:
- Use only the already local Binance USDT-M futures 5m files:
  - ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
- Source facts from prior validation:
  - each symbol has 573,984 loaded candles;
  - coverage is 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - sorted streams have gap_count=0 and duplicate_count=0;
  - zero_volume_count is BTC=66, ETH=47, SOL=47;
  - SOL had one physical non-monotonic row and was accepted only after sorting.
- Do not add symbols, source downloads, spot comparisons, derivatives context
  sources, exchange APIs, private endpoints, or broad mining.

Goal:
Create a documentation-only zero-trade audit brief or implementation plan for
BTC regime plus ETH/SOL context. Do not implement the audit. The brief must be
specific enough for a later session to implement only after explicit approval.

Expected doc:
- Create docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md.
- Add it to README.md's docs index if created.
- Keep the brief decision-complete but stop before code.

The brief must define:
- BTCUSDT role: market-regime context and diagnostic-only authority row.
- ETHUSDT/SOLUSDT role: possible authority rows only for a zero-trade context
  audit, not strategy promotion.
- Minimum audit question:
  Do BTC regime buckets known at closed decision-candle time improve separation
  of ETH/SOL usable, toxic, rotation, continuation, or no-trade range states?
- Allowed features at a high level:
  closed-candle BTC regime buckets plus closed-candle ETH/SOL local range state.
- Required anti-leakage rule:
  forward labels may appear only in label/cohort artifacts, never premise inputs.
- Required common-output rule:
  common summary/trades outputs must remain zero-trade compatible.
- Rejection criteria:
  closed-family reslice, broad mining, source gap, hidden future-label input,
  structured-compression rescue, ETH/SOL replay, BTC promotion, or any move
  toward entries/backtests.

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

Allowed outcome:
- Stop at:
  btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval.
- Refresh memory/NEXT_CODEX_BRIEF.md to wait for explicit user approval before
  any implementation.

Verification for this documentation-only closeout:
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the brief creates a durable boundary,
  no-go rule, or permission rule.
- Commit completed docs/memory updates after checks pass unless explicitly told
  not to commit.
```
