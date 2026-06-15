# Next Codex Brief: Hold-Inside Directional Edge Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md, docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md, docs/RANGE_REGIME_DURABILITY_REVIEW.md, docs/COMPRESSION_BREAKOUT_REVIEW.md, docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/ENTRY_READINESS_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Compression breakout audit was not entry-ready.
- Range regime durability review says the balanced detector regimes are not durable enough as context for future entry hypotheses.
- Detector durability sweep review says no current DefaultDetectorSweepProfiles profile is approved as future entry context; p30_c12_bollinger_on_adx_on is diagnostic only.
- Detector context refinement review says the delayed hold_3_inside and hold_6_inside context rules are the first refinement that materially and split-stably reduces quick invalidation and trend leakage with adequate candidates. They are the leading decision-candle context, but they are NOT promoted to entry context: the gain is a heavy survivorship/conditioning effect, residual 12 bar trend leakage stays material, and the label_* fields are regime-durability outcomes, not P&L.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest detector context refinement review:
- Durable report: docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md
- Inputs reviewed (existing artifacts, not rerun):
  - results/detector-context-refinement-audit/detector_context_refinement_summary.csv/json
  - results/detector-context-refinement-audit/detector_context_refinement_stability.csv/json
  - results/detector-context-refinement-audit/detector_context_refinement_candidates.csv (sampled)
- Audit size: profiles=8, rules=5, candidate_rows=113824, summary_rows=640, stability_rows=160
- Lead context rules: hold_3_inside, hold_6_inside (closed-candle knowable at the decision candle; labels start at decision_index+1).

Recommended next task:
Stay non-trading. Implement a compact hold-inside directional edge audit that conditions on the leading hold-inside decision-candle context and measures forward, split-stable DIRECTIONAL outcomes, not just range survival. Concretely:
- Reuse the existing context machinery: rangeRegimeDurabilityEpisodes, DetectorContextRefinementRule / detectorContextRulePasses, decisionClosePosition / decisionClosePositionBucket, classifyDetectorSweepProfile, and the split helpers.
- Focus on the balanced baseline detector p30_c12_bollinger_on_adx_off with hold_3_inside and hold_6_inside context (optionally one stricter mid_50 variant for comparison).
- At the decision candle, record decision-candle features only (close position bucket, distance to frozen episode high/low, width and width/ATR context). Add forward label_* outcomes that capture direction, for example reversion toward range mid versus continued trend out of the range, favorable-minus-adverse move by side, and max favorable/adverse excursion, starting at decision_index+1.
- Keep horizons 1, 3, 6, 12 and the 3-bar quick-invalidation window for comparability.
- Emit compact CSV/JSON candidate, summary, and split-stability outputs under results/, plus a CLI flag, mirroring the existing audit modes.
- Add focused tests for no-lookahead decision features, label-window start, support/resistance (toward-high vs toward-low) symmetry, hold-rule filtering, missing-future skipping, invalid config, deterministic sorting, and summary denominators.
Do not add entries, exits, scoring, sizing, or strategy replacement. Do not add a durable verdict doc in the same session as the audit build; the next session reviews the outputs for whether a directional edge survives split-stably and net of an assumed cost before any entry trigger.

Non-negotiables:
- Offline BTCUSDT 5m research only.
- Confirmed closed-candle decisions only.
- No entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.
- Keep generated CSV/JSON outputs under results/.
- Treat label_* columns as forward outcomes, not decision inputs.
- Keep helper modules behind adapters and limited to feature extraction or audit outputs.
- Update memory/PROGRESS.md with commands, result paths, and concise factual outcome after a completed milestone.
- Update memory/DECISIONS.md only if a durable constraint changes.
- After completing a brief or milestone, run closeout checks and commit the completed repo changes unless the user explicitly says not to commit.

Suggested verification for an audit build:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv -out-dir results/<new-audit-dir> -<new-audit-flag>
- wc -l results/<new-audit-dir>/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
```
