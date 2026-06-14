# Next Codex Brief: Review Range Regime Durability Stability

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/ENTRY_READINESS_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md, docs/COMPRESSION_BREAKOUT_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Compression breakout audit was not entry-ready.
- Range regime durability audit has been implemented but not reviewed for stability yet.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest range regime durability audit:
- CLI flag:
  - -range-regime-durability-audit
- Outputs:
  - results/range-regime-durability-audit/range_regime_durability_episodes.csv/json
  - results/range-regime-durability-audit/range_regime_durability_summary.csv/json
- Audit size:
  - episode_rows=11984
  - summary_rows=452
  - episode CSV lines including header: 11,985
  - summary CSV lines including header: 453
- Defaults:
  - detector_profile_id=p30_c12_bollinger_on_adx_off
  - horizons=1;3;6;12
  - quick_invalidation_bars=3
- Episode semantics:
  - episodes are contiguous RawActive detector runs that eventually become Active
  - episode high/low, width, length, and ATR context use only closed candles through the episode end
  - label windows start at episode_end_index + 1
  - all forward durability metrics are label_* fields and are labels only, not decision inputs
  - summary rows include period splits plus full_2021_2026 aggregate rows
- Latest smoke:
  - loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z
  - strategy=empty trades=0

Non-negotiables:
- Offline BTCUSDT 5m research only.
- Confirmed closed-candle decisions only.
- No entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.
- Keep generated CSV/JSON outputs under results/.
- Treat label_* columns as forward outcomes, not decision inputs.
- Update memory/PROGRESS.md with commands, result paths, and concise factual outcome after a completed milestone.
- Update memory/DECISIONS.md only if a durable constraint changes.
- After completing a brief or milestone, run closeout checks and commit the completed repo changes unless the user explicitly says not to commit.

Recommended next task:
Review the range regime durability outputs for split/time stability before testing another entry trigger. Decide whether the current range/compression regimes are durable, explainable, and narrow enough to use as future context. Do not add entries, exits, scoring, sizing, or strategy replacement in that review. Do not add a review verdict doc unless the audit outputs are actually reviewed in that same session.

Suggested verification for docs/memory-only closeouts:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- git diff --check
```
