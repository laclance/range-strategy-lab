# Next Codex Brief: Futures Hypothesis Pivot Inventory

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT 5m range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md
  - docs/FUTURES_DATA_IMPACT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active research target is Binance USDT-M futures BTCUSDT 5m, not spot.
- Full-history futures source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - market type: Binance USDT-M futures
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - source manifest: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- The futures midline touch prototype stop state was:
  prototype_failed_no_promotion.
- Failed exact template:
  - detector profile p30_c12_bollinger_on_adx_off
  - context rule hold_3_inside
  - first closed-candle mid_touch within 12 bars after hold decision
  - event close-position bucket mid_50
  - close-back side model
  - same-side boundary stop, opposite-boundary target, 6 bar time stop
  - next-bar-open entry
- Full prototype result:
  - 532 signal rows, 531 trades, one exact-mid skip
  - full gross P&L -95.54, full net P&L -418.99
  - full PF 0.3409, average net R -0.4276
  - all period splits and both sides failed
- Default runs still use lab.EmptyStrategy unless an explicit prototype/audit flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Stop mining the failed hold-inside/midline entry family and create a compact futures hypothesis pivot inventory before any new audit or entry work.

Implementation boundaries:
- Keep this task review/inventory only unless a clearly different non-trading audit is identified for a later brief.
- Do not add entries, exits, scoring, sizing changes, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, or two-exchange logic.
- Do not retune the failed midline-touch prototype.
- Do not broaden the failed surface to hold_6_inside, mid_close_across, side-specific cohorts, or old spot evidence.
- Do not use spot outputs as authority for futures promotion.

Required work:
- Create docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md.
- Inventory the major reviewed/closed families from the docs index, including:
  - SR rejection/confirmation/false-break timing
  - compression breakout
  - range regime durability
  - detector durability/context refinement
  - hold-inside directional edge
  - hold-inside midline transition/reaction
  - futures midline touch prototype
- For each family, mark:
  - current status: closed, diagnostic-only, or possibly reusable as infrastructure
  - what not to retry without a new premise
  - whether the conclusion is futures-authoritative or legacy spot-only
- If a materially different futures hypothesis is found, write the next brief for a non-trading audit only. If not, make the next brief ask the user for a new hypothesis or data premise.

Suggested verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Stop states:
- pivot_inventory_ready_for_new_non_trading_audit_brief
- pivot_inventory_no_viable_next_hypothesis
- pivot_inventory_needs_user_hypothesis
- pivot_inventory_review_only_no_strategy_change

Closeout:
- Update memory/PROGRESS.md with docs created, exact commands, and stop state.
- Update memory/DECISIONS.md only if a durable no-retry or next-hypothesis rule is added.
- Refresh memory/NEXT_CODEX_BRIEF.md with the exact next step based on the stop state.
- Commit completed repo changes after verification unless explicitly told not to.
```
