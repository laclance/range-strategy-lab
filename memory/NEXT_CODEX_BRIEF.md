# Next Codex Brief: Higher-Timeframe Range Premise Required

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md
  - docs/FUTURES_SCOPE_PIVOT_REVIEW.md
  - docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md
  - docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active project objective is range-strategy research, not generic futures strategy discovery.
- The active market target is Binance USDT-M futures. Spot evidence is historical context only unless explicitly rerun and reviewed on futures data.
- BTCUSDT 5m range mining is paused.
- Current authoritative parent source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - required product: -source-product binance-usdm-futures
  - market type: Binance USDT-M futures BTCUSDT 5m
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - accepted manifest facts: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- Higher-timeframe source spec status:
  higher_tf_range_source_spec_no_viable_range_premise.
- Candidate source shapes are closed UTC resamples from the accepted 5m parent:
  - 15m: 3 child bars, expected complete bars 191,328, first open 2021-01-01T00:00:00Z, last open 2026-06-16T23:45:00Z
  - 1h: 12 child bars, expected complete bars 47,832, first open 2021-01-01T00:00:00Z, last open 2026-06-16T23:00:00Z
  - 4h: 48 child bars, expected complete bars 11,958, first open 2021-01-01T00:00:00Z, last open 2026-06-16T20:00:00Z
- Closed UTC resampling contract:
  - higher-timeframe open time is the UTC bucket start;
  - each accepted bucket must contain every expected child 5m open exactly once;
  - open is first child open, high is max child high, low is min child low, close is last child close, volume is summed child volume;
  - no partial final bar, forward-filled bar, synthetic bar, duplicate timestamp, child gap, local-time bucket, spot/comparison parent, or future-looking close is allowed.
- Failed or closed families must not be reopened by renaming them onto higher-timeframe bars:
  - BTCUSDT 5m SR timing
  - compression breakout
  - hold-inside directional edge
  - hold-inside midline transition/reaction/prototype
  - abnormal OHLCV impulse absorption
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/prototype flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Do not start a new audit automatically.
- Ask the user for a materially different BTCUSDT higher-timeframe range premise.
- Once the premise is explicit, write a scoped next brief for a non-trading audit only.

Implementation boundaries:
- Keep this task planning/review-only unless the user supplies a clear new premise in the same session.
- Do not add code, CLI flags, audits, generated result directories, entries, exits, scoring, sizing, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, two-exchange logic, data downloads, spot comparison, symbol expansion, or sibling repo mutation.
- Do not reopen BTCUSDT 5m SR timing, compression, hold-inside/midline, or impulse surfaces by reslicing them onto 15m/1h/4h.
- Do not use spot outputs as authority for futures promotion.
- Do not expand to BTC/ETH or broader symbols unless the user explicitly changes scope.

Required next interaction:
- Ask for the new higher-timeframe range premise in concrete terms:
  - Which interval or interval set should be used: 15m, 1h, 4h, or a fixed comparison among them?
  - What range behavior should be audited?
  - Why should it differ from the closed 5m families?
  - Which closed-candle observable defines the candidate event?
  - What outcome would falsify it before any entry prototype?
- If the user provides a premise, create a new memory/NEXT_CODEX_BRIEF.md for a non-trading audit only. Include exact source, resampling contract, outputs, review gate, stop states, and verification commands.
- If the user does not provide a premise, stop with:
  higher_tf_range_source_spec_no_viable_range_premise.

Suggested verification after any brief refresh:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update memory/PROGRESS.md with the chosen premise or the lack of one, exact commands, and stop state.
- Update memory/DECISIONS.md only if a durable no-retry, source, or scope rule is added.
- Commit completed repo changes after verification unless explicitly told not to.
```
