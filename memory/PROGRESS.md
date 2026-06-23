# Progress

This file is the always-read project snapshot. Keep it compact: current state,
latest verification, result paths, and a milestone index. Detailed evidence
belongs in focused `docs/` reviews, generated artifacts under `results/`, and
git history.

## Current State

- Scope is now offline Binance USDT-M futures range-strategy discovery. The
  implemented CLI default remains BTCUSDT `5m`, but the next approved research
  scope is a local BTC/ETH/SOL range-universe discovery spec. BTCUSDT `5m`,
  buy/sell-touch, and single-candle reaction ideas remain available when they
  are materially reframed and compared inside a discovery-to-backtest funnel.
- Active market target is Binance USDT-M futures, not spot. Spot-generated
  audits/reviews are historical context only unless a futures rerun explicitly
  revalidates a specific conclusion.
- Default strategy remains `lab.EmptyStrategy`; trades remain `0` unless an
  explicit offline prototype flag is passed.
- Closed-candle decision semantics remain required.
- `cmd/rangelab` now enforces source identity before audits/backtests. The
  default `-csv` is the full-history Binance USDT-M futures file
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`; non-default CSVs
  require `-source-product`; spot CSVs require explicit comparison flags; every
  accepted run writes `source_manifest.json`.
- No live code, API keys, deploy scripts, grid, martingale, averaging down, or
  two-exchange execution is allowed.
- `memory/NEXT_CODEX_BRIEF.md` is the only canonical next-session prompt.
- The first minimal offline prototype for the futures-revalidated
  `hold_3_inside` + first `mid_touch` within `12` bars + event close-position
  bucket `mid_50` surface has been built and reviewed. The close-back
  boundary-target template failed P&L and is not promoted.
- Futures impulse absorption after abnormal OHLCV impulse candles has been
  implemented, run, and reviewed on Binance USDT-M futures data. It failed the
  non-trading review gate: continuation-first and quick continuation dominated
  midpoint reclaim-first across all period splits and horizons. No prototype or
  promotion is approved.
- Futures range candidate discovery has been implemented, run, and reviewed.
  It passed the discovery gate only for clean breakout continuation surfaces.
  The fixed-rule offline baseline backtest for the top non-duplicative `4h`
  up-breakout and `1h` all-side clean-breakout candidates has now been built
  and reviewed. Both candidates failed after costs; no optimization,
  portfolio-style stream, or automatic `15m` expansion is approved from that
  baseline.
- The current next task is a non-trading local universe discovery audit across
  BTCUSDT, ETHUSDT, and SOLUSDT Binance USDT-M futures `5m` sources. It must
  validate sources first, rank range-first surfaces, and route quickly to a
  fixed-rule baseline backtest only if one or two surfaces pass.

## 2026-06-26

Futures range universe discovery spec:

- Spec doc: `docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_SPEC.md`.
- Stop state: `range_universe_spec_ready_for_audit_implementation`.
- Outcome: the project explicitly broadens from BTCUSDT-only to a local
  Binance USDT-M futures range universe for discovery. The v1 universe is
  limited to local BTCUSDT, ETHUSDT, and SOLUSDT `5m` CSVs under
  `../binance-bot/data/`.
- Motivation: the BTCUSDT clean-breakout baseline had adequate trade counts
  but failed after costs, so the next useful step is source-validated
  cross-symbol range discovery rather than retuning the failed immediate
  breakout template.
- Local source quick checks:
  - `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`: `573,985`
    CSV lines including header, `573,984` rows, min open `1609459200000`,
    max open `1781654100000`, non-monotonic physical row count `0`
  - `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`: `573,985`
    CSV lines including header, `573,984` rows, min open `1609459200000`,
    max open `1781654100000`, non-monotonic physical row count `0`
  - `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`: `573,985`
    CSV lines including header, `573,984` rows, min open `1609459200000`,
    max open `1781654100000`, non-monotonic physical row count `1`
- The SOL caveat is recorded in the spec: the future implementation must
  validate and sort/fail explicitly instead of trusting filename coverage.
- Future candidate families:
  - breakout retest / acceptance after completed mature range;
  - boundary touch / rejection;
  - failed breakout re-entry;
  - mature balance rotation / persistence;
  - compression-to-expansion only if materially reframed from the failed
    immediate clean-breakout baseline.
- Required future outputs are scoped under
  `results/futures-range-universe-discovery-audit/` and remain zero-trade
  common outputs until a later baseline backtest is approved.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement
  `-futures-range-universe-discovery-audit` as a non-trading source/universe
  validation plus discovery audit.
- This was docs/memory-only: no code, CLI flags, generated results, entries,
  exits, scoring, sizing, optimization, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, or sibling
  repo mutation was added.
- Verification passed:
  - `wc -l ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`
  - `awk -F, 'NR>1{if(NR==2||$1<min)min=$1;if($1>max)max=$1;if(NR>2 && $1<=prev) nonmono++; prev=$1} END{print FILENAME, NR-1, min, max, nonmono+0}' <source>`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures clean breakout baseline backtest:

- Review doc: `docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md`.
- Stop state: `clean_breakout_baseline_failed_no_promotion`.
- Added explicit offline CLI flag:
  `-futures-clean-breakout-baseline-backtest`.
- Default runs still use `lab.EmptyStrategy`; this baseline runs only when the
  explicit flag is passed and rejects spot/comparison sources.
- First baseline evaluated candidates independently, per user direction. The
  user-approved portfolio-stream routing rule was checked and did not trigger:
  neither candidate was near-viable after costs.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT `5m`.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Resample coverage used by the baseline:
  - `1h`: `47,832` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:00:00Z`
  - `4h`: `11,958` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T20:00:00Z`
- Result dir:
  `results/futures-clean-breakout-baseline-backtest/`.
- Outputs:
  - `futures_clean_breakout_baseline_signals.csv/json`
  - `futures_clean_breakout_baseline_trades.csv/json`
  - `futures_clean_breakout_baseline_summary.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - signals `4,849`
  - trades `1,186`
  - baseline summary `25`
  - common summary `13`
