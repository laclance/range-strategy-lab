# Futures Range-First Strategy Construction Protocol

Date: 2026-06-26

## Verdict

Stop state:
`range_first_strategy_construction_protocol_ready_for_v1_spec`.

This is a documentation-only protocol milestone. It changes the project posture
from a no-automatic-implementation stop to a user-approved path for designing
fresh range-derived strategies from scratch through offline backtesting and
optimization.

It does not add a strategy, optimizer, replay, walk-forward run, CLI flag,
generated result directory, paper/testnet/live path, exchange API use,
credentials, deploy files, data downloads, broad symbol mining, martingale,
averaging down, or two-exchange logic.

## Intent

The lab should not stop research merely because the latest reviewed premises
failed. The correct boundary is narrower: do not keep rescuing failed premises
by retuning or renaming them.

The next research arc should use this repo as an offline strategy construction
factory:

1. define a range-derived strategy grammar;
2. run a fixed baseline backtest;
3. run a bounded optimization only after the baseline justifies it;
4. freeze and replay a selected configuration;
5. run walk-forward robustness;
6. review before any candidate strategy package.

This mirrors the useful research loop from `binance-bot` while keeping this lab
independent, offline, and free of copied strategy, scoring, live execution,
deploy, credential, or order-management code.

## Scope

The approved scope is range-first and BTCUSDT-first.

| Field | Contract |
| --- | --- |
| Market | Binance USDT-M futures |
| First symbol | `BTCUSDT` |
| Parent source | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Loaded rows | `573,984` |
| Coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Source status | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

The existing local BTC/ETH/SOL futures universe can remain review context, but
it is not the default first construction universe for this protocol. Any
multi-symbol transfer study, ETH/SOL authority stream, broad symbol mining, or
source expansion requires a later approved brief.

## Reusable Infrastructure

The next arc may reuse lab-owned infrastructure that is not itself a rejected
strategy premise:

- futures source guard and manifest discipline;
- closed UTC resampling from the accepted `5m` parent source;
- closed-candle decision semantics;
- next-bar-open entries;
- one open position max;
- engine costs, risk sizing, and notional cap;
- stop-first same-bar ambiguity handling;
- split metrics and side metrics;
- CSV/JSON artifact patterns under `results/`;
- detector, range episode, and event-label helpers as feature extraction only.

Default `cmd/rangelab` behavior should remain `lab.EmptyStrategy` until a later
implementation brief explicitly adds a new offline flag.

## Exclusion Boundary

The following reviewed families are exclusion evidence in their current forms:

- structured-compression ETH/SOL authority stream after fragile walk-forward;
- breakout-retest/acceptance baseline after negative P&L and failed transfer;
- clean-breakout baseline after costs;
- hold-inside/midline futures prototype after failed P&L;
- impulse absorption after continuation-dominant audit results;
- BTCUSDT higher-timeframe nested range rotation after the no-baseline event
  count failure;
- legacy spot-only SR timing and compression reviews as promotion evidence.

The next strategy construction work must not retune these exact families,
relax their gates, rename their rules, or rebuild them as entries without a
materially new structure premise and a fresh spec.

## Strategy Construction Ladder

Each future implementation step should advance only one rung:

| Rung | Purpose | Stop Condition |
| --- | --- | --- |
| v1 grammar spec | Define feature grammar, parameter bounds, outputs, gates, and stop states | Stop before code |
| baseline backtest | Prove a fixed, non-optimized strategy template can survive costs and splits | Stop before optimization |
| bounded optimization | Search only declared parameters and rank by predeclared criteria | Stop before replay |
| fixed replay | Freeze one selected config and reproduce it as normal outputs | Stop before walk-forward |
| walk-forward robustness | Test train-selection stability and forward results | Stop before packaging |
| package review | Decide whether the frozen stream is candidate-ready | Stop before any live-adjacent work |

No rung authorizes the next one implicitly after a failed result.

## Next Spec Requirements

The next brief should create:
`docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md`.

That spec should remain documentation-only and decide, before code:

- the first range-derived strategy grammar;
- allowed timeframes and feature primitives;
- entry, stop, target, max-hold, and invalid-geometry rules;
- parameter grid or search bounds;
- training, OOS, recent, and full-period split gates;
- ranking score and tie-breaks;
- artifact names and common-output behavior;
- CLI flag name for the later implementation brief;
- source/resample validation requirements;
- failure and promotion stop states.

Minimum stop states for the v1 spec:

- `range_first_strategy_v1_spec_ready_for_optimizer_implementation`
- `range_first_strategy_v1_spec_needs_user_premise_or_scope_input`
- `range_first_strategy_v1_spec_rejected_closed_family_reslice`

The v1 spec should prefer a broad range-derived grammar over a single rescued
micro-premise. It should be decision-complete enough that the following session
can implement the first bounded offline optimizer/backtester without inventing
research policy mid-flight.

## Review Rule

This protocol authorizes strategy construction research, not strategy
promotion. Any future candidate must earn promotion through documented source
validation, baseline P&L, bounded optimization, fixed replay, walk-forward
robustness, and review gates.
