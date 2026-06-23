# Next Codex Brief: Minimal Futures Midline Touch Prototype

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT 5m range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_DATA_IMPACT_REVIEW.md
  - docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active research target is Binance USDT-M futures BTCUSDT 5m, not spot.
- Full-history futures source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - market type: Binance USDT-M futures
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - source manifest: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- The futures data impact review stop state was:
  futures_reaction_gate_passed_needs_minimal_entry_brief.
- The only revalidated prototype surface is:
  - detector profile p30_c12_bollinger_on_adx_off
  - context rule hold_3_inside
  - first closed-candle mid_touch within 12 bars after the hold decision
  - event close-position bucket mid_50
  - any prototype entry no earlier than the next bar open after the event candle
- Still diagnostic only: old spot approval as evidence, hold_6_inside, mid_close_across, side-specific cohorts, and hold_3_inside_mid_50.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Build the first minimal offline entry prototype for the revalidated futures surface, then review its split-stable P&L and stress behavior.

Implementation boundaries:
- Keep the default smoke/backtest path on lab.EmptyStrategy unless an explicit prototype flag is passed.
- Add a separate offline CLI flag for this experiment, for example:
  -hold-inside-midline-touch-prototype
- Use the enforced futures source contract; do not add spot comparison unless a future task explicitly asks for it.
- Use closed-candle decisions only:
  - identify the hold_3_inside decision candle from data known at that close
  - find the first mid_touch within 12 closed bars after the hold decision
  - require the event close-position bucket mid_50
  - enter no earlier than the next candle open after the event candle
- Keep one open position max and existing stop-first ambiguity behavior.
- Do not add live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, or two-exchange logic.

Prototype evidence requirements:
- Report total trades, side splits, period splits, net/gross P&L after costs, profit factor, win rate, max drawdown, average R, and worst split.
- Include enough trade detail to inspect the event candle, entry candle, side, entry, stop, target/exit, fees/slippage, and split.
- Do not promote the strategy from aggregate full-period results alone.
- If either side or any period split is sparse or weak, document that directly instead of broadening the surface.

Suggested verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
    -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv \
    -source-product binance-usdm-futures \
    -hold-inside-midline-touch-prototype \
    -out-dir results/futures-hold-inside-midline-touch-prototype
- wc -l results/futures-hold-inside-midline-touch-prototype/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Stop states:
- prototype_source_gap
- prototype_codegen_or_test_blocked
- prototype_failed_no_promotion
- prototype_review_passed_needs_stricter_oos_review
- prototype_review_only_no_strategy_change

Closeout:
- Create or update a focused review doc for the prototype outcome.
- Update memory/PROGRESS.md with commands, result paths, manifest facts, row counts, and the factual outcome.
- Update memory/DECISIONS.md only if the prototype creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md with the exact next step based on the stop state.
- Commit completed repo changes after verification unless explicitly told not to.
```
