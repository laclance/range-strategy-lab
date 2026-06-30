# Progress Note: Backtest-First Candidate Packet

Date: 2026-06-30

This focused progress note records the docs-only candidate-packet milestone while
keeping the large always-read `memory/PROGRESS.md` stable for this connector-made
change.

## Outcome

Added docs-only candidate packet:

```text
docs/BACKTEST_FIRST_CANDIDATE_PACKET.md
```

Stop state:

```text
backtest_first_candidate_packet_selected_value_reversion_baseline_needs_implementation_approval
```

Selected next fixed baseline candidate:

```text
btc_5m_rolling_value_area_reversion_v1
```

## Candidate Set

The packet considered three materially different BTCUSDT Binance USDT-M futures
range-entry hypotheses:

1. `btc_5m_rolling_value_area_reversion_v1` — selected first baseline.
2. `btc_15m_previous_day_range_reversion_v1` — parked candidate.
3. `btc_15m_range_edge_exhaustion_fade_v1` — parked candidate.

Rejected lookalikes include post-compression breakout variants, clean breakout
continuation, breakout-retest acceptance, router-gated boundary reclaim,
occupancy rotation, midline/hold-inside variants, and derivatives-veto
interaction.

## Source Contract

Selected candidate uses the current accepted source:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
```

Source facts carried forward:

- Product: Binance USDT-M futures.
- Symbol/interval: BTCUSDT `5m`.
- Loaded candles: `573,984`.
- Coverage: `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`.
- `gap_count=0`.
- `duplicate_count=0`.
- `zero_volume_count=66`.
- `comparison_only=false`.
- `validation_status=accepted`.

## Next Gate

Next allowed task, only after explicit user approval:

```text
Implement exactly one fixed offline baseline backtest for btc_5m_rolling_value_area_reversion_v1.
```

No Go code, CLI flag, generated output, source download, optimizer grid, replay,
walk-forward, derivatives-veto interaction, paper/testnet/live path, exchange
API work, credentials, deploy file, martingale, averaging down, two-exchange
logic, or promotion was authorized by this docs-only packet.
