# Progress

This file is the always-read project snapshot. Keep it compact: current state,
latest verification, result paths, and a milestone index. Detailed evidence
belongs in focused `docs/` reviews, generated artifacts under `results/`, and
git history.

## Current State

- Scope remains offline BTCUSDT 5m range-strategy research only.
- Strategy remains `lab.EmptyStrategy`; trades remain `0`.
- Closed-candle decision semantics remain required.
- No live code, API keys, deploy scripts, grid, martingale, averaging down, or
  two-exchange execution is allowed.
- `memory/NEXT_CODEX_BRIEF.md` is the only canonical next-session prompt.
- Current next task: review the generated hold-inside midline reaction audit
  outputs and decide whether stable post-event evidence supports a first
  minimal entry prototype or closes this detector family.

## 2026-06-15

Hold-inside midline reaction audit milestone:

- Added CLI flag `-hold-inside-midline-reaction-audit`.
- Result directory: `results/hold-inside-midline-reaction-audit/`.
- Outputs:
  - `hold_inside_midline_reaction_candidates.csv/json`
  - `hold_inside_midline_reaction_funnel_summary.csv/json`
  - `hold_inside_midline_reaction_summary.csv/json`
  - `hold_inside_midline_reaction_stability.csv/json`
- Audit size:
  - profiles: `1`
  - rules: `3`
  - event types: `2`
  - candidate rows: `9,080`
  - funnel rows: `24`
  - summary rows: `1,296`
  - stability rows: `352`
  - CSV lines including header: `9,081` / `25` / `1,297` / `353`
  - horizons: `1`, `3`, `6`, `12`
  - max midline event delay: `12` bars
  - quick invalidation window: `3` bars
- Scope:
  - profile `p30_c12_bollinger_on_adx_off`
  - rules `hold_3_inside`, `hold_6_inside`, and diagnostic
    `hold_3_inside_mid_50`
  - event types `mid_touch` and `mid_close_across`
  - event candle is the reindexed decision candle; labels start at
    `event_index + 1`
  - no entries, exits, scoring, sizing, paper side, favorable/adverse fields,
    strategy replacement, or live wiring
- Last run loaded `569,451` candles through `2026-06-01T23:59:59Z` and printed
  `strategy=empty trades=0`.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv -hold-inside-midline-reaction-audit -out-dir results/hold-inside-midline-reaction-audit`
  - `wc -l results/hold-inside-midline-reaction-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Context budget cleanup milestone:

- Added a durable context-budget rule to `AGENTS.md`, `memory/README.md`, and
  `memory/DECISIONS.md`.
- Compacted `memory/PROGRESS.md` from transcript-style history into this
  rolling snapshot plus milestone index.
- Updated `README.md` and `memory/NEXT_CODEX_BRIEF.md` so future sessions use
  docs as a task-scoped index instead of reading every historical doc by
  default.
- Verification:
  - `wc -l AGENTS.md README.md memory/*.md docs/*.md`: `memory/PROGRESS.md`
    is now `199` lines; counted project context files total `1,999` lines
    versus `3,506` before compaction.
  - `rg -n "docs/\\*|NEXT_CODEX_BRIEF|CODEX_BRIEF" README.md docs memory AGENTS.md`:
    no handoff requires reading `docs/*.md`; remaining `docs/*.md` mention is
    the required-file inventory in `docs/VERIFICATION.md`.
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`:
    passed.
  - `git diff --check`: passed.

Hold-inside midline transition review milestone:

- Added durable review report:
  - `docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md`
- Updated `README.md` docs order to include the hold-inside midline transition
  review.
- Updated `memory/DECISIONS.md` with a durable no-promotion rule: current
  midline transition labels are not entry context or strategy inputs.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` with the next non-trading reindexed
  midline-event audit handoff.
- Review verdict:
  - hold-inside midline transition audit is not entry-ready
  - broad `hold_3_inside`/`hold_6_inside` rows show split-stable midline touch
    and close-across behavior by `12` bars
  - the midline is a promising observation point, but not a trade trigger
  - keep `lab.EmptyStrategy`
  - trades remain `0`
