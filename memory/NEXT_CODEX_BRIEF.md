# Next Codex Brief: Futures Impulse Absorption Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT 5m range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md
  - docs/FUTURES_DATA_IMPACT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active research target is Binance USDT-M futures BTCUSDT 5m, not spot.
- Full-history futures source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - market type: Binance USDT-M futures
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - source manifest: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- The futures hypothesis pivot inventory stop state was:
  pivot_inventory_needs_user_hypothesis.
- A new materially different futures premise has now been selected:
  futures impulse absorption after abnormal OHLCV impulse candles.
- This premise is not a retune or reslice of SR timing, compression breakout, detector/context promotion, hold-inside directional continuation, hold-inside/midline reaction, or the failed futures midline touch prototype.
- Default runs still use lab.EmptyStrategy unless an explicit prototype/audit flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Implement a non-trading audit behind a new explicit CLI flag:
  -futures-impulse-absorption-audit.
- Use only Binance USDT-M futures BTCUSDT 5m data for the audit.
- Keep the audit OHLCV-only and closed-candle-only.
- Do not add entries, exits, scoring, sizing changes, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, or two-exchange logic.

Hypothesis:
- After a liquidation-like abnormal futures impulse candle, BTCUSDT may show absorption: price reclaims the impulse candle midpoint before extending beyond the impulse extreme.
- Candidate event is a closed 5m candle after a 30-day prior rolling warmup where:
  - true-range percentile rank against prior 30 days is at least p99;
  - volume percentile rank against prior 30 days is at least p95;
  - close position within the event candle is extreme:
    - down impulse: close position <= 0.25;
    - up impulse: close position >= 0.75.
- Use prior rolling windows only; do not include the event candle in its own percentile reference.
- True range is max(high-low, abs(high-previous close), abs(low-previous close)).
- Close position is (close-low)/(high-low); skip zero-range events.

Implementation requirements:
- Add Go-native audit code; do not import strategy evidence or implementation from another repo.
- Add a CLI flag in cmd/rangelab:
  -futures-impulse-absorption-audit.
- Reject spot/comparison sources for this audit. The audit must require source_manifest product "Binance USDT-M futures" and comparison_only=false.
- Keep lab.EmptyStrategy as the default and keep the audit non-trading; common summary/trades outputs should remain zero-trade unless another explicit existing strategy/prototype flag is used.
- Use horizons 3, 6, 12, and 24 closed bars after the event.
- For an up impulse:
  - midpoint reclaim means a later candle low touches or crosses the event midpoint.
  - continuation beyond extreme means a later candle high exceeds the event high.
- For a down impulse:
  - midpoint reclaim means a later candle high touches or crosses the event midpoint.
  - continuation beyond extreme means a later candle low falls below the event low.
- Track which happened first within each horizon; if both happen on the same candle, record same_bar_ambiguous=true and count it separately from clean reclaim-first or clean continuation-first.
- Track quick continuation within the first 3 bars.
- Split by event close time using the existing default splits.

Required outputs under results/futures-impulse-absorption-audit:
- source_manifest.json
- summary.csv/json
- trades.json
- futures_impulse_absorption_candidates.csv/json
- futures_impulse_absorption_summary.csv/json
- futures_impulse_absorption_stability.csv/json

Candidate rows should include at minimum:
- event id, index, open/close time, split, direction, open/high/low/close/volume;
- previous close, true range, true-range percentile rank, volume percentile rank;
- event range pct, close position, event midpoint;
- label window start/end per horizon;
- midpoint reclaim, continuation beyond extreme, first outcome, same-bar ambiguity, quick continuation, missing future;
- bars to reclaim and bars to continuation when present.

Summary rows should group by split, direction, horizon, true-range percentile bucket, and volume percentile bucket, plus all-direction/all-bucket rollups. Include source event count, labeled event count, missing-future count, reclaim-first count/rate, continuation-first count/rate, same-bar ambiguity count/rate, quick-continuation count/rate, average bars to reclaim, and average bars to continuation.

Stability rows should compare period splits, excluding full_2021_2026 from the period set, and report weakest split event counts, min/max reclaim-first rate, min/max continuation-first rate, max quick-continuation rate, max same-bar ambiguity rate, and deltas.

Review gate:
- Stop with impulse_absorption_source_gap if source validation, manifest facts, or output completeness fails.
- Stop with impulse_absorption_codegen_or_test_blocked if implementation or verification cannot complete.
- Stop with impulse_absorption_no_viable_edge if the source is valid but any of these fail:
  - every period split has at least 100 candidates for the all-direction/all-bucket row;
  - midpoint reclaim-first beats continuation-first in every period split for at least one horizon;
  - quick continuation does not dominate the first 3 bars;
  - same-bar ambiguity does not dominate the apparent edge.
- Stop with impulse_absorption_audit_ready if outputs are complete and the non-trading review identifies a split-stable absorption surface worth reviewing in a dedicated verdict doc.
- Stop with impulse_absorption_needs_review_only if outputs are complete but the correct next task is only a review doc, not another implementation.

Suggested implementation tests:
- Candidate detection uses only prior 30-day percentile windows and excludes the event candle from thresholds.
- Up/down impulse direction and close-position cutoffs are exact.
- Zero-range events are skipped.
- Midpoint reclaim, continuation beyond extreme, same-bar ambiguity, and quick continuation labels are correct for both directions.
- Missing future windows are counted.
- Summary rollups and stability rows include expected split/direction/bucket combinations.
- CLI flag writes the impulse absorption artifacts while default runs remain empty-strategy only.
- The audit rejects spot/comparison sources.

Run:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
    -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv \
    -source-product binance-usdm-futures \
    -futures-impulse-absorption-audit \
    -out-dir results/futures-impulse-absorption-audit
- wc -l results/futures-impulse-absorption-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Closeout:
- Create docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md only after the audit outputs exist and have been reviewed.
- Add that review doc to the README docs index if created.
- Update memory/PROGRESS.md with commands, result paths, manifest facts, CSV counts, and stop state.
- Update memory/DECISIONS.md only if the review creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md from the chosen stop state.
- Commit completed repo changes after verification unless explicitly told not to.
```
