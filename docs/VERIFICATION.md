# Verification

This project should be runnable after moving the `range-strategy-lab` folder out
of the parent bot.

## Required Files

The folder contains everything needed to compile and run the lab code:

- `go.mod`
- `cmd/rangelab/main.go`
- `internal/lab/*.go`
- `docs/*.md`
- `memory/NEXT_CODEX_BRIEF.md`

The only external input needed for a real run is a candle CSV passed with
`-csv`. For current research, use Binance USDT-M futures BTCUSDT 5m candles,
not legacy spot candles, unless the task is explicitly a spot/futures
comparison.

`cmd/rangelab` defaults to the sibling full-history futures source:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
```

For any non-default CSV, pass `-source-product binance-usdm-futures` or
`-source-product binance-spot`. Spot paths require `-allow-spot-comparison` and
remain comparison-only.

The source guard requires BTCUSDT 5m identity, monotonic closed 5m candles, no
duplicate timestamps, no interval gaps, positive OHLC prices, finite
non-negative volume, and valid high/low containment. Official flat zero-volume
candles are allowed and counted in `source_manifest.json`.

## Local Verification

From inside `range-strategy-lab`:

```bash
/usr/local/go/bin/go test ./...

/usr/local/go/bin/go run ./cmd/rangelab \
  -source-product binance-usdm-futures \
  -out-dir results/smoke
```

Expected smoke result with the starter strategy:

- CSV loads successfully.
- Strategy name is `empty`.
- Trade count is `0`.
- `results/smoke/source_manifest.json`, `summary.csv`, `summary.json`, and
  `trades.json` are written.

Zero trades are correct until a real strategy replaces `lab.EmptyStrategy`.

## After Moving The Folder

If you move this folder elsewhere, either:

- copy a BTCUSDT 5m CSV into the new project, or
- pass an absolute CSV path with `-csv`.

Example:

```bash
/usr/local/go/bin/go run ./cmd/rangelab \
  -csv /absolute/path/to/btcusdt_futures_um_5m_2021_2026.csv \
  -source-product binance-usdm-futures \
  -out-dir results/smoke
```

## What To Check Before A Strategy Result Is Trusted

- Tests pass.
- Smoke run loads the expected number of candles.
- The CSV market type, path, row count, first candle, and last candle match the
  intended experiment in `source_manifest.json`.
- The manifest has `product` set to `Binance USDT-M futures`,
  `comparison_only=false`, `gap_count=0`, `duplicate_count=0`, and
  `validation_status=accepted` for current futures research. For the current
  full-history futures file, `zero_volume_count=66`.
- The strategy does not inspect future candles.
- Entries happen on the next bar open.
- Stop and target prices are on the correct side of entry.
- Same-bar stop/target ambiguity is stop-first.
- Results include all time splits and side splits.
- Generated outputs are kept under `results/`.
