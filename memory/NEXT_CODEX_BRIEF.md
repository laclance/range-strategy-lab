# Next Codex Brief: Futures Range-First Strategy Construction V1 Spec

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_PROTOCOL.md
  - docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with `lab.EmptyStrategy`
  unless an explicit offline audit/backtest flag is passed.
- The latest protocol milestone changed the posture from no automatic
  implementation to user-approved range-first strategy construction from
  scratch.
- Protocol stop state:
  range_first_strategy_construction_protocol_ready_for_v1_spec.
- Scope for the first construction arc is range-first and BTCUSDT-first.
- The first source remains the accepted Binance USDT-M futures BTCUSDT `5m`
  file:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Source facts:
  - loaded candles: 573,984;
  - coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - comparison_only=false, validation_status=accepted.
- Structured compression, breakout-retest/acceptance, clean-breakout,
  hold-inside/midline, impulse absorption, and higher-timeframe nested
  range-rotation are exclusion evidence in their reviewed forms.
- Research is not stopped, but retuning failed families under new names is not
  authorized.

Goal:
- Create a documentation-only v1 strategy construction spec:
  docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md
- Define the concrete first range-derived optimizer/backtester grammar and
  review gates, but do not implement code yet.
- The spec must be decision-complete enough that a later implementation brief
  can add the first bounded offline optimizer/backtester without inventing
  research policy during implementation.

Required spec content:
- State the selected v1 strategy grammar in plain language.
- Define allowed range-derived feature primitives and timeframes.
- Define entry, stop, target, max-hold, and invalid-geometry rules.
- Define parameter grid or search bounds.
- Define fixed baseline behavior versus optimization behavior.
- Define train/OOS/recent/full-period gates and walk-forward expectations.
- Define ranking score and tie-breaks.
- Define expected CLI flag name for the later implementation brief.
- Define result directory and artifact names for the later implementation.
- Define common-output behavior for `source_manifest.json`, `summary.csv/json`,
  and `trades.json`.
- Define source/resample validation requirements.
- Define stop states, including:
  - range_first_strategy_v1_spec_ready_for_optimizer_implementation
  - range_first_strategy_v1_spec_needs_user_premise_or_scope_input
  - range_first_strategy_v1_spec_rejected_closed_family_reslice

Boundaries:
- Documentation/spec only. Do not implement a strategy, optimizer, replay,
  walk-forward run, CLI flag, tests, generated result directory, or artifact
  writer in this task.
- Keep scope BTCUSDT-first and range-first.
- Do not expand to ETH/SOL or broad symbols in this spec unless framed only as
  a future optional review gate after BTCUSDT evidence.
- Do not retune or reopen structured compression, breakout-retest/acceptance,
  clean-breakout, hold-inside/midline, impulse absorption, or nested
  range-rotation as entries.
- Do not import old `binance-bot` strategy/scoring/live code.
- Do not add live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, data downloads, broad symbol mining, martingale, averaging down,
  or two-exchange logic.

Docs and memory:
- Add the new v1 spec doc to README.md.
- Update memory/PROGRESS.md with the spec decision, commands, and outcome.
- Update memory/DECISIONS.md only if the v1 spec creates a new durable rule.
- Replace memory/NEXT_CODEX_BRIEF.md based on the v1 spec stop state.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit completed docs/memory updates and verification evidence after checks
  pass unless explicitly told not to commit.
```
