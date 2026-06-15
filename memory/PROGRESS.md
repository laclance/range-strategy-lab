# Progress

## 2026-06-15

Hold-inside midline transition audit milestone:

- Added CLI flag `-hold-inside-midline-transition-audit`.
- Added compact non-trading hold-inside midline transition outputs:
  - `hold_inside_midline_transition_candidates.csv`
  - `hold_inside_midline_transition_candidates.json`
  - `hold_inside_midline_transition_summary.csv`
  - `hold_inside_midline_transition_summary.json`
  - `hold_inside_midline_transition_stability.csv`
  - `hold_inside_midline_transition_stability.json`
- Added exported lab API:
  - `DefaultHoldInsideMidlineTransitionAuditConfig`
  - `RunHoldInsideMidlineTransitionAudit`
  - `HoldInsideMidlineTransitionCandidateRow`
  - `HoldInsideMidlineTransitionSummaryRow`
  - `HoldInsideMidlineTransitionStabilityRow`
- Audit semantics:
  - uses balanced baseline detector `p30_c12_bollinger_on_adx_off`
  - uses context rules `hold_3_inside`, `hold_6_inside`, and diagnostic
    `hold_3_inside_mid_50`
  - emits one candidate row per passed source episode, context rule, and
    horizon
  - records only decision-candle-known fields: frozen high/low/mid, decision
    close position and bucket, decision mid side, distances to high/low/mid,
    episode width, ATR, and width/ATR context
  - records midline transition labels only; no `paper_side`, favorable/adverse
    fields, entries, exits, scoring, sizing, or trade-side interpretation
  - all forward `label_*` fields start at `decision_index + 1`
  - missing first-delay labels use `-1`
  - same-bar midline and boundary events do not satisfy "before"; ordering
    labels require strict earlier occurrence
  - summary rows aggregate by profile, context rule, split, horizon,
    `decision_mid_side`, and decision close position bucket, including `all`
  - source episode denominators remain independent of candidate rows across
    mid-side and bucket aggregations
  - stability rows compare only `2021_2022_stress`, `2023_2024_oos`, and
    `2025_2026_recent`
  - stability rows expose split min/max/delta rates for midline touch/cross,
    strict-before ordering, reentry, persistence, quick invalidation,
    invalidation up/down, and trend up/down labels
- This milestone did not add entries, exits, scoring, sizing, strategy
  replacement, live code, deploy scripts, API keys, grid, martingale,
  averaging down, or two-exchange execution.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.
- No durable verdict doc was added; the next step should review the generated
  outputs for split-stable midline transition evidence before any entry
  trigger.
- `memory/DECISIONS.md` was not changed because no durable constraint changed.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the next non-trading review
  handoff.

Latest hold-inside midline transition audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/hold-inside-midline-transition-audit \
  -hold-inside-midline-transition-audit

wc -l results/hold-inside-midline-transition-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```

Result:

- `go test ./...` passed.
- New hold-inside midline transition audit printed:
  - `hold_inside_midline_transition_audit profiles=1 rules=3 candidate_rows=7988 summary_rows=720 stability_rows=192 quick_invalidation_bars=3 horizons=1;3;6;12`
  - `loaded 569451 candles from 2021-01-01T00:00:00Z to
    2026-06-01T23:59:59Z`
  - `strategy=empty trades=0`
- New hold-inside midline transition CSV lines including header:
  - `hold_inside_midline_transition_candidates.csv`: `7,989`
  - `hold_inside_midline_transition_summary.csv`: `721`
  - `hold_inside_midline_transition_stability.csv`: `193`
  - base `summary.csv`: `13`
- Result paths:
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.csv`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.json`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.csv`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.json`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.csv`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.json`
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- `git diff --check` passed.

Hold-inside directional edge review milestone:

- Added durable review report:
  - `docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md`
- Updated `README.md` docs order to include the hold-inside directional edge
  review (now item `14`; `memory/NEXT_CODEX_BRIEF.md` is item `15`).
- Updated `memory/DECISIONS.md` with a durable no-promotion gate: the current
  hold-inside directional paper-side labels are not approved as entry context
  or strategy inputs.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the next non-trading
  hold-inside midline-transition audit handoff.
- Review verdict:
  - hold-inside directional edge audit is not entry-ready
  - `hold_3_inside` and `hold_6_inside` remain useful range-survival context,
    but do not show a split-stable directional edge toward the frozen range
    high or low
  - no primary all-bucket row passed the review gate of positive worst-split
    favorable-minus-adverse, worst-split favorable-greater-than-adverse above
    `50%`, and adequate candidate counts
  - no non-`all` decision-close-position bucket reached `100` candidates in
    every period split for `hold_3_inside`/`hold_6_inside`
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, strategy replacement, or convert
    `paper_side=toward_high`/`paper_side=toward_low` labels into trades
- Inputs reviewed:
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.json`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.json`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.json`
- Audit size:
  - profiles: `1`
  - context rules: `3`
  - paper sides: `2`
  - candidate rows: `15,976`
  - summary rows: `624`
  - stability rows: `168`
  - candidate CSV lines including header: `15,977`
  - summary CSV lines including header: `625`
  - stability CSV lines including header: `169`
  - horizons: `1`, `3`, `6`, `12`
  - quick invalidation window: `3` bars after the decision candle
- Compact evidence:
  - primary all-bucket `hold_3_inside`, `h3`, `toward_high`: `222` minimum
    split candidates, `+0.78bp` minimum favorable-minus-adverse, but only
    `46.56%` minimum favorable-greater-than-adverse
  - primary all-bucket `hold_3_inside`, `h6`, `toward_high`: `222` minimum
    split candidates, `+0.62bp` minimum favorable-minus-adverse, but only
    `45.51%` minimum favorable-greater-than-adverse
  - primary all-bucket `hold_6_inside`, `h6`, `toward_low`: `170` minimum split
    candidates and `50.18%` minimum favorable-greater-than-adverse, but
    `-0.03bp` minimum favorable-minus-adverse
  - all other primary all-bucket rows had negative minimum
    favorable-minus-adverse, minimum favorable-greater-than-adverse below
    `50%`, or both
  - non-`all` bucket rows for `hold_3_inside`/`hold_6_inside`: `0` rows with
    `>=100` candidates in every period split; `32` rows with `>=50` candidates
    in every period split, but the best relaxed positives were sparse or weak
    by FGTA

Latest hold-inside directional edge review verification:

