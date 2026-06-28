# Next Codex Brief: Derivatives Context Source Audit Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md only for the
  source-scope boundary.
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
- The derivatives source-audit brief stopped at:
  derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval.
- That brief is ready for explicit user approval, but it does not authorize
  source-audit implementation by itself.

Approval gate:
- If the current user request does not explicitly approve implementing the
  derivatives context zero-trade source audit, do not add code, CLI flags,
  generated result directories, source materialization, source parsing, source
  downloads, audit outputs, or context artifacts.
- In that no-approval case, report that the approved brief is ready at
  docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md and stop at:
  derivatives_context_zero_trade_source_audit_waiting_for_user_approval.
- If the current user request explicitly approves implementation of that
  zero-trade source audit, implement only the source/alignment audit described
  in docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md.

Allowed source-audit scope only after explicit approval:
- First source family only:
  Binance USDT-M futures mark-price, index-price, or premium-index klines.
- Symbols:
  BTCUSDT, ETHUSDT, and SOLUSDT only.
- Candle anchors only:
  - ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
- Candle anchor facts:
  each file has 573,984 loaded candles; coverage is
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z; sorted streams have
  gap_count=0 and duplicate_count=0; zero_volume_count is BTC=66, ETH=47,
  SOL=47; SOL had one physical non-monotonic row and was accepted only after
  sorting.
- Derivatives source files must be durable local/offline files under
  ../binance-bot/data/ or a documented subdirectory, or the run must stop at
  source gap.
- Equivalent source path names are acceptable only when a manifest maps the
  actual files to product, symbol, interval, source family, coverage, timestamp
  semantics, finality rule, and provenance.

Implementation question only after explicit approval:
Can durable local/offline Binance USDT-M mark-price, index-price, or
premium-index kline rows for BTCUSDT, ETHUSDT, and SOLUSDT be validated,
provenance-bound, and anti-lookahead-aligned to the approved 5m candle anchors
well enough to justify a separate later zero-trade context-audit brief?

Required source-audit rules:
- Record durable path, source family, source owner/archive family, symbol,
  product, native interval, timestamp semantics, first/last timestamp, row
  count, duplicate count, gap or missing-interval count, non-monotonic physical
  row count, timezone assumption, publication/finality rule, max staleness,
  missing-data policy, checksum or provenance identifier, comparison-only
  status, and validation status.
- Use UTC closed-candle semantics.
- Align source rows only when source_close_time + publication_lag is less than
  or equal to the decision candle close.
- If publication lag cannot be proven, reject the source or report only a
  conservative one-native-interval-lag alignment view.
- Missing context must produce missingness or skip rows, not silent defaults.
- Common outputs must remain zero-trade compatible.

Allowed generated directory only after explicit approval:
- results/futures-derivatives-context-source-audit/

Allowed source-audit artifacts only after explicit approval:
- futures_derivatives_context_source_audit_sources.csv/json
- futures_derivatives_context_source_audit_candle_anchors.csv/json
- futures_derivatives_context_source_audit_external_coverage.csv/json
- futures_derivatives_context_source_audit_timestamp_alignment.csv/json
- futures_derivatives_context_source_audit_publication_lag.csv/json
- futures_derivatives_context_source_audit_missingness.csv/json
- futures_derivatives_context_source_audit_provenance.csv/json
- futures_derivatives_context_source_audit_skips.csv/json
- futures_derivatives_context_source_audit_summary.csv/json
- common summary.json, summary.csv, and trades.json with 0 trades

Forbidden scope:
- No funding, aggregate trades, taker flow, open interest, long/short ratios,
  liquidations, force-order rows, order book/depth, spot sources,
  cross-exchange sources, broad symbol mining, or adjacent research result CSVs
  as strategy evidence.
- No /tmp caches as durable inputs.
- No source downloads, live probes, private endpoints, signed requests,
  exchange API keys, credentials, WebSockets as historical substitutes, paper,
  testnet, live execution, deploy files, martingale, averaging down, or
  two-exchange logic.
- No context-gain features, labels, cohorts, rankings, entries, exits, P&L
  backtests, optimizer grids, replay, walk-forward, portfolio construction,
  strategy promotion, or closed-family rescue.

Allowed implementation stop states, only after explicit approval:
- derivatives_context_zero_trade_source_audit_rejected_source_gap
- derivatives_context_zero_trade_source_audit_rejected_live_or_private_api_path
- derivatives_context_zero_trade_source_audit_rejected_timestamp_or_finality_gap
- derivatives_context_zero_trade_source_audit_rejected_alignment_gap
- derivatives_context_zero_trade_source_audit_rejected_closed_family_rescue
- derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief

Verification if implementation is explicitly approved:
- rg --files ../binance-bot/data | rg -i "(mark[_-]?price|index[_-]?price|premium[_-]?index|premiumIndex|basis)"
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-source-audit -out-dir results/futures-derivatives-context-source-audit
- wc -l results/futures-derivatives-context-source-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout after any approved implementation:
- Create a focused review doc with source facts, artifact paths, command
  outcomes, common zero-trade output status, and stop state.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the implementation creates a durable
  boundary, no-go rule, or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next approval-gated task.
- Commit completed docs/memory/code updates after checks pass unless explicitly
  told not to commit.
```
