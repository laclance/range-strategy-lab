# Backtest-First Candidate Packet Review

Date: 2026-06-30

## Verdict

Stop state:

```text
backtest_first_candidate_packet_selected_value_reversion_baseline_needs_implementation_approval
```

The docs-only packet selected exactly one next fixed baseline candidate:

```text
btc_5m_rolling_value_area_reversion_v1
```

Review authority is `docs/BACKTEST_FIRST_CANDIDATE_PACKET.md`.

## Scope

This was a docs-only milestone. It added no Go code, CLI flag, generated output,
source download, optimizer grid, replay, walk-forward, derivatives-veto
interaction, paper/testnet/live path, exchange API work, credentials, deploy
file, martingale, averaging down, two-exchange logic, or promotion.

## Selected Candidate

The selected baseline uses native closed `5m` BTCUSDT Binance USDT-M futures
candles. It tests whether excursions into the outer zones of a rolling `288`-bar
value area revert toward a rolling VWAP anchor before any optimizer, veto,
replay, walk-forward, or broader source work is considered.

## Files Changed

- Added `docs/BACKTEST_FIRST_CANDIDATE_PACKET.md`.
- Updated `README.md` to index the backtest-first lane and candidate packet.
- Updated `memory/NEXT_CODEX_BRIEF.md` to point the next task at the selected
  fixed baseline, only after explicit user implementation approval.

## Next Gate

The next allowed task is explicit user approval to implement exactly one fixed
offline baseline backtest for:

```text
btc_5m_rolling_value_area_reversion_v1
```

If implemented and failed, the baseline must not be rescued by alternate VWAP
windows, outer-zone percentages, target changes, time-stop changes, side
selection, volume filters, derivatives-veto interaction, replay, walk-forward,
or optimizer grids.