```bash
wc -l results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- CSV line counts matched expected values (`15,977` / `625` / `169`).
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- Verification passed.

Hold-inside directional edge audit milestone:

- Added CLI flag `-hold-inside-directional-edge-audit`.
- Added compact non-trading hold-inside directional edge outputs:
  - `hold_inside_directional_edge_candidates.csv`
  - `hold_inside_directional_edge_candidates.json`
  - `hold_inside_directional_edge_summary.csv`
  - `hold_inside_directional_edge_summary.json`
  - `hold_inside_directional_edge_stability.csv`
  - `hold_inside_directional_edge_stability.json`
- Added exported lab API:
  - `DefaultHoldInsideDirectionalEdgeAuditConfig`
  - `RunHoldInsideDirectionalEdgeAudit`
  - `HoldInsideDirectionalEdgeCandidateRow`
  - `HoldInsideDirectionalEdgeSummaryRow`
  - `HoldInsideDirectionalEdgeStabilityRow`
- Audit semantics:
  - uses only balanced baseline detector `p30_c12_bollinger_on_adx_off`
  - uses context rules `hold_3_inside`, `hold_6_inside`, and
    `hold_3_inside_mid_50`
  - emits two paper-side labels per passed source episode:
    `toward_high` and `toward_low`
  - records only decision-candle-known features: frozen high/low/mid, decision
    close position and bucket, distance to high/low/mid, episode width, ATR,
    and width/ATR context
  - all forward `label_*` fields start at `decision_index + 1`
  - summary rows aggregate by profile, context rule, split, horizon, paper
    side, and decision close position bucket, including `all`
  - stability rows compare only `2021_2022_stress`, `2023_2024_oos`, and
    `2025_2026_recent`
- This milestone did not add entries, exits, scoring, sizing, strategy
  replacement, live code, deploy scripts, API keys, grid, martingale,
  averaging down, or two-exchange execution.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.
- No durable verdict doc was added; the next step should review the generated
  outputs for split-stable directional edge before any entry trigger.
- `memory/DECISIONS.md` was not changed because no durable constraint changed.

Latest hold-inside directional edge audit verification:

```bash
git status --short

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/hold-inside-directional-edge-audit \
  -hold-inside-directional-edge-audit

wc -l results/hold-inside-directional-edge-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```

Result:

- Initial `git status --short` was clean.
- `go test ./...` passed.
- New hold-inside directional edge audit printed:
  - `hold_inside_directional_edge_audit profiles=1 rules=3 paper_sides=2 candidate_rows=15976 summary_rows=624 stability_rows=168 quick_invalidation_bars=3 horizons=1;3;6;12`
  - `loaded 569451 candles from 2021-01-01T00:00:00Z to
    2026-06-01T23:59:59Z`
  - `strategy=empty trades=0`
- New hold-inside directional edge CSV lines including header:
  - `hold_inside_directional_edge_candidates.csv`: `15,977`
  - `hold_inside_directional_edge_summary.csv`: `625`
  - `hold_inside_directional_edge_stability.csv`: `169`
  - base `summary.csv`: `13`
- Result paths:
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.json`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.json`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv`
  - `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.json`
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- `git diff --check` passed.

Detector context refinement review milestone:

- Added durable review report:
  - `docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md`
- Updated `README.md` docs order to include the detector context refinement
  review (now item `13`; `memory/NEXT_CODEX_BRIEF.md` is item `14`).
- Updated `memory/DECISIONS.md` with a durable gate: the delayed
  `hold_3_inside`/`hold_6_inside` context rules are the leading context
  refinement and materially improve split-stable durability, but remain
  regime-durability filters and are not approved as entry context.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the next non-trading
  hold-inside directional edge handoff.
- Review verdict:
  - the delayed hold-inside context rules are the first refinement that
    materially and split-stably reduces quick invalidation and trend leakage
    with adequate candidates; lead rules are `hold_3_inside` and
    `hold_6_inside`
  - not promoted to entry context; the gain is heavy survivorship/conditioning,
    residual `12` bar trend leakage is still material, and the labels are
    regime-durability outcomes, not P&L
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this review
- Inputs reviewed:
  - `results/detector-context-refinement-audit/detector_context_refinement_summary.csv`
  - `results/detector-context-refinement-audit/detector_context_refinement_summary.json`
  - `results/detector-context-refinement-audit/detector_context_refinement_stability.csv`
  - `results/detector-context-refinement-audit/detector_context_refinement_stability.json`
  - `results/detector-context-refinement-audit/detector_context_refinement_candidates.csv` (sampled)
- Audit size:
  - profiles: `8`
  - context rules: `5`
  - candidate rows: `113,824`
  - summary rows: `640`
  - stability rows: `160`
  - candidate CSV lines including header: `113,825`
  - summary CSV lines including header: `641`
  - stability CSV lines including header: `161`
  - horizons: `1`, `3`, `6`, `12`
  - quick invalidation window: `3` bars after the decision candle
- Compact evidence (balanced baseline `p30_c12_bollinger_on_adx_off`, `12` bar
  horizon, worst split across the three period splits):
  - `episode_end`: `13.61%` min persisted, `70.43%` max quick invalidated,
    `51.79%` max trended, `742` min split candidates
  - `hold_3_inside`: `40.54%` min persisted, `25.32%` max quick invalidated,
    `40.09%` max trended, `222` min split candidates, `~30%` candidate rate
  - `hold_6_inside`: `44.71%` min persisted, `21.46%` max quick invalidated,
    `35.88%` max trended, `170` min split candidates, `~22%-24%` candidate rate
  - `hold_3_inside_mid_50`: `47.83%` min persisted, `16.28%` max quick
    invalidated, `33.06%` max trended, `94` min split candidates, `~14%`
    candidate rate
  - the same monotonic improvement held across all `8` profiles; the strongest
    worst-split `12` bar rows with at least `150` candidates per split were all
    delayed-hold rules
- Existing generated artifacts were used; the audit was not rerun because files
  were present and matched expected counts.

Latest detector context refinement review verification:

```bash
wc -l results/detector-context-refinement-audit/detector_context_refinement_candidates.csv results/detector-context-refinement-audit/detector_context_refinement_summary.csv results/detector-context-refinement-audit/detector_context_refinement_stability.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- CSV line counts matched expected values (`113,825` / `641` / `161`).
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- Verification passed.

Detector context refinement audit milestone:

- Added CLI flag `-detector-context-refinement-audit`.
- Added compact non-trading detector/context refinement outputs:
  - `detector_context_refinement_candidates.csv`
  - `detector_context_refinement_candidates.json`
  - `detector_context_refinement_summary.csv`
  - `detector_context_refinement_summary.json`
  - `detector_context_refinement_stability.csv`
  - `detector_context_refinement_stability.json`
- Added exported lab API:
  - `DefaultDetectorContextRefinementAuditConfig`
  - `DefaultDetectorContextRefinementProfiles`
  - `RunDetectorContextRefinementAudit`
  - `DetectorContextRefinementCandidateRow`
  - `DetectorContextRefinementSummaryRow`
  - `DetectorContextRefinementStabilityRow`
- Audit semantics:
  - uses 8 detector profiles from the detector durability review, including
    the ADX diagnostic variants
  - uses 5 closed-candle context rules: `episode_end`, `hold_1_inside`,
    `hold_3_inside`, `hold_6_inside`, and `hold_3_inside_mid_50`
  - freezes episode high/low at the original raw-active episode end
  - delayed rules set `decision_index = episode_end_index + hold_bars`
  - all `label_*` fields start at `decision_index + 1` and are forward
    outcomes only, not decision inputs
  - summary rows include source episode counts, candidate counts, candidate
    rates, and label rates by profile, context rule, split, and horizon
  - stability rows compare only `2021_2022_stress`, `2023_2024_oos`, and
    `2025_2026_recent`
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.
- No detector promotion verdict was added in this milestone; the next step
  should review the detector context refinement outputs before any detector
  promotion or entry-trigger work.

Latest detector context refinement audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-context-refinement-audit \
  -detector-context-refinement-audit

wc -l results/detector-context-refinement-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```

