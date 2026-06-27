# Next Codex Brief: Futures Range Context Triage Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_SPEC.md
  - docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md
  - docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_PROTOCOL.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/FUTURES_SCOPE_PIVOT_REVIEW.md
  - docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project remains offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- Research is not stopped, but automatic reuse or retuning of failed premises
  is stopped.
- The latest executable grammar, range_occupancy_rotation_v1, failed its
  bounded optimizer review.
- V1 stop state:
  range_first_strategy_v1_optimizer_failed_no_replay.
- Source and closed UTC resample validation passed for V1:
  - source:
    ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv;
  - loaded candles: 573,984;
  - open-time coverage:
    2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - 15m resample: 191,328 rows through 2026-06-16T23:45:00Z;
  - 1h resample: 47,832 rows through 2026-06-16T23:00:00Z.
- The fixed V1 baseline lost after costs and the declared 1,152-row grid had
  0 passing configs. A fixed replay spec, walk-forward, package review,
  retune, gate relaxation, symbol expansion, or live-adjacent path is not
  authorized from V1.
- The new user-approved premise is a non-trading range-context triage audit
  that evaluates range quality, UTC session behavior, and failure-mode taxonomy
  in parallel before any new strategy grammar.
- Spec stop state:
  range_context_triage_spec_ready_for_audit_implementation.

Goal:
- Implement the non-trading BTCUSDT futures range-context triage audit behind:
  -futures-range-context-triage-audit
- Follow docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_SPEC.md exactly.
- This is an audit only: no entries, exits, strategy P&L backtest, optimizer,
  replay, walk-forward, strategy package, source expansion, symbol expansion,
  or live-adjacent work.

Default audit config:
- approved source:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
- expected source facts:
  - rows: 573,984;
  - first open: 2021-01-01T00:00:00Z;
  - last open: 2026-06-16T23:55:00Z;
  - gaps: 0;
  - duplicates: 0;
  - zero-volume candles: 66;
  - product: Binance USDT-M futures;
  - symbol: BTCUSDT;
  - interval: 5m;
  - comparison_only=false;
  - validation_status=accepted.
- selected closed UTC resamples:
  - 15m: expected 191,328 rows through 2026-06-16T23:45:00Z;
  - 1h: expected 47,832 rows through 2026-06-16T23:00:00Z;
  - 4h: expected 11,958 rows through 2026-06-16T20:00:00Z.
- detector:
  - profile p30_c12_bollinger_on_adx_off;
  - percentile 0.30;
  - min consecutive bars 12;
  - lookback days 20;
  - Bollinger on;
  - ADX off.
- horizons and thresholds:
  - outcome horizons: 12, 24, 48 closed bars per timeframe;
  - quick failure horizon: 6 closed bars;
  - reentry window: 6 closed bars;
  - clean expansion threshold: 0.75 frozen range widths;
  - drift threshold: final close outside range with less than 0.50
    range-width extension;
  - boundary chop transition threshold: 3 inside/outside state transitions;
  - low width threshold: width_pct < 0.0015 or width_to_atr_ratio < 0.75;
  - wide volatile threshold: width_pct >= 0.0150 or
    width_to_atr_ratio >= 4.00;
  - choppy pre-mature threshold: at least 3 midpoint crosses and at least 4
    boundary touches;
  - minimum full cohort count: 300;
  - minimum split cohort count: 50;
  - minimum session split cohort count: 30;
  - minimum usable context rate: 0.55 full and 0.45 in every period split;
  - maximum toxic context rate: 0.45 full and 0.55 in every period split;
  - maximum missing future rate: 0.02 full.

Audit behavior:
- Add internal/lab/futures_range_context_triage_audit.go with exported
  config/result/row types, default config, runner, source/resample validation,
  range candidate builder, quality lens, session lens, failure-mode labeling,
  cohort summary, cohort ranking, and stop-state helper.
- Build range candidates from each raw-active detector run that reaches a
  mature Active candle.
- Freeze candidate geometry at the first mature close using only bars from
  raw-active start through mature close.
- Record raw-active end metadata as descriptive context only; do not use it
  for future eligibility or cohort ranking.
- Label future windows starting at the next closed bar after mature close.
- Emit exactly one primary context label per candidate/horizon:
  contained_rotation, clean_expansion_up, clean_expansion_down,
  false_break_reentry_up, false_break_reentry_down, boundary_chop,
  drift_through_up, drift_through_down, low_width_noise, no_resolution,
  or missing_future.
