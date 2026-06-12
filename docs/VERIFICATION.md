# Verification

This project should be runnable after moving the `range-strategy-lab` folder out
of the parent bot.

## Required Files

The folder contains everything needed to compile and run the lab code:

- `go.mod`
- `cmd/rangelab/main.go`
- `internal/lab/*.go`
- `docs/*.md`
- `CODEX_BRIEF.md`

The only external input needed for a real run is a candle CSV passed with
`-csv`.

## Local Verification

From inside `range-strategy-lab`:

```bash
/usr/local/go/bin/go test ./...

/usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/smoke
```

Expected smoke result with the starter strategy:

- CSV loads successfully.
- Strategy name is `empty`.
- Trade count is `0`.
- `results/smoke/summary.csv`, `summary.json`, and `trades.json` are written.

Zero trades are correct until a real strategy replaces `lab.EmptyStrategy`.

## After Moving The Folder

If you move this folder elsewhere, either:

- copy a BTCUSDT 5m CSV into the new project, or
- pass an absolute CSV path with `-csv`.

Example:

```bash
/usr/local/go/bin/go run ./cmd/rangelab \
  -csv /absolute/path/to/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/smoke
```

## What To Check Before A Strategy Result Is Trusted

- Tests pass.
- Smoke run loads the expected number of candles.
- The strategy does not inspect future candles.
- Entries happen on the next bar open.
- Stop and target prices are on the correct side of entry.
- Same-bar stop/target ambiguity is stop-first.
- Results include all time splits and side splits.
- Generated outputs are kept under `results/`.
