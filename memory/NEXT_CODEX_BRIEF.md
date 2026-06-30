# Next Codex Brief: BTCUSDT 15m Post-Compression Directional Expansion Zero-Trade Audit Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do not implement anything unless the user explicitly approves the zero-trade
audit for:

btc_15m_post_compression_directional_expansion_v1

If the user has not approved that exact audit, make no docs edits, no Go code
changes, no CLI flag, no generated result directory, no audit run, no source
download, no network request, no data write, and no strategy/P&L work. Report
that the project is waiting for explicit zero-trade audit approval and stop.

Before any nontrivial work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_PREMISE_SPEC.md.
- Read docs/FUTURES_INDEPENDENT_ENTRY_PREMISE_AND_HYPOTHESIS_MAP.md only for
  the prior user-choice boundary.
- Read docs/COMPRESSION_BREAKOUT_REVIEW.md and
  docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md only
  if exact closed-family boundaries are needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline audit/backtest flag is passed.
- The latest docs-only spec stopped at:
  independent_entry_premise_spec_ready_for_user_approval.
- The selected independent entry-premise candidate for possible later audit is:
  btc_15m_post_compression_directional_expansion_v1.
- The source is the accepted local BTCUSDT Binance USDT-M futures 5m CSV:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
- Source facts to carry forward: 573,984 loaded candles;
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z; gap_count=0;
  duplicate_count=0; zero_volume_count=66; comparison_only=false;
  validation_status=accepted.
- The derivatives veto candidate remains parked as future skip/retain evidence
  only:
  btc_15m_basis_discount_no_trade_veto_v1.
- The veto may not shape candidate rows, create entries, choose side, rank,
  score P&L, replay, walk forward, optimize, promote a strategy, or reopen
  closed families. Any future veto interaction audit requires a separate
  approval after an independent entry audit exists.

If the user explicitly approves the zero-trade audit implementation:
- Add one offline audit mode only, preferably behind:
  -futures-btc-15m-post-compression-directional-expansion-audit
- Default output directory:
  results/futures-btc-15m-post-compression-directional-expansion-audit
- Use only the accepted BTCUSDT local futures 5m source.
- Resample to closed UTC 15m candles from exact 5m children.
- Keep summary.json, summary.csv, and trades.json zero-trade compatible with
  trades=0.
- Do not simulate fills, fees, slippage, stops, targets, P&L, replay,
  walk-forward, portfolio construction, strategy promotion, or veto
  interaction.

Required event math:
- Decision candle is a closed UTC 15m candle d.
- Compression lookbacks: 48, 96, 192 prior closed 15m bars.
- For lookback L, prior range is [d-L, d-1].
- range_width_pct_L(d) = (range_high_L(d) - range_low_L(d)) / close[d-1].
- Compression thresholds: current range_width_pct_L(d) at or below the bottom
  20%, 30%, or 40% rolling prior percentile.
- Percentile reference set: previous 1,920 valid range_width_pct_L observations,
  ending at d-1. The decision value is excluded.
- ATR: repo's Wilder-style ATR(candles, 14) on 15m candles; use only
  atr14[d-1].
- Upside expansion: close[d] >= range_high_L(d) + M * atr14[d-1].
- Downside expansion: close[d] <= range_low_L(d) - M * atr14[d-1].
- Breakout multiples M: 0.1, 0.2, 0.3.
- Volume modes: none, volume[d] > prior 96-bar median, or volume[d] > prior
  96-bar 60th percentile. Current volume is not included in the threshold.
- Side/timing labels: long after upside expansion at next 15m bar open; short
  after downside expansion at next 15m bar open.
- Candidate viability counts must use de-duplicated (decision_close, side) rows
  so repeated parameter-cell hits cannot inflate evidence.

Required zero-trade labels:
- Horizons: 16, 32, and 48 closed 15m bars after the decision candle.
- Label anchor: open[d+1], diagnostic only.
- Emit intended-side forward close return, intended-side favorable excursion,
  intended-side adverse excursion, favorable_gt_adverse, and
  favorable_minus_adverse.
- Compare every candidate cell against the unconditional eligible 15m baseline
  by split, side, and horizon.
- Primary period splits are 2021_2022_stress, 2023_2024_oos, and
  2025_2026_recent, assigned by decision candle close time.

Required falsification gates:
- Reject if source validation or 15m resampling fails.
- Reject if any feature uses data after decision candle close, or any prior
  threshold includes the decision value.
- Reject if fewer than 300 de-duplicated (decision_close, side) candidates exist
  in the full sample.
- Reject if any primary period split has fewer than 50 de-duplicated
  candidates.
- Reject if no side/horizon/cell cluster separates beyond the unconditional
  eligible 15m baseline.
- Reject if only one isolated parameter cell passes.
- Reject if full-sample separation is not supported in every primary period
  split.
- Reject if the result is a closed-family reslice of compression breakout,
  structured compression, clean breakout, router rotation, occupancy rotation,
  midline/hold-inside, higher-timeframe nested rotation, BTC/ETH/SOL context, or
  derivatives veto work.
- Reject if derivatives veto rows, basis/premium facts, or no-trade filter
  labels shape entries, side, candidate selection, or pass/fail scoring.

Expected outputs:
- source_manifest.json
- summary.json, summary.csv, trades.json with trades=0
- btc_15m_post_compression_directional_expansion_sources.csv
- btc_15m_post_compression_directional_expansion_resample_coverage.csv
- btc_15m_post_compression_directional_expansion_parameter_cells.csv
- btc_15m_post_compression_directional_expansion_candidates.csv
- btc_15m_post_compression_directional_expansion_dedup_events.csv
- btc_15m_post_compression_directional_expansion_baseline.csv
- btc_15m_post_compression_directional_expansion_split_summary.csv
- btc_15m_post_compression_directional_expansion_adjacency.csv
- btc_15m_post_compression_directional_expansion_missingness.csv
- btc_15m_post_compression_directional_expansion_falsification.json

Expected implementation review doc:
- docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md

Later audit stop states:
- btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review
- btc_15m_post_compression_directional_expansion_zero_trade_audit_failed_no_usable_entry_premise
- independent_entry_premise_spec_rejected_closed_family_reslice

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

Closeout if a future approved audit is completed:
- Add/update the implementation review doc.
- Update README.md docs index and intro if the current next gate changes.
- Update memory/PROGRESS.md with the stop state, commands, result paths, source
  facts, row counts, and short factual outcome.
- Update memory/DECISIONS.md only for durable boundaries.
- Refresh this memory/NEXT_CODEX_BRIEF.md to the new canonical next state.
- Run gofmt on changed Go files.
- Run:
  env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
  wc -l results/futures-btc-15m-post-compression-directional-expansion-audit/*.csv
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- After staging, run:
  git diff --cached --check
- Commit completed docs/memory/code changes after checks pass unless told not
  to.
```
