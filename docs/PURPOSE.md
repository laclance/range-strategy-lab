# Purpose

`range-strategy-lab` is a clean research workbench for trying to build a
BTCUSDT range strategy from the ground up.

It exists because the parent `binance-bot` project has a mature trend/scoring
system, portfolio research history, live execution code, deploy scripts, and
presets. Those are useful for that bot, but they can bias a new range-strategy
experiment. This lab keeps the reusable research plumbing and removes the
strategy baggage.

## What This Project Includes

- 5m candle CSV loading for Binance archive-style data.
- A small one-position backtest engine.
- Confirmed-candle decisions with next-bar-open entries.
- Stop-first ambiguity when stop and target are both touched in the same bar.
- Fee and slippage modeling.
- Risk-at-stop sizing with a notional cap.
- Split metrics by trade close time.
- JSON and CSV outputs.
- Unit tests around loading, execution, costs, sizing, split metrics, and
  drawdown.

## What This Project Excludes

- No live trading.
- No API keys or exchange clients.
- No deploy scripts.
- No borrowed scoring engine from the parent bot.
- No copied live presets.
- No portfolio coordinator.
- No existing range strategy hidden in the starter.

## Starting Assumptions

- Build BTCUSDT-only first.
- Use 5m candles first.
- Make decisions only from confirmed closed candles.
- Enter on the next bar open until a later experiment explicitly justifies a
  different fill model.
- Keep one open position max while searching for the first real edge.
- Require gross profitability before trusting cost-sensitive variants.

## Success Criteria

A candidate is worth further work only if it survives all of these:

- Positive gross expectancy over the full available sample.
- Positive net expectancy after realistic fees and slippage.
- Not isolated to a single lucky period.
- Reasonable drawdown for the return.
- Enough trades to be more than noise.
- Simple enough to explain and re-run.

The first milestone is not live trading. The first milestone is a reproducible
offline result that still looks good after you try to disprove it.
