# Next Codex Brief: New Futures Hypothesis Or Data Premise After Impulse Failure

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT 5m range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
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
- The futures impulse absorption audit stop state was:
  impulse_absorption_no_viable_edge.
- The tested impulse premise was:
  - 30-day prior rolling warmup
  - true-range percentile rank >= p99
  - volume percentile rank >= p95
  - close position <=0.25 for down impulses or >=0.75 for up impulses
  - horizons 3, 6, 12, and 24 closed bars after the event
  - midpoint reclaim before continuation beyond the event extreme as the desired absorption behavior
- Result:
  - 3,434 distinct impulse events
  - 13,736 event/horizon candidate rows
  - all period splits had at least 1,003 all-direction/all-bucket events at every horizon
  - continuation-first beat midpoint reclaim-first in every period split and horizon
  - quick continuation reached roughly 78%-81%
  - no missing future windows
  - strategy remained empty with 0 trades
- This futures-authoritative result blocks converting the abnormal OHLCV impulse absorption surface into an entry prototype.
- The prior closed families remain an exclusion map:
  - SR rejection/confirmation/false-break timing: closed legacy spot-only evidence.
  - compression breakout: closed legacy spot-only evidence.
  - range regime durability: diagnostic/infrastructure, not an entry surface.
  - detector durability/context refinement: reusable infrastructure, not entry promotion.
  - hold-inside directional edge: closed.
  - hold-inside midline transition/reaction: diagnostic only after the failed prototype.
  - futures midline touch prototype: futures-authoritative failure.
  - futures impulse absorption: futures-authoritative no viable edge.
- Default runs still use lab.EmptyStrategy unless an explicit prototype/audit flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Do not start a new audit or prototype automatically.
- Ask the user for a materially new futures hypothesis or data premise.
- Once the premise is explicit, write a scoped next brief for a non-trading audit only.

Implementation boundaries:
- Keep this task planning/review-only unless the user supplies a clear new premise in the same session.
- Do not add entries, exits, scoring, sizing changes, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, or two-exchange logic.
- Do not retune or broaden the failed hold-inside/midline touch prototype.
- Do not reopen SR timing reslices, compression reslices, detector/context promotion, hold-inside directional continuation, hold-inside/midline reaction, abnormal OHLCV impulse absorption, or old spot evidence without a materially new futures premise.
- Do not use spot outputs as authority for futures promotion.

Required next interaction:
- Ask for the new hypothesis or data premise in concrete terms:
  - What market-structure behavior should be audited?
  - Why should it be different from the closed families?
  - Which closed-candle observable defines the candidate event?
  - What would falsify it before any entry prototype?
- If the user provides a premise, create a new memory/NEXT_CODEX_BRIEF.md for a non-trading audit only. Include exact source, outputs, review gate, stop states, and verification commands.
- If the user does not provide a premise, stop with:
  hypothesis_premise_needed_after_impulse_failure.

Suggested verification after any brief refresh:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Closeout:
- Update memory/PROGRESS.md with the chosen premise or the lack of one, exact commands, and stop state.
- Update memory/DECISIONS.md only if a durable no-retry or new-premise rule is added.
- Commit completed repo changes after verification unless explicitly told not to.
```
