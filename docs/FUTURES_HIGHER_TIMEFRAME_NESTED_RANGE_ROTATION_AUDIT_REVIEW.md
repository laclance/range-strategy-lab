# Futures Higher-Timeframe Nested Range Rotation Audit Review

Date: 2026-06-26

## Scope

This review covers the bounded non-trading audit behind:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-higher-tf-nested-range-rotation-audit -out-dir results/futures-higher-tf-nested-range-rotation-audit
```

The audit uses the accepted BTCUSDT Binance USDT-M futures `5m` source,
resamples it into closed UTC `4h` parent bars and `1h` child bars, then checks
whether mature child ranges nested inside frozen mature parent ranges produce
internal rotation events toward the parent midpoint and far quartile.

No entries, exits, P&L, optimizer, replay, walk-forward, source expansion,
strategy package, paper/testnet/live path, exchange API, deploy file,
martingale, averaging down, or two-exchange logic was added.

## Source And Resample Validation

The source guard and closed UTC resample validation passed.

| Source | Rows | First open | Last open | Gaps | Duplicates | Zero volume | Status |
| --- | ---: | --- | --- | ---: | ---: | ---: | --- |
| `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | 573,984 | 2021-01-01T00:00:00Z | 2026-06-16T23:55:00Z | 0 | 0 | 66 | accepted |

| Role | Timeframe | Rows | First open | Last open | Complete | Status |
| --- | --- | ---: | --- | --- | --- | --- |
| parent | 4h | 11,958 | 2021-01-01T00:00:00Z | 2026-06-16T20:00:00Z | true | accepted |
| child | 1h | 47,832 | 2021-01-01T00:00:00Z | 2026-06-16T23:00:00Z | true | accepted |

Both resample rows had `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, `complete=true`, and
`validation_status=accepted`.

## Audit Rules

The frozen audit config used:

- detector profile: `p30_c12_bollinger_on_adx_off`;
- parent timeframe: closed UTC `4h`;
- child timeframe: closed UTC `1h`;
- parent maturity: `12` parent bars;
- child maturity: `12` child bars;
- maximum child width: `40%` of parent width;
- outcome horizon: `24` closed child bars;
- quick invalidation horizon: `6` closed child bars.

Parent ranges were frozen at the first mature close using the raw active run's
high, low, midpoint, and quartiles. Parent eligibility ended on the first later
closed `4h` candle that closed outside the frozen range.

Child ranges were eligible only when they had positive width, were fully inside
the latest still-valid mature parent, were no wider than `40%` of parent width,
and had midpoint geometry consistent with a rotation from the parent half
toward the parent interior.

Only the first closed `1h` break after child maturity could become the child
event. Invalid first breaks and later duplicate breaks were recorded as skipped
event rows.

## Artifacts

The run wrote:

- result directory:
  `results/futures-higher-tf-nested-range-rotation-audit/`
- common zero-trade outputs:
  `source_manifest.json`, `summary.csv/json`, `trades.json`
- specific CSV line counts including headers:
  - child ranges: `283`
  - coverage: `3`
  - events: `20`
  - parent ranges: `69`
  - sources: `2`
  - summary: `13`
  - common summary: `13`

Common `trades.json` contains zero trades, and common `summary.csv` remains
the zero-trade compatibility summary.

## Results

The audit found `68` parent ranges and `282` child ranges.

Child eligibility collapsed before event review:

| Child status | Count |
| --- | ---: |
| eligible | 11 |
| no valid parent | 262 |
| child width above 40% parent | 8 |
| child not inside parent | 1 |

Event review produced `19` event rows, but only `3` valid events:

| Event status | Count |
| --- | ---: |
| valid | 3 |
| duplicate child event | 9 |
| event beyond parent midpoint | 6 |
| event wrong rotation side | 1 |

All valid events were upside rotations. No downside event passed the nested
geometry and first-break rules.

| Split | Valid events | Outcome read |
| --- | ---: | --- |
| 2021_2022_stress | 1 | adverse parent invalidation, also adverse child invalidation |
| 2023_2024_oos | 1 | favorable midpoint and far quartile |
| 2025_2026_recent | 1 | no resolution |
| full_2021_2026 | 3 | 1 favorable, 1 adverse parent/child, 1 no resolution |

The full sample had:

- favorable midpoint rate: `33.33%`;
- favorable far-quartile rate: `33.33%`;
- adverse child invalidation rate: `33.33%`;
- adverse parent invalidation rate: `33.33%`;
- no-resolution rate: `33.33%`;
- quick invalidation rate: `0%`;
- missing future rate: `0%`.

The result failed the declared review gates because full-period events were far
below the `100` event gate, every split was far below the `25` event split
gate, downside had zero valid events, and the favorable/outcome tests were too
sparse and unstable to justify a baseline.

## Verdict

Stop state:

```text
higher_tf_nested_range_rotation_audit_failed_no_baseline
```

This is not a source-gap stop. It is exclusion evidence for this exact
BTCUSDT `4h` parent / `1h` child nested range-rotation premise with the frozen
`40%` width, `24` bar outcome, and `6` bar quick-invalidation review contract.

The audit does not authorize a baseline backtest, optimizer, replay,
walk-forward, retuning around the failed result, source expansion, symbol
expansion, strategy package, paper/testnet/live wiring, exchange API use,
deployment, martingale, averaging down, or two-exchange work.

## Verification

Passed:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-higher-tf-nested-range-rotation-audit -out-dir results/futures-higher-tf-nested-range-rotation-audit
wc -l results/futures-higher-tf-nested-range-rotation-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