Result:

- `go test ./...` passed.
- New detector context refinement audit printed:
  - `detector_context_refinement_audit profiles=8 rules=5 candidate_rows=113824 summary_rows=640 stability_rows=160`
  - `quick_invalidation_bars=3`
  - `horizons=1;3;6;12`
  - `loaded 569451 candles from 2021-01-01T00:00:00Z to
    2026-06-01T23:59:59Z`
  - `strategy=empty trades=0`
- New detector context refinement CSV lines including header:
  - `detector_context_refinement_candidates.csv`: `113,825`
  - `detector_context_refinement_summary.csv`: `641`
  - `detector_context_refinement_stability.csv`: `161`
  - base `summary.csv`: `13`
- Result paths:
  - `results/detector-context-refinement-audit/detector_context_refinement_candidates.csv`
  - `results/detector-context-refinement-audit/detector_context_refinement_candidates.json`
  - `results/detector-context-refinement-audit/detector_context_refinement_summary.csv`
  - `results/detector-context-refinement-audit/detector_context_refinement_summary.json`
  - `results/detector-context-refinement-audit/detector_context_refinement_stability.csv`
  - `results/detector-context-refinement-audit/detector_context_refinement_stability.json`
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- `git diff --check` passed.

Detector durability sweep review milestone:

- Added durable review report:
  - `docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md`
- Updated `README.md` docs order to include the detector durability sweep
  review.
- Updated `memory/DECISIONS.md` with a durable gate: no profile in the current
  `DefaultDetectorSweepProfiles` detector durability sweep is approved as
  future entry context; `p30_c12_bollinger_on_adx_on` is diagnostic only.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the next non-trading
  detector/context refinement handoff.
- Review verdict:
  - no current sweep profile is approved as future entry context
  - `p30_c12_bollinger_on_adx_on` improves short-horizon persistence and quick
    invalidation, but is not promoted
  - the best broad `12` bar persistence floors remain below `18%`
  - the best fully specified `12` bar slices remain too weak or sparse
  - next implementation should refine or reframe detector/context first
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this review
- Inputs reviewed:
  - `results/detector-durability-sweep/detector_durability_sweep.csv`
  - `results/detector-durability-sweep/detector_durability_sweep.json`
  - `results/detector-durability-sweep/detector_durability_slices.csv`
  - `results/detector-durability-sweep/detector_durability_slices.json`
  - `results/detector-durability-sweep/detector_durability_stability.csv`
  - `results/detector-durability-sweep/detector_durability_stability.json`
- Audit size:
  - profiles: `19`
  - broad rows: `304`
  - slice rows: `9,088`
  - stability rows: `76`
  - broad CSV lines including header: `305`
  - slice CSV lines including header: `9,089`
  - stability CSV lines including header: `77`
  - horizons: `1`, `3`, `6`, `12`
  - quick invalidation window: `3` bars after episode end
- Compact evidence:
  - balanced baseline `p30_c12_bollinger_on_adx_off` had only `13.61%`
    minimum `12` bar persistence, `70.43%` maximum quick invalidation, and
    `51.79%` maximum trended rate across period splits
  - ADX comparison `p30_c12_bollinger_on_adx_on` improved `1` bar minimum
    persistence to `59.28%` and maximum quick invalidation to `40.72%`, but
    still had only `16.47%` minimum `12` bar persistence and `50.75%` maximum
    trended rate
  - the best broad `12` bar persistence floors were `17.57%`, `17.43%`, and
    `17.30%`, all with high quick-invalidation or trended rates
  - fully specified `12` bar slice counts by minimum episodes in every period
    split were `203`, `80`, `45`, `13`, `0`, and `0` at thresholds `10`,
    `25`, `50`, `100`, `250`, and `500`
  - the best fully specified `12` bar slice at the `100` episodes-per-split
    threshold had `26.95%` minimum persistence and `60.28%` maximum quick
    invalidation
- Existing generated artifacts were used; the detector durability sweep was not
  rerun because files were present and matched expected counts.

Latest detector durability sweep review verification:

```bash
wc -l results/detector-durability-sweep/detector_durability_sweep.csv results/detector-durability-sweep/detector_durability_slices.csv results/detector-durability-sweep/detector_durability_stability.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- CSV line counts matched expected values.
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- Verification passed.

Detector durability sweep milestone:

- Added CLI flag `-detector-durability-sweep`.
- Added compact non-trading detector durability sweep outputs:
  - `detector_durability_sweep.csv`
  - `detector_durability_sweep.json`
  - `detector_durability_slices.csv`
  - `detector_durability_slices.json`
  - `detector_durability_stability.csv`
  - `detector_durability_stability.json`
- Added exported lab API:
  - `RunDetectorDurabilitySweep`
  - `DetectorDurabilitySweepRow`
  - `DetectorDurabilitySliceRow`
  - `DetectorDurabilityStabilityRow`
- Sweep semantics:
  - reuses the existing `DefaultDetectorSweepProfiles` 19-profile grid
  - applies existing range-regime durability label semantics to each profile
  - broad rows are one row per profile, split, and horizon
  - slice rows use existing raw length, active length, width, and width/ATR
    buckets
  - stability rows compare only `2021_2022_stress`, `2023_2024_oos`, and
    `2025_2026_recent`
  - all `label_*` fields are forward outcomes only, not decision inputs
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.
- No detector promotion verdict was added in this milestone; the next step
  should review the detector durability sweep outputs before any detector
  promotion or entry-trigger work.

Latest detector durability sweep verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-durability-sweep \
  -detector-durability-sweep

wc -l results/detector-durability-sweep/detector_durability_sweep.csv results/detector-durability-sweep/detector_durability_slices.csv results/detector-durability-sweep/detector_durability_stability.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```

Result:

- `go test ./...` passed.
- New detector durability sweep printed:
  - `detector_durability_sweep profiles=19 broad_rows=304 slice_rows=9088 stability_rows=76`
  - `quick_invalidation_bars=3`
  - `horizons=1;3;6;12`
  - `loaded 569451 candles from 2021-01-01T00:00:00Z to
    2026-06-01T23:59:59Z`
  - `strategy=empty trades=0`
- New detector durability sweep CSV lines including header:
  - `detector_durability_sweep.csv`: `305`
  - `detector_durability_slices.csv`: `9,089`
  - `detector_durability_stability.csv`: `77`
- Result paths:
  - `results/detector-durability-sweep/detector_durability_sweep.csv`
  - `results/detector-durability-sweep/detector_durability_sweep.json`
  - `results/detector-durability-sweep/detector_durability_slices.csv`
  - `results/detector-durability-sweep/detector_durability_slices.json`
  - `results/detector-durability-sweep/detector_durability_stability.csv`
  - `results/detector-durability-sweep/detector_durability_stability.json`
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- `git diff --check` passed.

Range regime durability review milestone:

- Added durable review report:
  - `docs/RANGE_REGIME_DURABILITY_REVIEW.md`
