# Next Codex Brief: BTCUSDT 15m Post-Compression Strategy-Premise Spec Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do not implement anything unless the user explicitly approves a docs-only
strategy-premise spec for the passed zero-trade audit:

btc_15m_post_compression_directional_expansion_v1

If the user has not approved that exact docs-only strategy-premise spec, make no
docs edits, no Go code changes, no CLI flag, no generated result directory, no
audit/backtest run, no source download, no network request, no data write, and
no strategy/P&L work. Report that the project is waiting for explicit
strategy-premise spec approval and stop.

Before any nontrivial work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_PREMISE_SPEC.md.
- Read docs/FUTURES_INDEPENDENT_ENTRY_PREMISE_AND_HYPOTHESIS_MAP.md only if the
  independent-entry route boundary is needed.
- Read docs/COMPRESSION_BREAKOUT_REVIEW.md and
  docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md only
  if exact closed-family boundaries are needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline research flag is passed.
- The selected independent entry-premise family is:
  btc_15m_post_compression_directional_expansion_v1.
- The approved zero-trade audit passed at:
  btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review.
- Source: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  Binance USDT-M futures BTCUSDT 5m; 573,984 loaded candles;
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z; gap_count=0;
  duplicate_count=0; zero_volume_count=66; comparison_only=false;
  validation_status=accepted.
- Exact closed UTC 15m resample facts: 191,328 rows; first open
  2021-01-01T00:00:00Z; last open 2026-06-16T23:45:00Z; last close
  2026-06-16T23:59:59Z; 3 expected child bars; 0 missing child opens;
  validation accepted.
- Audit artifacts live under:
  results/futures-btc-15m-post-compression-directional-expansion-audit/
  They are Git-ignored but should be referenced by path.
- Audit counts: parameter_cells=81, candidate_rows=386,694,
  dedup_events=4,677, baseline_rows=24, split_summary_rows=1,944,
  adjacency_rows=486, missingness_rows=4, passing_cells=9,
  adjacent_pass_clusters=9, trades=0.
- Falsification gates passed: source/resample, leakage, candidate size,
  split size, baseline separation, adjacent-cell cluster, split stability,
  closed-family protection, derivatives-veto contamination, and zero-trade
  common outputs.
- Passing evidence is narrow:
  long side only; horizon 48 closed 15m bars only; compression lookback 192;
  compression threshold bottom 20%; breakout thresholds 0.1, 0.2, and 0.3
  prior-bar ATR(14); volume modes none, above_prior_96_median, and
  above_prior_96_p60. No short-side, 16-bar, 32-bar, lookback 48/96, or
  compression threshold 30%/40% surface passed the full gate.
- Common outputs stayed zero-trade: summary.json, summary.csv, and trades.json
  report 0 trades.
- The derivatives veto candidate remains parked as future skip/retain evidence
  only:
  btc_15m_basis_discount_no_trade_veto_v1.
- The veto may not shape candidate rows, create entries, choose side, rank,
  score P&L, replay, walk forward, optimize, promote a strategy, or reopen
  closed families. Any future veto interaction audit requires a separate
  approval after an independent backtest candidate exists.

If the user explicitly approves the docs-only strategy-premise spec:
- Add one docs-only artifact:
  docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_STRATEGY_PREMISE_SPEC.md
- Do not add Go code, CLI flags, generated outputs, audit/backtest runs, source
  downloads, P&L, replay, walk-forward, optimizer selection, paper/testnet/live
  path, or derivatives veto interaction.
- The spec must decide whether the passed zero-trade evidence is sufficient to
  request a later offline backtest spec.
- If it selects a later backtest candidate, it must define exactly one bounded
  candidate from the passing pocket only. It should not reopen the full 81-cell
  grid as an optimizer.

Required spec content if approved:
- Restate that the evidence is zero-trade label separation only, not P&L.
- Identify the only eligible evidence pocket:
  long, h48, lookback 192, q20 compression, breakout ATR multiple 0.1/0.2/0.3,
  volume mode none/median/p60.
- Decide whether a later backtest should use:
  a single conservative representative cell, a small predeclared robustness
  cluster, or no backtest because the evidence is too narrow.
- If a later backtest candidate is selected, predeclare candidate construction,
  side/timing, allowed parameters, fixed risk/exits question, source facts,
  no-lookahead rules, exact stop/fail gates, and output expectations for a
  future separate approval.
- Explicitly reject any strategy premise that depends on short-side evidence,
  16/32-bar labels, non-passing lookbacks, non-passing compression thresholds,
  derivatives veto facts, old detector episodes, router rows, occupancy
  rotation, midline/hold-inside states, BTC/ETH/SOL context, or source expansion.
- Preserve the derivatives veto as parked future skip/retain evidence only.

Allowed docs-only stop states:
- post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval
- post_compression_directional_expansion_strategy_premise_spec_failed_too_narrow
- post_compression_directional_expansion_strategy_premise_spec_needs_user_choice

Forbidden unless a later explicit brief changes scope:
- Source downloads, network requests, source materialization, data writes under
  ../binance-bot/data/derivatives/, entries, exits, P&L backtests, optimizer
  grids, replay, walk-forward, portfolio construction, paper/testnet/live paths,
  exchange API, credentials, deploy files, martingale, averaging down,
  two-exchange logic, strategy promotion, or derivatives veto interaction.
- Treating btc_15m_basis_discount_no_trade_veto_v1 as an entry signal, a
  basis-tradability claim, a basis-fade rule, a rotation-entry rule, a
  continuation-entry rule, or evidence that P&L would improve.
- Reopening, retuning, renaming, gate-relaxing, or promoting reviewed closed
  families.

Closeout if the docs-only spec is completed:
- Update README.md intro and docs index.
- Update memory/PROGRESS.md with the stop state, commands, paths, and short
  factual outcome.
- Update memory/DECISIONS.md only for durable boundaries.
- Refresh this memory/NEXT_CODEX_BRIEF.md to the new canonical next state:
  backtest-spec approval gate, user-choice gate, or no selected next
  implementation.
- Run:
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- After staging, run:
  git diff --cached --check
- Commit completed docs/memory changes after checks pass unless told not to.
```
