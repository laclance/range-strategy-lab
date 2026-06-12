# Strategy Workflow

Use this workflow to avoid building another parameter soup.

## 1. Build The Detector First

Before adding entries, build a detector and report:

- active bars
- duty cycle
- average active-run length
- median active-run length
- number of range episodes
- split-by-period duty cycle

The detector should answer: "When do we believe BTCUSDT is in a range?"

## 2. Add One Entry Template

Before implementing entries, use helper modules only to create inspectable
features or audit outputs. For example, support/resistance zones can define the
boundary under test, while indicators or candle patterns can be recorded as
diagnostic columns. Do not let helper modules turn the first entry into a bundle
of unrelated filters.

Start with one entry type only. Examples:

- boundary rejection
- false-break reclaim
- compression breakout
- midpoint reclaim after failed continuation
- volatility contraction followed by expansion

Do not optimize many filters at once. Add the minimum rule, then inspect trade
quality.

## 3. Add One Exit Model

Keep the first exit simple:

- stop at a structural invalidation point
- fixed R target or opposite range boundary
- time stop

Use stop-first ambiguity. If a strategy only works when same-bar ambiguity is
optimistic, it is not ready.

## 4. Evaluate Gross Before Net

A candidate that is negative before costs is done. Costs should not be used to
explain away a missing signal.

For candidates that are gross positive, compare:

- gross
- realistic net
- fee/slippage stress

## 5. Split Everything

Always review:

- `2021_2022_stress`
- `2023_2024_oos`
- `2025_2026_recent`
- `full_2021_2026`
- long side
- short side

Do not promote a candidate that only works in one tiny split.

## 6. Keep Promotion Gates Hard

A candidate must be:

- net positive after costs
- profitable or at least non-fragile in both 2023-2024 and 2025-2026
- PF comfortably above 1.2
- not reliant on a tiny number of trades
- drawdown acceptable versus return
- simple enough to explain

## 7. Preserve Failed Results

When a branch fails, write down:

- the exact hypothesis
- the commands
- best gross result
- best net result
- why it failed
- what not to try next

Dead ends are useful if they prevent repeat work.