- Updated `README.md` docs order to include the range regime durability review.
- Updated `memory/DECISIONS.md` with a durable gate: the current
  `p30_c12_bollinger_on_adx_off` detector is not approved as context for future
  entry hypotheses until detector/context refinement is reviewed.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the new detector/context
  refinement handoff.
- Review verdict:
  - current balanced detector regimes are not durable enough to use as context
    for future entry hypotheses
  - split stability mostly reflects weakness repeating across splits, not high
    regime quality
  - next implementation should refine detector/context first
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this audit
- Inputs reviewed:
  - `results/range-regime-durability-audit/range_regime_durability_summary.csv`
  - `results/range-regime-durability-audit/range_regime_durability_episodes.csv`
- Audit size:
  - episode rows: `11,984`
  - unique episodes: `2,996`
  - summary rows: `452`
  - episode CSV lines including header: `11,985`
  - summary CSV lines including header: `453`
  - detector profile: `p30_c12_bollinger_on_adx_off`
  - horizons: `1`, `3`, `6`, `12`
  - quick invalidation window: `3` bars after the episode end
- Compact evidence:
  - full-sample persistence fell from `44.66%` at `1` bar to `14.95%` at
    `12` bars
  - full-sample quick invalidation reached `70.06%` by the `3` bar horizon and
    stayed there
  - the period splits had similar broad weakness: `12` bar persistence was
    `15.68%`, `15.07%`, and `13.61%`
  - no fully specified bucket slice had at least `100` episodes in every period
    split
  - the best `12` bar fully specified slice with at least `25` episodes in
    every split still had only `23.93%` minimum persistence and `58.12%`
    maximum quick invalidation
  - the broad `12` bar width-only `gt_50bp` slice still had only `16.80%`
    minimum persistence and `68.07%` maximum quick invalidation
- Existing generated artifacts were used; the audit smoke was not rerun because
  files were present and matched expected counts.

Latest range regime durability review verification:

```bash
wc -l results/range-regime-durability-audit/range_regime_durability_episodes.csv results/range-regime-durability-audit/range_regime_durability_summary.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- CSV line counts matched expected values.
- `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt.
- Verification passed.

## 2026-06-14

Range regime durability audit milestone:

- Added CLI flag `-range-regime-durability-audit`.
- Added compact non-trading range regime durability outputs:
  - `range_regime_durability_episodes.csv`
  - `range_regime_durability_episodes.json`
  - `range_regime_durability_summary.csv`
  - `range_regime_durability_summary.json`
- Defaults:
  - balanced compression detector profile:
    `p30_c12_bollinger_on_adx_off`
  - horizons: `1`, `3`, `6`, `12` bars after the detected episode end
  - quick invalidation window: `3` bars after the detected episode end
- Episode semantics:
  - episodes are contiguous `RawActive` detector runs that eventually become
    `Active`
  - episode high/low, width, length, and ATR context use only closed candles
    through the episode end
  - all forward durability metrics are `label_*` fields and start at
    `episode_end_index + 1`
  - summary rows include period splits plus `full_2021_2026` aggregate rows
- Added focused tests for no-lookahead raw-run episode construction,
  deterministic output, label-window start, split/full-summary handling,
  invalid config, missing-future skipping, and summary denominators.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.
- No durable review verdict document was added in this milestone; the next step
  should review regime durability stability before any entry trigger work.

Latest range regime durability audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/range-regime-durability-audit \
  -range-regime-durability-audit

wc -l results/range-regime-durability-audit/range_regime_durability_episodes.csv results/range-regime-durability-audit/range_regime_durability_summary.csv
git diff --check
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- New range regime durability audit printed:
  - `range_regime_durability_audit episode_rows=11984 summary_rows=452`
  - `quick_invalidation_bars=3`
  - `horizons=1;3;6;12`
  - `detector_profile_id=p30_c12_bollinger_on_adx_off`
  - `loaded 569451 candles from 2021-01-01T00:00:00Z to
    2026-06-01T23:59:59Z`
  - `strategy=empty trades=0`
- New range regime durability audit CSV lines including header:
  - `range_regime_durability_episodes.csv`: `11,985`
  - `range_regime_durability_summary.csv`: `453`
- Result paths:
  - `results/range-regime-durability-audit/range_regime_durability_episodes.csv`
  - `results/range-regime-durability-audit/range_regime_durability_episodes.json`
  - `results/range-regime-durability-audit/range_regime_durability_summary.csv`
  - `results/range-regime-durability-audit/range_regime_durability_summary.json`

Compression breakout review milestone:

- Added durable review report:
  - `docs/COMPRESSION_BREAKOUT_REVIEW.md`
- Updated `README.md` docs order to include the compression breakout review.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the new verdict and next
  non-trading handoff.
- Review verdict:
  - compression breakout audit is not entry-ready
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this audit
  - do not continue narrow compression-breakout slicing unless the next
    hypothesis changes materially
- Inputs reviewed:
  - `results/compression-breakout-audit/compression_breakout_summary.csv`
  - `results/compression-breakout-audit/compression_breakout_summary.json`
  - `results/compression-breakout-audit/compression_breakout_candidates.csv`
  - `results/compression-breakout-audit/compression_breakout_candidates.json`
- Audit size:
  - candidate rows: `5,096`
  - summary rows: `24`
  - candidate CSV lines including header: `5,097`
  - summary CSV lines including header: `25`
  - breakout decisions across all splits: `2,548` per horizon
  - one-bar side counts: `1,290` up breakouts and `1,258` down breakouts
- Compact evidence:
  - broad side/horizon favorable-minus-adverse was positive, about `+2.09bp`
    to `+2.97bp`, but every broad FGTA row was below `50%`
  - `2025_2026_recent` up breakouts were negative at `3`, `6`, and `12` bars
  - down breakouts were negative in `2021_2022_stress` at `12` bars and in
    `2025_2026_recent` at `1` bar
  - range re-entry rose to about `65.97%` for up and `70.11%` for down by the
    `12` bar horizon
  - fully sliced, no-delay, and breakout-shape cohorts had `0` stable positive
    rows at every threshold from `25` to `500` candidates per split
  - no positive stable cohort survived the `50` candidates-per-split threshold
    in any reviewed grouping
- Existing generated artifacts were used; the audit smoke was not rerun because
  files were present and matched expected counts.

Latest compression breakout review verification:

