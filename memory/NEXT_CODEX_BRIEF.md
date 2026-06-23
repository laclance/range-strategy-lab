# Next Codex Brief: Futures Scope Pivot Review, Range-Only

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_SCOPE_PIVOT_REVIEW_SPEC.md
  - docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md
  - docs/FUTURES_DATA_IMPACT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active project objective is still range-strategy research, not generic futures strategy discovery.
- The active market target is Binance USDT-M futures. Spot evidence is historical context only unless explicitly rerun and reviewed on futures data.
- Current authoritative lab source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - market type: Binance USDT-M futures BTCUSDT 5m
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - accepted manifest facts: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- BTCUSDT 5m range work is not close to a strategy:
  - futures midline touch prototype failed full sample, every period split, and both sides;
  - futures impulse absorption failed because continuation-first and quick continuation dominated midpoint reclaim-first;
  - SR timing, compression, range durability, detector/context, hold-inside directional edge, and hold-inside midline/reaction are closed, diagnostic-only, or infrastructure.
- The scope pivot must remain range-only unless the user explicitly changes the objective.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/prototype flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Create docs/FUTURES_SCOPE_PIVOT_REVIEW.md as a review-only decision document.
- Decide whether the next range-only scope should be:
  - stop BTCUSDT 5m range work until a materially new premise appears;
  - specify a higher-timeframe BTCUSDT futures range source/premise;
  - specify a narrow BTC/ETH futures range-only source/premise;
  - declare no viable range-only lane from current evidence;
  - ask the user to explicitly change scope if they want non-range futures strategy work.

Sibling-repo boundaries:
- Inspect only range-relevant sibling context from ~/binance-bot and ~/crypto-trading-bot.
- Use sibling repos as process/source-contract/exclusion context only.
- Do not import strategy results, parameters, Python implementation, Go runtime code, live/deploy behavior, or multi-pair selectors as evidence for this lab.
- Treat known sibling evidence carefully:
  - ~/binance-bot Binance-only BTCUSDT 5m range sleeve failed even gross on legacy spot data; use as caution only.
  - ~/binance-bot cross-exchange spread seed failed local data gates and remains out of scope because this lab forbids two-exchange execution.
  - ~/binance-bot multi-pair selector/alt overlay work is a no-go for broad symbol selection.
  - ~/binance-bot daily ATR contraction evidence is non-range trend/volatility work and is outside this range-only review.
  - ~/crypto-trading-bot BTC/ETH 4h range-reversal/re-entry workbench has failed as a frozen family; use it as exclusion evidence, not a port target.
  - ~/crypto-trading-bot BTC/ETH USD-M 1h-to-4h/1d source-contract discipline may be useful if a higher-timeframe or BTC/ETH range source spec is selected.

Implementation boundaries:
- This task is docs/memory-only.
- Do not add code, CLI flags, audits, generated results, entries, exits, scoring, sizing, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, or two-exchange logic.
- Do not reopen closed BTCUSDT 5m hold-inside/midline, SR timing, compression, or impulse surfaces by reslicing.
- Do not run or mutate sibling repos.
- Do not use spot outputs as authority for futures promotion.

Required work:
- Create docs/FUTURES_SCOPE_PIVOT_REVIEW.md using docs/FUTURES_SCOPE_PIVOT_REVIEW_SPEC.md as the controlling spec.
- Add the new review doc to the README docs index.
- Update memory/PROGRESS.md with docs created, exact commands, sibling repo facts used, and stop state.
- Update memory/DECISIONS.md only if the review adds a durable range-scope or no-retry rule.
- Refresh memory/NEXT_CODEX_BRIEF.md from the chosen stop state.

Review requirements:
- Include current BTCUSDT 5m futures source facts.
- Include a closed-family table for lab 5m range work.
- Include a sibling-repo evidence table with process-only and exclusion-only labels.
- Include a scope-lane table covering:
  - BTCUSDT 5m range continuation;
  - BTCUSDT higher-timeframe range;
  - BTC/ETH range-only source expansion;
  - larger multi-pair range universe;
  - cross-exchange range spread;
  - non-range higher-timeframe trend/volatility strategy.
- Recommend exactly one next route or stop with no viable range-only lane.

Stop states:
- range_scope_pivot_ready_for_higher_timeframe_source_spec
- range_scope_pivot_ready_for_btc_eth_range_source_spec
- range_scope_pivot_no_viable_range_lane
- range_scope_pivot_needs_user_scope_change
- range_scope_pivot_review_only_no_strategy_change

Suggested verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit completed repo changes after verification unless explicitly told not to.
```
