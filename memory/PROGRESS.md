# Progress

This file is the always-read project snapshot. Keep it compact: current state,
latest verification, result paths, and a milestone index. Detailed evidence
belongs in focused `docs/` reviews, generated artifacts under `results/`, and
git history.

## Current State

- Scope remains offline BTCUSDT 5m range-strategy research only.
- Active market target is Binance USDT-M futures, not spot. Prior generated
  audits/reviews were run on spot CSVs and are no longer authoritative for
  promotion or entry decisions until rerun/reviewed on futures data.
- Strategy remains `lab.EmptyStrategy`; trades remain `0`.
- Closed-candle decision semantics remain required.
- No live code, API keys, deploy scripts, grid, martingale, averaging down, or
  two-exchange execution is allowed.
- `memory/NEXT_CODEX_BRIEF.md` is the only canonical next-session prompt.
- Current next task: assess the spot-to-futures data-source impact before any
  entries. The previously planned hold-inside midline touch prototype is paused.

## 2026-06-16

Futures data-source correction:

- User clarified the real trading target is Binance futures rather than spot.
- Previous audit/review outputs used the spot path
  `../binance-bot/data/btcusdt_spot_5m_2021_2026.csv`; treat all promotion
  implications from those results as suspended until futures revalidation.
- Local sibling data currently visible:
  - `../binance-bot/data/btcusdt_spot_5m_2021_2026.csv`: `573,697` CSV lines
    including header, spanning `2021-01-01T00:00:00Z` through
    `2026-06-15T23:59:59Z`
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-13_2026-06-15.csv`:
    `865` CSV lines including header
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-14_2026-06-15.csv`:
    `577` CSV lines including header
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-15.csv`: `289` CSV
    lines including header
- Updated memory and docs so the next brief is an impact assessment, not entry
  implementation. If full-history futures data is unavailable, the next verdict
  should be "data gap first," not a prototype.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Memory context-budget wording adjustment:

- Updated `AGENTS.md`, `memory/README.md`, and `memory/DECISIONS.md` so the
  size guidance for all always-read memory files is a soft `300-350` line
  judgment band, not a hard threshold. Compact or split memory when it feels
  bulky or repetitive, not merely because one file crosses `300` lines.

## 2026-06-15

Hold-inside midline reaction review milestone:

- Review doc: `docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md`.
- Inputs: `results/hold-inside-midline-reaction-audit/`.
- Original spot-data verdict: one first minimal offline entry prototype was
  justified, but only for `hold_3_inside` + first `mid_touch` within `12` bars
  + event close-position bucket `mid_50`.
- Current status after futures correction: suspended pending futures rerun and
  impact review.
- Not promoted: live use, strategy promotion, broad detector-family entries,
  `hold_6_inside`, `mid_close_across`, side-specific cohorts, and
  `hold_3_inside_mid_50`.
- Key evidence:
  - funnel pass: `hold_3_inside` + `mid_touch` has weakest-split event rate
    `52.25%` with `116` event candidates
  - `hold_3_inside` + `mid_touch`, h6 all-bucket: weakest split `116`
    candidates, `55.17%` minimum close-back, `45.69%` minimum mid-rejection
    before boundary, `28.25%` maximum boundary-before-rejection, `18.39%`
    maximum quick invalidation, `22.22%` maximum trend
  - `hold_3_inside` + `mid_touch` + event-close `mid_50`, h6: weakest split
    `104` candidates, `58.65%` minimum close-back, `50.96%` minimum
    mid-rejection before boundary, `22.12%` maximum boundary-before-rejection,
    `11.54%` maximum quick invalidation, `19.23%` maximum trend
- Updated `README.md`, `memory/DECISIONS.md`, and
  `memory/NEXT_CODEX_BRIEF.md` so the next task is the bounded offline
  prototype.
- Verification passed:
  - `wc -l results/hold-inside-midline-reaction-audit/*.csv`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

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

- Review doc: `docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md`.
- Inputs: `results/hold-inside-midline-transition-audit/`; `7,988`
  candidate rows, `720` summary rows, `192` stability rows.
- Verdict: not entry-ready; broad `hold_3_inside`/`hold_6_inside` rows show
  split-stable midline touch/close-across behavior by `12` bars, but current
  midline labels are not entry context or strategy inputs.

Hold-inside midline transition audit milestone:

- Commit: `277ce49 Add hold-inside midline transition audit`.
- Added `-hold-inside-midline-transition-audit`; result dir
  `results/hold-inside-midline-transition-audit/`; `7,988` candidate rows,
  `720` summary rows, `192` stability rows; `strategy=empty trades=0`.

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