```bash
wc -l results/compression-breakout-audit/compression_breakout_candidates.csv results/compression-breakout-audit/compression_breakout_summary.csv
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- CSV line counts matched expected values.
- Verification passed.

Compression breakout audit milestone:

- Added CLI flag `-compression-breakout-audit`.
- Added compact non-trading compression breakout outputs:
  - `compression_breakout_candidates.csv`
  - `compression_breakout_candidates.json`
  - `compression_breakout_summary.csv`
  - `compression_breakout_summary.json`
- Defaults:
  - balanced compression detector profile:
    `p30_c12_bollinger_on_adx_off`
  - max closed-candle breakout delay: `12` bars after the raw-active
    compression episode ends
  - horizons: `1`, `3`, `6`, `12` bars after the breakout candle
- Decision semantics:
  - episodes are contiguous `RawActive` detector runs that eventually become
    `Active`
  - episode high/low are frozen using only closed candles through the episode
    end
  - the first close above the frozen episode high or below the frozen episode
    low within `12` bars is the decision breakout candle
  - all forward outcome metrics remain `label_*` fields and start after the
    breakout candle
- Added focused tests for raw-run episode construction, frozen episode bounds,
  first close-break selection, up/down symmetry, forward label-window start,
  no-break/missing-future skipping, invalid config, deterministic sorting, and
  aggregation denominators.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.
- No durable promotion or no-promotion review document was added in this
  milestone; the next step should review split and cohort stability before any
  entries.

Latest compression breakout audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/compression-breakout-audit \
  -compression-breakout-audit

wc -l results/compression-breakout-audit/compression_breakout_candidates.csv results/compression-breakout-audit/compression_breakout_summary.csv
awk -F, 'NR>1 {counts[$3]+=$5} END {for (h in counts) print h, counts[h]}' results/compression-breakout-audit/compression_breakout_summary.csv
awk -F, 'NR>1 && $3==1 {counts[$2]+=$5} END {for (s in counts) print s, counts[s]}' results/compression-breakout-audit/compression_breakout_summary.csv
```

Result:

- `go test ./...` passed.
- New compression breakout audit printed:
  - `compression_breakout_audit candidate_rows=5096 summary_rows=24`
  - `max_breakout_delay=12`
  - `horizons=1;3;6;12`
  - `detector_profile_id=p30_c12_bollinger_on_adx_off`
  - `strategy=empty trades=0`
- New compression breakout audit CSV lines including header:
  - `compression_breakout_candidates.csv`: `5,097`
  - `compression_breakout_summary.csv`: `25`
- Result paths:
  - `results/compression-breakout-audit/compression_breakout_candidates.csv`
  - `results/compression-breakout-audit/compression_breakout_candidates.json`
  - `results/compression-breakout-audit/compression_breakout_summary.csv`
  - `results/compression-breakout-audit/compression_breakout_summary.json`
- Compact aggregate read of `compression_breakout_summary.csv`:
  - breakout decisions across all splits: `2,548` per horizon
  - one-bar horizon side counts: `1,290` up breakouts and `1,258` down
    breakouts
  - summary rows are split by `2021_2022_stress`, `2023_2024_oos`, and
    `2025_2026_recent`

False-break reclaim timing review milestone:

- Added durable review report:
  - `docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md`
- Updated `README.md` docs order to include the review note.
- Updated `memory/DECISIONS.md` to record that completed briefs/milestones
  should be checked and committed automatically unless the user says not to
  commit.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the new verdict and next
  non-trading handoff.
- Review verdict:
  - false-break reclaim timing audit is not entry-ready
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this audit
  - do not continue narrow SR timing slices unless the next hypothesis changes
    materially
- Inputs reviewed:
  - `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_summary.csv`
  - `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_candidates.csv`
- Audit size:
  - candidate rows: `17,652`
  - summary rows: `24`
  - candidate CSV lines including header: `17,653`
  - summary CSV lines including header: `25`
  - support reclaim decisions across all splits: `4,120` per horizon
  - resistance reclaim decisions across all splits: `4,150` per horizon
- Compact evidence:
  - broad side/horizon favorable-minus-adverse was small-positive, topping at
    `+2.41bp`
  - broad FGTA was mostly near coin-flip, about `49.27%` to `51.54%`
  - `2025_2026_recent` support was negative across all horizons
  - `2023_2024_oos` resistance was negative at `3`, `6`, and `12` bars
  - fully sliced cohorts had `0` stable rows at every threshold from `25` to
    `500` candidates per split
  - coarse cohorts had no stable rows at `500` candidates per split; close
    shape cohorts survived at `500`, but the best row had only `+0.51bp` min
    split diff and `46.69%` min FGTA
- Existing generated artifacts were used; the audit smoke was not rerun because
  files were present and matched expected counts.

Latest false-break reclaim review verification:

```bash
wc -l results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_candidates.csv results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_summary.csv
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- CSV line counts matched expected values.
- Verification passed.

False-break reclaim timing audit milestone:

- Added CLI flag `-sr-false-break-reclaim-timing-audit`.
- Added compact non-trading false-break reclaim outputs:
  - `sr_false_break_reclaim_timing_candidates.csv`
  - `sr_false_break_reclaim_timing_candidates.json`
  - `sr_false_break_reclaim_timing_summary.csv`
  - `sr_false_break_reclaim_timing_summary.json`
- Defaults:
  - max closed-candle break delay: `3` bars after the anchor candle
  - max closed-candle reclaim delay: `12` bars after the break candle
  - horizons: `1`, `3`, `6`, `12` bars after the reclaim candle
  - `detector_active=true` anchor rows only
- Decision semantics:
  - support false break closes below the frozen anchor support zone bottom,
    then reclaims with a close at or above the frozen anchor support level
  - resistance false break closes above the frozen anchor resistance zone top,
    then reclaims with a close at or below the frozen anchor resistance level
  - the reclaim candle is the decision candle
  - all forward outcome metrics remain `label_*` fields and start after the
    reclaim candle
- Added focused tests for no-lookahead decision features, label-window start,
  support/resistance symmetry, detector-active filtering, no-break/no-reclaim
  skipping, end-of-data skipping, invalid config, candidate aggregation, and
  summary denominators.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest false-break reclaim audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-false-break-reclaim-timing-audit \
  -sr-false-break-reclaim-timing-audit
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- New false-break reclaim audit printed:
  - `sr_false_break_reclaim_timing_audit candidate_rows=17652 summary_rows=24`
  - `max_break_delay=3`
  - `max_reclaim_delay=12`
  - `horizons=1;3;6;12`
  - `detector_active_only=true`
  - `strategy=empty trades=0`
- New false-break reclaim audit CSV lines including header:
  - `sr_false_break_reclaim_timing_candidates.csv`: `17,653`
  - `sr_false_break_reclaim_timing_summary.csv`: `25`
- Result paths:
  - `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_candidates.csv`
  - `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_candidates.json`
  - `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_summary.csv`
  - `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_summary.json`
- Compact aggregate read of `sr_false_break_reclaim_timing_summary.csv`:
  - support reclaim decisions across all splits: `4,120` per horizon
  - resistance reclaim decisions across all splits: `4,150` per horizon
  - broad side/horizon favorable-minus-adverse was small-positive:
    - support: about `+1.03bp`, `+1.76bp`, `+1.61bp`, `+0.75bp` at
      horizons `1`, `3`, `6`, `12`
    - resistance: about `+0.84bp`, `+1.28bp`, `+1.85bp`, `+2.41bp` at
      horizons `1`, `3`, `6`, `12`
  - broad favorable-greater-than-adverse rates were about `49.3%` to `51.5%`
- No durable promotion or no-promotion review document was added in this
  milestone; the next step should review split and cohort stability before any
  entries.

## 2026-06-13

SR confirmation timing review milestone:

- Added durable review report:
  - `docs/SR_CONFIRMATION_TIMING_REVIEW.md`
- Updated `README.md` docs order to include the review note.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` for the next materially different
  non-trading hypothesis.
