# Next Codex Brief: BTC Regime Plus ETH/SOL Context Audit Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md only for the
  approval boundary.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md only for parked context.
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
- The BTC regime plus ETH/SOL zero-trade audit brief stopped at:
  btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval.
- That brief is ready for user approval, but it does not authorize audit
  implementation by itself.

Approval gate:
- If the current user request does not explicitly approve implementing the
  BTC regime plus ETH/SOL zero-trade context audit, do not add code, CLI flags,
  generated result directories, or audit outputs.
- In that no-approval case, report that the approved brief is ready at
  docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md and stop at:
  btc_regime_eth_sol_context_zero_trade_audit_waiting_for_user_approval.
- If the current user request explicitly approves implementation of that
  zero-trade audit, implement only the audit described in
  docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md.

Allowed source scope if implementation is explicitly approved:
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

Implementation scope only after explicit approval:
- BTCUSDT role: market-regime context and diagnostic-only authority row.
- ETHUSDT/SOLUSDT role: possible authority rows only for a zero-trade context
  audit, not strategy promotion.
- Minimum audit question:
  Do BTC regime buckets known at closed decision-candle time improve separation
  of ETH/SOL usable, toxic, rotation, continuation, or no-trade range states?
- Allowed features at a high level:
  closed-candle BTC regime buckets plus closed-candle ETH/SOL local range state.
- Required anti-leakage rule:
  forward labels may appear only in label/cohort/ranking/summary artifacts,
  never premise, state-ID, router, gating, or feature-bucket inputs.
- Required common-output rule:
  common summary/trades outputs must remain zero-trade compatible.
- Rejection criteria:
  closed-family reslice, broad mining, source gap, hidden future-label input,
  structured-compression rescue, ETH/SOL replay, BTC promotion, or any move
  toward entries/backtests.

Boundaries:
- Do not add strategy code, entries, exits, P&L strategy backtests, optimizer
  grids, replay, walk-forward logic, strategy packages, paper/testnet/live
  paths, exchange API, credentials, deploy files, broad mining, martingale,
  averaging down, or two-exchange logic.
- Do not retune, rename, relax gates for, or directly repackage the failed
  router-gated boundary-reclaim rotation premise.
- Do not rescue the structured-compression branch or reuse its ETH/SOL
  authority result as promotion evidence.
- Do not import old binance-bot strategy/scoring/live code.
- Do not use future labels as premise inputs.

Allowed implementation stop states, only after explicit approval:
- btc_regime_eth_sol_context_zero_trade_audit_source_gap
- btc_regime_eth_sol_context_zero_trade_audit_rejected_closed_family_reslice
- btc_regime_eth_sol_context_zero_trade_audit_rejected_future_label_leak
- btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context
- btc_regime_eth_sol_context_zero_trade_audit_passed_needs_strategy_premise_spec

Verification if implementation is explicitly approved:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-regime-eth-sol-context-audit -out-dir results/futures-btc-regime-eth-sol-context-audit
- wc -l results/futures-btc-regime-eth-sol-context-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout after any approved implementation:
- Update or create a focused review doc with source facts, artifact paths,
  command outcomes, common zero-trade output status, and stop state.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the implementation creates a durable
  boundary, no-go rule, or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next approval-gated task.
- Commit completed docs/memory/code updates after checks pass unless explicitly
  told not to commit.
```
