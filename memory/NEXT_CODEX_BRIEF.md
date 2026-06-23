# Next Codex Brief: Futures Higher-Timeframe Range Source Spec

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_SCOPE_PIVOT_REVIEW.md
  - docs/FUTURES_SCOPE_PIVOT_REVIEW_SPEC.md
  - docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md
  - docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active project objective is still range-strategy research, not generic futures strategy discovery.
- The active market target is Binance USDT-M futures. Spot evidence is historical context only unless explicitly rerun and reviewed on futures data.
- Current authoritative lab source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - market type: Binance USDT-M futures BTCUSDT 5m
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - accepted manifest facts: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- BTCUSDT 5m range work is paused:
  - futures midline touch prototype failed full sample, every period split, and both sides;
  - futures impulse absorption failed because continuation-first and quick continuation dominated midpoint reclaim-first;
  - SR timing, compression, range durability, detector/context, hold-inside directional edge, and hold-inside midline/reaction are closed, diagnostic-only, or infrastructure.
- The futures scope pivot review stop state was:
  range_scope_pivot_ready_for_higher_timeframe_source_spec.
- The selected next lane is BTCUSDT higher-timeframe futures range source/premise specification.
- BTC/ETH range expansion is deferred until the BTCUSDT higher-timeframe source/premise review is complete or the user explicitly changes scope.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/prototype flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Create docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md as a docs/source-spec-only milestone.
- Define the accepted parent source and closed UTC resampling contract for BTCUSDT higher-timeframe range research.
- Keep the next step source/premise review only: no audit implementation, no generated results, no entry prototype, and no strategy change.

Implementation boundaries:
- This task is docs/memory-only.
- Do not add code, CLI flags, audits, generated result directories, entries, exits, scoring, sizing, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, two-exchange logic, data downloads, or sibling repo mutation.
- Do not reopen BTCUSDT 5m hold-inside/midline, SR timing, compression, or impulse surfaces by reslicing.
- Do not use spot outputs as authority for futures promotion.
- Do not expand to BTC/ETH or broader symbols in this spec; record those as deferred unless the user explicitly changes scope.

Required work:
- Create docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md.
- Specify the parent source:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  with -source-product binance-usdm-futures facts carried from the accepted manifest.
- Specify candidate higher-timeframe bars:
  - 15m
  - 1h
  - 4h
- Define closed UTC resampling rules:
  - parent bars are 5m open-time candles;
  - higher-timeframe open time is the UTC bucket start;
  - a higher-timeframe bar is complete only when all expected child 5m opens exist;
  - open is first child open, high is max child high, low is min child low, close is last child close, volume is summed child volume;
  - no partial final bar, forward-filled bar, synthetic bar, duplicate timestamp, or gap is allowed.
- Define source acceptance gates:
  - parent manifest accepted, futures product, BTCUSDT, 5m, comparison_only=false;
  - generated 15m/1h/4h coverage and row counts are documented before any audit;
  - no missing child bars inside accepted higher-timeframe buckets;
  - timestamps remain UTC and closed-candle finality is explicit.
- Define premise acceptance gates:
  - the future audit premise must be range-only;
  - it must be materially different from failed 5m SR timing, compression, hold-inside/midline, and impulse surfaces;
  - it must define a closed-candle candidate event and falsification rule before any entry prototype;
  - it must remain BTCUSDT-only unless the user explicitly changes scope.
- Add the new spec doc to the README docs index.
- Update memory/PROGRESS.md with docs created, exact verification commands, source facts, and stop state.
- Update memory/DECISIONS.md only if the spec adds a durable source or no-retry rule.
- Refresh memory/NEXT_CODEX_BRIEF.md from the chosen stop state.

Suggested stop states:
- higher_tf_range_source_spec_ready_for_non_trading_audit_brief
- higher_tf_range_source_spec_blocked_by_source_or_resampling_gap
- higher_tf_range_source_spec_rejected_as_5m_reslice
- higher_tf_range_source_spec_no_viable_range_premise
- higher_tf_range_source_spec_review_only_no_strategy_change

Suggested verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit completed repo changes after verification unless explicitly told not to.
```