- Review verdict:
  - delayed confirmation after SR rejection is not entry-ready
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this audit
  - do not continue broad rejection-confirmation slicing unless the next
    hypothesis changes materially
- Inputs reviewed:
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.csv`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.csv`
- Audit size:
  - candidate rows: `9,692`
  - summary rows: `72`
  - candidate CSV lines including header: `9,693`
  - summary CSV lines including header: `73`
- Compact evidence:
  - broad decision-confirmation cohorts were about `48.9%` to `49.6%` of seed
    rejection candidates
  - side/delay/horizon aggregate favorable-minus-adverse topped near `+0.96bp`
  - recent split support turned flat or negative in several rows
  - FGTA was mostly below `50%`
  - with at least `500` rows in every split, best minimum split diff fell to
    about `+0.16bp`
  - with at least `750` rows in every split, no stable cohorts remained

Latest confirmation-review verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result:

- Verification passed.
- Existing generated artifacts were used; the audit smoke was not rerun because
  `results/sr-confirmation-timing-audit/` was present and current.

SR confirmation timing audit milestone:

- Added CLI flag `-sr-confirmation-timing-audit`.
- Added compact non-trading delayed-confirmation outputs:
  - `sr_confirmation_timing_candidates.csv`
  - `sr_confirmation_timing_candidates.json`
  - `sr_confirmation_timing_summary.csv`
  - `sr_confirmation_timing_summary.json`
- Defaults:
  - confirmation delays: `1`, `2`, `3` bars after the seed rejection candle
  - horizons: `1`, `3`, `6`, `12` bars after the confirmation candle
  - `detector_active=true` seed rows only
- Decision semantics:
  - seed candle must be an existing SR rejection candidate
  - confirmation candle is the decision candle
  - all forward outcome metrics remain `label_*` fields and start after the
    confirmation candle
- Added focused tests for no-lookahead decision features, label-window start,
  support/resistance symmetry, seed filtering, end-of-data skipping, invalid
  config, candidate aggregation, and summary denominators.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest confirmation-audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-confirmation-timing-audit \
  -sr-confirmation-timing-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-confirmation-combined-compat-check \
  -sr-audit \
  -sr-boundary-audit \
  -sr-boundary-inspect \
  -sr-rejection-timing-audit \
  -sr-confirmation-timing-audit
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- New confirmation audit printed:
  - `sr_confirmation_timing_audit candidate_rows=9692 summary_rows=72`
  - `delays=1;2;3`
  - `horizons=1;3;6;12`
  - `detector_active_only=true`
  - `strategy=empty trades=0`
- New confirmation audit CSV lines including header:
  - `sr_confirmation_timing_candidates.csv`: `9,693`
  - `sr_confirmation_timing_summary.csv`: `73`
- Result paths:
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.csv`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.json`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.csv`
  - `results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.json`
- Combined SR compatibility check preserved existing counts:
  - SR audit rows: `569,313`
  - near-support rows: `126,861`
  - near-resistance rows: `128,085`
  - boundary events: `281,080`
  - boundary quality rows: `192`
  - boundary inspect comparison rows: `192`
  - rejection timing rows: `968` candidates / `24` summary rows
  - confirmation timing rows: `9,692` candidates / `72` summary rows
- Compact aggregate read of `sr_confirmation_timing_summary.csv`:
  - broad side/delay/horizon decision-confirmation cohorts were about
    `48.9%` to `49.6%` of seed rejection candidates
  - decision-candidate favorable-minus-adverse was small-positive across the
    aggregate side/delay/horizon grid, topping out at about `+0.96bp`
  - favorable-greater-than-adverse rates were still mostly below `50%`
- No durable promotion or no-promotion review document was added in this
  milestone; the next step should review split and cohort stability before any
  entries.

SR rejection timing review milestone:

- Added durable review report:
  - `docs/SR_REJECTION_TIMING_REVIEW.md`
- Updated `README.md` docs order to include the review note.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` for the next non-trading
  confirmation-audit milestone.
- Review verdict:
  - boundary-rejection timing audit is not entry-ready
  - keep `lab.EmptyStrategy`
  - trades remain `0`
  - do not add entries, exits, scoring, sizing, or strategy replacement from
    this audit alone
- Compact evidence:
  - broad support decision candidates were flat-to-negative by
    favorable-minus-adverse across most horizons
  - resistance was better but still tiny in aggregate, topping out at
    `+0.45bp` decision-candidate favorable-minus-adverse at `12` bars
  - top OOS/recent cohorts were split-specific or side-specific rather than one
    simple support/resistance-symmetric template
  - common h12 in-zone strength-2 shape was unstable across pierced state and
    splits
- Result paths reviewed:
  - `results/sr-rejection-timing-audit/sr_rejection_timing_summary.csv`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_candidates.csv`
- Next recommended work:
  - stay non-trading
  - test delayed confirmation after an SR rejection candle, re-indexed so the
    confirmation candle is the decision candle and future `label_*` outcomes
    start after that confirmation candle

Latest review verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-rejection-timing-audit \
  -sr-rejection-timing-audit
```

Timing audit smoke result:

- `sr_rejection_timing_audit candidate_rows=968 summary_rows=24`
- `strategy=empty trades=0`

Boundary-rejection timing audit milestone:

- Added CLI flag `-sr-rejection-timing-audit`.
- Added compact non-trading outputs:
  - `sr_rejection_timing_candidates.csv`
  - `sr_rejection_timing_candidates.json`
  - `sr_rejection_timing_summary.csv`
  - `sr_rejection_timing_summary.json`
- Candidate cohorts group decision-candle features separately from forward
  labels:
  - side, horizon, close location, touched/pierced/closed-back state
  - wick-beyond, strength, and distance buckets
  - balanced detector context
  - all forward outcome metrics use `label_` prefixes
- Added tests for no-lookahead decision features, support/resistance symmetry,
  touch/pierce/closed-back behavior, bucket behavior, detector-active
  filtering, missing-future skipping, and summary denominators.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest timing-audit verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection-check \
  -sr-boundary-inspect

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-rejection-timing-audit \
  -sr-rejection-timing-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-combined-compat-check \
  -sr-audit \
  -sr-boundary-audit \
  -sr-boundary-inspect \
  -sr-rejection-timing-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-both-check \
  -detector \
  -detector-sweep
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- Existing boundary inspect check printed:
  - `sr_boundary_inspect events=281080 comparison_rows=192`
  - `strategy=empty trades=0`
- New timing audit printed:
  - `sr_rejection_timing_audit candidate_rows=968 summary_rows=24`
  - `strategy=empty trades=0`
- New timing audit CSV lines including header:
  - `sr_rejection_timing_candidates.csv`: `969`
  - `sr_rejection_timing_summary.csv`: `25`
- Result paths:
  - `results/sr-rejection-timing-audit/sr_rejection_timing_candidates.csv`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_candidates.json`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_summary.csv`
  - `results/sr-rejection-timing-audit/sr_rejection_timing_summary.json`
- Combined SR compatibility check preserved existing counts:
  - SR audit rows: `569,313`
  - near-support rows: `126,861`
  - near-resistance rows: `128,085`
  - boundary events: `281,080`
  - boundary quality rows: `192`
  - boundary inspect comparison rows: `192`
  - timing audit rows: `968` candidates / `24` summary rows
