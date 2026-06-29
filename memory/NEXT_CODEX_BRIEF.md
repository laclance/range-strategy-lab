# Next Codex Brief: Derivatives Context Zero-Trade Context-Audit Implementation Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do NOT implement anything in this task without explicit user approval. The
context-audit brief is written and gated; the next step is the user approving or
rejecting the context-audit implementation scope. If approval is not explicitly
given in the session, make no Go code changes, no CLI flag, no result directory,
no audit run, no source download, no network request, and no data write; report
the waiting state and stop.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md as the
  decision-complete plan to implement.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md for the passed
  source/alignment facts and the loader/alignment code to reuse.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md only for
  the source/provenance boundary.
- Check git status before editing or writing outputs.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline audit/backtest flag is passed.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed (0 cohorts), closed at:
  btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context.
- Derivatives source materialization passed; 729 checksum-verified zips, 9
  normalized CSVs, 5 manifests under ../binance-bot/data/derivatives/.
- The derivatives zero-trade source audit passed at:
  derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief.
  It validated the 9 mark/index/premium 5m source files, SHA-256-bound their
  provenance, and proved anti-lookahead alignment to the 5m candle anchors under
  a conservative one-interval lag; all 6 required mark/index streams cleared the
  0.99 coverage bar (required floor 0.994472, index BTCUSDT). Flag:
  -futures-derivatives-context-source-audit.
- The derivatives zero-trade context-audit brief is written and stopped at:
  derivatives_context_zero_trade_context_audit_brief_ready_for_user_approval.
  This task is its implementation approval gate.

Task (only after explicit user approval):
- Implement the zero-trade derivatives context audit exactly as scoped in
  docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md. Answer only
  the separation question: whether mark-minus-index basis, premium-index level,
  or basis-change context buckets known at closed decision-candle time improve
  separation of BTCUSDT/ETHUSDT/SOLUSDT local range states (usable, toxic,
  rotation, continuation, no-trade) beyond the local price/volume state alone.
- Use only the 9 validated derivatives CSVs under ../binance-bot/data/derivatives/
  plus the 3 candle anchors. Reuse the passed source audit's loader and its
  conservative one-5m-interval anti-lookahead alignment
  (source_close_time + 5m <= decision_candle_close_time). No forward fill,
  interpolation, or nearest-future joins; bounded recorded missingness only with
  the 0.994472 required basis context coverage floor.
- Build derivatives context (basis level, basis sign/change, optional premium
  level) only from lagged closed source rows. Local range states come from the
  existing range-state construction infrastructure over the candle anchors.
- Enforce the anti-leakage rules: forward labels only in label/cohort/ranking/
  summary artifacts; missing context produces missingness/skip rows, never silent
  defaults; summary.json, summary.csv, and trades.json stay zero-trade
  compatible.
- Enforce the orthogonality gate: a basis/premium bucket must add separation
  beyond the local price/volume state; collinear/dominated buckets fail. Preserve
  every closed-family verdict; do not reopen, retune, rename, gate-relax, or
  promote any closed family, and do not rescue a closed family with derivatives
  context.
- Wire a single CLI flag (suggested -futures-derivatives-context-audit) with
  default out-dir results/futures-derivatives-context-audit, matching the
  source-product and flag-conflict guards of the other zero-trade audits. Write
  the brief's nine artifact families (CSV+JSON) plus the zero-trade common
  outputs. Add unit tests (pass path plus rejection paths) and a CLI flag test.

Forbidden in this task:
- Any implementation before explicit user approval.
- Basis tradability tests, entries, exits, P&L backtests, optimizer grids,
  replay, walk-forward, portfolio construction, paper/testnet/live paths,
  exchange API, credentials, deploy files, martingale, averaging down,
  two-exchange logic, broad symbol/interval/source mining, source downloads,
  network requests, data writes under ../binance-bot/data/derivatives/, or
  closed-family rescue.

Allowed implementation stop states (only after explicit user approval):
- derivatives_context_zero_trade_context_audit_source_gap;
- derivatives_context_zero_trade_context_audit_rejected_future_label_leak;
- derivatives_context_zero_trade_context_audit_rejected_closed_family_rescue;
- derivatives_context_zero_trade_context_audit_failed_no_usable_context;
- derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec.

A passing implementation stop state would still not authorize entries, exits,
P&L backtests, replay, walk-forward, packaging, source downloads,
paper/testnet/live paths, exchange API work, credentials, deploy files, strategy
promotion, or deployment.

Closeout (only after an approved implementation):
- Add a review doc docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md.
- Update README.md docs index and intro.
- Update memory/PROGRESS.md with the implementation result and stop state.
- Update memory/DECISIONS.md only if the result creates a durable boundary.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next gate.
- Run:
  env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
  env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-audit -out-dir results/futures-derivatives-context-audit
  wc -l results/futures-derivatives-context-audit/*.csv
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- Commit completed code, docs, and memory changes after checks pass unless told
  not to. Generated results/ stay Git-ignored; derivatives data stays outside
  the repo.
```
