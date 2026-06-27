# Next Codex Brief: Futures Range-First Premise Stop

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to the user-supplied premise or question.
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- The accepted BTCUSDT futures 5m source remains:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
- Source facts:
  - loaded candles: 573,984;
  - coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - comparison_only=false, validation_status=accepted.
- The latest non-trading range-context triage audit passed source/resample
  validation, evaluated closed UTC 15m/1h/4h range quality, UTC session
  behavior, and failure-mode cohorts, and failed with no gated strategy premise.
- Current effective stop state:
  range_context_triage_failed_no_strategy_premise.

Closed or failed premises:
- structured compression;
- breakout-retest/acceptance;
- clean breakout;
- hold-inside/midline;
- impulse absorption;
- higher-timeframe nested range rotation;
- range_occupancy_rotation_v1;
- range-context triage-derived quality/session/failure-mode cohorts.

Goal:
- Review-only stop unless the user explicitly supplies a materially different
  offline range-first premise.
- Do not implement a strategy, optimizer, replay, walk-forward, grid, source
  expansion, symbol expansion, or result-producing run from the closed families
  without user input.
- If the user asks what is available, inventory the current docs and present a
  short choice set of materially different offline range-first premises,
  clearly marking closed families as exclusion evidence.

Boundaries:
- Do not retune or rename failed reviewed families under new labels.
- Do not relax review gates around failed results.
- Do not import old binance-bot strategy/scoring/live code.
- Do not add live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, data downloads, broad symbol mining, martingale, averaging down,
  or two-exchange logic.

Expected shape if a new premise is supplied:
- Start from current docs only.
- Treat all reviewed failures as exclusion evidence.
- Separate reusable infrastructure from rejected strategy premises.
- First produce a documentation-only spec unless the user's prompt explicitly
  supplies an implementation-ready brief with source requirements, artifacts,
  gates, tests, and stop states.

Verification for any documentation-only update:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- If files are changed, update memory/PROGRESS.md with exact commands and
  factual outcomes.
- Update memory/DECISIONS.md only for a durable user-approved rule or
  no-promotion/no-strategy-change result.
- Commit completed docs/memory updates and verification evidence after checks
  pass unless explicitly told not to commit.
```
