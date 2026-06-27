package lab

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRangeFirstOccupancyRotationGridAndBaselineConfig(t *testing.T) {
	cfg := DefaultFuturesRangeFirstOccupancyRotationV1OptimizationConfig()
	grid := rangeFirstOccupancyRotationGridConfigs(cfg)
	if len(grid) != 1152 {
		t.Fatalf("grid rows=%d, want 1152", len(grid))
	}
	baseline := rangeFirstOccupancyRotationBaselineGridConfig()
	if baseline.ConfigID != RangeFirstOccupancyRotationV1BaselineConfigID {
		t.Fatalf("baseline id=%q", baseline.ConfigID)
	}
	found := false
	for _, row := range grid {
		if row.ConfigID == baseline.ConfigID {
			found = true
			if row.Timeframe != RangeDiscoveryTimeframe1h || row.LookbackHours != 48 || row.MaxHoldBars != 12 || row.StopBufferWidth != 0.05 {
				t.Fatalf("baseline grid params mismatch: %+v", row)
			}
		}
	}
	if !found {
		t.Fatalf("baseline config not present in declared grid")
	}
}

func TestRangeFirstOccupancyRotationSourceAndResampleAcceptance(t *testing.T) {
	parent := make([]Candle, 12)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	manifest := SourceManifest{
		Path:             "btcusdt_futures_um_5m_fixture.csv",
		Venue:            "Binance",
		Product:          "Binance USDT-M futures",
		Symbol:           RangeUniverseSymbolBTCUSDT,
		Interval:         "5m",
		RowCount:         len(parent),
		FirstOpenTime:    parent[0].OpenTime.UTC().Format(timeLayout),
		LastOpenTime:     parent[len(parent)-1].OpenTime.UTC().Format(timeLayout),
		ValidationStatus: "accepted",
	}
	cfg := testRangeFirstOccupancyRotationConfig()
	cfg.Timeframes = []string{RangeDiscoveryTimeframe15m, RangeDiscoveryTimeframe1h}
	cfg.LookbackHours = []int{1, 48}
	cfg.OccupancyWindows = []int{2, 12}
	result, err := RunFuturesRangeFirstOccupancyRotationV1Optimization(parent, manifest, cfg, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SourceRows) != 1 || !result.SourceRows[0].SourceFactsPass {
		t.Fatalf("bad source rows: %+v", result.SourceRows)
	}
	if len(result.CoverageRows) != 2 {
		t.Fatalf("coverage rows=%d, want 2", len(result.CoverageRows))
	}
	for _, row := range result.CoverageRows {
		if !row.CoverageFactsPass || row.ValidationStatus != "accepted" {
			t.Fatalf("coverage should pass: %+v", row)
		}
	}

	manifest.ComparisonOnly = true
	manifest.Product = "Binance spot"
	if _, err := RunFuturesRangeFirstOccupancyRotationV1Optimization(parent, manifest, cfg, BacktestConfig{}, []Split{{Name: fullSplitName}}); err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures") {
		t.Fatalf("expected spot/comparison rejection, got %v", err)
	}
}

