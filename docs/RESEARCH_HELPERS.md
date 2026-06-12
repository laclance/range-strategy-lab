# Research Helper Modules

This project may use small, pure-Go helper modules for feature extraction and
audit outputs. These modules do not change the project boundary: strategy
hypotheses, entries, exits, scoring, sizing, and backtest behavior stay inside
`range-strategy-lab`.

## Added Modules

- `github.com/laclance/go-sr v1.0.0`
  - Use for closed-candle support/resistance detection.
  - Best first candidate for the next audit milestone: SR zones, nearest
    support/resistance, and boundary-touch diagnostics.
  - Keep it behind a lab adapter that converts from `lab.Candle` and records
    explainable SR fields in generated results.

- `github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f`
  - Use as a general indicator toolbox when an experiment needs standard TA
    features such as RSI, ATR, Bollinger bands, ADX, stochastic, MFI, or
    regression slope.
  - Prefer adapters and focused parity tests before replacing existing local
    indicator helpers.
  - Do not add many indicators at once; each indicator should answer a specific
    hypothesis.

- `nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc`
  - Use for candlestick-pattern audit or a single confirmation filter.
  - Treat pattern labels as weak evidence until split-tested on BTCUSDT 5m.
  - Do not use candle-pattern combinations as a broad optimizer surface.

## Integration Rules

- Keep all usage offline and deterministic.
- Use confirmed closed-candle prefixes only.
- Enter on the next bar open when entries are eventually implemented.
- Keep generated helper outputs under `results/`.
- Add focused tests for adapters, especially candle conversion, warmup behavior,
  and no-lookahead assumptions.
- Do not import live execution, exchange clients, deployment code, strategy
  scoring, grid, martingale, averaging down, or two-exchange execution logic.

## Preferred Next Adapter

Add an SR audit mode with `github.com/laclance/go-sr` before adding trades:

- fixed BTCUSDT 5m input
- balanced detector baseline available for context
- outputs such as `sr_zones.csv/json` or `sr_touch_audit.csv/json`
- no trade entries in the audit milestone

After SR audit outputs are inspectable, choose one entry template only, such as
boundary rejection or false-break reclaim.