- Apply deterministic label precedence from the spec.
- Summarize cohorts by:
  - split, timeframe, horizon;
  - split, timeframe, horizon, quality bucket;
  - split, timeframe, horizon, mature session;
  - split, timeframe, horizon, primary context label;
  - split, timeframe, horizon, quality bucket, mature session.
- Rank only non-trading context cohorts. A ranked cohort is not an executable
  trade config.

Stop states:
- range_context_triage_source_gap
- range_context_triage_no_range_episodes
- range_context_triage_no_usable_cohorts
- range_context_triage_failed_no_strategy_premise
- range_context_triage_ready_for_strategy_spec
- range_context_triage_rejected_closed_family_reslice

Use range_context_triage_ready_for_strategy_spec only if source/resample
validation passes and at least one non-closed-family context cohort passes all
review gates. That stop state authorizes only a later documentation-only
strategy spec, not a baseline backtest or optimizer.

CLI and artifacts:
- Wire cmd/rangelab/main.go with a config hook for tests:
  futuresRangeContextTriageAuditConfigForRun.
- Add the flag:
  -futures-range-context-triage-audit
- Reject spot/comparison sources.
- Reject combinations with trade-producing prototype/baseline/optimization/
  replay/walk-forward flags.
- Keep common outputs zero-trade compatible:
  source_manifest.json, summary.csv/json, trades.json.
- Write triage artifacts under:
  results/futures-range-context-triage-audit/
- Write all of these as CSV and JSON:
  - futures_range_context_triage_sources
  - futures_range_context_triage_coverage
  - futures_range_context_triage_episodes
  - futures_range_context_triage_quality
  - futures_range_context_triage_sessions
  - futures_range_context_triage_failure_modes
  - futures_range_context_triage_cohorts
  - futures_range_context_triage_rankings
  - futures_range_context_triage_summary
- Console summary should include source row count, coverage row count, episode
  rows, quality rows, session rows, failure-mode rows, cohort rows, ranking
  rows, passing cohort count, and stop state.

Closed-family boundaries:
- Do not retune range_occupancy_rotation_v1.
- Do not retune structured compression.
- Do not retune breakout_retest_acceptance.
- Do not reopen clean breakout, hold-inside/midline, impulse absorption,
  higher-timeframe nested range rotation, boundary touch rejection,
  single-candle wick rejection, failed breakout re-entry, or mature balance
  persistence as entries.
- Do not import old binance-bot strategy/scoring/live code.
- Do not add live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, data downloads, broad symbol mining, martingale, averaging down,
  or two-exchange logic.

Tests:
- Add focused lab tests for:
  - source acceptance and source-gap rejection;
  - closed UTC 15m, 1h, and 4h coverage validation;
  - mature range freeze using only raw-start through mature close;
  - raw-active end metadata not affecting eligibility;
  - invalid candidate skip reasons;
  - UTC session classification boundaries;
  - quality bucket assignment;
  - contained rotation label;
  - clean expansion up/down labels;
  - false-break reentry labels;
  - boundary chop transition label;
  - drift-through labels;
  - missing future label;
  - deterministic label precedence;
  - cohort summary gates;
  - ranking score and tie-breaks;
  - stop-state selection.
- Add CLI tests proving:
  - default runs do not write range-context triage artifacts;
  - the new flag writes all required artifacts;
  - spot comparison is rejected;
  - combinations with trade-producing flags are rejected.

Review and memory after the run:
- Add docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md with source/resample
  facts, artifact line counts, episode counts, quality/session/failure-mode
  findings, cohort ranking, caveats, and stop state.
- Add the review doc to README.md.
- Update memory/PROGRESS.md with exact commands, result path, source/resample
  facts, CSV line counts, cohort outcomes, and stop state.
- Update memory/DECISIONS.md only if the audit creates a durable promotion,
  no-promotion, or no-strategy-change rule.
- Replace memory/NEXT_CODEX_BRIEF.md based on the final stop state.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-context-triage-audit -out-dir results/futures-range-context-triage-audit
- wc -l results/futures-range-context-triage-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit the completed implementation, generated review doc, README/memory
  updates, refreshed next brief, and verification evidence after checks pass
  unless explicitly told not to commit.
```
