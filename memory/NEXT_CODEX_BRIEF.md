# Next Codex Brief: Build First Hold-Inside Midline Touch Prototype

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this prototype:
  - docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md
  - docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md
  - docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md
  - docs/STRATEGY_WORKFLOW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current verdict:
- The hold-inside midline reaction review approved exactly one first minimal offline entry prototype surface:
  - context rule: hold_3_inside
  - event: first mid_touch within 12 bars after the hold decision
  - event close-position bucket: mid_50
- This is prototype permission only, not strategy promotion or live approval.
- Keep all work offline. Do not add live code, API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.

Prototype goal:
- Build the smallest closed-candle, next-bar-open prototype that tests whether the reviewed mid_touch reaction surface has real P&L.
- Keep the default strategy empty unless the new prototype flag is explicitly passed.
- Report side splits and period splits so the review can tell whether the combined surface hides long/short asymmetry.

Implementation boundaries:
- Add a disabled-by-default CLI flag for the prototype, for example:
  - -hold-inside-midline-touch-prototype
- Use the existing balanced detector profile p30_c12_bollinger_on_adx_off.
- Use only hold_3_inside.
- Search for only the first mid_touch event within 12 bars after the hold decision.
- Require the event close-position bucket to be mid_50.
- The event candle is the signal candle; any entry must occur through the existing next-bar-open engine path.
- Skip at_mid cases unless the implementation has an explicit symmetric rule for them.
- Do not add scoring, sizing changes, live wiring, or multiple competing detector families.

Suggested first strategy shape:
- Direction should be symmetric and based on the event close side of the frozen mid:
  - event close below mid -> long
  - event close above mid -> short
- Use the frozen range boundaries and midline for stop/target design.
- Keep max hold anchored to the review's cleanest horizon, 6 bars, unless tests show the existing engine requires a different default.
- Add enough metadata in Signal.Reason or outputs to identify rule, event type, bucket, and event delay.

Expected outputs:
- Result directory:
  - results/hold-inside-midline-touch-prototype/
- Existing backtest outputs:
  - trades.csv/json if supported by current CLI outputs, otherwise keep existing trades.json plus summary.csv/json
  - summary.csv/json
- Add any compact prototype-specific CSV/JSON only if needed to inspect skipped signals or event counts.

Review after run:
- If net P&L, side splits, stress split, or drawdown fail, write a no-promotion review and stop mining this detector family.
- If it works, require a separate review before any promotion or live discussion.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...
- Run the prototype on BTCUSDT 5m:
  - env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go run ./cmd/rangelab \
      -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
      -hold-inside-midline-touch-prototype \
      -out-dir results/hold-inside-midline-touch-prototype
- wc -l results/hold-inside-midline-touch-prototype/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Closeout:
- Update memory/PROGRESS.md with commands, paths, row counts/trade counts, and factual outcome.
- Update memory/DECISIONS.md only if a durable promotion/no-promotion rule changes.
- Refresh memory/NEXT_CODEX_BRIEF.md with the next review task.
- Commit completed repo changes after verification unless explicitly told not to.
```
