# Next Codex Brief: Review Hold-Inside Midline Reaction Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this review:
  - docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md
  - docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md
  - docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The hold-inside midline reaction audit has already been implemented and generated.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Result directory:
- results/hold-inside-midline-reaction-audit/

Generated outputs:
- hold_inside_midline_reaction_candidates.csv/json
- hold_inside_midline_reaction_funnel_summary.csv/json
- hold_inside_midline_reaction_summary.csv/json
- hold_inside_midline_reaction_stability.csv/json

Audit facts:
- profile: p30_c12_bollinger_on_adx_off
- rules: hold_3_inside, hold_6_inside, diagnostic hold_3_inside_mid_50
- event types: mid_touch and mid_close_across
- max midline event delay: 12 bars after hold decision
- horizons: 1, 3, 6, 12
- quick invalidation: 3 bars
- candidate rows: 9,080
- funnel rows: 24
- reaction summary rows: 1,296
- stability rows: 352
- CSV lines including headers:
  - candidates: 9,081
  - funnel summary: 25
  - reaction summary: 1,297
  - stability: 353
- last run loaded 569,451 candles through 2026-06-01T23:59:59Z and printed strategy=empty trades=0.

Review task:
- Create docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md.
- Review the funnel first:
  - source hold candidates
  - event occurrence rate
  - missing events
  - missing future
  - average event delay
  - rule/event-type/split stability
- Then review reaction summaries and stability:
  - start with all-bucket rows by rule, event type, horizon, and split
  - inspect event_mid_side and event close-position cohorts only if weakest split counts are adequate
  - compare range persistence, quick invalidation, trend leakage, high/low touch, opposite-half touch, close-back-across-mid, mid rejection before boundary touch, and boundary touch before mid rejection
- Require split-stable and cohort-stable evidence before recommending any entry prototype.
- If evidence is stable enough, make the next brief a first minimal non-live entry prototype task.
- If evidence is not stable enough, write a clear no-promotion verdict and make the next brief stop mining this detector family or pivot to a materially different non-trading hypothesis.

Closeout:
- Update README.md docs order if the review doc is added.
- Update memory/PROGRESS.md with review doc path, commands, row counts used, and factual verdict.
- Update memory/DECISIONS.md only for a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md with the next task after the review.
- Commit completed repo changes after closeout unless explicitly told not to.

Verification:
- wc -l results/hold-inside-midline-reaction-audit/*.csv
- env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
```
