# Next Codex Brief: Futures Range-Universe Post-Compression Fragility Pivot

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_REPLAY_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/FUTURES_SCOPE_PIVOT_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy discovery.
- The frozen ETH/SOL 4h structured-compression stream
  sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00 passed fixed replay but
  failed walk-forward robustness.
- Walk-forward stop state:
  structured_compression_walk_forward_fragile_needs_review.
- Walk-forward result:
  - source and closed UTC 4h resample validation passed for BTCUSDT, ETHUSDT,
    and SOLUSDT;
  - grid size remained 216 declared configs;
  - common trades.json and summary.* stayed frozen ETH/SOL authority-only;
  - fold 1 selected sc4h_btc_diagnostic_eth_sol_cw2_h12_t0_75_sb0_10, but its
    test net P&L 92.68 was worse than frozen 229.02;
  - fold 2 selected no config because the combined train period had 97
    authority trades, below the 100 multi-split train gate;
  - fold 3 selected the exact frozen config and passed, with 32 test trades,
    net P&L 193.06, PF 1.9121.
- Durable decision:
  the walk-forward result does not authorize candidate strategy packaging,
  retuning, BTCUSDT promotion, new grid dimensions, live/paper/testnet,
  exchange API, deployment, martingale, averaging down, or two-exchange work.

Goal:
- Do a review-only post-fragility pivot pass.
- Create a concise doc:
  docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md
- The doc should treat the structured-compression walk-forward result as
  exclusion evidence, summarize why the stream is not package-ready, and define
  the next materially different offline range-strategy premise or conclude that
  no next implementation is authorized without user input.

Boundaries:
- Do not implement a new strategy, optimizer, or replay in this task.
- Do not rerun the structured-compression walk-forward unless files are missing
  or counts mismatch.
- Do not tune target, stop, confirmation, max-hold, symbol set, or training
  gates around the failed walk-forward result.
- Do not reopen the failed 1h structured-compression surface.
- Do not add live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, data downloads, broad symbol mining, grid/martingale/averaging
  down, or two-exchange logic.
- BTCUSDT remains diagnostic-only for the completed ETH/SOL structured-
  compression stream.

Expected review shape:
- Inventory the still-open range-first premise space from current docs only.
- Exclude already failed or fragile families unless a materially new data or
  structure premise is explicit.
- If a next premise is available, write a bounded next implementation brief in
  memory/NEXT_CODEX_BRIEF.md with exact reads, source requirements, outputs,
  tests, review gates, and stop states.
- If no next premise is available, make memory/NEXT_CODEX_BRIEF.md a short
  review-only stop brief that says no implementation is authorized until the
  user chooses a new premise.
- Update README.md if a new review doc is added.
- Update memory/PROGRESS.md with the review command/checks and factual outcome.
- Update memory/DECISIONS.md only if the pivot review creates a durable rule.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit the completed review doc, README/memory updates, refreshed next brief,
  and verification evidence after checks pass unless explicitly told not to.
```