- Inputs reviewed:
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.csv`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.json`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.csv`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.json`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.csv`
  - `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.json`
- Audit size:
  - candidate rows: `7,988`
  - summary rows: `720`
  - stability rows: `192`
  - CSV lines including header: `7,989` / `721` / `193`
- Compact evidence:
  - `hold_3_inside`, h12 all-bucket: `222` minimum split candidates,
    `52.25%` minimum mid touch, `45.05%` minimum close-across, `37.84%`
    minimum cross-before-boundary-break, `25.32%` maximum quick invalidation,
    and `40.09%` maximum trend leakage
  - `hold_6_inside`, h12 all-bucket: `170` minimum split candidates,
    `52.35%` minimum mid touch, `46.47%` minimum close-across, `40.00%`
    minimum cross-before-boundary-break, `21.46%` maximum quick invalidation,
    and `35.88%` maximum trend leakage
  - mid-position h12 rows are cleaner but diagnostic: weakest split counts are
    `94` for `hold_3_inside` and `98` for `hold_6_inside`
- Verification:
  - `wc -l results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.csv results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.csv results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.csv`:
    `7,989` / `721` / `193` lines including headers
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`:
    `memory/NEXT_CODEX_BRIEF.md` remains the only canonical next-session prompt
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`:
    passed
  - `git diff --check`: passed

Hold-inside midline transition audit milestone:

- Commit: `277ce49 Add hold-inside midline transition audit`.
- Added CLI flag `-hold-inside-midline-transition-audit`.
- Result directory: `results/hold-inside-midline-transition-audit/`.
- Outputs:
  - `hold_inside_midline_transition_candidates.csv/json`
  - `hold_inside_midline_transition_summary.csv/json`
  - `hold_inside_midline_transition_stability.csv/json`
- Audit size:
  - profiles: `1`
  - rules: `3`
  - candidate rows: `7,988`
  - summary rows: `720`
  - stability rows: `192`
  - CSV lines including header: `7,989` / `721` / `193`
  - horizons: `1`, `3`, `6`, `12`
  - quick invalidation window: `3` bars
- Scope:
  - profile `p30_c12_bollinger_on_adx_off`
  - primary rules `hold_3_inside` and `hold_6_inside`
  - diagnostic rule `hold_3_inside_mid_50`
  - no entries, exits, scoring, sizing, paper side, favorable/adverse fields,
    strategy replacement, or live wiring
- Last run printed `strategy=empty trades=0`.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv -out-dir results/hold-inside-midline-transition-audit -hold-inside-midline-transition-audit`
  - `wc -l results/hold-inside-midline-transition-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Hold-inside directional edge review milestone:

- Commit: `0601952 Review hold-inside directional edge audit`.
- Review doc: `docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md`.
- Inputs: `results/hold-inside-directional-edge-audit/`.
- Verdict: not entry-ready. `hold_3_inside` and `hold_6_inside` remain useful
  range-survival context but do not show a split-stable directional edge toward
  frozen range high or low.
- Durable rule added to `memory/DECISIONS.md`: do not promote
  `paper_side=toward_high` or `paper_side=toward_low` into entry context.
- Verification passed:
  - `wc -l results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `git diff --check`

Hold-inside directional edge audit milestone:

- Commit: `00b8844 Add hold-inside directional edge audit`.
- Added CLI flag `-hold-inside-directional-edge-audit`.
- Result directory: `results/hold-inside-directional-edge-audit/`.
- Audit size: `15,976` candidate rows, `624` summary rows, `168` stability
  rows; CSV lines including header `15,977` / `625` / `169`.
- Scope stayed non-trading; `strategy=empty trades=0`.
- Verification passed with `go test ./...`, audit run, `wc -l`, handoff
  reference check, and `git diff --check`.

Detector context refinement review milestone:

- Commit: `a4698bf Review detector context refinement audit`.
- Review doc: `docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md`.
- Inputs: `results/detector-context-refinement-audit/`.
- Verdict: delayed `hold_3_inside` and `hold_6_inside` are the leading context
  refinement and materially reduce quick invalidation and trend leakage, but
  are not promoted to entry context because the gain is conditioned survival
  and labels are regime-durability outcomes, not P&L.
- Verification passed with CSV line counts `113,825` / `641` / `161`,
  `go test ./...`, and `git diff --check`.

Detector context refinement audit milestone:

- Commit: `b980095 Add detector context refinement audit`.
- Added CLI flag `-detector-context-refinement-audit`.
- Result directory: `results/detector-context-refinement-audit/`.
- Audit size: `113,824` candidate rows, `640` summary rows, `160` stability
  rows; horizons `1`, `3`, `6`, `12`.
- Scope stayed non-trading; `strategy=empty trades=0`.

Detector durability sweep review milestone:

- Commit: `7bd2852 Review detector durability sweep`.
- Review doc: `docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md`.
- Inputs: `results/detector-durability-sweep/`.
- Verdict: no current `DefaultDetectorSweepProfiles` profile is approved as
  future entry context; `p30_c12_bollinger_on_adx_on` is diagnostic only.
- Verification passed with CSV line counts, `go test ./...`, and
  `git diff --check`.

Detector durability sweep milestone:

- Commit: `d70236b Add detector durability sweep`.
- Added CLI flag `-detector-durability-sweep`.
- Result directory: `results/detector-durability-sweep/`.
- Audit size: `304` broad rows, `9,088` slice rows, `76` stability rows.
- Scope stayed non-trading; `strategy=empty trades=0`.

Range regime durability review milestone:

- Commit: `bc81d9d Review range regime durability`.
- Review doc: `docs/RANGE_REGIME_DURABILITY_REVIEW.md`.
- Inputs: `results/range-regime-durability-audit/`.
- Verdict: current balanced detector regimes are not durable enough as context
  for future entry hypotheses. Refine detector/context before trigger work.
- Verification passed with CSV line counts, `go test ./...`, and
  `git diff --check`.

## Historical Milestone Index

- `7e1494b Add range regime durability audit`: added
  `-range-regime-durability-audit`; result dir
  `results/range-regime-durability-audit/`; `11,984` episode rows and `452`
  summary rows; `strategy=empty trades=0`.
- `b92207a Review compression breakout audit`: added
  `docs/COMPRESSION_BREAKOUT_REVIEW.md`; verdict not entry-ready; result dir
  `results/compression-breakout-audit/`.
- `be1e8d5 Add compression breakout audit mode`: added
  `-compression-breakout-audit`; `5,096` candidate rows and `24` summary rows;
  `strategy=empty trades=0`.
- `0256301 Review SR false-break reclaim timing`: added
  `docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md`; verdict not entry-ready;
  result dir `results/sr-false-break-reclaim-timing-audit/`.
- `1264a2f Add false-break reclaim timing audit`: added
  `-sr-false-break-reclaim-timing-audit`; `17,652` candidate rows and `24`
  summary rows; `strategy=empty trades=0`.
- `4e1a682 Add SR confirmation timing audit and review`: added
  `-sr-confirmation-timing-audit` and `docs/SR_CONFIRMATION_TIMING_REVIEW.md`;
  `9,692` candidate rows and `72` summary rows; verdict not entry-ready.
- `ed82756 Add SR rejection timing audit and review`: added
  `-sr-rejection-timing-audit` and `docs/SR_REJECTION_TIMING_REVIEW.md`; `968`
  candidate rows and `24` summary rows; verdict not entry-ready.
- `f60a8d9 Add entry readiness review gate`: added
  `docs/ENTRY_READINESS_REVIEW.md`; `internal/lab` coverage reached `99.8%`;
  first entries remained blocked pending non-trading timing evidence.
- `04edccc Use memory brief as canonical handoff`: removed duplicate root
  `CODEX_BRIEF.md`; `memory/NEXT_CODEX_BRIEF.md` became canonical.
- `1e26695 Add SR boundary inspection mode`: added `-sr-boundary-inspect`;
  `281,080` boundary events inspected and `192` comparison rows; `strategy=empty
  trades=0`.
- `9833525 Add SR boundary audit mode`: added `-sr-boundary-audit`; `281,080`
  boundary event rows and `192` boundary quality rows.
- `bdf1398 Add mods and next plan;`: pinned helper modules and documented
  adapter-only research helper boundaries in `docs/RESEARCH_HELPERS.md`.
- `6483e8c Add detector sweep audit mode`: added `-detector-sweep`; `19`
  profiles and `76` sweep rows; balanced baseline
  `p30_c12_bollinger_on_adx_off`.
- `e902a9e Initialize range strategy lab`: initialized standalone lab,
  tracked memory, and empty-strategy smoke path.

## Standard Closeout Checks

Use these unless the task has a narrower verifier:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```