- Baseline result:
  - `4,848` signal rows
  - `1,185` executed trades
  - `clean_breakout_4h_up_h12`: `496` signals, `121` trades, gross P&L
    `2.55`, net P&L `-69.29`, PF `0.8386`
  - `clean_breakout_1h_all_h12`: `4,352` signals, `47` skipped signals,
    `1,064` trades, gross P&L `272.63`, net P&L `-320.25`, PF `0.8846`
  - aggregate compatibility summary: `1,185` trades, gross P&L `275.18`,
    net P&L `-389.54`, PF `0.8784`
- Failure facts:
  - both candidates had negative full net P&L after costs;
  - both had full PF below `1.2`;
  - both were negative in `2023_2024_oos` and `2025_2026_recent`;
  - the `1h` all-side candidate was negative on both long and short sides.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to request a new futures premise or
  explicit scope choice before another backtest; do not optimize this clean
  breakout baseline.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-clean-breakout-baseline-backtest -out-dir results/futures-clean-breakout-baseline-backtest`
  - `wc -l results/futures-clean-breakout-baseline-backtest/*.csv`

Futures range candidate discovery audit:

- Review doc: `docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md`.
- Stop state: `range_discovery_audit_ready`.
- Added explicit non-trading CLI flag:
  `-futures-range-candidate-discovery-audit`.
- Default runs still use `lab.EmptyStrategy`; this audit writes labels,
  coverage, rankings, summaries, and stability rows only and produced `0`
  trades.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT `5m`.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Resample coverage:
  - `5m`: `573,984` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:55:00Z`
  - `15m`: `191,328` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:45:00Z`
  - `1h`: `47,832` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T23:00:00Z`
  - `4h`: `11,958` rows, first open `2021-01-01T00:00:00Z`, last open
    `2026-06-16T20:00:00Z`
  - every timeframe had `gap_count=0`, `duplicate_count=0`,
    `missing_child_open_count=0`, `complete=true`,
    `validation_status=accepted`
- Result dir:
  `results/futures-range-candidate-discovery-audit/`.
- Outputs:
  - `futures_range_discovery_coverage.csv/json`
  - `futures_range_discovery_candidates.csv/json`
  - `futures_range_discovery_summary.csv/json`
  - `futures_range_discovery_rankings.csv/json`
  - `futures_range_discovery_stability.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - candidates `144,637`
  - coverage `5`
  - rankings `121`
  - stability `121`
  - discovery summary `481`
  - common summary `13`
- Audit result:
  - `144,636` event/horizon rows
  - `48,212` distinct events
  - `120` ranked surfaces
  - `24` passing rows, all from `clean_breakout_continuation`
  - passing rows by timeframe: `15m=9`, `1h=9`, `4h=6`
- Top non-duplicative candidates for the next baseline:
  - `4h clean_breakout_continuation up h12`: full count `496`, weakest split
    count `112`, full favorable rate `86.69%`, weakest favorable rate
    `82.38%`, weakest cost buffer `4.2910%`
  - `1h clean_breakout_continuation all h12`: full count `4,352`, weakest
    split count `1,213`, full favorable rate `90.37%`, weakest favorable rate
    `89.01%`, weakest cost buffer `2.5907%`
- Rejected for this baseline: boundary touch rejection, single-candle wick
  rejection, failed breakout re-entry, and mature balance persistence.
- Added README docs-index entry for the new review doc.
- Added a durable decision that only clean breakout continuation is authorized
  for the next fixed-rule baseline; the other discovery families are not.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement
  `-futures-clean-breakout-baseline-backtest`.
- No entries, exits, scoring, sizing, optimization, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, spot
  comparison, symbol expansion, or sibling repo mutation was added by this
  audit.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-range-candidate-discovery-audit -out-dir results/futures-range-candidate-discovery-audit`
  - `wc -l results/futures-range-candidate-discovery-audit/*.csv`

Futures range candidate discovery spec:

- Spec doc: `docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_SPEC.md`.
- Stop state: `range_discovery_spec_ready_for_audit_implementation`.
- Outcome: the repo scope is corrected from overly narrow range-only gating to
  range-first broad discovery. `5m`, buy/sell touch, single-candle reactions,
  boundary rejection, failed-break re-entry, and breakout continuation remain
  open when materially reframed and tied to a fast baseline-backtest gate.
- Current lab authority remains Binance USDT-M futures BTCUSDT `5m`:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - required CLI product: `-source-product binance-usdm-futures`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Discovery menu for the next audit:
  - mature balance persistence on `1h`/`4h` with optional `15m` cross-check;
  - boundary touch rejection on `5m`/`15m`/`1h`;
  - single-candle wick rejection on `5m`/`15m`/`1h`;
  - failed breakout re-entry on `5m`/`15m`/`1h`;
  - clean breakout continuation after mature range compression on
    `15m`/`1h`/`4h`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to implement
  `-futures-range-candidate-discovery-audit` as a non-trading discovery audit.
  If one or two candidates clear the gate, the following brief should be a
  baseline offline backtest, not another review-only loop.
- Added README docs-index entry for the new spec.
- Added a durable decision that closed families are not permanently banned;
  exact failed templates are banned, but materially reframed range-first ideas
  may compete in a broad futures discovery funnel.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, spot
  comparison, symbol expansion, or sibling repo mutation was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures higher-timeframe range source spec:

- Spec doc: `docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md`.
- Stop state: `higher_tf_range_source_spec_no_viable_range_premise`.
- Outcome: the parent source and closed UTC resampling contract are specified,
  but no higher-timeframe audit is ready until the user supplies a materially
  different range premise, closed-candle candidate event, and falsification
  rule.
- Current lab authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - required CLI product: `-source-product binance-usdm-futures`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Candidate higher-timeframe bars from the current parent:
  - `15m`: `3` child bars, expected complete bars `191,328`, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T23:45:00Z`
  - `1h`: `12` child bars, expected complete bars `47,832`, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T23:00:00Z`
  - `4h`: `48` child bars, expected complete bars `11,958`, first open
    `2021-01-01T00:00:00Z`, last open `2026-06-16T20:00:00Z`
- Resampling contract: UTC bucket start as higher-timeframe open; complete
  buckets only; open first child open, high max child high, low min child low,
  close last child close, volume sum; reject partial, missing-child,
  duplicate, synthetic, forward-filled, spot/comparison, non-UTC, or
  future-looking bars.
- Added README docs-index entry for the new source spec.
- Added a durable decision that higher-timeframe range work must derive from
  the accepted BTCUSDT futures 5m parent source until an explicit source/scope
  change.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` into a premise-request brief, not an
  audit brief.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, data download, or sibling
  repo mutation was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures scope pivot review, range-only:

- Review doc: `docs/FUTURES_SCOPE_PIVOT_REVIEW.md`.
- Stop state: `range_scope_pivot_ready_for_higher_timeframe_source_spec`.
- Decision: pause BTCUSDT 5m range continuation work. The next range-only
  route is a BTCUSDT higher-timeframe futures source/premise spec, not an audit
  or prototype.
- Added README docs-index entry for the new review doc.
- Added a durable decision that BTCUSDT higher-timeframe range-source work is
  the next authorized lane; BTC/ETH expansion remains deferred until that
  source/premise review is complete or the user explicitly changes scope.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to create
  `docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md`.
- Current lab authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Sibling repo facts used as context only:
  - `~/binance-bot` BTCUSDT 5m range sleeve failed on legacy spot data; caution
    only, not futures authority.
  - `~/binance-bot` cross-exchange spread and broad selector work remain
    exclusion evidence for this lab.
  - `~/binance-bot` daily ATR contraction is non-range trend/volatility context
    and cannot promote range work.
  - `~/crypto-trading-bot` BTC/ETH USD-M source-contract discipline is process
    context only.
  - `~/crypto-trading-bot` BTC/ETH 4h range-reversal/re-entry is closed
    exclusion evidence, not a port target.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, two-exchange logic, or sibling repo mutation was
  added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

Futures scope pivot review spec, range-only:

- Spec doc: `docs/FUTURES_SCOPE_PIVOT_REVIEW_SPEC.md`.
- Stop state: `range_scope_pivot_spec_ready`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to make the next task a docs-only
  `docs/FUTURES_SCOPE_PIVOT_REVIEW.md`, focused on range strategy scope only.
- Added README docs-index entry for the new spec.
- Added a durable decision: a scope pivot remains range-strategy-only unless
  the user explicitly changes the project objective. Higher-timeframe,
  multi-symbol, or sibling-repo context may be used only as range-source,
  range-premise, process, or exclusion evidence.
- Current lab authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- Sibling repo facts recorded for next-review context only:
  - `~/binance-bot` range seed: Binance-only BTCUSDT 5m range sleeve on
    legacy spot data failed even gross; use only as caution.
  - `~/binance-bot` cross-exchange spread seed failed local gates and remains
    out of scope because this lab forbids two-exchange execution.
  - `~/binance-bot` multi-pair selector/alt overlay work is no-go evidence for
    broad symbol selection.
  - `~/binance-bot` daily ATR contraction is positive-looking but non-range
    trend/volatility work, so it is outside this range-only pivot.
  - `~/crypto-trading-bot` BTC/ETH 4h range-reversal/re-entry workbench failed
    as a frozen family; use only as exclusion evidence.
  - `~/crypto-trading-bot` BTC/ETH USD-M source-contract discipline may be
    useful if the review selects higher-timeframe or BTC/ETH range-source
    specification.
- This was docs/memory-only: no code, CLI flags, audits, generated results,
  source mutation, entries, exits, scoring, sizing, strategy replacement,
  paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, or two-exchange logic was added.
- Sibling repo status note:
  - `~/binance-bot` was clean during inspection.
  - `~/crypto-trading-bot` had pre-existing dirty docs/memory/research files
    and was treated as read-only.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` showed only intended doc/memory changes before
    commit.

## 2026-06-25

Futures impulse absorption audit:

- Review doc: `docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md`.
- Stop state: `impulse_absorption_no_viable_edge`.
- Added explicit non-trading CLI flag:
  `-futures-impulse-absorption-audit`.
- Default runs still use `lab.EmptyStrategy`; this audit writes labels and
  summaries only and produced `0` trades.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT 5m.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Result dir:
  `results/futures-impulse-absorption-audit/`.
- Outputs:
  - `futures_impulse_absorption_candidates.csv/json`
  - `futures_impulse_absorption_summary.csv/json`
  - `futures_impulse_absorption_stability.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - candidates `13,737`
  - summary `949`
  - stability `241`
  - common summary `13`
- Audit result:
  - `3,434` distinct impulse events
  - `13,736` event/horizon candidate rows
  - `1,838` up impulse events and `1,596` down impulse events
  - no missing future windows
  - all period splits have at least `1,003` all-direction/all-bucket events
    at every horizon
- Gate failure:
  - midpoint reclaim-first never beats continuation-first in any period split
    at horizons `3`, `6`, `12`, or `24`
  - all-direction/all-bucket continuation-first rates are roughly
    `57.74%` to `59.58%`
  - reclaim-first rates are only roughly `24.13%` to `25.43%`
  - maximum quick-continuation rate is `80.76%`
  - minimum reclaim-minus-continuation margin stays negative at every horizon
- No entry, exit, scoring, sizing, strategy replacement, paper/testnet/live
  wiring, exchange API use, deploy files, grid, martingale, averaging down, or
  two-exchange logic was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-impulse-absorption-audit -out-dir results/futures-impulse-absorption-audit`
  - `wc -l results/futures-impulse-absorption-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Futures impulse absorption next-brief refresh:

- Stop state: `impulse_absorption_audit_ready`.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` from the open-ended
  new-hypothesis prompt into a scoped non-trading audit brief for
  `-futures-impulse-absorption-audit`.
- Chosen premise: after a liquidation-like abnormal Binance USDT-M futures 5m
  impulse candle, BTCUSDT may show absorption by reclaiming the event candle
  midpoint before extending beyond the event extreme.
- Candidate event definition for the next audit:
  - closed BTCUSDT 5m futures candle after a `30` day prior rolling warmup
  - true-range percentile rank at least `p99`
  - volume percentile rank at least `p95`
  - close position `<=0.25` for down impulses or `>=0.75` for up impulses
- Required next-audit source remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - expected accepted manifest facts: `gap_count=0`, `duplicate_count=0`,
    `zero_volume_count=66`, `comparison_only=false`,
    `validation_status=accepted`
- This was docs/memory-only handoff work: no code, CLI flags, audits, result
  directories, entries, exits, scoring, sizing, paper/testnet/live wiring,
  exchange API use, deploy files, grid, martingale, averaging down, or
  two-exchange logic was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` before commit showed only intended memory changes.

Futures hypothesis pivot inventory:

- Inventory doc: `docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md`.
- Stop state: `pivot_inventory_needs_user_hypothesis`.
- Outcome: reviewed material is now an exclusion map plus reusable
  infrastructure, not a new strategy queue.
- Source authority remains Binance USDT-M futures BTCUSDT 5m:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Classification:
  - SR timing and compression breakout are closed legacy spot-only evidence.
  - Range durability is diagnostic/infrastructure, not an entry surface.
  - Detector durability/context refinement is reusable infrastructure; futures
    review revalidated context shape only, not entry promotion.
  - Hold-inside directional edge is closed.
  - Hold-inside midline transition/reaction is diagnostic after the failed
    prototype; do not broaden into `hold_6_inside`, `mid_close_across`, side
    cohorts, or old spot authority.
  - Futures midline touch prototype is a futures-authoritative failure.
- No code, CLI flags, audits, result directories, entry/exit logic, scoring,
  sizing, paper/testnet/live wiring, exchange API use, deploy files, grid,
  martingale, averaging down, or two-exchange logic was added.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short` before commit showed only intended doc/memory changes.

Minimal futures midline touch prototype:

- Review doc: `docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md`.
- Stop state: `prototype_failed_no_promotion`.
- Added explicit offline CLI flag:
  `-hold-inside-midline-touch-prototype`.
- Default runs still use `lab.EmptyStrategy`; the prototype runs only when the
  flag is passed.
- Prototype source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT 5m.
- Source manifest:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Result dir:
  `results/futures-hold-inside-midline-touch-prototype/`.
- Prototype outputs:
  - `hold_inside_midline_touch_prototype_signals.csv/json`
  - `hold_inside_midline_touch_prototype_trades.csv/json`
  - `hold_inside_midline_touch_prototype_summary.csv/json`
  - common `summary.csv/json`, `trades.json`, and `source_manifest.json`
- CSV line counts including headers:
  - signals `533`
  - prototype trades `532`
  - prototype summary `13`
  - common summary `13`
- Run output:
  - `signal_rows=532`
  - `trades=531`
  - `summary_rows=12`
  - one exact-mid skip
  - exit reasons: `140` stop losses, `82` take profits, `309` time stops
- Full-sample result:
  - `531` trades, win rate `29.94%`
  - gross P&L `-95.54`, net P&L `-418.99`
  - profit factor `0.3409`, average net R `-0.4276`
  - max drawdown `42.11%`
- Period splits all failed:
  - `2021_2022_stress`: net P&L `-226.03`, PF `0.3824`
  - `2023_2024_oos`: net P&L `-116.58`, PF `0.3165`
  - `2025_2026_recent`: net P&L `-76.38`, PF `0.2298`
- Side splits both failed:
  - long: `248` trades, net P&L `-234.50`, PF `0.2805`
  - short: `283` trades, net P&L `-184.49`, PF `0.4045`
- Added engine entry-geometry guard so next-bar entries are skipped when the
  stop/target would be on the wrong side of the actual entry.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -hold-inside-midline-touch-prototype -out-dir results/futures-hold-inside-midline-touch-prototype`
  - `wc -l results/futures-hold-inside-midline-touch-prototype/*.csv`

Futures data impact review:

- Review doc: `docs/FUTURES_DATA_IMPACT_REVIEW.md`.
- Stop state: `futures_reaction_gate_passed_needs_minimal_entry_brief`.
- Source path:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`.
- Market type: Binance USDT-M futures BTCUSDT 5m.
- Source manifests in all futures result dirs:
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Result dirs:
  - `results/futures-detector-context-refinement-audit/`
  - `results/futures-hold-inside-midline-transition-audit/`
  - `results/futures-hold-inside-midline-reaction-audit/`
- Audit sizes:
  - detector context refinement: `117,848` candidate rows, `640` summary rows,
    `160` stability rows
  - hold-inside midline transition: `8,600` candidate rows, `672` summary rows,
    `168` stability rows
  - hold-inside midline reaction: `10,172` candidate rows, `24` funnel rows,
    `1,240` summary rows, `336` stability rows
- Futures reaction gate:
  - `hold_3_inside + mid_touch`: weakest split `129` event candidates,
    `55.60%` minimum event rate, `44.40%` maximum missing-event rate
  - h6 all-bucket: weakest split `129` candidates, `56.59%` minimum
    close-back, `48.06%` minimum mid-rejection before boundary, `29.46%`
    maximum boundary-before-rejection, `18.60%` maximum quick invalidation,
    `22.48%` maximum trend
  - h6 `mid_50`: weakest split `114` candidates, `57.89%` minimum close-back,
    `52.63%` minimum mid-rejection before boundary, `23.68%` maximum
    boundary-before-rejection, `13.16%` maximum quick invalidation, `21.93%`
    maximum trend
- Revalidated for the next task: first minimal offline prototype around
  `hold_3_inside` + first `mid_touch` within `12` bars + event close-position
  bucket `mid_50`.
- Still diagnostic only: old spot approval as evidence, `hold_6_inside`,
  `mid_close_across`, side-specific cohorts, and `hold_3_inside_mid_50`.
- No entry, exit, scoring, sizing, strategy replacement, paper, testnet, live,
  exchange-key, deploy, grid, martingale, averaging-down, or two-exchange work
  was added. The runs reported `strategy=empty trades=0`.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -detector-context-refinement-audit -out-dir results/futures-detector-context-refinement-audit`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -hold-inside-midline-transition-audit -out-dir results/futures-hold-inside-midline-transition-audit`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -hold-inside-midline-reaction-audit -out-dir results/futures-hold-inside-midline-reaction-audit`
  - `wc -l results/futures-*/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Futures source guard and next-brief refresh:

