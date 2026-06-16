# Next Codex Brief: Assess Futures Data Impact Before Entries

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this impact review:
  - docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md
  - docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md
  - docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current issue:
- The user clarified the actual trading target is Binance futures, not spot.
- Prior generated audits and review verdicts were based on spot CSV data, especially:
  - ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv
- The previously approved first minimal hold-inside midline touch prototype is paused.
- No entries, exits, scoring, sizing, or strategy replacement should be built in this task.

Known local data snapshot from 2026-06-16:
- Spot full-history file visible:
  - ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv
  - 573,697 CSV lines including header
  - 2021-01-01T00:00:00Z through 2026-06-15T23:59:59Z
- Futures files visible, but only short recent slices:
  - ../binance-bot/data/btcusdt_futures_um_5m_2026-06-13_2026-06-15.csv
  - ../binance-bot/data/btcusdt_futures_um_5m_2026-06-14_2026-06-15.csv
  - ../binance-bot/data/btcusdt_futures_um_5m_2026-06-15.csv
- First verify current files again; the user may have added a longer futures CSV after this brief was written.

Goal:
- Produce a futures data impact review that answers:
  - what prior spot-based conclusions are now invalid, suspended, or probably reusable as code only
  - whether a full-history Binance USDT-M futures 5m dataset is available
  - whether the key hold-inside midline reaction evidence survives on futures data
  - what must happen before entries are allowed again

Key changes:
- Create docs/FUTURES_DATA_IMPACT_REVIEW.md.
- Inventory available BTCUSDT 5m CSVs by path, market type inferred from filename, row count, first candle, last candle, and schema.
- If a full-history futures CSV exists, rerun at minimum:
  - -detector-context-refinement-audit on futures data
  - -hold-inside-midline-transition-audit on futures data
  - -hold-inside-midline-reaction-audit on futures data
- Put futures rerun outputs under distinct result directories, for example:
  - results/futures-detector-context-refinement-audit/
  - results/futures-hold-inside-midline-transition-audit/
  - results/futures-hold-inside-midline-reaction-audit/
- Compare the futures reaction funnel/stability against the prior spot review gate, especially:
  - hold_3_inside + mid_touch
  - event close-position bucket mid_50
  - weakest split candidate count
  - event rate, missing events, missing future
  - h6 close-back, mid-rejection-before-boundary, boundary-before-rejection
  - quick invalidation and trend leakage
- If only short futures slices exist, do not rerun promotion audits as if they are comparable. Write the review verdict as data gap first, and make the next brief about obtaining/building full-history futures BTCUSDT 5m data.

Review gate:
- The old spot-based approval does not carry forward automatically.
- A first minimal entry prototype may return to the next brief only if futures data has comparable coverage and the futures rerun passes the same evidence gate previously used for spot:
  - reaction-eligible candidates in every period split
  - all-bucket stability has at least 100 candidates in the weakest split
  - any prototype cohort has at least 100 candidates in every period split
  - event occurrence is split-stable enough that missing events do not dominate
  - post-event behavior points to one coherent prototype shape
  - quick invalidation and trend leakage do not contradict that behavior
- If the futures data is too short, the verdict is not no-promotion; it is insufficient data.
- If full-history futures data exists and the rerun fails, write a no-promotion or pivot verdict for this detector family.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...
- If a full-history futures CSV exists, run the required futures audits with that CSV, for example:
  - env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go run ./cmd/rangelab \
      -csv PATH_TO_FULL_HISTORY_FUTURES_5M_CSV \
      -hold-inside-midline-reaction-audit \
      -out-dir results/futures-hold-inside-midline-reaction-audit
- wc -l results/futures-*/*.csv when futures result dirs exist
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Closeout:
- Update memory/PROGRESS.md with data paths, row counts, coverage, commands, result paths, and factual outcome.
- Update memory/DECISIONS.md if the active dataset path/coverage becomes durable, if spot evidence is permanently retired, or if the futures rerun creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md with either:
  - a data acquisition/build brief if full futures history is missing
  - a futures reaction review brief if audits ran but still need review
  - the first minimal entry prototype brief only if futures evidence clears the gate
- Commit completed repo changes after verification unless explicitly told not to.
```
