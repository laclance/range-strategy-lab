# Next Codex Brief: Derivatives Context Zero-Trade Source Audit Brief

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md only for the
  parked derivatives-context boundary.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md only
  for the latest closed context-audit exclusion.
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed with 0 passing cohorts and is
  closed at:
  btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context.
- The derivatives source scope review stopped at:
  derivatives_context_source_scope_review_approved_needs_zero_trade_source_audit_brief.
- That review approved only a later docs-only zero-trade source-audit brief, not
  implementation.

Current source facts:
- The durable local source directory checked from this repo is
  ../binance-bot/data/.
- It contains durable candle CSVs only, not approved derivatives market-data
  rows.
- ../binance-bot/data/raw/ had no files in the reviewed search depth.
- Existing approved candle alignment anchors are:
  - ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
- Those candle anchors each have 573,984 loaded candles covering
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z, with sorted gap_count=0
  and duplicate_count=0; zero-volume counts are BTC=66, ETH=47, SOL=47; SOL is
  accepted only after sorting one physical non-monotonic row.
- Adjacent ../binance-bot/research/ source-proof artifacts may be referenced
  only as process/source evidence, not as lab input data or strategy evidence.

Task:
- Create docs-only brief:
  docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md.
- Do not add Go code, CLI flags, generated result directories, source
  downloads, source materialization, source parsing, context features, entries,
  exits, P&L backtests, optimizer grids, replay, walk-forward, paper/testnet,
  live paths, exchange API, credentials, deploy files, martingale, averaging
  down, or two-exchange logic.
- The brief must decide the exact later source-audit question and stop before
  implementation.

Allowed brief-writing scope:
- The later source audit may use the approved BTC/ETH/SOL futures 5m candle
  files only as alignment anchors.
- The later source audit may consider at most one first derivatives source
  family.
- Recommended order:
  1. mark/index/premium basis first;
  2. funding rate history second;
  3. aggregate trades/taker flow only as a high-volume secondary candidate.
- The brief must define durable local/offline source-file requirements before
  implementation. If no durable local files or explicitly approved offline
  materialization path can be named, stop the later brief at source gap.

Forbidden source scope:
- No /tmp caches as durable source inputs.
- No adjacent research result CSVs as strategy evidence.
- No source downloads or live probes.
- No private account endpoints, signed requests, exchange API keys, or
  credentials.
- No live WebSocket stream as a historical substitute.
- No spot or cross-exchange source expansion.
- No broad symbol mining.
- No open-interest metrics from the current evidence set.
- No long/short ratios without full-era public archive proof.
- No liquidation or force-order history from current evidence.
- No order-book/depth source without a separate historical archive proof.
- No aggregate-trade implementation before a narrower funding or basis source
  brief explicitly rejects the lighter family.

Required brief content:
- Verdict and stop state.
- Authority chain from the derivatives source scope review.
- Exact proposed source family.
- Durable local/offline source path requirements.
- Candle alignment anchor contract.
- Timestamp, finality, publication-lag, missing-data, checksum/provenance, and
  max-staleness requirements.
- Allowed and forbidden artifacts for the later implementation.
- Anti-lookahead join rules.
- Rejection criteria.
- Verification commands for the future implementation if explicitly approved.
- Current docs-only closeout commands.

Allowed stop states for this brief-writing task:
- derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval
- derivatives_context_zero_trade_source_audit_brief_rejected_source_gap
- derivatives_context_zero_trade_source_audit_brief_rejected_live_or_private_api_path
- derivatives_context_zero_trade_source_audit_brief_rejected_closed_family_rescue

Recommended stop state if the brief can define a narrow source-audit plan:
derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval.

Verification for this docs-only closeout:
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Add the new brief to README.md docs index after
  docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md and renumber
  following entries.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if a durable boundary or permission rule
  changes.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next approval gate.
- Commit completed docs/memory changes after checks pass unless explicitly told
  not to commit.
```
