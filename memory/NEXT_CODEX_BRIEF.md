# Next Codex Brief: Derivatives Context Source Materialization Execution Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md only
  for the later source-audit boundary.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md only for the
  source-scope boundary.
- Check git status before editing or writing outputs.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed with 0 passing cohorts and is
  closed at:
  btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context.
- The derivatives source-audit brief stopped at:
  derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval.
- The derivatives source materialization plan stopped at:
  derivatives_context_source_materialization_plan_ready_for_execution_approval.
- That materialization plan records user approval for the plan, but it does not
  by itself execute downloads, write durable files, run the source audit, or
  authorize context/strategy work.

Approval gate:
- If the current user request does not explicitly approve executing the
  derivatives context source materialization, make no repo mutations, no network
  requests, no file writes under ../binance-bot/data/derivatives/, no generated
  result directories, and no audit/source parsing runs.
- In that no-execution-approval case, report that the materialization plan is
  ready at docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md and
  stop at:
  derivatives_context_source_materialization_waiting_for_execution_approval.
- If the current user request explicitly approves executing source
  materialization, perform only the public-archive materialization described in
  docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md.

Allowed materialization scope only after explicit execution approval:
- Public archive owner:
  Binance Data Vision.
- Product:
  Binance USDT-M futures.
- Required archive families:
  markPriceKlines and indexPriceKlines.
- Optional archive family:
  premiumIndexKlines as cross-check only.
- Symbols:
  BTCUSDT, ETHUSDT, and SOLUSDT only.
- Interval:
  5m only.
- Era:
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z only.
- Object set:
  monthly objects from 2021-01 through 2026-05 plus daily tail objects from
  2026-06-01 through 2026-06-16.
- Target root:
  ../binance-bot/data/derivatives/

Known planning facts:
- Adjacent local inventory for required mark/index scope recorded 486 official
  archive objects and 90,067,194 compressed bytes.
- markPriceKlines contributed 243 objects and 42,833,309 bytes.
- indexPriceKlines contributed 243 objects and 47,233,885 bytes.
- Optional premiumIndexKlines is estimated, not proven in this lab, at another
  same-shaped 243 objects and roughly 45 MB compressed.
- No durable derivatives files were present under ../binance-bot/data/derivatives/
  when the plan was written.

Required execution behavior:
- Use deterministic public Data Vision archive object URLs only:
  https://data.binance.vision/data/futures/um/{monthly|daily}/{archive_family}/{symbol}/5m/{symbol}-5m-{key}.zip
- Write raw zips, normalized CSVs, and manifests only under
  ../binance-bot/data/derivatives/.
- Record object URL, durable path, source owner, product, archive family, source
  family, symbol, interval, monthly/daily scope, key, captured timestamp, bytes,
  SHA-256, ETag if returned, Last-Modified if returned, Content-Length if
  returned, rows_total, rows_used, rows_out_of_range, parse_errors,
  duplicate_same, duplicate_conflict, first_open, last_open, and validation
  status.
- Normalize each family/symbol CSV with:
  open_time,open,high,low,close,close_time,source_object_id
- Sort rows by open_time, count physical non-monotonic rows before sorting,
  collapse duplicate identical rows, reject duplicate conflicting rows, and fail
  closed on gaps, schema ambiguity, checksum mismatch, or required source gaps.
- Missing optional premium-index objects may produce an optional-family skip
  state, but missing required mark/index objects must reject materialization.
- Do not run the derivatives source audit unless the user separately approves
  that later after materialization.

Forbidden scope:
- No source families outside markPriceKlines, indexPriceKlines, and optional
  premiumIndexKlines.
- No REST tail endpoints, WebSockets, signed/private endpoints, exchange API
  keys, credentials, live probes, archive broad discovery, funding, aggregate
  trades, taker flow, open interest, long/short ratios, liquidations,
  force-order rows, order book/depth, spot sources, cross-exchange sources,
  broad symbol mining, paper/testnet/live paths, deploy files, martingale,
  averaging down, or two-exchange logic.
- No context-gain features, labels, cohorts, rankings, entries, exits, P&L
  backtests, optimizer grids, replay, walk-forward, portfolio construction,
  strategy promotion, or closed-family rescue.

Allowed materialization stop states:
- derivatives_context_source_materialization_rejected_public_archive_gap
- derivatives_context_source_materialization_rejected_checksum_or_schema_gap
- derivatives_context_source_materialization_rejected_unapproved_source_path
- derivatives_context_source_materialization_rejected_unapproved_live_or_private_path
- derivatives_context_source_materialization_passed_ready_for_source_audit_approval
- derivatives_context_source_materialization_passed_with_optional_premium_gap_ready_for_source_audit_approval

Verification if execution is explicitly approved:
- find ../binance-bot/data/derivatives -type f | sort
- wc -l ../binance-bot/data/derivatives/*.csv
- wc -l ../binance-bot/data/derivatives/manifests/*.csv
- sha256sum ../binance-bot/data/derivatives/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout after approved execution:
- Create or update a focused materialization review doc with source facts,
  durable paths, manifest paths, measured object counts/bytes, command outcomes,
  and stop state.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if execution creates a durable boundary,
  no-go rule, or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next approval-gated task, likely a
  separate zero-trade derivatives context source-audit implementation approval
  gate if materialization passes.
- Commit completed docs/memory changes after checks pass unless explicitly told
  not to commit. Do not commit generated data outside this repository unless
  the user explicitly requests a commit in that other repository.
```