- Detector compatibility check preserved existing counts:
  - detector active bars: `77,231`
  - detector total bars: `569,451`
  - detector episodes: `2,996`
  - detector sweep profiles: `19`
  - detector sweep rows: `76`
- Every smoke/compatibility run printed `strategy=empty trades=0`.

Follow-up:

- Superseded by the SR rejection timing review milestone above.
- The review found this boundary-rejection timing shape is not entry-ready.

Entry-readiness review gate:

- Added durable review report:
  - `docs/ENTRY_READINESS_REVIEW.md`
- Added focused tests across CSV loading, detector/default behavior, detector
  sweep, engine, indicators, SR audit helpers, SR boundary audit helpers, and
  empty strategy behavior.
- Fixed two defaulting bugs found during review:
  - zero-value `RangeDetectorConfig` now returns
    `DefaultCompressionRangeDetectorConfig()`, including `UseBollinger=true`
  - zero-value `SRBoundaryAuditConfig` now returns
    `DefaultSRBoundaryAuditConfig()`, including `DetectorActiveOnly=true`
- Review verdict:
  - the codebase passes the review gate for the next non-trading timing audit
  - first trade entries should still wait until the timing audit proves a
    closed-candle rejection signal can be identified without future bars
- Updated `memory/NEXT_CODEX_BRIEF.md` for the boundary-rejection timing audit.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest review verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./cmd/rangelab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./internal/lab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -coverprofile=/tmp/range-strategy-lab-internal.cover ./internal/lab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go tool cover -func=/tmp/range-strategy-lab-internal.cover
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection-check \
  -sr-boundary-inspect
```

Result:

- `go test ./...` passed.
- `git diff --check` passed.
- `go test ./cmd/rangelab` passed.
- `internal/lab` coverage: `99.8%`.
- Remaining uncovered internal statements are justified defensive SR audit error
  propagation branches in `RunSRAudit`.
- `go test -cover ./cmd/rangelab` and `go test -cover ./...` still fail with
  the known snap-confine issue for `cmd/rangelab` coverage.
- SR boundary inspection check printed:
  - `sr_boundary_inspect events=281080 comparison_rows=192`
  - `loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z`
  - `strategy=empty trades=0`

Repository setup milestone:

- Set up tracked project memory in `memory/`.
- Added root `AGENTS.md` so future Codex sessions read and maintain memory.
- Replaced the starter Codex brief with the detector sweep/audit brief.

Detector-only milestone completed before repository setup:

- Added indicator helpers:
  - normalized ATR14
  - Donchian20 width
  - Bollinger20 width
  - optional ADX14
- Added `CompressionRangeDetector` with detector-only diagnostics.
- Added detector outputs:
  - `detector_duty_cycle.csv`
  - `detector_duty_cycle.json`
  - `range_episodes.csv`
  - `range_episodes.json`
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest detector smoke run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-smoke \
  -detector
```

Observed `full_2021_2026` detector metrics:

- Active bars: `77,231`
- Total bars: `569,451`
- Duty cycle: `13.5624%`
- Episodes: `2,996`

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
```

Result: passed.

Detector sweep/audit milestone:

- Added CLI flag `-detector-sweep`.
- Added a compact detector sweep:
  - percentile: `0.20`, `0.30`, `0.40`
  - min consecutive bars: `6`, `12`, `24`
  - Bollinger on/off
  - ADX off for the grid, plus one balanced ADX-on comparison
- Added outputs:
  - `detector_sweep.csv`
  - `detector_sweep.json`
- Marked the balanced baseline:
  - `p30_c12_bollinger_on_adx_off`
  - percentile: `0.30`
  - min consecutive bars: `12`
  - Bollinger: on
  - ADX: off
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest detector sweep run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep \
  -detector-sweep
```

Observed detector sweep facts:

- Output rows: `76` (`19` profiles x `4` splits)
- Balanced baseline full split:
  - Active bars: `77,231`
  - Total bars: `569,451`
  - Duty cycle: `13.5624%`
  - Episodes: `2,996`
- All profiles had nonzero episodes in every period split.
- Profiles that roughly fit the first-pass usability screen (`5%`-`25%` full duty, nonzero episodes in every split, no obviously unstable split duty):
  - all ADX-off profiles except `p20_c24_bollinger_on_adx_off`, `p40_c06_bollinger_on_adx_off`, and `p40_c06_bollinger_off_adx_off`
  - `p40_c12_bollinger_off_adx_off` is near the upper edge
- The balanced ADX-on comparison was too restrictive for the first-pass screen:
  - `p30_c12_bollinger_on_adx_on`
  - full duty cycle: `4.36%`

Latest detector compatibility check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector
```

Result:

- Existing detector outputs still write.
- `strategy=empty trades=0`.

Latest combined detector/sweep check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-both \
  -detector \
  -detector-sweep
```

Result:

- Existing detector outputs and detector sweep outputs write in one run.
- `strategy=empty trades=0`.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
```

Result: passed.

Research helper dependency milestone:

- Added pinned pure-Go helper modules:
  - `github.com/laclance/go-sr v1.0.0`
  - `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
  - `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`
- Added `docs/RESEARCH_HELPERS.md`.
- Updated docs to keep helper modules behind adapters and audit outputs.
- No strategy entries, exits, scoring, live code, or generated result artifacts
  were added.

Dependency add command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go get \
  github.com/laclance/go-sr@v1.0.0 \
  github.com/markcheno/go-talib@v0.0.0-20250114000313-ec55a20c902f \
  nproject.io/gitlab/libraries/talib-cdl-go@v0.0.0-20211217160304-2ed8176448cc
```

Verification:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

SR audit milestone:

- Added CLI flag `-sr-audit`.
- Added `go-sr` zone-mode adapter behind lab-owned audit output:
  - timeframe: `5m`
  - lookback bars: `120`
  - min strength: `2`
  - warmup bars: `138`
- Added outputs:
  - `sr_touch_audit.csv`
  - `sr_touch_audit.json`
- Included balanced detector context on each SR row:
  - `detector_profile_id=p30_c12_bollinger_on_adx_off`
  - `detector_raw_active`
  - `detector_active`
- Refreshed the next Codex brief for the next SR boundary-inspection step.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest SR audit smoke run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-audit-smoke \
  -sr-audit
```

Observed SR audit facts:

- Output rows: `569,313`
- CSV lines including header: `569,314`
- Near-support rows: `126,861`
- Near-resistance rows: `128,085`
- Result paths:
  - `results/sr-audit-smoke/sr_touch_audit.csv`
  - `results/sr-audit-smoke/sr_touch_audit.json`
  - `results/sr-audit-smoke/summary.csv`
  - `results/sr-audit-smoke/summary.json`
  - `results/sr-audit-smoke/trades.json`
- `strategy=empty trades=0`.

Latest detector compatibility check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector
```

Result:

- Existing detector outputs still write.
- Full split active bars: `77,231`
- Full split total bars: `569,451`
- Full split duty cycle: `13.56%`
- Episodes: `2,996`
- `strategy=empty trades=0`.

Latest detector sweep compatibility check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep-check \
  -detector-sweep
```