- Added CLI source enforcement:
  - default CSV:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - valid `-source-product` values: `binance-usdm-futures` and `binance-spot`
  - non-default CSV paths must pass `-source-product`
  - spot paths require `-source-product binance-spot` plus
    `-allow-spot-comparison`
  - accepted runs write `source_manifest.json`
- Added Go-native source validation inspired by `crypto-trading-bot` source
  contract discipline, without importing its Python helpers or strategy
  evidence.
- Real futures smoke manifest:
  - path:
    `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`
  - market type: Binance USDT-M futures
  - CSV lines including header: `573,985`
  - loaded candles / manifest `row_count`: `573,984`
  - open-time coverage: `2021-01-01T00:00:00Z` through
    `2026-06-16T23:55:00Z`
  - close-time end from smoke output: `2026-06-16T23:59:59Z`
  - `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
    `comparison_only=false`, `validation_status=accepted`
- Updated `README.md`, `docs/VERIFICATION.md`, `memory/DECISIONS.md`, and
  `memory/NEXT_CODEX_BRIEF.md`. The next brief now starts from the full-history
  futures CSV and asks for `docs/FUTURES_DATA_IMPACT_REVIEW.md` plus the three
  paused futures audit reruns.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -out-dir results/source-guard-smoke`
  - `wc -l results/source-guard-smoke/*.json results/source-guard-smoke/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

## 2026-06-16

Futures data-source correction:

- User clarified the real trading target is Binance futures rather than spot.
- Previous audit/review outputs used the spot path
  `../binance-bot/data/btcusdt_spot_5m_2021_2026.csv`; treat all promotion
  implications from those results as suspended until futures revalidation.
- Local sibling data currently visible:
  - `../binance-bot/data/btcusdt_spot_5m_2021_2026.csv`: `573,697` CSV lines
    including header, spanning `2021-01-01T00:00:00Z` through
    `2026-06-15T23:59:59Z`
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-13_2026-06-15.csv`:
    `865` CSV lines including header
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-14_2026-06-15.csv`:
    `577` CSV lines including header
  - `../binance-bot/data/btcusdt_futures_um_5m_2026-06-15.csv`: `289` CSV
    lines including header
- Updated memory and docs so the next brief is an impact assessment, not entry
  implementation. If full-history futures data is unavailable, the next verdict
  should be "data gap first," not a prototype.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Memory context-budget wording adjustment:

- Updated `AGENTS.md`, `memory/README.md`, and `memory/DECISIONS.md` so the
  size guidance for all always-read memory files is a soft `300-350` line
  judgment band, not a hard threshold. Compact or split memory when it feels
  bulky or repetitive, not merely because one file crosses `300` lines.

## 2026-06-15

Hold-inside midline reaction review milestone:

- Review doc: `docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md`.
- Inputs: `results/hold-inside-midline-reaction-audit/`.
- Original spot-data verdict: one first minimal offline entry prototype was
  justified, but only for `hold_3_inside` + first `mid_touch` within `12` bars
  + event close-position bucket `mid_50`.
- Current status after futures correction: suspended pending futures rerun and
  impact review.
- Not promoted: live use, strategy promotion, broad detector-family entries,
  `hold_6_inside`, `mid_close_across`, side-specific cohorts, and
  `hold_3_inside_mid_50`.
- Key evidence:
  - funnel pass: `hold_3_inside` + `mid_touch` has weakest-split event rate
    `52.25%` with `116` event candidates
  - `hold_3_inside` + `mid_touch`, h6 all-bucket: weakest split `116`
    candidates, `55.17%` minimum close-back, `45.69%` minimum mid-rejection
    before boundary, `28.25%` maximum boundary-before-rejection, `18.39%`
    maximum quick invalidation, `22.22%` maximum trend
  - `hold_3_inside` + `mid_touch` + event-close `mid_50`, h6: weakest split
    `104` candidates, `58.65%` minimum close-back, `50.96%` minimum
    mid-rejection before boundary, `22.12%` maximum boundary-before-rejection,
    `11.54%` maximum quick invalidation, `19.23%` maximum trend
- Updated `README.md`, `memory/DECISIONS.md`, and
  `memory/NEXT_CODEX_BRIEF.md` so the next task is the bounded offline
  prototype.
- Verification passed:
  - `wc -l results/hold-inside-midline-reaction-audit/*.csv`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Hold-inside midline reaction audit milestone:

- Added CLI flag `-hold-inside-midline-reaction-audit`.
- Result directory: `results/hold-inside-midline-reaction-audit/`.
- Outputs:
  - `hold_inside_midline_reaction_candidates.csv/json`
  - `hold_inside_midline_reaction_funnel_summary.csv/json`
  - `hold_inside_midline_reaction_summary.csv/json`
  - `hold_inside_midline_reaction_stability.csv/json`
- Audit size:
  - profiles: `1`
  - rules: `3`
  - event types: `2`
  - candidate rows: `9,080`
  - funnel rows: `24`
  - summary rows: `1,296`
  - stability rows: `352`
  - CSV lines including header: `9,081` / `25` / `1,297` / `353`
  - horizons: `1`, `3`, `6`, `12`
  - max midline event delay: `12` bars
  - quick invalidation window: `3` bars
- Scope:
  - profile `p30_c12_bollinger_on_adx_off`
  - rules `hold_3_inside`, `hold_6_inside`, and diagnostic
    `hold_3_inside_mid_50`
  - event types `mid_touch` and `mid_close_across`
  - event candle is the reindexed decision candle; labels start at
    `event_index + 1`
  - no entries, exits, scoring, sizing, paper side, favorable/adverse fields,
    strategy replacement, or live wiring
- Last run loaded `569,451` candles through `2026-06-01T23:59:59Z` and printed
  `strategy=empty trades=0`.
- Verification passed:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv -hold-inside-midline-reaction-audit -out-dir results/hold-inside-midline-reaction-audit`
  - `wc -l results/hold-inside-midline-reaction-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`

Context budget cleanup milestone:

- Added a durable context-budget rule to `AGENTS.md`, `memory/README.md`, and
  `memory/DECISIONS.md`.
- Compacted `memory/PROGRESS.md` from transcript-style history into this
  rolling snapshot plus milestone index.
- Updated `README.md` and `memory/NEXT_CODEX_BRIEF.md` so future sessions use
  docs as a task-scoped index instead of reading every historical doc by
  default.
- Verification:
  - `wc -l AGENTS.md README.md memory/*.md docs/*.md`: `memory/PROGRESS.md`
    is now `199` lines; counted project context files total `1,999` lines
    versus `3,506` before compaction.
  - `rg -n "docs/\\*|NEXT_CODEX_BRIEF|CODEX_BRIEF" README.md docs memory AGENTS.md`:
    no handoff requires reading `docs/*.md`; remaining `docs/*.md` mention is
    the required-file inventory in `docs/VERIFICATION.md`.
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`:
    passed.
  - `git diff --check`: passed.

Hold-inside midline transition review milestone:

- Review doc: `docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md`.
- Inputs: `results/hold-inside-midline-transition-audit/`; `7,988`
  candidate rows, `720` summary rows, `192` stability rows.
- Verdict: not entry-ready; broad `hold_3_inside`/`hold_6_inside` rows show
  split-stable midline touch/close-across behavior by `12` bars, but current
  midline labels are not entry context or strategy inputs.

Hold-inside midline transition audit milestone:

- Commit: `277ce49 Add hold-inside midline transition audit`.
- Added `-hold-inside-midline-transition-audit`; result dir
  `results/hold-inside-midline-transition-audit/`; `7,988` candidate rows,
  `720` summary rows, `192` stability rows; `strategy=empty trades=0`.

Hold-inside directional edge review milestone:

- Commit: `0601952 Review hold-inside directional edge audit`.
- Review doc: `docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md`.
- Inputs: `results/hold-inside-directional-edge-audit/`.
- Verdict: not entry-ready. `hold_3_inside` and `hold_6_inside` remain useful
  range-survival context but do not show a split-stable directional edge toward
  frozen range high or low.
- Durable rule added to `memory/DECISIONS.md`: do not promote
  `paper_side=toward_high` or `paper_side=toward_low` into entry context.
- Verification passed:
  - `wc -l results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `git diff --check`

Hold-inside directional edge audit milestone:

- Commit: `00b8844 Add hold-inside directional edge audit`.
- Added CLI flag `-hold-inside-directional-edge-audit`.
- Result directory: `results/hold-inside-directional-edge-audit/`.
- Audit size: `15,976` candidate rows, `624` summary rows, `168` stability
  rows; CSV lines including header `15,977` / `625` / `169`.
- Scope stayed non-trading; `strategy=empty trades=0`.
- Verification passed with `go test ./...`, audit run, `wc -l`, handoff
  reference check, and `git diff --check`.

Detector context refinement review milestone:

- Commit: `a4698bf Review detector context refinement audit`.
- Review doc: `docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md`.
- Inputs: `results/detector-context-refinement-audit/`.
- Verdict: delayed `hold_3_inside` and `hold_6_inside` are the leading context
  refinement and materially reduce quick invalidation and trend leakage, but
  are not promoted to entry context because the gain is conditioned survival
  and labels are regime-durability outcomes, not P&L.
- Verification passed with CSV line counts `113,825` / `641` / `161`,
  `go test ./...`, and `git diff --check`.

Detector context refinement audit milestone:

- Commit: `b980095 Add detector context refinement audit`.
- Added CLI flag `-detector-context-refinement-audit`.
- Result directory: `results/detector-context-refinement-audit/`.
- Audit size: `113,824` candidate rows, `640` summary rows, `160` stability
  rows; horizons `1`, `3`, `6`, `12`.
- Scope stayed non-trading; `strategy=empty trades=0`.

Detector durability sweep review milestone:

- Commit: `7bd2852 Review detector durability sweep`.
- Review doc: `docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md`.
- Inputs: `results/detector-durability-sweep/`.
- Verdict: no current `DefaultDetectorSweepProfiles` profile is approved as
  future entry context; `p30_c12_bollinger_on_adx_on` is diagnostic only.
- Verification passed with CSV line counts, `go test ./...`, and
  `git diff --check`.

Detector durability sweep milestone:

- Commit: `d70236b Add detector durability sweep`.
- Added CLI flag `-detector-durability-sweep`.
- Result directory: `results/detector-durability-sweep/`.
- Audit size: `304` broad rows, `9,088` slice rows, `76` stability rows.
- Scope stayed non-trading; `strategy=empty trades=0`.

Range regime durability review milestone:

- Commit: `bc81d9d Review range regime durability`.
- Review doc: `docs/RANGE_REGIME_DURABILITY_REVIEW.md`.
- Inputs: `results/range-regime-durability-audit/`.
- Verdict: current balanced detector regimes are not durable enough as context
  for future entry hypotheses. Refine detector/context before trigger work.
- Verification passed with CSV line counts, `go test ./...`, and
  `git diff --check`.

## Historical Milestone Index

- `7e1494b Add range regime durability audit`: added
  `-range-regime-durability-audit`; result dir
  `results/range-regime-durability-audit/`; `11,984` episode rows and `452`
  summary rows; `strategy=empty trades=0`.
- `b92207a Review compression breakout audit`: added
  `docs/COMPRESSION_BREAKOUT_REVIEW.md`; verdict not entry-ready; result dir
  `results/compression-breakout-audit/`.
- `be1e8d5 Add compression breakout audit mode`: added
  `-compression-breakout-audit`; `5,096` candidate rows and `24` summary rows;
  `strategy=empty trades=0`.
- `0256301 Review SR false-break reclaim timing`: added
  `docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md`; verdict not entry-ready;
  result dir `results/sr-false-break-reclaim-timing-audit/`.
- `1264a2f Add false-break reclaim timing audit`: added
  `-sr-false-break-reclaim-timing-audit`; `17,652` candidate rows and `24`
  summary rows; `strategy=empty trades=0`.
- `4e1a682 Add SR confirmation timing audit and review`: added
  `-sr-confirmation-timing-audit` and `docs/SR_CONFIRMATION_TIMING_REVIEW.md`;
  `9,692` candidate rows and `72` summary rows; verdict not entry-ready.
- `ed82756 Add SR rejection timing audit and review`: added
  `-sr-rejection-timing-audit` and `docs/SR_REJECTION_TIMING_REVIEW.md`; `968`
  candidate rows and `24` summary rows; verdict not entry-ready.
- `f60a8d9 Add entry readiness review gate`: added
  `docs/ENTRY_READINESS_REVIEW.md`; `internal/lab` coverage reached `99.8%`;
  first entries remained blocked pending non-trading timing evidence.
- `04edccc Use memory brief as canonical handoff`: removed duplicate root
  `CODEX_BRIEF.md`; `memory/NEXT_CODEX_BRIEF.md` became canonical.
- `1e26695 Add SR boundary inspection mode`: added `-sr-boundary-inspect`;
  `281,080` boundary events inspected and `192` comparison rows; `strategy=empty
  trades=0`.
- `9833525 Add SR boundary audit mode`: added `-sr-boundary-audit`; `281,080`
  boundary event rows and `192` boundary quality rows.
- `bdf1398 Add mods and next plan;`: pinned helper modules and documented
  adapter-only research helper boundaries in `docs/RESEARCH_HELPERS.md`.
- `6483e8c Add detector sweep audit mode`: added `-detector-sweep`; `19`
  profiles and `76` sweep rows; balanced baseline
  `p30_c12_bollinger_on_adx_off`.
- `e902a9e Initialize range strategy lab`: initialized standalone lab,
  tracked memory, and empty-strategy smoke path.

## Standard Closeout Checks

Use these unless the task has a narrower verifier:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```
