# Strategy Workflow

Use this workflow to avoid building another parameter soup while also avoiding
excessive docs-only gatekeeping before the first executable test.

## 1. Keep The Safety Rails

Every strategy result must preserve the project invariants:

- accepted source manifest;
- confirmed closed-candle decisions;
- next-bar-open entries;
- no lookahead in features, filters, sizing, stops, targets, or selection;
- one open position max unless a later engine spec changes that;
- stop-first same-bar ambiguity;
- fees and slippage;
- split metrics by `2021_2022_stress`, `2023_2024_oos`, `2025_2026_recent`, and
  `full_2021_2026`;
- side-separated results when both long and short are tested;
- generated outputs under `results/`.

Do not add paper/testnet/live execution, exchange clients, credentials, deploy
files, martingale, averaging down, or two-exchange logic.

## 2. Prefer Backtest-First For New Simple Ideas

For materially different BTCUSDT range-entry ideas, the default path is now:

1. write a short candidate packet;
2. implement one fixed baseline;
3. run the offline backtest;
4. record the result;
5. kill failures quickly or review a pass.

Do not require a long premise-spec -> zero-trade audit -> strategy-premise spec
-> backtest spec chain for every simple candidate. Use heavier docs-only gates
only when the idea changes source, timeframe, engine mechanics, source family,
portfolio scope, paper/live readiness, or promotion state.

See `docs/BACKTEST_FIRST_RESEARCH_LANE.md` for the detailed lane contract.

## 3. Short Candidate Packet

Before implementation, capture only what is needed to make the first backtest
fixed and reproducible:

- hypothesis;
- why it is materially different from recently closed failures;
- source and timeframe;
- closed-candle entry rule;
- fixed stop, target, time stop, sizing, fee, and slippage assumptions;
- expected output path and artifacts;
- pass/fail gates;
- explicit no-rescue boundaries.

A candidate packet should be compact. It can live in one focused doc, a review
section, or a Codex task prompt.

## 4. Build A Detector First Only When Needed

Detector-only work is useful when the candidate depends on defining a range
regime. In that case, report:

- active bars;
- duty cycle;
- average active-run length;
- median active-run length;
- number of range episodes;
- split-by-period duty cycle.

Do not run detector/audit work by default when the entry rule is already fixed
enough for a baseline backtest.

## 5. Implement One Entry Template

Start with one entry type only. Examples:

- boundary rejection;
- false-break reclaim;
- compression breakout;
- midpoint reclaim after failed continuation;
- volatility contraction followed by expansion.

Do not optimize many filters at once. Do not add filters after seeing the first
result. Add the minimum rule, then inspect trade quality.

## 6. Use One Simple Exit Model

Keep the first exit simple:

- stop at a structural invalidation point or ATR-based invalidation;
- fixed R target or opposite range boundary;
- time stop.

Use stop-first ambiguity. If a strategy only works when same-bar ambiguity is
optimistic, it is not ready.

## 7. Evaluate Gross Before Net

A candidate that is clearly negative before costs is done. Costs should not be
used to explain away a missing signal.

For candidates that are gross positive, compare:

- gross;
- realistic net;
- fee/slippage stress when the candidate is close enough to merit it.

## 8. Split Everything

Always review:

- `2021_2022_stress`;
- `2023_2024_oos`;
- `2025_2026_recent`;
- `full_2021_2026`;
- long side;
- short side.

Do not promote a candidate that only works in one tiny split or only in
`2021_2022_stress`.

## 9. Kill Failures Quickly

Close a fixed baseline when gross edge fails, costs kill the edge, recent/OOS
splits fail, trade count is too small, drawdown is unacceptable, or the rule
needs post-result retuning to survive.

A failed candidate should not be rescued by adjacent-cell P&L selection, exit
retuning, side changes, filter additions, derivatives-veto interaction, replay,
or walk-forward. Move to a materially different idea instead.

## 10. Become Stricter After A Pass

A promising fixed baseline may request heavier review, such as:

- replay or walk-forward robustness;
- limited predeclared sensitivity checks;
- conservative optimization;
- derivatives no-trade veto interaction;
- broader source or symbol scope;
- paper/shadow readiness review.

A passing backtest still does not authorize paper, testnet, live trading,
exchange API work, credentials, deployment, or promotion. It only creates
evidence for the next review gate.

## 11. Keep Promotion Gates Hard

A candidate must be:

- net positive after costs;
- profitable or at least non-fragile in both `2023_2024_oos` and
  `2025_2026_recent`;
- PF comfortably above `1.2` before promotion;
- not reliant on a tiny number of trades;
- drawdown acceptable versus return;
- simple enough to explain.

## 12. Preserve Failed Results

When a branch fails, write down:

- the exact hypothesis;
- the commands;
- source path, market type, coverage, and row count;
- best gross result;
- best net/stress result;
- split result;
- why it failed;
- what not to try next.

Dead ends are useful if they prevent repeat work.
