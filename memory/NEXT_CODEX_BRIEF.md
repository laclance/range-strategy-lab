# Next Codex Brief: Independent Entry Premise User Scope Choice Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

There is no selected implementation from the current state. Do not implement an
entry audit, derivatives veto interaction, strategy, P&L run, replay,
walk-forward, source expansion, or generated result directory unless the user
first supplies and explicitly approves one concrete next route.

If the user has not chosen exactly one next route, make no docs edits, no Go
code changes, no CLI flag, no generated result directory, no audit run, no
source download, no network request, no data write, and no strategy/P&L work.
Report the waiting state and stop.

Before any nontrivial work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_INDEPENDENT_ENTRY_PREMISE_AND_HYPOTHESIS_MAP.md for the
  current user-choice stop state.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_INTEGRATION_SPEC.md
  only if exact veto boundaries are needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline audit/backtest flag is passed.
- The latest docs-only review stopped at:
  independent_entry_premise_and_hypothesis_map_needs_user_scope_choice.
- No single BTCUSDT 15m local-source independent entry premise is selected from
  the current reviewed evidence.
- The current reviewed candidates either collapse into closed families, remain
  filter/context evidence only, or require a fresh user-supplied event premise
  or explicit scope change.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed with 0 cohorts and remains
  closed.
- Derivatives source materialization, source audit, context audit, strategy-
  premise spec, no-trade filter premise audit, and no-trade filter integration
  spec all completed in reviewed zero-trade form.
- The canonical derivatives veto candidate is preserved as future skip/retain
  evidence only:
  btc_15m_basis_discount_no_trade_veto_v1.
- The veto premise facts remain: 1,823 de-duplicated veto rows; 1,241 no-trade
  toxic rows; toxic rate 0.680746; minimum split toxic rate 0.665485; weakest
  split rows 387; full toxic improvement versus local-only baseline 0.046269;
  311 rotation-useful and 271 continuation-useful full-sample collateral rows
  were reported.
- There is no independently approved candidate-entry stream for the veto to
  filter, so no veto interaction audit is selected.

The user must choose exactly one next route before work starts:
1. New BTCUSDT 15m local-source premise: requires a closed-candle candidate
   event, intended side/timing, why it is not a closed-family reslice, and a
   falsification rule. First task should be a docs-only premise spec unless the
   user provides a complete spec and explicitly approves a zero-trade audit.
2. Higher-timeframe premise: requires interval set and materially different
   range behavior, not the failed nested rotation premise. First task should be
   a docs-only premise spec or zero-trade audit brief.
3. Spread-range / pair-range: requires explicit source and engine scope
   acceptance. First task should be a docs-only source/engine scope spec.
4. New source family: requires source family, provenance plan, and anti-leakage
   alignment premise. First task should be a docs-only source-scope or
   source-audit brief.
5. No further audit: update docs/memory only if the user explicitly asks for a
   closeout note.

If the user asks to implement the derivatives veto alone:
- Do not implement it.
- Explain that the latest approved entry map stopped at
  independent_entry_premise_and_hypothesis_map_needs_user_scope_choice.
- Explain that the veto can only be tested against a separate independently
  approved entry premise.
- Keep the worktree unchanged unless the user explicitly asks for a docs-only
  brief or supplies a new route choice.

If the user supplies a new BTCUSDT 15m local-source premise:
- Ensure the derivatives veto is not the source of the entry premise.
- The premise must define its candidate-event stream before any veto is applied.
- The premise must define side/timing, decision candle, source facts,
  closed-family separation, and falsification criteria.
- The next likely milestone is docs-only:
  independent_entry_premise_spec_ready_for_user_approval.
- Any later zero-trade premise audit still needs separate explicit approval.
- Any later veto interaction audit still needs another separate explicit
  approval after the independent entry premise exists.

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
- Reopening, retuning, renaming, gate-relaxing, or promoting reviewed closed
  families.

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