Result:

- Existing detector sweep outputs still write.
- Profiles: `19`
- Rows: `76`
- Balanced baseline active bars: `77,231`
- Balanced baseline total bars: `569,451`
- Balanced baseline duty cycle: `13.56%`
- Balanced baseline episodes: `2,996`
- `strategy=empty trades=0`.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

Implementation verification rerun:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-audit-smoke \
  -sr-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep-check \
  -detector-sweep
```

Result:

- Tests passed and `git diff --check` was clean.
- SR audit wrote:
  - `results/sr-audit-smoke/sr_touch_audit.csv`
  - `results/sr-audit-smoke/sr_touch_audit.json`
- SR audit rows: `569,313`.
- SR audit CSV lines including header: `569,314`.
- Near-support rows: `126,861`.
- Near-resistance rows: `128,085`.
- Detector compatibility still writes `results/detector-check/detector_duty_cycle.csv/json`.
- Detector full split: `77,231` active bars / `569,451` total bars, `13.56%` duty cycle, `2,996` episodes.
- Detector sweep compatibility still writes `results/detector-sweep-check/detector_sweep.csv/json`.
- Detector sweep full baseline: `77,231` active bars / `569,451` total bars, `13.56%` duty cycle, `2,996` episodes.
- Every smoke/compatibility run printed `strategy=empty trades=0`.

SR boundary-quality audit milestone:

- Added CLI flag `-sr-boundary-audit`.
- Added non-trading SR boundary quality outputs:
  - `sr_boundary_events.csv`
  - `sr_boundary_events.json`
  - `sr_boundary_quality.csv`
  - `sr_boundary_quality.json`
- Defaults:
  - horizons: `1`, `3`, `6`, `12` bars
  - `detector_active=true` rows only
  - one event per near boundary side
  - skip event/horizon pairs without enough future candles
- Metrics include favorable/adverse forward move, wick break, close break,
  reclaim-after-break, rejection, strength bucket, and distance bucket.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest SR boundary-quality smoke run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-quality \
  -sr-boundary-audit
```

Observed SR boundary-quality facts:

- Boundary event rows: `281,080`.
- Boundary event CSV lines including header: `281,081`.
- Boundary quality rows: `192`.
- Boundary quality CSV lines including header: `193`.
- Result paths:
  - `results/sr-boundary-quality/sr_boundary_events.csv`
  - `results/sr-boundary-quality/sr_boundary_events.json`
  - `results/sr-boundary-quality/sr_boundary_quality.csv`
  - `results/sr-boundary-quality/sr_boundary_quality.json`
- `strategy=empty trades=0`.

Combined SR audit/boundary check:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-combined-check \
  -sr-audit \
  -sr-boundary-audit
```

Result:

- `sr_touch_audit.csv` lines including header: `569,314`.
- `sr_boundary_events.csv` lines including header: `281,081`.
- `sr_boundary_quality.csv` lines including header: `193`.
- `strategy=empty trades=0`.

Latest compatibility checks after SR boundary-quality work:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-audit-smoke \
  -sr-audit

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-check \
  -detector

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/detector-sweep-check \
  -detector-sweep
```

Result:

- SR audit still writes `sr_touch_audit.csv/json`.
- SR audit rows: `569,313`.
- Near-support rows: `126,861`.
- Near-resistance rows: `128,085`.
- Detector full split remains `77,231` active bars / `569,451` total bars,
  `13.56%` duty cycle, `2,996` episodes.
- Detector sweep full baseline remains `77,231` active bars / `569,451` total
  bars, `13.56%` duty cycle, `2,996` episodes.
- Every compatibility run printed `strategy=empty trades=0`.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

SR boundary-inspection milestone:

- Added CLI flag `-sr-boundary-inspect`.
- Added compact non-trading candidate comparison outputs:
  - `sr_boundary_candidate_comparison.csv`
  - `sr_boundary_candidate_comparison.json`
- Grouping:
  - split
  - side
  - horizon bars
  - strength bucket
  - distance bucket
- Metrics include counts and rates for close breaks, rejections, reclaimed
  breaks, and favorable-vs-adverse cohorts for all, rejected, and reclaimed
  events.
- This milestone did not add entries, exits, scoring, sizing, or strategy
  replacement.
- Strategy remains `lab.EmptyStrategy`.
- Trades remain `0`.

Latest SR boundary-inspection run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection \
  -sr-boundary-inspect
```

Observed SR boundary-inspection facts:

- Boundary event rows inspected in memory: `281,080`.
- Candidate comparison rows: `192`.
- Candidate comparison CSV lines including header: `193`.
- Result paths:
  - `results/sr-boundary-inspection/sr_boundary_candidate_comparison.csv`
  - `results/sr-boundary-inspection/sr_boundary_candidate_comparison.json`
- Compact inspect mode did not write `sr_boundary_events.*` or
  `sr_boundary_quality.*`; only the standard empty-strategy summary/trades
  files were also written.
- `strategy=empty trades=0`.

Factual comparison outcome:

- Boundary rejection has better current evidence than false-break reclaim.
- By side/horizon across all splits:
  - rejection rates ranged from `10.52%` to `36.77%`
  - rejected-cohort favorable-minus-adverse ranged from `14.75bp` to `27.73bp`
  - rejected-cohort favorable-greater-than-adverse was `96.34%` to `98.16%`
- Reclaim-after-break was mostly a longer-horizon conditional outcome:
  - reclaim event rate reached about `20%` at `12` bars
  - reclaim given close break reached about `43%` at `12` bars
  - reclaimed-cohort favorable-minus-adverse was negative at `3`, `6`, and
    `12` bars for both support and resistance in this event definition
- The `10_20bp` distance bucket remains sparse:
  - resistance events across all splits/horizons: `208`
  - support events across all splits/horizons: `96`
- Do not treat rejection as a ready entry rule yet; it is still an ex-post
  audit label. False-break reclaim needs a later post-reclaim timing audit
  before becoming a first entry template.

Latest test command:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```

Result: passed.

Next implementation:

- Start with an entry-readiness review before adding trades.
- Use `memory/NEXT_CODEX_BRIEF.md` as the next prompt.
- Review code correctness, lookahead safety, coverage gaps, docs, and memory
  before the first entry implementation.
- After the review is clean, consider one more compact timing audit that asks
  whether rejection can be identified on the decision candle without using
  future bars.

Entry-readiness handoff:

- Committed SR boundary inspection mode:
  - commit: `1e26695 Add SR boundary inspection mode`
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with a paste-ready prompt for:
  - whole-project review before entries
  - test coverage gap investigation
  - docs/memory readiness
  - next non-trading boundary-rejection timing audit
- Removed duplicate root `CODEX_BRIEF.md`; `memory/NEXT_CODEX_BRIEF.md` is the
  canonical next-session prompt.
- Coverage check:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./internal/lab`
  - result: `84.5%` statement coverage
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./...`
  - result: failed in `cmd/rangelab` with the known snap-confine issue, while
    `/usr/local/go/bin/go test ./cmd/rangelab` passed