func TestRangeFirstOccupancyRotationEnvelopeExcludesSignalAndLongTrade(t *testing.T) {
	candles := []Candle{
		testCandle(0, 105, 110, 100, 105),
		testCandle(1, 106, 108, 101, 106),
		testCandle(2, 101, 104, 100, 101),
		testCandle(3, 102, 105, 101, 102),
		testCandle(4, 104, 200, 103, 104),
		testCandle(5, 103, 107, 102, 106),
	}
	grid := testRangeFirstOccupancyRotationGrid(Long)
	frame := testRangeFirstOccupancyRotationFrame(t, candles, grid)
	run, err := rangeFirstOccupancyRotationRunConfig(grid, frame, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, []Split{{Name: fullSplitName}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(run.signals) != 1 {
		t.Fatalf("signals=%d, want 1 skips=%+v", len(run.signals), run.skips)
	}
	signal := run.signals[0]
	if signal.Side != Long || signal.RangeHigh != 110 || signal.RangeLow != 100 {
		t.Fatalf("signal should use previous bars only: %+v", signal)
	}
	if signal.Stop != 99.5 || signal.Target != 106.6 || signal.MaxHoldBars != 2 {
		t.Fatalf("unexpected long geometry: stop=%v target=%v max_hold=%d", signal.Stop, signal.Target, signal.MaxHoldBars)
	}
	if len(run.trades) != 1 || run.trades[0].Reason != "take_profit" {
		t.Fatalf("trade mismatch: %+v", run.trades)
	}
}

func TestRangeFirstOccupancyRotationShortTrade(t *testing.T) {
	candles := []Candle{
		testCandle(0, 105, 110, 100, 105),
		testCandle(1, 104, 109, 101, 104),
		testCandle(2, 109, 110, 107, 109),
		testCandle(3, 108, 110, 106, 108),
		testCandle(4, 106, 107, 104, 106),
		testCandle(5, 107, 108, 103, 104),
	}
	grid := testRangeFirstOccupancyRotationGrid(Short)
	frame := testRangeFirstOccupancyRotationFrame(t, candles, grid)
	run, err := rangeFirstOccupancyRotationRunConfig(grid, frame, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, []Split{{Name: fullSplitName}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(run.signals) != 1 || run.signals[0].Side != Short {
		t.Fatalf("short signal mismatch: %+v", run.signals)
	}
	if run.signals[0].Stop != 110.5 || run.signals[0].Target != 103.4 {
		t.Fatalf("unexpected short geometry: %+v", run.signals[0])
	}
	if len(run.trades) != 1 || run.trades[0].Reason != "take_profit" {
		t.Fatalf("short trade mismatch: %+v", run.trades)
	}
}

func TestRangeFirstOccupancyRotationDualSideAndInvalidGeometrySkips(t *testing.T) {
	candles := []Candle{
		testCandle(0, 105, 110, 100, 105),
		testCandle(1, 101, 110, 100, 101),
		testCandle(2, 109, 110, 100, 109),
		testCandle(3, 105, 106, 104, 105),
		testCandle(4, 105, 106, 104, 105),
	}
	grid := testRangeFirstOccupancyRotationGrid(Long)
	grid.LookbackHours = 2
	grid.OccupancyWindow = 3
	grid.OccupancyMinFraction = 0.33
	grid.RecaptureLevel = 0.25
	grid.ConfigID = rangeFirstOccupancyRotationConfigID(grid)
	frame := testRangeFirstOccupancyRotationFrame(t, candles, grid)
	run, err := rangeFirstOccupancyRotationRunConfig(grid, frame, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, []Split{{Name: fullSplitName}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !rangeFirstOccupancyRotationHasSkip(run.skips, "ambiguous_dual_side_signal") {
		t.Fatalf("expected dual-side skip, skips=%+v", run.skips)
	}

	invalid := []Candle{
		testCandle(0, 105, 110, 100, 105),
		testCandle(1, 106, 108, 101, 106),
		testCandle(2, 101, 104, 100, 101),
		testCandle(3, 102, 105, 101, 102),
		testCandle(4, 104, 105, 103, 104),
		testCandle(5, 108, 109, 107, 108),
	}
	grid = testRangeFirstOccupancyRotationGrid(Long)
	frame = testRangeFirstOccupancyRotationFrame(t, invalid, grid)
	run, err = rangeFirstOccupancyRotationRunConfig(grid, frame, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, []Split{{Name: fullSplitName}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !rangeFirstOccupancyRotationHasSkip(run.skips, "entry_stop_target_invalid") {
		t.Fatalf("expected invalid geometry skip, skips=%+v", run.skips)
	}
}

func TestRangeFirstOccupancyRotationAlreadyInPositionSkip(t *testing.T) {
	candles := []Candle{
		testCandle(0, 105, 110, 100, 105),
		testCandle(1, 106, 108, 101, 106),
		testCandle(2, 101, 104, 100, 101),
		testCandle(3, 101, 105, 100, 101),
		testCandle(4, 104, 105, 103, 104),
		testCandle(5, 102, 105, 101, 102),
		testCandle(6, 104, 105, 103, 104),
		testCandle(7, 104, 107, 103, 106),
	}
	grid := testRangeFirstOccupancyRotationGrid(Long)
	grid.RecaptureLevel = 0.25
	frame := testRangeFirstOccupancyRotationFrame(t, candles, grid)
	run, err := rangeFirstOccupancyRotationRunConfig(grid, frame, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, []Split{{Name: fullSplitName}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(run.signals) < 2 {
		t.Fatalf("expected overlapping signal candidates, got %+v", run.signals)
	}
	if !rangeFirstOccupancyRotationHasSkip(run.skips, "already_in_position") {
		t.Fatalf("expected already-in-position skip, skips=%+v signals=%+v trades=%+v", run.skips, run.signals, run.trades)
	}
}

func TestRangeFirstOccupancyRotationRankingAndStopStates(t *testing.T) {
	rows := []FuturesRangeFirstOccupancyRotationV1GridRow{
		{ConfigID: "pass_15m", Timeframe: RangeDiscoveryTimeframe15m, MaxHoldBars: 8, PassesGate: true, RankScore: 10, TrainNetPnL: 100, TrainProfitFactor: 2, FullTrades: 250},
		{ConfigID: "pass_1h", Timeframe: RangeDiscoveryTimeframe1h, MaxHoldBars: 12, PassesGate: true, RankScore: 10, TrainNetPnL: 100, TrainProfitFactor: 2, FullTrades: 250},
	}
	rankings := rangeFirstOccupancyRotationRankingRows(rows)
	if rankings[0].ConfigID != "pass_1h" {
		t.Fatalf("1h tie-break should win: %+v", rankings)
	}
	source := []FuturesRangeFirstOccupancyRotationV1SourceRow{{SourceFactsPass: true, ValidationStatus: "accepted"}}
	coverage := []FuturesRangeFirstOccupancyRotationV1CoverageRow{{CoverageFactsPass: true, FuturesRangeDiscoveryCoverageRow: FuturesRangeDiscoveryCoverageRow{Complete: true, ValidationStatus: "accepted"}}}
	grid := []FuturesRangeFirstOccupancyRotationV1GridRow{{ConfigID: RangeFirstOccupancyRotationV1BaselineConfigID, SignalCount: 10}, {ConfigID: "pass", SignalCount: 10, PassesGate: true}}
	if got := FuturesRangeFirstOccupancyRotationV1OptimizationStopState(source, coverage, grid, []FuturesRangeFirstOccupancyRotationV1RankingRow{{ConfigID: "pass", PassesGate: true}}); got != RangeFirstStrategyV1StopStatePassedNeedsFixedReplaySpec {
		t.Fatalf("stop=%s, want passed", got)
	}
	grid = []FuturesRangeFirstOccupancyRotationV1GridRow{{ConfigID: RangeFirstOccupancyRotationV1BaselineConfigID, SignalCount: 10, PassesGate: false}}
	if got := FuturesRangeFirstOccupancyRotationV1OptimizationStopState(source, coverage, grid, nil); got != RangeFirstStrategyV1StopStateOptimizerFailedNoReplay {
		t.Fatalf("stop=%s, want optimizer failed", got)
	}
	grid = []FuturesRangeFirstOccupancyRotationV1GridRow{{ConfigID: RangeFirstOccupancyRotationV1BaselineConfigID, SignalCount: 0}}
	if got := FuturesRangeFirstOccupancyRotationV1OptimizationStopState(source, coverage, grid, nil); got != RangeFirstStrategyV1StopStateNoValidSignals {
		t.Fatalf("stop=%s, want no valid signals", got)
	}
}

func testRangeFirstOccupancyRotationConfig() FuturesRangeFirstOccupancyRotationV1OptimizationConfig {
	cfg := DefaultFuturesRangeFirstOccupancyRotationV1OptimizationConfig()
	cfg.SkipSourceFactCheck = true
	cfg.SkipCoverageCountCheck = true
	cfg.Timeframes = []string{RangeDiscoveryTimeframe1h}
	cfg.LookbackHours = []int{1, 48}
	cfg.MaxWidthPcts = []float64{0.50}
	cfg.OccupancyWindows = []int{2, 12}
	cfg.OccupancyZoneLevels = []float64{0.25}
	cfg.OccupancyMinFractions = []float64{0.50}
	cfg.RecaptureLevels = []float64{0.33}
	cfg.TargetLevels = []float64{0.66}
	cfg.MaxHoldBars = []int{2}
	cfg.StopBufferWidths = []float64{0.05}
	cfg.SideModes = []string{RangeDiscoverySideAll}
	cfg.MinTrainTrades = 1
	cfg.MinOOSTrades = 1
	cfg.MinRecentTrades = 1
	cfg.MinFullTrades = 1
	return cfg
}

func testRangeFirstOccupancyRotationGrid(side Direction) FuturesRangeFirstOccupancyRotationV1GridConfig {
	grid := FuturesRangeFirstOccupancyRotationV1GridConfig{
		Timeframe:            RangeDiscoveryTimeframe1h,
		LookbackHours:        4,
		MaxWidthPct:          0.50,
		OccupancyWindow:      3,
		OccupancyZoneLevel:   0.25,
		OccupancyMinFraction: 0.60,
		RecaptureLevel:       0.33,
		TargetLevel:          0.66,
		MaxHoldBars:          2,
		StopBufferWidth:      0.05,
		SideMode:             RangeDiscoverySideAll,
	}
	if side == Short {
		grid.ConfigID = "test_short"
	} else {
		grid.ConfigID = "test_long"
	}
	return grid
}

func testRangeFirstOccupancyRotationFrame(t *testing.T, candles []Candle, grid FuturesRangeFirstOccupancyRotationV1GridConfig) rangeFirstOccupancyFrameCache {
	t.Helper()
	envelopes, err := rangeFirstOccupancyRotationPrecomputeEnvelopes(candles, grid.Timeframe, grid.LookbackHours)
	if err != nil {
		t.Fatal(err)
	}
	frame := rangeFirstOccupancyFrameCache{
		candles:         candles,
		splitsByIndex:   make([]string, len(candles)),
		envelopes:       map[int][]rangeFirstOccupancyEnvelope{grid.LookbackHours: envelopes},
		occupancyCounts: map[string][]rangeFirstOccupancyCounts{},
	}
	for i := range candles {
		frame.splitsByIndex[i] = fullSplitName
	}
	frame.occupancyCounts[rangeFirstOccupancyCountsKey(grid.LookbackHours, grid.OccupancyWindow)] = rangeFirstOccupancyRotationPrecomputeOccupancy(candles, envelopes, grid.OccupancyWindow, grid.OccupancyZoneLevel)
	return frame
}

func rangeFirstOccupancyRotationHasSkip(rows []FuturesRangeFirstOccupancyRotationV1SkipRow, reason string) bool {
	for _, row := range rows {
		if row.Reason == reason && row.Count > 0 {
			return true
		}
	}
	return false
}

func TestRangeFirstOccupancyRotationSourcePathRejection(t *testing.T) {
	dir := t.TempDir()
	parent := make([]Candle, 12)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	path := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_v1.csv", parent)
	candles, manifest, err := LoadResearchSourceCSV(path, SourceValidationOptions{Product: SourceProductBinanceUSDMFutures})
	if err != nil {
		t.Fatal(err)
	}
	cfg := testRangeFirstOccupancyRotationConfig()
	cfg.SkipSourceFactCheck = false
	cfg.ApprovedSourcePath = filepath.Join(dir, "other_btcusdt_futures_um_5m.csv")
	if _, err := RunFuturesRangeFirstOccupancyRotationV1Optimization(candles, manifest, cfg, BacktestConfig{}, []Split{{Name: fullSplitName}}); err == nil || !strings.Contains(err.Error(), "approved path") {
		t.Fatalf("expected approved-path rejection, got %v", err)
	}
}
