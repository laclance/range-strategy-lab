# Next Codex Brief: BTCUSDT 15m Post-Compression Offline Backtest Spec Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do not implement anything unless the user explicitly approves a docs-only
offline backtest spec for the selected post-compression representative
candidate:

btc_15m_post_compression_l192_q20_m020_none_long_h48_v1

If the user has not approved that exact docs-only backtest spec, make no docs
edits, no Go code changes, no CLI flag, no generated result directory, no
audit/backtest run, no source download, no network request, no data write, no
P&L artifact, and no strategy/veto work. Report that the project is waiting for
explicit backtest-spec approval and stop.

Before any nontrivial work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_STRATEGY_PREMISE_SPEC.md.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_PREMISE_SPEC.md only if exact event-construction details are needed.
- Read docs/FUTURES_INDEPENDENT_ENTRY_PREMISE_AND_HYPOTHESIS_MAP.md,
  docs/COMPRESSION_BREAKOUT_REVIEW.md, or
  docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md only
  if exact independent-entry or closed-family boundaries are needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline research flag is passed.
- The selected independent entry-premise family is:
  btc_15m_post_compression_directional_expansion_v1.
- Its approved zero-trade audit passed at:
  btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review.
- The docs-only strategy-premise spec then stopped at:
  post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval.
- The only selected later backtest-spec candidate is:
  btc_15m_post_compression_l192_q20_m020_none_long_h48_v1.
- Source: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  Binance USDT-M futures BTCUSDT 5m; 573,984 loaded candles;
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z; gap_count=0;
  duplicate_count=0; zero_volume_count=66; comparison_only=false;
  validation_status=accepted.
- Exact closed UTC 15m resample facts: 191,328 rows; first open
  2021-01-01T00:00:00Z; last open 2026-06-16T23:45:00Z; last close
  2026-06-16T23:59:59Z; 3 expected child bars; 0 missing child opens;
  validation accepted.
- Zero-trade audit artifacts live under:
  results/futures-btc-15m-post-compression-directional-expansion-audit/
  They are Git-ignored but may be referenced by path.
- Audit counts: parameter_cells=81, candidate_rows=386,694,
  dedup_events=4,677, baseline_rows=24, split_summary_rows=1,944,
  adjacency_rows=486, missingness_rows=4, passing_cells=9,
  adjacent_pass_clusters=9, trades=0.
- Passing evidence was narrow:
  long side only; horizon 48 closed 15m bars only; compression lookback 192;
  compression threshold bottom 20%; breakout thresholds 0.1, 0.2, and 0.3
  prior-bar ATR(14); volume modes none, above_prior_96_median, and
  above_prior_96_p60. No short-side, 16-bar, 32-bar, lookback 48/96, or
  compression threshold 30%/40% surface passed the full gate.
- The selected representative candidate is the center/no-extra-filter cell from
  that pocket: long only, L=192, q20 compression, M=0.2 prior-bar ATR(14),
  volume=none, evidence horizon h48.
- Common zero-trade outputs from the audit stayed zero-trade:
  summary.json, summary.csv, and trades.json report 0 trades.
- The derivatives veto candidate remains parked as future skip/retain evidence
  only:
  btc_15m_basis_discount_no_trade_veto_v1.
- The veto may not shape candidate rows, create entries, choose side, rank,
  score P&L, replay, walk forward, optimize, promote a strategy, or reopen
  closed families. Any future veto interaction audit requires a separate
  approval after an independent backtest candidate exists.

If the user explicitly approves the docs-only offline backtest spec:
- Add one docs-only artifact:
  docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_SPEC.md
- Do not add Go code, CLI flags, generated outputs, audit/backtest runs, source
  downloads, P&L artifacts, replay, walk-forward, optimizer selection,
  paper/testnet/live path, or derivatives veto interaction.
- The spec must decide whether the representative candidate is sufficiently
  defined for a later offline backtest implementation approval gate.
- The spec must lock one fixed risk/exit model before implementation or stop at
  a user-choice gate. It may not request a stop grid, target grid, holding
  period grid, parameter-cell grid, volume-filter grid, side optimizer, or veto
  interaction.

Required spec content if approved:
- Restate that the current evidence is zero-trade label separation only, not
  P&L.
- Define exactly one candidate stream:
  btc_15m_post_compression_l192_q20_m020_none_long_h48_v1.
- Candidate construction:
  closed UTC BTCUSDT 15m candles resampled from exact local 5m children; prior
  range [d-192,d-1]; q20 compression against prior 1,920 valid range-width
  observations; decision close above prior range high by 0.2 * ATR(14)[d-1];
  volume mode none; long side only; next 15m open entry timing for any later
  approved implementation.
- Predeclare one fixed risk/exit model for the future implementation brief, or
  stop at:
  post_compression_directional_expansion_backtest_spec_needs_user_exit_choice.
- Define no-lookahead rules, source/resample validation, skip/missingness
  handling, expected generated artifacts, P&L/split metrics, pass/fail gates,
  and exact later implementation stop states.
- Explicitly reject any backtest premise that depends on short-side evidence,
  16/32-bar labels, non-passing lookbacks, non-passing compression thresholds,
  volume-filter selection, neighboring-cell P&L rescue, derivatives veto facts,
  old detector episodes, router rows, occupancy rotation, midline/hold-inside
  states, BTC/ETH/SOL context, source expansion, or the full 81-cell grid.
- Preserve the derivatives veto as parked future skip/retain evidence only.

Allowed docs-only stop states:
- post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval
- post_compression_directional_expansion_backtest_spec_needs_user_exit_choice
- post_compression_directional_expansion_backtest_spec_rejected_too_narrow
- post_compression_directional_expansion_backtest_spec_rejected_optimizer_contamination
- post_compression_directional_expansion_backtest_spec_rejected_closed_family_reslice

Forbidden unless a later explicit brief changes scope:
- Source downloads, network requests, source materialization, data writes under
  ../binance-bot/data/derivatives/, Go code, CLI flags, entries, exits, P&L
  backtest runs, optimizer grids, replay, walk-forward, portfolio construction,
  paper/testnet/live paths, exchange API, credentials, deploy files, martingale,
  averaging down, two-exchange logic, strategy promotion, or derivatives veto
  interaction.
- Treating btc_15m_basis_discount_no_trade_veto_v1 as an entry signal, a
  basis-tradability claim, a basis-fade rule, a rotation-entry rule, a
  continuation-entry rule, or evidence that P&L would improve.
- Reopening, retuning, renaming, gate-relaxing, or promoting reviewed closed
  families.

Closeout if the docs-only backtest spec is completed:
- Update README.md intro and docs index.
- Update memory/PROGRESS.md with the stop state, commands, paths, and short
  factual outcome.
- Update memory/DECISIONS.md only for durable boundaries.
- Refresh this memory/NEXT_CODEX_BRIEF.md to the new canonical next state:
  backtest implementation approval gate, user-choice gate, or no selected next
  implementation.
- Run:
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- After staging, run:
  git diff --cached --check
- Commit completed docs/memory changes after checks pass unless told not to.
```
