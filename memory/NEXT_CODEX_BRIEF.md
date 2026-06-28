# Next Codex Brief: Derivatives Context Zero-Trade Source Audit Implementation Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md as the
  decision-complete source-audit plan.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md for the
  durable source files and their recorded gap facts.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md and
  docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md only for source
  boundaries.
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
- The derivatives source materialization was explicitly approved and executed.
  It passed at:
  derivatives_context_source_materialization_passed_ready_for_source_audit_approval.
- Durable Binance public Data Vision USDT-M futures source files now exist under
  ../binance-bot/data/derivatives/: 729 checksum-verified raw zips, 9 normalized
  CSVs, and 5 manifests. The 9 normalized streams are mark/index/premium 5m for
  BTCUSDT/ETHUSDT/SOLUSDT, schema
  open_time,open,high,low,close,close_time,source_object_id, covering
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z. Each stream has 0
  duplicate-conflict, 0 parse-error, 0 out-of-range, 0 non-monotonic rows, and
  recorded internal gaps (required mark/index missing 9,820 of 3,443,904,
  0.285%). Manifests record per-stream gap_count and missing_interval_count under
  a no-imputation policy.
- Passing materialization created durable candidate source inputs only. It does
  NOT by itself authorize the zero-trade source audit, context features, labels,
  cohorts, rankings, entries, exits, P&L, replay, walk-forward, or promotion.

Approval gate:
- If the current user request does not explicitly approve implementing the
  zero-trade derivatives context source audit, make no repo mutations, no new
  Go code or CLI flags, no generated result directories, no source-audit parsing
  runs, and no network requests.
- In that no-approval case, report that the materialized derivatives source
  files are ready at ../binance-bot/data/derivatives/ and the source-audit plan
  is ready at docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md,
  then stop at:
  derivatives_context_source_audit_waiting_for_implementation_approval.
- If the current user request explicitly approves implementing the source audit,
  implement only the zero-trade source audit described in
  docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md.

Allowed source-audit scope only after explicit implementation approval:
- Question: whether the durable mark/index/premium 5m source rows can be
  validated, provenance-bound, and anti-lookahead-aligned to the existing
  Binance USDT-M futures 5m candle anchors well enough to justify a separate
  later zero-trade context-audit brief. It must NOT ask whether basis can be
  traded.
- Inputs: only the materialized files under ../binance-bot/data/derivatives/
  (plus their manifests) as candidate sources, and the three candle anchors
  ../binance-bot/data/{btcusdt,ethusdt,solusdt}_futures_um_5m_2021_2026.csv as
  alignment anchors only (573,984 candles each; gap_count=0; duplicate_count=0;
  zero-volume BTC 66, ETH 47, SOL 47; SOL one physical non-monotonic row accepted
  only after sorting).
- New CLI flag: -futures-derivatives-context-source-audit, default output
  results/futures-derivatives-context-source-audit/.
- Allowed artifacts: sources, candle_anchors, external_coverage,
  timestamp_alignment, publication_lag, missingness, provenance, skips, and
  summary CSV/JSON, plus zero-trade-compatible summary.json/summary.csv and an
  empty trades.json. Mark-minus-index basis / premium-index level rows may be
  derived only for source validation and alignment reporting.
- Required source facts per candidate file: durable path, source family, source
  owner/archive family, symbol, product scope, native interval, timestamp
  semantics, first/last timestamp, row count, duplicate count, gap/missing count,
  non-monotonic physical count, timezone assumption, publication/finality rule,
  max staleness, missing-data policy, checksum or provenance identifier,
  comparison-only status, and validation status.
- Anti-lookahead: align by symbol/product/interval/UTC; use only source rows with
  source_close_time + publication_lag <= decision_candle_close_time; require
  exact closed 5m alignment or report the row missing; preferred max staleness is
  one closed 5m interval; no future fill, no interpolation, no nearest-future
  joins, no labels/future returns/future revisions as validation inputs; report
  skipped/missing/stale/duplicate/rejected rows as first-class artifacts. The
  recorded materialization gaps must surface as bounded missingness here, not be
  silently filled.

Forbidden scope:
- No funding, aggregate trades, taker flow, open interest, long/short ratios,
  liquidations, force-order rows, order book/depth, spot, or cross-exchange
  sources.
- No source downloads, REST tail endpoints, WebSockets, signed/private
  endpoints, exchange API keys, credentials, live probes, /tmp caches, or
  adjacent research-result CSVs as inputs.
- No context-gain features, cohorts, rankings, labels, forward returns,
  adverse/favorable excursion rows, trade-outcome rows, entries, exits, P&L
  backtests, optimizer grids, replay, walk-forward, portfolio construction,
  paper/testnet/live paths, deploy files, martingale, averaging down,
  two-exchange logic, or closed-family rescue (structured compression, BTC regime
  plus ETH/SOL, router-gated boundary reclaim, breakout-retest/acceptance, clean
  breakout, hold-inside/midline, impulse absorption, nested range rotation).

Allowed source-audit stop states (only after explicit approval):
- derivatives_context_zero_trade_source_audit_rejected_source_gap
- derivatives_context_zero_trade_source_audit_rejected_live_or_private_api_path
- derivatives_context_zero_trade_source_audit_rejected_timestamp_or_finality_gap
- derivatives_context_zero_trade_source_audit_rejected_alignment_gap
- derivatives_context_zero_trade_source_audit_rejected_closed_family_rescue
- derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief

A passing source-audit stop state still would not authorize context-gain
implementation, entries, exits, P&L backtests, optimizer grids, replay,
walk-forward, packaging, source downloads, paper/testnet/live paths, exchange
API work, credentials, deploy files, strategy promotion, martingale, averaging
down, or two-exchange logic.

Verification if implementation is explicitly approved:
- rg --files ../binance-bot/data | rg -i "(mark[_-]?price|index[_-]?price|premium[_-]?index|premiumIndex|basis)"
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-source-audit -out-dir results/futures-derivatives-context-source-audit
- wc -l results/futures-derivatives-context-source-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout after approved implementation:
- Create or update a focused source-audit review doc with source paths, product,
  symbol, source family, interval, timestamp semantics, finality rule, coverage,
  row counts, duplicate/gap/missing/stale counts, checksums or provenance
  identifiers, generated artifact paths, zero-trade common-output status, stop
  state, and exact command outcomes.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if implementation creates a durable boundary,
  no-go rule, or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next approval-gated task (a separate
  zero-trade derivatives context-audit brief if the source audit passes, or the
  next parked direction if it rejects).
- Commit completed repo changes after checks pass unless explicitly told not to
  commit. Do not commit generated data outside this repository.
```
