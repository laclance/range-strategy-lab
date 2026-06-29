# Next Codex Brief: Independent Entry Premise Required Before Derivatives Veto Integration

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

There is no selected derivatives filter integration implementation from the
current state. Do not implement the derivatives no-trade veto by itself. If the
user has not supplied and explicitly approved a separate independent entry
premise for interaction testing, make no docs edits, no Go code changes, no CLI
flag, no generated result directory, no audit run, no source download, no
network request, no data write, and no strategy/P&L work. Report the waiting
state and stop.

Before any nontrivial work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_INTEGRATION_SPEC.md for
  the current deferred stop state.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_PREMISE_AUDIT_REVIEW.md
  only if exact veto evidence is needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline audit/backtest flag is passed.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed with 0 cohorts and remains
  closed.
- Derivatives source materialization, source audit, context audit, strategy-
  premise spec, and no-trade filter premise audit all completed in reviewed
  zero-trade form.
- The zero-trade derivatives no-trade filter premise audit passed at:
  derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec.
- The docs-only derivatives no-trade filter integration spec stopped at:
  derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise.
- The canonical veto candidate is preserved as future filter evidence only:
  btc_15m_basis_discount_no_trade_veto_v1.
- The veto premise facts remain: 1,823 de-duplicated veto rows; 1,241 no-trade
  toxic rows; toxic rate 0.680746; minimum split toxic rate 0.665485; weakest
  split rows 387; full toxic improvement versus local-only baseline 0.046269;
  311 rotation-useful and 271 continuation-useful full-sample collateral rows
  were reported.
- There is no independently approved candidate-entry stream for the veto to
  filter, so no integration audit is selected.

If the user asks to implement the veto alone:
- Do not implement it.
- Explain that the latest approved stop state is
  derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise.
- Explain that a no-trade veto can only be tested against a separate,
  independently approved entry premise.
- Keep the worktree unchanged unless the user explicitly asks for a new
  docs-only brief or a new independent entry-premise task.

If the user supplies a materially new independent entry premise:
- First decide whether the requested task is a docs-only premise spec, a
  zero-trade premise audit, or a later veto interaction audit.
- The derivatives veto must not be the source of the entry premise.
- The entry premise must define its own candidate-event stream before the veto
  is applied.
- The veto may only annotate already-approved candidate rows as skipped or
  retained. It may not create entries, choose side, change entry logic, act as
  an exit, rank trades, score P&L, optimize, replay, walk forward, or promote a
  strategy.
- Any later interaction audit still needs separate explicit user approval.

Forbidden unless a later explicit brief changes scope:
- Go code, CLI flags, generated result directories, audit runs, source
  downloads, network requests, source materialization, data writes under
  ../binance-bot/data/derivatives/, entries, exits, P&L backtests, optimizer
  grids, replay, walk-forward, portfolio construction, paper/testnet/live paths,
  exchange API, credentials, deploy files, martingale, averaging down,
  two-exchange logic, strategy promotion, or closed-family rescue.
- Treating btc_15m_basis_discount_no_trade_veto_v1 as an entry signal, a
  basis-tradability claim, a basis-fade rule, a rotation-entry rule, a
  continuation-entry rule, or evidence that P&L would improve.
- Expanding this veto to ETHUSDT, SOLUSDT, other timeframes, new source
  families, BTC-regime context, router rotation, spread-range, volatility-aware
  exits, or closed BTCUSDT price-only families.

Closeout if a future approved milestone is completed:
- Update README.md docs index and intro if the current next gate changes.
- Update memory/PROGRESS.md with the stop state, commands, result paths if any,
  and short factual outcome.
- Update memory/DECISIONS.md only for durable boundaries.
- Refresh this memory/NEXT_CODEX_BRIEF.md to the new canonical next state.
- Run:
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- Also run Go tests and artifact row-count checks only when code or generated
  outputs changed.
- Commit completed docs/memory/code changes after checks pass unless told not
  to.
```
