# Next Codex Brief: Futures Data Impact Review

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT 5m range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this review:
  - docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md
  - docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md
  - docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active research target is Binance USDT-M futures BTCUSDT 5m, not spot.
- `cmd/rangelab` now enforces source identity before running:
  - default CSV: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - default source product: binance-usdm-futures
  - non-default CSVs require -source-product
  - spot CSVs require -source-product binance-spot plus -allow-spot-comparison
  - every accepted run writes source_manifest.json
- Full-history futures source currently visible:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - market type: Binance USDT-M futures
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - source manifest: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- Prior generated audits/reviews were based on spot data, especially:
  - ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv
- The spot-based first minimal hold-inside midline touch prototype remains paused.
- No entries, exits, scoring, sizing, strategy replacement, paper, testnet, live, exchange keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution should be built in this task.

Goal:
- Create docs/FUTURES_DATA_IMPACT_REVIEW.md and answer whether the prior spot-based conclusions survive on Binance USDT-M futures data.

Required runs:
- Run tests first:
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- Run the three paused audits on the full-history futures CSV:
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
      -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv \
      -source-product binance-usdm-futures \
      -detector-context-refinement-audit \
      -out-dir results/futures-detector-context-refinement-audit
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
      -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv \
      -source-product binance-usdm-futures \
      -hold-inside-midline-transition-audit \
      -out-dir results/futures-hold-inside-midline-transition-audit
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
      -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv \
      -source-product binance-usdm-futures \
      -hold-inside-midline-reaction-audit \
      -out-dir results/futures-hold-inside-midline-reaction-audit
- Each result directory should contain source_manifest.json; verify the manifest confirms Binance USDT-M futures, BTCUSDT, 5m, 573984 rows, no gaps, no duplicates, zero_volume_count=66, and comparison_only=false.

Review focus:
- Inventory the source path, market type, row count, first open, last open, schema, timestamp semantics, gap count, duplicate count, and finality rule from source_manifest.json.
- Compare the futures reaction funnel/stability against the old spot review gate, especially:
  - hold_3_inside + mid_touch
  - event close-position bucket mid_50
  - weakest split candidate count
  - event rate, missing events, missing future
  - h6 close-back, mid-rejection-before-boundary, boundary-before-rejection
  - quick invalidation and trend leakage
- State clearly what prior spot-based conclusions are invalid, suspended, reusable as code only, or revalidated on futures.

Stop states:
- futures_data_gap
- futures_reaction_gate_passed_needs_minimal_entry_brief
- futures_reaction_gate_failed_no_promotion
- futures_review_only_no_strategy_change

Review gate:
- The old spot-based approval does not carry forward automatically.
- A first minimal entry prototype may return to the next brief only if the futures rerun clears the same evidence gate previously used for spot:
  - reaction-eligible candidates in every period split
  - all-bucket stability has at least 100 candidates in the weakest split
  - any prototype cohort has at least 100 candidates in every period split
  - event occurrence is split-stable enough that missing events do not dominate
  - post-event behavior points to one coherent prototype shape
  - quick invalidation and trend leakage do not contradict that behavior

Closeout:
- Update memory/PROGRESS.md with commands, result paths, row counts, manifest facts, and the factual outcome.
- Update memory/DECISIONS.md only if the futures review creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md with the exact next step based on the stop state.
- Run:
  - wc -l results/futures-*/*.csv
  - rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  - git diff --check
- Commit completed repo changes after verification unless explicitly told not to.
```
