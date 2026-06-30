# Backtest-First Research Lane

Date: 2026-06-30

## Purpose

This lane reduces workflow friction after the lab's long sequence of zero-trade
audits and failed premises. It keeps the hard systematic-trading safeguards, but
moves simple, materially different ideas to fixed baseline backtests faster.

The goal is not to be less careful with results. The goal is to stop requiring a
large docs-only evidence chain before the first executable test of a simple
hypothesis.

## When This Lane Applies

Use this lane when the user asks to continue offline strategy research and the
candidate is:

- BTCUSDT Binance USDT-M futures unless a separate source-scope decision changes
  the active market;
- closed-candle and local/offline;
- materially different from closed failed families;
- simple enough to express as one fixed baseline backtest; and
- not a paper, live, optimizer, portfolio, or deployment task.

This lane does not reopen closed failures. It does not authorize rescue retuning
of failed post-compression, router rotation, occupancy rotation, structured
compression, clean breakout, hold-inside/midline, higher-timeframe nested
rotation, BTC-regime/ETH/SOL context, or derivatives-veto paths.

## Non-Negotiable Safeguards

Every backtest-first candidate must keep these safeguards:

- accepted source manifest with market, symbol, interval, coverage, row count,
  gaps, duplicates, zero-volume count, comparison-only status, and validation
  status recorded;
- confirmed closed-candle decisions only;
- next-bar-open entries;
- no feature, filter, sizing, stop, target, or candidate rule using future data;
- one open position max unless a separate engine spec changes it;
- stop-first same-bar ambiguity;
- fees and slippage included;
- split metrics for `2021_2022_stress`, `2023_2024_oos`, `2025_2026_recent`,
  and `full_2021_2026`;
- long and short side results separated when both sides are tested;
- generated outputs under `results/`; and
- no paper/testnet/live path, exchange API, credentials, deploy files,
  martingale, averaging down, or two-exchange execution.

## Short Candidate Packet

A new idea no longer needs several docs-only gates before its first fixed
backtest. A short candidate packet is enough if it states:

1. the hypothesis in one or two sentences;
2. why it is materially different from recently closed failures;
3. the exact source and timeframe;
4. the closed-candle entry rule;
5. the fixed stop, target, time stop, sizing, fee, and slippage assumptions;
6. expected output path and artifacts;
7. pass/fail gates; and
8. explicit no-rescue boundaries.

Keep the packet small. It may live in one focused doc, a section of a review doc,
or a Codex task prompt. Do not create a separate premise spec, audit brief,
auditor review, backtest spec, and implementation gate unless the idea changes
source, timeframe, source family, engine mechanics, or promotion state.

## First Baseline Backtest Contract

For each materially different idea, implement exactly one simple fixed baseline
first:

- one entry template;
- one fixed exit model;
- one sizing model;
- no optimizer grid;
- no adjacent-cell P&L selection;
- no adding filters after seeing results;
- no derivatives-veto interaction until an independent entry baseline passes;
- no replay or walk-forward until the fixed baseline passes.

Failure is useful. A failed fixed backtest should be recorded and closed, not
rescued by parameter drift.

## Fast Failure Rules

Close the candidate and move on when any of these occur:

- full-sample gross P&L is negative;
- `2023_2024_oos` or `2025_2026_recent` gross P&L is clearly negative;
- costs/slippage turn a weak gross edge into a clearly negative net edge;
- the result depends only on `2021_2022_stress`;
- trade count is too small to trust;
- drawdown is unacceptable relative to return;
- the result only works through same-bar optimism, source drift, or leakage;
- the explanation becomes too complex for the first baseline.

Record the exact hypothesis, command, output path, gross result, net/stress
result, split result, failure reason, and what not to retry.

## When To Become Stricter

Only after a fixed baseline is promising should the lab consider heavier gates:

- replay or walk-forward robustness;
- limited predeclared sensitivity checks;
- conservative optimization;
- derivatives no-trade veto interaction;
- broader source or symbol scope;
- paper/shadow readiness review.

A passing baseline still does not authorize paper, testnet, live trading,
exchange API work, credentials, deployment, martingale, averaging down, or
promotion. It only creates evidence for a later review.

## Default Next Research Shape

When there is no selected next strategy, prefer this shape:

1. list `3` to `5` materially different BTCUSDT range-entry hypotheses;
2. reject anything that is only a retuned closed-family rescue;
3. select the simplest first baseline;
4. write the short candidate packet;
5. implement and run the fixed offline backtest only when the user explicitly
   asks for implementation;
6. preserve the result and either close it or promote it to the next review gate.
