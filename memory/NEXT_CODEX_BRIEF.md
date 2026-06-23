# Next Codex Brief: Futures Range Candidate Discovery Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_SPEC.md
  - docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md
  - docs/FUTURES_SCOPE_PIVOT_REVIEW.md
  - docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md
  - docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active project objective is BTCUSDT futures range-strategy research.
- Scope is range-first broad, not narrow range-only:
  - 5m, buy/sell touch, single-candle reactions, boundary rejection, failed-break re-entry, and breakout continuation remain open when materially reframed.
  - Do not rerun exact failed templates under new names.
- The active source is Binance USDT-M futures. Spot evidence is historical context only unless explicitly rerun and reviewed on futures data.
- Current authoritative parent source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - required product: -source-product binance-usdm-futures
  - market type: Binance USDT-M futures BTCUSDT 5m
  - CSV lines including header: 573,985
  - loaded candles: 573,984
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - accepted manifest facts: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- Candidate source shapes:
  - native 5m parent candles;
  - closed UTC 15m, 1h, and 4h resamples from docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/prototype flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Implement `-futures-range-candidate-discovery-audit` as a non-trading discovery audit.
- Compare a compact menu of BTCUSDT futures range-adjacent candidate families across 5m, 15m, 1h, and 4h.
- Rank candidates and decide whether one or two should move to a baseline offline backtest brief.
- Keep common `summary.*` and `trades.json` zero-trade.

Implementation boundaries:
- Do not add entries, exits, scoring, sizing, optimization, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, two-exchange logic, data downloads, spot comparison, symbol expansion, or sibling repo mutation.
- Do not use spot outputs as authority for futures promotion.
- Do not broaden to BTC/ETH or other symbols unless the user explicitly changes scope.
- Do not optimize before a baseline candidate clears this discovery gate.

Candidate families:
- Mature balance persistence:
  - timeframes: 1h and 4h, with optional 15m cross-check;
  - candidate: compact multi-bar range with adequate inside closes and no fresh expansion;
  - labels: inside persistence, internal rotation, expansion failure.
- Boundary touch rejection:
  - timeframes: 5m, 15m, 1h;
  - candidate: closed candle touches a mature range boundary without closing decisively beyond it;
  - labels: reject inward, accept outside, stall/none.
- Single-candle wick rejection:
  - timeframes: 5m, 15m, 1h;
  - candidate: large wick at a mature range boundary with close back inside the range;
  - labels: inward follow-through, boundary break, no follow-through.
- Failed breakout re-entry:
  - timeframes: 5m, 15m, 1h;
  - candidate: closed break outside a mature range followed by closed re-entry within a fixed window;
  - labels: re-entry continuation, second break, no follow-through.
- Clean breakout continuation:
  - timeframes: 15m, 1h, 4h;
  - candidate: mature range compression followed by a decisive closed break;
  - labels: continuation beyond break, failed break/re-entry, no extension.

Required outputs:
- Result dir: results/futures-range-candidate-discovery-audit/
- Write compact CSV/JSON artifacts:
  - futures_range_discovery_candidates.csv/json
  - futures_range_discovery_summary.csv/json
  - futures_range_discovery_rankings.csv/json
  - futures_range_discovery_stability.csv/json
  - normal source_manifest.json, summary.* and trades.json
- Include source/resample coverage facts for every timeframe used.
- Include split, timeframe, family, direction/boundary side where applicable, candidate count, favorable/adverse/neutral rates, excursion summaries, rough cost buffer proxy, and ranking fields.

Review gate:
- A candidate can move to baseline backtest only if:
  - source manifest is accepted Binance USDT-M futures BTCUSDT, comparison_only=false;
  - resample coverage is complete for any 15m/1h/4h candidate;
  - each period split has adequate candidate counts;
  - favorable outcome rate beats adverse outcome rate in every period split;
  - adverse move and quick invalidation do not dominate;
  - rough cost buffer is plausible before sizing or optimization;
  - full sample and weakest split are coherent enough for a fixed-rule prototype.

Stop states:
- range_discovery_audit_ready: implementation and outputs complete; at least one family clears the review gate and the next brief is a baseline backtest.
- range_discovery_no_backtest_candidate: implementation and outputs complete; no family clears the review gate.
- range_discovery_source_or_resample_gap: source validation or resampling completeness fails.
- range_discovery_codegen_or_test_blocked: implementation or verification cannot complete.
- range_discovery_review_only_no_strategy_change: outputs document context but do not change next implementation direction.

Review and memory:
- Create docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md after outputs exist.
- Add the review doc to README.
- Update memory/PROGRESS.md with commands, result paths, manifest/resample facts, row counts, candidate/ranking counts, and stop state.
- Update memory/DECISIONS.md only if the review creates a durable promotion/no-promotion/source rule.
- Refresh memory/NEXT_CODEX_BRIEF.md from the chosen stop state.
- If the gate passes, the refreshed next brief should be a baseline offline backtest for the top one or two candidates, not another review-only loop.

Suggested verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-range-candidate-discovery-audit -out-dir results/futures-range-candidate-discovery-audit
- wc -l results/futures-range-candidate-discovery-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit completed implementation, generated review/doc memory updates, and refreshed next brief after verification unless explicitly told not to.
```
