# Next Codex Brief: Derivatives Context Zero-Trade Context-Audit Brief-Writing Task

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md (the passed source
  audit and its validated/aligned source facts).
- Read docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md and
  docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md only for the
  source/provenance boundary.
- Read docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md as the
  shape template for a zero-trade context-audit brief.
- Check git status before editing or writing outputs.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline audit/backtest flag is passed.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed (0 cohorts), closed at:
  btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context.
- Derivatives source materialization passed at:
  derivatives_context_source_materialization_passed_ready_for_source_audit_approval
  (729 checksum-verified zips, 9 normalized CSVs, 5 manifests under
  ../binance-bot/data/derivatives/).
- The derivatives zero-trade source audit passed at:
  derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief.
  It validated the 9 mark/index/premium 5m source files, SHA-256-bound their
  provenance, and proved anti-lookahead alignment to the 5m candle anchors under
  a conservative one-interval lag; all 6 required mark/index streams cleared the
  0.99 coverage bar (min 0.994472, index BTCUSDT) with recorded missingness and
  no forward fill. Audit flag: -futures-derivatives-context-source-audit;
  artifacts under results/futures-derivatives-context-source-audit/.

Task (docs-only, no implementation):
- Write a single new docs-only brief
  docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md that turns
  the passed source audit into a decision-complete plan for a LATER zero-trade
  derivatives context audit. The brief decides the exact later question, the
  derived context inputs, the anti-leakage rules, the artifacts, the gates, and
  the stop states. It does not implement anything.
- The later audit's question must be a zero-trade separation question only, e.g.
  whether mark-minus-index basis / premium-index level / basis-change context
  buckets known at closed decision-candle time improve separation of BTC/ETH/SOL
  usable, toxic, rotation, continuation, or no-trade local range states. It must
  not ask whether basis can be traded and must not measure entry/exit/P&L.
- Allowed source inputs for the later audit are only the 9 validated derivatives
  CSVs under ../binance-bot/data/derivatives/ plus the 3 candle anchors. Carry
  forward the proven facts: SHA-256 provenance, conservative one-5m-interval
  publication lag, exact closed-interval joins, no forward fill / interpolation /
  nearest-future joins, and bounded recorded missingness (required min lag
  coverage 0.994472).
- Required anti-leakage rules to specify: forward labels may appear only in
  label/cohort/ranking/summary artifacts, never in premise/state/context/gating/
  feature-bucket inputs; basis/premium context features must be built only from
  source rows whose source_close_time + 5m <= decision_candle_close_time; missing
  context must produce missingness/skip rows, not silent defaults.
- Required common-output rule: summary.json, summary.csv, and trades.json must
  stay zero-trade compatible.
- Preserve all closed-family verdicts; the brief must explain why basis/premium
  context is materially different from the closed price-only and BTC-regime
  families and must not be a rescue of any closed family.

Forbidden in this task:
- No Go code, CLI flags, generated result directories, audit runs, source
  downloads, network requests, or data writes under
  ../binance-bot/data/derivatives/.
- No context features, labels, cohorts, rankings, entries, exits, P&L backtests,
  optimizer grids, replay, walk-forward, portfolio construction, paper/testnet/
  live paths, exchange API, credentials, deploy files, martingale, averaging
  down, two-exchange logic, or closed-family rescue.
- Do not implement the context audit; this is brief-writing only. The brief must
  stop before code and require explicit user approval before any implementation.

Stop state for this task:
- derivatives_context_zero_trade_context_audit_brief_ready_for_user_approval.

The later implementation (only after explicit user approval) should reuse the
passed source audit's loader/alignment facts and stay zero-trade.

Closeout:
- Update README.md docs index and intro.
- Update memory/PROGRESS.md with the new brief and stop state.
- Update memory/DECISIONS.md only if the brief creates a durable boundary or
  rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the context-audit implementation
  approval gate.
- Run: rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md;
  git diff --check; git status --short.
- Commit completed docs/memory changes after checks pass unless told not to.
```
